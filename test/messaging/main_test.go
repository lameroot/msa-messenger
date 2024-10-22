package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Mock handlers for testing
func mockMessagingReadinessHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "messaging module ready",
	})
}

func mockMessagingLivenessHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "messaging module alive",
	})
}

func TestMessagingReadinessHandler(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Setup the router with the mock handler
	r := gin.Default()
	r.GET("/ready", mockMessagingReadinessHandler)

	// Create a mock request to pass to our handler
	req, _ := http.NewRequest("GET", "/ready", nil)
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check to see if the response was what you expected
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d but got %d", http.StatusOK, w.Code)
	}

	expected := `{"status":"messaging module ready"}`
	if w.Body.String() != expected {
		t.Errorf("Expected body %s but got %s", expected, w.Body.String())
	}
}

func TestMessagingLivenessHandler(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Setup the router with the mock handler
	r := gin.Default()
	r.GET("/health", mockMessagingLivenessHandler)

	// Create a mock request to pass to our handler
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check to see if the response was what you expected
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d but got %d", http.StatusOK, w.Code)
	}

	expected := `{"status":"messaging module alive"}`
	if w.Body.String() != expected {
		t.Errorf("Expected body %s but got %s", expected, w.Body.String())
	}
}
