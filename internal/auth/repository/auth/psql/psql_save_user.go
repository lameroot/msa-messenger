package auth_repository_psql

import auth_models "github.com/lameroot/msa-messenger/internal/auth/models"

func (s *PostgresAuthRepository) SaveUser(user *auth_models.User) (*auth_models.User, error) {
	_, err := s.db.Exec("INSERT INTO users (id, email, password_hash, nickname, avatar_url, info) VALUES ($1, $2, $3, $4, $5, $6)",
		user.ID, user.Email, user.Password, user.Nickname, user.AvatarURL, user.Info)
	if err != nil {
		return user, err
	}

	return user, nil
}
