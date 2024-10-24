package user_repository_psql

import (
	"database/sql"

	user_usecase "github.com/lameroot/msa-messenger/internal/user/usecase"
	_ "github.com/lib/pq"
)

/*
*
CREATE TABLE friendships (

	id SERIAL PRIMARY KEY,
	user_id INTEGER REFERENCES users(id),
	friend_id INTEGER REFERENCES users(id),
	status VARCHAR(20) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(user_id, friend_id)

);

CREATE INDEX idx_friendships_user_id ON friendships(user_id);
CREATE INDEX idx_friendships_friend_id ON friendships(friend_id);
*/
type PosgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(dbURL string) (user_usecase.PersistentRepository, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PosgresUserRepository{db: db}, nil
}
