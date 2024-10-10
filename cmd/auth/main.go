package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Auth module main function
	r := gin.Default()

	// Readiness probe
	r.GET("/ready", authReadinessHandler)

	// Liveness probe
	r.GET("/health", authLivenessHandler)

	// Запуск сервера на порту 8080
	r.Run(":8080")
}

func authReadinessHandler(c *gin.Context) {
	// Здесь должна быть логика проверки готовности модуля Auth
	// Например, проверка подключения к базе данных пользователей
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
