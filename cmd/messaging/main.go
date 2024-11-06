package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/ds248a/closer"
	"github.com/joho/godotenv"
	messaging_notification "github.com/lameroot/msa-messenger/internal/messaging/adapters/notification"
	messaging_http "github.com/lameroot/msa-messenger/internal/messaging/controller"
	messaging_repository_psql "github.com/lameroot/msa-messenger/internal/messaging/repository/messaging/psql"
	messaging_usecase "github.com/lameroot/msa-messenger/internal/messaging/usecase"
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
	closer.NewCloser()
	// Init config
	log.Default().Print("Start load configs for messaging")
	loadEnv()
	log.Default().Print("Loaded configs: ", os.Getenv("DB_POSTGRES_URL"))

	// Init database
	dbURL := os.Getenv("DB_POSTGRES_URL")
	persistentRepository, err := messaging_repository_psql.NewPostgresMessagingRepository(dbURL)
	closer.Add(persistentRepository.Close)
	if err != nil {
		log.Fatalf("Failed to initialize auth service: %v", err)
	}

	// Init notification service
	notificationService := messaging_notification.NewInMemmoryNotificationService()
	closer.Add(notificationService.Close)

	// Create grpc client
	auth_verify_service, err := auth_verify_service.NewAuthVerifyService(nil)
	if err != nil {
		log.Fatalf("Failed to create AuthVerifyService: %v", err)
	}
	closer.Add(func() {
		log.Default().Println("Close auth_verify_service")
		auth_verify_service.Close()
	})

	messagingService := messaging_usecase.NewMessagingService(persistentRepository, notificationService)
	messagingHandler := messaging_http.NewMessagingHandler(messagingService)
	messagingRouter := messaging_http.NewRouter(messagingHandler, auth_verify_service)

	// Start the server
	hostPortMessagingHttpServer := os.Getenv("MESSAGING_HTTP_HOST_PORT")
	httpSrv := &http.Server {
		Addr: hostPortMessagingHttpServer,
		Handler: messagingRouter,
	}
	go func ()  {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start messaing server: %v", err)
		}
	}()
	
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	closer.Add(func() {
		log.Default().Println("Close http")
		httpSrv.Shutdown(ctxWithTimeout)
	})
	closer.ListenSignal(syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	log.Default().Println("Close all")
}
