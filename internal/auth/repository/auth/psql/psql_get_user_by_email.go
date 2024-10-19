package auth_repository_psql

import (
	"database/sql"

	auth_models "github.com/lameroot/msa-messenger/internal/auth/models"
)

func (s *PostgresAuthRepository) GetUserByEmail(email string) (*auth_models.User, error) {
	user := &auth_models.User{}
	err := s.db.QueryRow("SELECT id, email, password_hash, nickname, avatar_url, info FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Email, &user.Password, &user.Nickname, &user.AvatarURL, &user.Info)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
