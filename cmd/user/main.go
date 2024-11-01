package main

import (
	"log"
	"os"
	"path/filepath"

	//user_http "msa-messenger/internal/user/controller/http"

	"github.com/joho/godotenv"

	user_http "github.com/lameroot/msa-messenger/internal/user/controller/http"
	user_repository_psql "github.com/lameroot/msa-messenger/internal/user/repository/user/psql"
	user_usecase "github.com/lameroot/msa-messenger/internal/user/usecase"
	auth_verify_service "github.com/lameroot/msa-messenger/pkg/auth"
)

func loadEnv() {
	// Init config
	dir, _ := os.Getwd()
	envPath := filepath.Join(dir, "configs", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Default().Print("Error loading .env file: ", err)
	}
}

func main() {
	// Init config
	log.Default().Print("Start load configs")
	loadEnv()
	log.Default().Print("Loaded configs: ", os.Getenv("DB_POSTGRES_URL"))

	// Init database
	dbURL := os.Getenv("DB_POSTGRES_URL")
	persistentRepository, err := user_repository_psql.NewPostgresUserRepository(dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize auth service: %v", err)
	}

	userService := user_usecase.NewUserService(persistentRepository)

	// Create grpc client
	auth_verify_service, err := auth_verify_service.NewAuthVerifyService(nil)
	if err != nil {
		log.Fatalf("Failed to create AuthVerifyService: %v", err)
	}
	defer auth_verify_service.Close()

	// Create user handler
	userHandler := user_http.NewUserHandler(userService)

	// Create server
	router := user_http.NewRouter(userHandler, auth_verify_service)

	// Start the server
	hostPortUserHttpServer := os.Getenv("USER_HTTP_HOST_PORT")
	if err := router.Run(hostPortUserHttpServer); err != nil {
		log.Fatalf("Failed to start user server: %v", err)
	}
}
