package auth

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	db                   *sql.DB
	jwtSecret            []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewAuthService(db *sql.DB, jwtSecret string, accessTokenDuration, refreshTokenDuration time.Duration) *AuthService {
	return &AuthService{
		db:                   db,
		jwtSecret:            []byte(jwtSecret),
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (s *AuthService) Register(email, password string) (*User, error) {
	existingUser, _ := GetUserByEmail(s.db, email)
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	return CreateUser(s.db, email, password)
}

func (s *AuthService) Authenticate(email, password string) (*TokenPair, error) {
	user, err := GetUserByEmail(s.db, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !CheckPasswordHash(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	return s.generateTokenPair(strconv.FormatInt(user.ID, 10))
}

func (s *AuthService) RefreshToken(refreshToken string) (*TokenPair, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	userID := claims.Subject

	return s.generateTokenPair(userID)
}

func (s *AuthService) VerifyToken(accessToken string) (string, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid access token")
	}

	return claims.Subject, nil
}

func (s *AuthService) generateTokenPair(userID string) (*TokenPair, error) {
	accessToken, err := s.generateToken(userID, s.accessTokenDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(userID, s.refreshTokenDuration)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) generateToken(userID string, duration time.Duration) (string, error) {
	claims := jwt.StandardClaims{
		Subject:   userID,
		ExpiresAt: time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
