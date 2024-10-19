package user_http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewRouter(userHandler *UserHandler) *gin.Engine {
	var engine = gin.New()
	// Options
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(AuthRequiredMiddleware())

	// Routes
	engine.GET("/friends", userHandler.GetFriends)
	engine.POST("/friends", userHandler.AddUserToFriends)
	engine.DELETE("/friends", userHandler.DeleteUserFromFriends)
	engine.POST("/friendship", userHandler.AcceptFriend)
	engine.DELETE("/friendship", userHandler.RejectFriend)

	// Readiness probe
	engine.GET("/ready", readinessHandler)

	// Liveness probe
	engine.GET("/health", livenessHandler)

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

func readinessHandler(c *gin.Context) {
	// Здесь вы можете добавить логику для проверки готовности вашего приложения
	// Например, проверка подключения к базе данных, доступность внешних сервисов и т.д.

	// В этом примере мы просто возвращаем успешный статус
	c.JSON(http.StatusOK, gin.H{
		"status": "user module ready",
	})
}

func livenessHandler(c *gin.Context) {
	// Здесь вы можете добавить логику для проверки жизнеспособности вашего приложения
	// Например, проверка критических компонентов, отсутствие deadlock'ов и т.д.

	// В этом примере мы просто возвращаем успешный статус
	c.JSON(http.StatusOK, gin.H{
		"status": "user module alive",
	})
}
