package messaging_repository_psql

import (
	"database/sql"

	messaging_usecase "github.com/lameroot/msa-messenger/internal/messaging/usecase"
	_ "github.com/lib/pq"
)

/*
*
create table messages (

	id VARCHAR(36) PRIMARY KEY,
	sender_id VARCHAR(36) REFERENCES users(id),
	receiver_id VARCHAR(36) REFERENCES users(id),
	content TEXT NOT NULL,
	is_read BOOLEAN DEFAULT FALSE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
*/
type PostgresMessagingRepository struct {
	db *sql.DB
}

func NewPostgresMessagingRepository(dbURL string) (messaging_usecase.PersistentRepository, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresMessagingRepository{db}, nil
}
