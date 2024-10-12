package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	UnimplementedAuthServiceServer
	authService *AuthService
}

func NewGRPCServer(authService *AuthService) *GRPCServer {
	return &GRPCServer{authService: authService}
}

func (s *GRPCServer) VerifyToken(ctx context.Context, req *VerifyTokenRequest) (*VerifyTokenResponse, error) {
	userID, err := s.authService.VerifyToken(req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token: %v", err)
	}

	return &VerifyTokenResponse{
		Valid:  true,
		UserId: userID,
	}, nil
}
