package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles HTTP requests for authentication
// @title Authentication Service API
// @version 1.0
// @description This is the API for the Authentication Service
// @host localhost:8080
// @BasePath /auth
type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags authentication
// @Accept json
// @Produce json
// @Param user body RegisterInput true "User Registration Info"
// @Success 201 {object} User
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.service.Register(input.Email, input.Password)
	if err != nil {
		if err == ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, ErrorResponse{Error: "User already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to register user"})
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Authenticate godoc
// @Summary Authenticate a user
// @Description Authenticate a user with email and password and return JWT tokens
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body LoginInput true "User Login Info"
// @Success 200 {object} TokenPair
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /login [post]
func (h *AuthHandler) Authenticate(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	tokens, err := h.service.Authenticate(input.Email, input.Password)
	if err != nil {
		if err == ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to authenticate"})
		}
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// RefreshToken godoc
// @Summary Refresh authentication tokens
// @Description Refresh the access token using a valid refresh token
// @Tags authentication
// @Accept json
// @Produce json
// @Param refresh_token body RefreshTokenInput true "Refresh Token"
// @Success 200 {object} TokenPair
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var input RefreshTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	tokens, err := h.service.RefreshToken(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Authorization header is missing"})
			c.Abort()
			return
		}

		tokenString := authHeader[7:] // Remove "Bearer " prefix
		userID, err := h.service.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid access token"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func SetupRoutes(router *gin.Engine, handler *AuthHandler) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Authenticate)
		auth.POST("/refresh", handler.RefreshToken)
	}

	// Example of a protected route
	protected := router.Group("/protected")
	protected.Use(handler.AuthMiddleware())
	{
		protected.GET("/", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			c.JSON(http.StatusOK, gin.H{"message": "This is a protected route", "user_id": userID})
		})
	}
}

// RegisterInput represents the input for user registration
type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginInput represents the input for user login
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenInput represents the input for token refresh
type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}
