package user_http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	auth_verify_service "github.com/lameroot/msa-messenger/pkg/auth"
)

type HttpRoute struct {
	engine *gin.Engine
}

func NewRouter(userHandler *UserHandler, authClient *auth_verify_service.AuthVerifyService) *gin.Engine {
	var engine = gin.New()
	// Options
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	notAuthGroup := engine.Group("/")
	{
		// Readiness probe
		notAuthGroup.GET("/ready", readinessHandler)	
		// Liveness probe
		notAuthGroup.GET("/health", livenessHandler)
	}
	
	authGroup := engine.Group("/")
	{
		authGroup.Use(AuthRequiredMiddleware(authClient))

		// Routes
		authGroup.GET("/friends", userHandler.GetFriends)
		authGroup.POST("/friends", userHandler.AddUserToFriends)
		authGroup.DELETE("/friends", userHandler.DeleteUserFromFriends)
		authGroup.POST("/friendship", userHandler.AcceptFriend)
		authGroup.DELETE("/friendship", userHandler.RejectFriend)
	}

	return engine
}

func AuthRequiredMiddleware(authClient *auth_verify_service.AuthVerifyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}
		IDUser, err := authClient.VerifyToken(context.Background(), token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error})
			return
		}

		c.Set("user_id", IDUser)
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
