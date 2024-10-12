package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/lameroot/msa-messenger/internal/auth"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Database connection
	db, err := sql.Open("postgres", "postgres://messenger:messenger@localhost/messenger?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize AuthService
	jwtSecret := "your-secret-key" // In production, use a secure method to manage this secret
	authService := auth.NewAuthService(db, jwtSecret, 15*time.Minute, 24*time.Hour)

	// Initialize AuthHandler
	authHandler := auth.NewAuthHandler(authService)

	// Set up Gin router
	r := gin.Default()

	// Set up auth routes
	auth.SetupRoutes(r, authHandler)

	// Serve Swagger UI
	r.StaticFile("/docs/swagger.json", "./docs/swagger.json")
	r.StaticFile("/swagger-ui", "./docs/swagger-ui.html")

	// Readiness probe
	r.GET("/ready", authReadinessHandler)

	// Liveness probe
	r.GET("/health", authLivenessHandler)

	// Start HTTP server
	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	auth.RegisterAuthServiceServer(grpcServer, auth.NewGRPCServer(authService))

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	go func() {
		log.Println("Starting gRPC server on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// Set up gRPC-Gateway
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = auth.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC-Gateway: %v", err)
	}

	// Start gRPC-Gateway server
	gwServer := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	log.Println("Starting HTTP server on :8080, gRPC server on :50051, and gRPC-Gateway server on :8081")
	log.Println("Swagger UI available at http://localhost:8080/swagger-ui")
	if err := gwServer.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start gRPC-Gateway server: %v", err)
	}
}

func authReadinessHandler(c *gin.Context) {
	// Here should be the logic to check if the Auth module is ready
	// For example, checking the database connection
	c.JSON(http.StatusOK, gin.H{
		"status": "auth module ready",
	})
}

func authLivenessHandler(c *gin.Context) {
	// Here should be the logic to check if the Auth module is alive
	c.JSON(http.StatusOK, gin.H{
		"status": "auth module alive",
	})
}
