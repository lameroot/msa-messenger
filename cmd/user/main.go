package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"time"

	//user_http "msa-messenger/internal/user/controller/http"

	"github.com/ds248a/closer"
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
	closer.NewCloser()
	// Init config
	log.Default().Print("Start load configs")
	loadEnv()
	log.Default().Print("Loaded configs: ", os.Getenv("DB_POSTGRES_URL"))

	// Init database
	dbURL := os.Getenv("DB_POSTGRES_URL")
	persistentRepository, err := user_repository_psql.NewPostgresUserRepository(dbURL)
	closer.Add(persistentRepository.Close)

	if err != nil {
		log.Fatalf("Failed to initialize auth service: %v", err)
	}

	userService := user_usecase.NewUserService(persistentRepository)

	// Create grpc client
	auth_verify_service, err := auth_verify_service.NewAuthVerifyService(nil)
	if err != nil {
		log.Fatalf("Failed to create AuthVerifyService: %v", err)
	}
	closer.Add(func() {
		err := auth_verify_service.Close()
		if err != nil {
			log.Default().Printf(err.Error())
		}
	})

	// Create user handler
	userHandler := user_http.NewUserHandler(userService)

	// Create server
	router := user_http.NewRouter(userHandler, auth_verify_service)

	// Start the server
	hostPortUserHttpServer := os.Getenv("USER_HTTP_HOST_PORT")
	httpSrv := &http.Server{
		Addr: hostPortUserHttpServer,
		Handler: router.Handler(),
	}
	go func ()  {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start user server: %v", err)	
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
