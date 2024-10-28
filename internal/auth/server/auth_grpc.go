package auth_grpc

import (
	"context"

	auth_usecase "github.com/lameroot/msa-messenger/internal/auth/usecase"
	auth_proto "github.com/lameroot/msa-messenger/pkg/api/auth"
)

type server struct {
	auth_proto.UnimplementedTokenVerifyServiceServer
	tokenRepository auth_usecase.TokenRepository
}

func NewServer(tokenRepository auth_usecase.TokenRepository) *server {
	return &server{tokenRepository: tokenRepository}
}

func (s *server) VerifyToken(context context.Context, tokenRequest *auth_proto.TokenVerificationRequest) (*auth_proto.TokenVerificationResponse, error) {
	JWTClaims, err := s.tokenRepository.ValidateAccessToken(tokenRequest.Token)
	if err != nil {
		return &auth_proto.TokenVerificationResponse{Verified: false}, err
	}
	return &auth_proto.TokenVerificationResponse{
		UserId:         JWTClaims.UserID.String(),
		Verified:       true,
		ExpirationTime: JWTClaims.ExpiresAt.Unix(),
	}, nil
}
