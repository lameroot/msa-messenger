package auth_verify_service

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	auth_proto "github.com/lameroot/msa-messenger/pkg/api/auth"
	"github.com/sony/gobreaker/v2"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type AuthVerifyService struct {
	authClient auth_proto.TokenVerifyServiceClient
	conn       *grpc.ClientConn
}

type AuthVerifyConfig struct {
	RateLimitCountRequestInDuration time.Duration
	RateLimitBurst                  int
	СircuitBreakerName              string
	СircuitBreakerMaxRequests       uint32
	СircuitBreakerInterval          time.Duration
	СircuitBreakerTimeout           time.Duration
	СircuitBreakerReadyToTrip       func(counts gobreaker.Counts) bool
	СircuitBreakerOnStateChange     func(name string, from gobreaker.State, to gobreaker.State)
	СircuitBreakerIsSuccessful      func(err error) bool
	RetryMaxRetries                 int
	RetryBaseDelay                  time.Duration
}

var DefaultAuthVerifyConfig = AuthVerifyConfig{
	RateLimitCountRequestInDuration: 10 * time.Second,
	RateLimitBurst:                  1,
	СircuitBreakerName:              "DefaultСircuitBreaker",
	СircuitBreakerMaxRequests:       3,
	СircuitBreakerInterval:          5 * time.Second,
	СircuitBreakerTimeout:           10 * time.Second,
	СircuitBreakerReadyToTrip: func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	},
	СircuitBreakerOnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
		fmt.Printf("Circuit Breaker '%s' changed from '%s' to '%s'\n", name, from, to)
	},
	RetryMaxRetries: 3,
	RetryBaseDelay: 100 * time.Millisecond,
}

func NewAuthVerifyService(authVerifyConfig *AuthVerifyConfig) (*AuthVerifyService, error) {
	if authVerifyConfig == nil {
		authVerifyConfig = &DefaultAuthVerifyConfig
	}
	hostPortAuthGrpcServer, exists := os.LookupEnv("AUTH_GRPC_SERVER")
	if !exists {
		log.Fatalf("Not found variable AUTH_GRPC_SERVER")
		return nil, fmt.Errorf("not found variable AUTH_GRPC_SERVER")
	}

	conn, err := grpc.NewClient(hostPortAuthGrpcServer,
		buildRateLimiter(authVerifyConfig),
		buildCircuitBreaker(authVerifyConfig),
		buildRetry(authVerifyConfig),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil, fmt.Errorf("failed to connect: %v", err)
	}
	log.Default().Println("Grpc auth client successfully created and connected to host ", hostPortAuthGrpcServer)
	authClient := auth_proto.NewTokenVerifyServiceClient(conn)

	return &AuthVerifyService{
		authClient: authClient,
		conn:       conn,
	}, nil
}

func (a *AuthVerifyService) Close() error {
	return a.conn.Close()
}

func buildRateLimiter(authVerifyConfig *AuthVerifyConfig) grpc.DialOption {
	limiter := rate.NewLimiter(rate.Every(authVerifyConfig.RateLimitCountRequestInDuration), authVerifyConfig.RateLimitBurst)
	rateLimiterInterceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if err := limiter.Wait(ctx); err != nil {
			return err
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
	return grpc.WithUnaryInterceptor(rateLimiterInterceptor)
}

func buildCircuitBreaker(authVerifyConfig *AuthVerifyConfig) grpc.DialOption {
	st := gobreaker.Settings{
		Name:          authVerifyConfig.СircuitBreakerName,
		MaxRequests:   authVerifyConfig.СircuitBreakerMaxRequests,
		Interval:      authVerifyConfig.СircuitBreakerInterval,
		Timeout:       authVerifyConfig.СircuitBreakerTimeout,
		ReadyToTrip:   authVerifyConfig.СircuitBreakerReadyToTrip,
		OnStateChange: DefaultAuthVerifyConfig.СircuitBreakerOnStateChange,
	}
	cb := gobreaker.NewCircuitBreaker[interface{}](st)

	unary := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		_, err := cb.Execute(func() (any, error) {
			err := invoker(ctx, method, req, reply, cc, opts...)
			var zero any
			if err != nil {
				// Преобразуем gRPC ошибки в обычные ошибки для Circuit Breaker
				if s, ok := status.FromError(err); ok {
					switch s.Code() {
					case codes.DeadlineExceeded, codes.Unavailable, codes.Internal:
						return zero, err
					}
				}
			}
			return zero, err
		})
		return err
	}

	return grpc.WithUnaryInterceptor(unary)
}

func buildRetry(authVerifyConfig *AuthVerifyConfig) grpc.DialOption {
	unary := func (ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var lastErr error
		for attempt := 0; attempt < authVerifyConfig.RetryMaxRetries; attempt++ {
			err := invoker(ctx, method, req, reply, cc, opts...)
			if err == nil {
				return nil
			}

			lastErr = err
			st, ok := status.FromError(err)
			if !ok {
				// Если ошибка не является gRPC статусом, не повторяем
				return err
			}

			// Повторяем только для определенных кодов ошибок
			if st.Code() != codes.Unavailable && st.Code() != codes.DeadlineExceeded {
				return err
			}

			// Экспоненциальная задержка перед повторной попыткой
			delay := authVerifyConfig.RetryBaseDelay * time.Duration(1<<uint(attempt))
			select {
				case <-ctx.Done():
				return ctx.Err()
				case <-time.After(delay):
			}
		}

		return lastErr
	}
	return grpc.WithUnaryInterceptor(unary)
}

func (a *AuthVerifyService) VerifyToken(ctx context.Context, token string) (*uuid.UUID, *ErrorVerifyToken) {
	req := &auth_proto.TokenVerificationRequest{
		Token: token,
	}

	verifyResponse, err := a.authClient.VerifyToken(ctx, req)
	if err != nil {
		return nil, &ErrorVerifyToken{Error: "Unauthorized :" + err.Error()}
	}
	if !verifyResponse.Verified {
		return nil, &ErrorVerifyToken{Error: "Unauthorized"}
	}
	IDUser, err := uuid.Parse(verifyResponse.UserId)
	if err != nil {
		return nil, &ErrorVerifyToken{Error: "Unauthorized: " + err.Error()}
	}
	return &IDUser, nil
}

type ErrorVerifyToken struct {
	Error string
}
