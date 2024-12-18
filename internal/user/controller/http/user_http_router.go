package user_http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	auth_proto "github.com/lameroot/msa-messenger/pkg/api/auth"
)

func NewRouter(userHandler *UserHandler, authClient *auth_proto.TokenVerifyServiceClient) *gin.Engine {
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
		authGroup.Use(AuthRequiredMiddleware(*authClient))

		// Routes
		authGroup.GET("/friends", userHandler.GetFriends)
		authGroup.POST("/friends", userHandler.AddUserToFriends)
		authGroup.DELETE("/friends", userHandler.DeleteUserFromFriends)
		authGroup.POST("/friendship", userHandler.AcceptFriend)
		authGroup.DELETE("/friendship", userHandler.RejectFriend)
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
