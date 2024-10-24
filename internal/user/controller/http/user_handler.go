package user_http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	user_models "github.com/lameroot/msa-messenger/internal/user/models"
	user_usecase "github.com/lameroot/msa-messenger/internal/user/usecase"
)

type UserHandler struct {
	userService *user_usecase.UserService
}

func NewUserHandler(userService *user_usecase.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// @Summary Add user to friends
// @Description Add a user to the current user's friends list
// @Tags friends
// @Accept json
// @Produce json
// @Param request body user_models.AddUserToFriendsRequest true "Add user to friends request"
// @Success 200 {object} user_models.AddUserToFriendsResponse
// @Failure 400 {object} ErrorResponse
// @Router /friends [post]
func (h *UserHandler) AddUserToFriends(c *gin.Context) {
	var req user_models.AddUserToFriendsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	userID, exists := c.MustGet("user_id").(uuid.UUID)
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "User not found"})
		return
	}
	res, err := h.userService.AddUserToFriend(userID, req.FriendUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// @Summary Delete user from friends
// @Description Remove a user from the current user's friends list
// @Tags friends
// @Accept json
// @Produce json
// @Param request body user_models.DeleteUserFromFriendsRequest true "Delete user from friends request"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Router /friends [delete]
func (h *UserHandler) DeleteUserFromFriends(c *gin.Context) {
	var req user_models.DeleteUserFromFriendsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	userID, exists := c.MustGet("user_id").(uuid.UUID)
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "User not found"})
		return
	}
	err := h.userService.DeleteUserFromFriends(userID, req.FriendUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Accept friend request
// @Description Accept a friend request from another user
// @Tags friends
// @Accept json
// @Produce json
// @Param request body user_models.AcceptFriendRequest true "Accept friend request"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Router /friends/accept [post]
func (h *UserHandler) AcceptFriend(c *gin.Context) {
	var req user_models.AcceptFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	userID, exists := c.MustGet("user_id").(uuid.UUID)
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "User not found"})
		return
	}
	err := h.userService.AcceptFriend(userID, req.FriendUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Reject friend request
// @Description Reject a friend request from another user
// @Tags friends
// @Accept json
// @Produce json
// @Param request body user_models.RejectFriendRequest true "Reject friend request"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Router /friends/reject [post]
func (h *UserHandler) RejectFriend(c *gin.Context) {
	var req user_models.RejectFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	userID, exists := c.MustGet("user_id").(uuid.UUID)
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "User not found"})
		return
	}
	err := h.userService.RejectFriend(userID, req.FriendUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Get friends
// @Description Get the list of friends for the current user
// @Tags friends
// @Produce json
// @Success 200 {array} user_models.Friend
// @Failure 400 {object} ErrorResponse
// @Router /friends [get]
func (h *UserHandler) GetFriends(c *gin.Context) {
	userID, exists := c.MustGet("user_id").(uuid.UUID)
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "User not found"})
		return
	}
	friends, err := h.userService.GetFriends(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, friends)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}
