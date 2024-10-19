package auth_repository_psql

import (
	"database/sql"

	auth_models "github.com/lameroot/msa-messenger/internal/auth/models"
)

func (s *PostgresAuthRepository) GetUserByNickname(nickname string) (*auth_models.User, error) {
	user := &auth_models.User{}
	err := s.db.QueryRow("SELECT id, email, password_hash, nickname, avatar_url, info FROM users WHERE nickname = $1", nickname).
		Scan(&user.ID, &user.Email, &user.Password, &user.Nickname, &user.AvatarURL, &user.Info)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
