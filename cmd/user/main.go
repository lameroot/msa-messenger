package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Readiness probe
	r.GET("/ready", readinessHandler)

	// Liveness probe
	r.GET("/health", livenessHandler)

	// Запуск сервера на порту 8080
	r.Run(":8080")
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
