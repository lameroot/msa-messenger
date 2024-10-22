package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	messaging_notification "github.com/lameroot/msa-messenger/internal/messaging/adapters/notification"
	messaging_http "github.com/lameroot/msa-messenger/internal/messaging/controller"
	messaging_repository_psql "github.com/lameroot/msa-messenger/internal/messaging/repository/messaging/psql"
	messaging_usecase "github.com/lameroot/msa-messenger/internal/messaging/usecase"
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
	persistentRepository, err := messaging_repository_psql.NewPostgresMessagingRepository(dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize auth service: %v", err)
	}

	// Init notification service
	notificationService := messaging_notification.NewInMemmoryNotificationService()

	// Create grpc client
	hostPortAuthGrpcServer := os.Getenv("AUTH_GRPC_HOST_PORT")
	conn, err := grpc.NewClient(hostPortAuthGrpcServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	log.Default().Println("Grpc auth client successfully created and connected to host ", hostPortAuthGrpcServer)
	authClient := auth_proto.NewTokenVerifyServiceClient(conn)

	messagingService := messaging_usecase.NewMessagingService(persistentRepository, notificationService)
	messagingHandler := messaging_http.NewMessagingHandler(messagingService)
	messagingRouter := messaging_http.NewRouter(messagingHandler, &authClient)

	// Start the server
	if err := messagingRouter.Run(":8082"); err != nil {
		log.Fatalf("Failed to start messaing server: %v", err)
	}

}
