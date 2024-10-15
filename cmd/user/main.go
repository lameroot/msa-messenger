package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lameroot/msa-messenger/internal/auth"
	"github.com/lameroot/msa-messenger/internal/user"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

var db *sql.DB
var authClient auth.AuthServiceClient

func main() {
	var err error
	db, err = sql.Open("postgres", "postgres://messenger:messenger@localhost/messenger?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}
	defer conn.Close()

	authClient = auth.NewAuthServiceClient(conn)

	r := gin.Default()

	r.Use(authMiddleware)

	r.POST("/friends/add/:id", addFriend)
	r.DELETE("/friends/remove/:id", removeFriend)
	r.GET("/friends", getFriends)
	r.POST("/friends/accept/:id", acceptFriendRequest)
	r.POST("/friends/reject/:id", rejectFriendRequest)

	r.GET("/ready", readinessHandler)
	r.GET("/health", livenessHandler)

	r.Run(":8082")
}

func authMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
		c.Abort()
		return
	}

	resp, err := authClient.VerifyToken(context.Background(), &auth.VerifyTokenRequest{Token: token})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	if !resp.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	c.Set("user_id", resp.UserId)
	c.Next()
}

func addFriend(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.GetString("user_id"), 10, 64)
	friendID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID"})
		return
	}

	err = user.AddFriend(db, userID, friendID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add friend"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request sent"})
}

func removeFriend(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.GetString("user_id"), 10, 64)
	friendID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID"})
		return
	}

	err = user.RemoveFriend(db, userID, friendID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove friend"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend removed"})
}

func getFriends(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.GetString("user_id"), 10, 64)

	friends, err := user.GetFriends(db, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get friends"})
		return
	}

	c.JSON(http.StatusOK, friends)
}

func acceptFriendRequest(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.GetString("user_id"), 10, 64)
	friendID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID"})
		return
	}

	err = user.AcceptFriendRequest(db, userID, friendID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept friend request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request accepted"})
}

func rejectFriendRequest(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.GetString("user_id"), 10, 64)
	friendID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID"})
		return
	}

	err = user.RejectFriendRequest(db, userID, friendID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject friend request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request rejected"})
}

func readinessHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "user module ready"})
}

func livenessHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "user module alive"})
}
