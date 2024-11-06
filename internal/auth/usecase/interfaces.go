package auth_usecase

import (
	"github.com/google/uuid"
	auth_models "github.com/lameroot/msa-messenger/internal/auth/models"
)

type (
	TokenRepository interface {
		CreateAccessToken(user *auth_models.User) (string, error)
		CreateRefreshToken(user *auth_models.User) (string, error)
		ValidateAccessToken(tokenString string) (*auth_models.JWTClaims, error)
		ValidateRefreshToken(tokenString string) (*auth_models.JWTClaims, error)
		Close()
	}
	PersistentRepository interface {
		GetUserByEmail(email string) (*auth_models.User, error)
		GetUserByID(id uuid.UUID) (*auth_models.User, error)
		GetUserByNickname(nickname string) (*auth_models.User, error)
		UpdateUser(id uuid.UUID, avatart_url, info string) (*auth_models.User, error)
		SaveUser(user *auth_models.User) (*auth_models.User, error)
		Close()
	}
)
