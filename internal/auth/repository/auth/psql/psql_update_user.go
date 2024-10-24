package auth_repository_psql

import (
	"log"

	"github.com/google/uuid"
	auth_models "github.com/lameroot/msa-messenger/internal/auth/models"
)

func (s *PostgresAuthRepository) UpdateUser(id uuid.UUID, avatart_url, info string) (*auth_models.User, error) {
	_, err := s.db.Exec("UPDATE users set avatar_url = $1, info = $2 where id = $3", avatart_url, info, id)
	if err != nil {
		log.Default().Printf("failed to update user: %v", err)
		return nil, err
	}

	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
