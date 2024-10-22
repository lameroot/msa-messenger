package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	user_http "github.com/lameroot/msa-messenger/internal/user/controller/http"
	user_repository_psql "github.com/lameroot/msa-messenger/internal/user/repository/user/psql"
	user_usecase "github.com/lameroot/msa-messenger/internal/user/usecase"
	auth_proto "github.com/lameroot/msa-messenger/pkg/api/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	// Create grpc client
	hostPortAuthGrpcServer := os.Getenv("AUTH_GRPC_HOST_PORT")
	conn, err := grpc.NewClient(hostPortAuthGrpcServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	log.Default().Println("Grpc auth client successfully created and connected to host ", hostPortAuthGrpcServer)
	authClient := auth_proto.NewTokenVerifyServiceClient(conn)

	// Create user handler
	userHandler := user_http.NewUserHandler(userService)

	// Create server
	router := user_http.NewRouter(userHandler, &authClient)

	// Start the server
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to start user server: %v", err)
	}
}
