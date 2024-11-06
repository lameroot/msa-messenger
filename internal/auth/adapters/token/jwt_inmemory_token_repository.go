package adapters_token

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	auth_models "github.com/lameroot/msa-messenger/internal/auth/models"
	auth_usecase "github.com/lameroot/msa-messenger/internal/auth/usecase"
)

type JwtInMemoryTokenRepository struct {
	tokenConfig auth_models.TokenConfig
}

func NewJwtInMemoryTokenRepository(tokenConfig auth_models.TokenConfig) auth_usecase.TokenRepository {
	return &JwtInMemoryTokenRepository{
		tokenConfig: tokenConfig,
	}
}

func (r *JwtInMemoryTokenRepository) CreateAccessToken(user *auth_models.User) (string, error) {
	return r.createToken(user, r.tokenConfig.AccessTokenExpiry, r.tokenConfig.AccessTokenSecret)
}

func (r *JwtInMemoryTokenRepository) CreateRefreshToken(user *auth_models.User) (string, error) {
	return r.createToken(user, r.tokenConfig.RefreshTokenExpiry, r.tokenConfig.RefreshTokenSecret)
}

func (r *JwtInMemoryTokenRepository) createToken(user *auth_models.User, expiry time.Duration, secret string) (string, error) {
	claims := &auth_models.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (r *JwtInMemoryTokenRepository) ValidateAccessToken(tokenString string) (*auth_models.JWTClaims, error) {
	return r.validateToken(tokenString, r.tokenConfig.AccessTokenSecret)
}
func (r *JwtInMemoryTokenRepository) ValidateRefreshToken(tokenString string) (*auth_models.JWTClaims, error) {
	return r.validateToken(tokenString, r.tokenConfig.RefreshTokenSecret)
}

func (r *JwtInMemoryTokenRepository) validateToken(tokenString, secret string) (*auth_models.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &auth_models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*auth_models.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (r *JwtInMemoryTokenRepository) Close() {
	log.Default().Println("Close JwtInMemoryTokenRepository")
}
