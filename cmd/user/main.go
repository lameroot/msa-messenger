package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	user_http "github.com/lameroot/msa-messenger/internal/user/controller/http"
	user_repository_psql "github.com/lameroot/msa-messenger/internal/user/repository/user/psql"
	user_usecase "github.com/lameroot/msa-messenger/internal/user/usecase"
)

func main() {
	// Init config
	dir, _ := os.Getwd()
	envPath := filepath.Join(dir, "configs", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	// Init database
	dbURL := os.Getenv("DB_POSTGRES_URL")
	persistentRepository, err := user_repository_psql.NewPostgresUserRepository(dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize auth service: %v", err)
	}

	userService := user_usecase.NewUserService(persistentRepository)

	// Create user handler
	userHandler := user_http.NewUserHandler(userService)

	// Create server
	router := user_http.NewRouter(userHandler)

	// Start the server
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to start user server: %v", err)
	}
}
