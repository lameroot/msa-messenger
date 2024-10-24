package auth_http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	auth_models "github.com/lameroot/msa-messenger/internal/auth/models"
	auth_usecase "github.com/lameroot/msa-messenger/internal/auth/usecase"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	service *auth_usecase.AuthService
}

func NewAuthHandler(service *auth_usecase.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth_models.RegisterRequest true "Registration information"
// @Success 201 {object} auth_models.UserResponse "User successfully registered"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 409 {object} ErrorResponse "User already exists or nickname taken"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req auth_models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.service.Register(req.Email, req.Password, req.Nickname)
	if err != nil {
		if err.Error() == "user already exists" || err.Error() == "nickname already taken" {
			c.JSON(http.StatusConflict, ErrorResponse{Error: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to register user"})
		}
		return
	}

	c.JSON(http.StatusCreated, auth_models.UserResponse{User: *user})
}

// Login godoc
// @Summary User login
// @Description Authenticate a user and return access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth_models.LoginRequest true "Login credentials"
// @Success 200 {object} auth_models.AuthResponse "Authentication successful"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth_models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	resp, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RefreshToken godoc
// @Summary Refresh authentication token
// @Description Generate a new access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth_models.RefreshRequest true "Refresh token request"
// @Success 200 {object} auth_models.TokenResponse "New access token"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req auth_models.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	token, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, auth_models.TokenResponse{Token: *token})
}

// UpdateUser godoc
// @Summary Update user information
// @Description Update the avatar URL and info for an authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth_models.UpdateUserRequest true "Update user request"
// @Success 200 {object} auth_models.User "Updated user information"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/update [put]
func (h *AuthHandler) UpdateUser(c *gin.Context) {
	var req auth_models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Default().Printf("Error in binding request: %v", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	user, err := h.service.CheckPassword(req.Email, req.Password)
	if err != nil {
		log.Default().Printf("Error in authentication: %v", err.Error())
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid email or password"})
		return
	}

	updatedUser, err := h.service.UpdateUser(user.ID, req.AvatarURL, req.Info)
	if err != nil {
		log.Default().Printf("Error in updating user info: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}
