package auth_http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(authHandler *AuthHandler) *gin.Engine {
	var engine = gin.New()
	// Options
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// Define the routes
	engine.POST("/register", authHandler.Register)
	engine.POST("/login", authHandler.Login)
	engine.POST("/refresh", authHandler.RefreshToken)
	engine.PUT("/update", authHandler.UpdateUser)

	// Swagger documentation route
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return engine
}
