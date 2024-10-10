package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Mock handlers for testing
func mockReadinessHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "user module ready",
	})
}

func mockLivenessHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "user module alive",
	})
}

func TestReadinessHandler(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Setup the router with the mock handler
	r := gin.Default()
	r.GET("/ready", mockReadinessHandler)

	// Create a mock request to pass to our handler
	req, _ := http.NewRequest("GET", "/ready", nil)
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check to see if the response was what you expected
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d but got %d", http.StatusOK, w.Code)
	}

	expected := `{"status":"user module ready"}`
	if w.Body.String() != expected {
		t.Errorf("Expected body %s but got %s", expected, w.Body.String())
	}
}

func TestLivenessHandler(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Setup the router with the mock handler
	r := gin.Default()
	r.GET("/health", mockLivenessHandler)

	// Create a mock request to pass to our handler
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check to see if the response was what you expected
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d but got %d", http.StatusOK, w.Code)
	}

	expected := `{"status":"user module alive"}`
	if w.Body.String() != expected {
		t.Errorf("Expected body %s but got %s", expected, w.Body.String())
	}
}
