package auth_http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(authHandler *AuthHandler) *gin.Engine {
	var engine = gin.New()
	// Options
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	notAuthGroup := engine.Group("/")
	{
		// Readiness probe
		notAuthGroup.GET("/ready", authReadinessHandler)

		// Liveness probe
		notAuthGroup.GET("/health", authLivenessHandler)
	}

	// Define the routes
	engine.POST("/register", authHandler.Register)
	engine.POST("/login", authHandler.Login)
	engine.POST("/refresh", authHandler.RefreshToken)
	engine.PUT("/update", authHandler.UpdateUser)

	// Swagger documentation route
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return engine
}

func authReadinessHandler(c *gin.Context) {
	// Здесь должна быть логика проверки готовности модуля Auth
	// Например, проверка подключения к брокеру сообщений
	c.JSON(http.StatusOK, gin.H{
		"status": "auth module ready",
	})
}

func authLivenessHandler(c *gin.Context) {
	// Здесь должна быть логика проверки жизнеспособности модуля Auth
	c.JSON(http.StatusOK, gin.H{
		"status": "auth module alive",
	})
}
