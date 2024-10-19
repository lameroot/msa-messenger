package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	adapters_token "github.com/lameroot/msa-messenger/internal/auth/adapters/token"
	auth_http "github.com/lameroot/msa-messenger/internal/auth/controller/http"
	auth_models "github.com/lameroot/msa-messenger/internal/auth/models"
	auth_repository_psql "github.com/lameroot/msa-messenger/internal/auth/repository/auth/psql"
	auth_usecase "github.com/lameroot/msa-messenger/internal/auth/usecase"

	"github.com/joho/godotenv"

	_ "github.com/lameroot/msa-messenger/docs" // This is where Swag has generated docs.go
)

// @title Authentication Service API
// @version 1.0
// @description This is the API for the authentication service of MSA Messenger.
// @host localhost:8080
// @BasePath /
func main() {
	// Init config
	dir, _ := os.Getwd()
	envPath := filepath.Join(dir, "configs", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	// Initialize the auth service
	tokenConfig := auth_models.TokenConfig{
		AccessTokenSecret:  os.Getenv("JWT_ACCESS_TOKEN_SECRET"),
		RefreshTokenSecret: os.Getenv("JWT_REFRESH_TOKEN_SECRET"),
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
	}
	tokenRepository := adapters_token.NewJwtInMemoryTokenRepository(tokenConfig)

	// Init database
	dbURL := os.Getenv("DB_POSTGRES_URL")
	persistentRepository, err := auth_repository_psql.NewPostgresAuthRepository(dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize auth service: %v", err)
	}

	authService := auth_usecase.NewAuthUseCase(tokenRepository, persistentRepository)

	// Initialize the auth handler
	authHandler := auth_http.NewAuthHandler(authService)

	// Create http router
	router := auth_http.NewRouter(authHandler)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start auth server: %v", err)
	}
}
