package auth_usecase

import (
	"errors"
	"log"

	"github.com/google/uuid"
	auth_models "github.com/lameroot/msa-messenger/internal/auth/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	tokenRepository      TokenRepository
	persistentRepository PersistentRepository
}

func NewAuthUseCase(tokenRepository TokenRepository, persistentRepository PersistentRepository) *AuthService {
	return &AuthService{
		tokenRepository:      tokenRepository,
		persistentRepository: persistentRepository,
	}
}

func (a *AuthService) Register(email, password, nickname string) (*auth_models.User, error) {
	existingUser, err := a.persistentRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Check if nickname is already taken
	existingUser, err = a.persistentRepository.GetUserByNickname(nickname)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("nickname already taken")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	// Create user
	user := &auth_models.User{
		ID:       uuid.New(),
		Email:    email,
		Password: string(hashedPassword),
		Nickname: nickname,
	}

	// Save user to database
	user, err = a.persistentRepository.SaveUser(user)
	if err != nil {
		log.Default().Printf("Error save user %s", err)
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Authenticate(email, password string) (*auth_models.AuthResponse, error) {
	user, err := s.CheckPassword(email, password)
	if err != nil {
		return nil, err
	}
	token, err := s.GenerateTokens(user)
	if err != nil {
		return nil, err
	}
	return &auth_models.AuthResponse{
		User:  *user,
		Token: *token,
	}, nil
}

func (s *AuthService) CheckPassword(email, password string) (*auth_models.User, error) {
	user, err := s.persistentRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*auth_models.Token, error) {
	claims, err := s.tokenRepository.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, err
	}

	user, err := s.persistentRepository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return s.GenerateTokens(user)
}

func (s *AuthService) GenerateTokens(user *auth_models.User) (*auth_models.Token, error) {
	accessToken, err := s.tokenRepository.CreateAccessToken(user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.tokenRepository.CreateRefreshToken(user)
	if err != nil {
		return nil, err
	}
	return &auth_models.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) GetUserByEmail(email string) (*auth_models.User, error) {
	return s.persistentRepository.GetUserByEmail(email)
}

func (s *AuthService) GetUserByID(id uuid.UUID) (*auth_models.User, error) {
	return s.persistentRepository.GetUserByID(id)
}

func (s *AuthService) GetUserByNickname(nickname string) (*auth_models.User, error) {
	return s.persistentRepository.GetUserByNickname(nickname)
}

func (s *AuthService) UpdateUser(id uuid.UUID, avatart_url, info string) (*auth_models.User, error) {
	return s.persistentRepository.UpdateUser(id, avatart_url, info)
}
