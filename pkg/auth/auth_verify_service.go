package auth_verify_service

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	auth_proto "github.com/lameroot/msa-messenger/pkg/api/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthVerifyService struct {
	authClient auth_proto.TokenVerifyServiceClient
	conn       *grpc.ClientConn
}

func NewAuthVerifyService() (*AuthVerifyService, error) {
	hostPortAuthGrpcServer, exists := os.LookupEnv("AUTH_GRPC_SERVER")
	if !exists {
		log.Fatalf("Not found variable AUTH_GRPC_SERVER")
		return nil, fmt.Errorf("not found variable AUTH_GRPC_SERVER")
	}
	conn, err := grpc.NewClient(hostPortAuthGrpcServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil, fmt.Errorf("failed to connect: %v", err)
	}
	log.Default().Println("Grpc auth client successfully created and connected to host ", hostPortAuthGrpcServer)
	authClient := auth_proto.NewTokenVerifyServiceClient(conn)

	return &AuthVerifyService{
		authClient: authClient,
		conn: conn,
	}, nil
}

func (a *AuthVerifyService) Close() error {
	return a.conn.Close()
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