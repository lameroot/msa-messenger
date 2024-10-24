package auth_repository_psql

import (
	"database/sql"

	_ "github.com/lib/pq"

	auth_usecase "github.com/lameroot/msa-messenger/internal/auth/usecase"
)

type PostgresAuthRepository struct {
	db *sql.DB
}

func NewPostgresAuthRepository(dbURL string) (auth_usecase.PersistentRepository, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresAuthRepository{db: db}, nil
}
