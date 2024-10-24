package main

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	adapters_token "github.com/lameroot/msa-messenger/internal/auth/adapters/token"
	auth_http "github.com/lameroot/msa-messenger/internal/auth/controller/http"
	auth_models "github.com/lameroot/msa-messenger/internal/auth/models"
	auth_repository_psql "github.com/lameroot/msa-messenger/internal/auth/repository/auth/psql"
	auth_grpc "github.com/lameroot/msa-messenger/internal/auth/server"
	auth_usecase "github.com/lameroot/msa-messenger/internal/auth/usecase"
	auth_proto "github.com/lameroot/msa-messenger/pkg/api/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/joho/godotenv"

	_ "github.com/lameroot/msa-messenger/docs" // This is where Swag has generated docs.go
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

// @title Authentication Service API
// @version 1.0
// @description This is the API for the authentication service of MSA Messenger.
// @host localhost:8080
// @BasePath /
func main() {
	// Init config
	loadEnv()
	log.Default().Print("Loaded configs: ", os.Getenv("DB_POSTGRES_URL"))

	// Initialize the auth service
	tokenConfig := auth_models.TokenConfig{
		AccessTokenSecret:  os.Getenv("JWT_ACCESS_TOKEN_SECRET"),
		RefreshTokenSecret: os.Getenv("JWT_REFRESH_TOKEN_SECRET"),
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
	}
	tokenRepository := adapters_token.NewJwtInMemoryTokenRepository(tokenConfig)
	log.Default().Print("Loaded JWT config")

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

	var wg sync.WaitGroup
	wg.Add(1)
	// Start the server
	hostPortAuthHttpServer := os.Getenv("AUTH_HTTP_HOST_PORT")
	go func() {
		defer wg.Done()

		if err := router.Run(hostPortAuthHttpServer); err != nil {
			log.Fatalf("Failed to start auth server: %v", err)
		}
	}()

	// Init gRPC Server
	hostPortAuthGrpcServer := os.Getenv("AUTH_GRPC_HOST_PORT")
	wg.Add(1)
	go func() {
		defer wg.Done()

		lis, err := net.Listen("tcp", hostPortAuthGrpcServer)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		impl := auth_grpc.NewServer(tokenRepository)
		server := grpc.NewServer()
		auth_proto.RegisterTokenVerifyServiceServer(server, impl)
		reflection.Register(server)
		log.Printf("server listening at %v", lis.Addr())
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	wg.Wait()

}
