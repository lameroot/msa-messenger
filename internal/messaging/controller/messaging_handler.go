package messaging_http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	messaging_models "github.com/lameroot/msa-messenger/internal/messaging/models"
	messaging_usecase "github.com/lameroot/msa-messenger/internal/messaging/usecase"
)

type MessagingHandler struct {
	messagingService *messaging_usecase.MessagingService
}

func NewMessagingHandler(messagingService *messaging_usecase.MessagingService) *MessagingHandler {
	return &MessagingHandler{
		messagingService: messagingService,
	}
}

// SendMessage godoc
// @Summary Send a message
// @Description Send a message from the authenticated user to another user
// @Tags messaging
// @Accept json
// @Produce json
// @Param request body messaging_models.SendMessageRequest true "Message details"
// @Security BearerAuth
// @Success 200 {object} messaging_models.SendMessageResponse "Message sent successfully"
// @Failure 400 {object} ErrorResponse "Bad request or user not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /messaging/send [post]
func (handler *MessagingHandler) SendMessage(c *gin.Context) {
	var req messaging_models.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userID, exists := c.MustGet("user_id").(uuid.UUID)
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "User not found"})
		return
	}
	res, err := handler.messagingService.SendMessage(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetMessages godoc
// @Summary Get messages for a user
// @Description Retrieve messages for the authenticated user based on provided criteria
// @Tags messaging
// @Accept json
// @Produce json
// @Param request body messaging_models.GetMessagesRequest true "Message retrieval criteria"
// @Security BearerAuth
// @Success 200 {object} messaging_models.Message "Messages retrieved successfully"
// @Failure 400 {object} ErrorResponse "Bad request or user not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /messaging/messages [post]
func (handler *MessagingHandler) GetMessages(c *gin.Context) {
	var req messaging_models.GetMessagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	userID, exists := c.MustGet("user_id").(uuid.UUID)
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "User not found"})
		return
	}
	res, err := handler.messagingService.GetMessages(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}
