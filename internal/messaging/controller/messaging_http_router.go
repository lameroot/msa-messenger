package messaging_http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	auth_proto "github.com/lameroot/msa-messenger/pkg/api/auth"
)

func NewRouter(messagingHandler *MessagingHandler, authClient *auth_proto.TokenVerifyServiceClient) *gin.Engine {
	var engine = gin.New()
	// Options
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	notAuthGroup := engine.Group("/")
	{
		// Readiness probe
		notAuthGroup.GET("/ready", messagingReadinessHandler)

		// Liveness probe
		notAuthGroup.GET("/health", messagingLivenessHandler)
	}

	authGroup := engine.Group("/")
	{
		authGroup.Use(AuthRequiredMiddleware(*authClient))

		authGroup.POST("/send", messagingHandler.SendMessage)
		authGroup.GET("/messages", messagingHandler.GetMessages)
	}


	return engine
}

func AuthRequiredMiddleware(authClient auth_proto.TokenVerifyServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}
		req := &auth_proto.TokenVerificationRequest{
			Token: token,
		}
		verifyResponse, err := authClient.VerifyToken(context.Background(), req)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		if !verifyResponse.Verified {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		// if c.GetHeader("Authorization") != "123" { //todo auth
		// 	c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		// 	return
		// }
		// IDUser, err := uuid.Parse("cf197cc1-3e93-476c-8b38-08f52cbe5a46")

		IDUser, err := uuid.Parse(verifyResponse.UserId)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		c.Set("user_id", IDUser)
		c.Next()
	}
}

func messagingReadinessHandler(c *gin.Context) {
	// Здесь должна быть логика проверки готовности модуля Messaging
	// Например, проверка подключения к брокеру сообщений
	c.JSON(http.StatusOK, gin.H{
		"status": "messaging module ready",
	})
}

func messagingLivenessHandler(c *gin.Context) {
	// Здесь должна быть логика проверки жизнеспособности модуля Messaging
	c.JSON(http.StatusOK, gin.H{
		"status": "messaging module alive",
	})
}
