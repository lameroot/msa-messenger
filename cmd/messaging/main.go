package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Messaging module main function
	r := gin.Default()

	// Readiness probe
	r.GET("/ready", messagingReadinessHandler)

	// Liveness probe
	r.GET("/health", messagingLivenessHandler)

	// Запуск сервера на порту 8080
	r.Run(":8080")
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
