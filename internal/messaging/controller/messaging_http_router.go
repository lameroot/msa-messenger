package messaging_http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewRouter(messagingHandler *MessagingHandler) *gin.Engine {
	var engine = gin.New()
	// Options
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(AuthRequiredMiddleware())

	engine.POST("/send", messagingHandler.SendMessage)
	engine.GET("/messages", messagingHandler.GetMessages)

	// Readiness probe
	engine.GET("/ready", messagingReadinessHandler)

	// Liveness probe
	engine.GET("/health", messagingLivenessHandler)

	return engine
}

func AuthRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") != "123" { //todo auth
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		IDUser, err := uuid.Parse("cf197cc1-3e93-476c-8b38-08f52cbe5a46")
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		c.Set("user_id", IDUser) //todo auth
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
