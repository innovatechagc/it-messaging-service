package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/company/microservice-template/internal/auth"
	"github.com/company/microservice-template/internal/services"
	"github.com/company/microservice-template/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	healthService := services.NewHealthService()
	messagingService := services.NewMessagingService(nil, nil, nil)
	fileService := services.NewFileService(nil, nil)
	jwtManager := auth.NewJWTManager("test-secret", "test-issuer", 24)
	logger := logger.NewLogger("debug")
	
	SetupRoutes(router, healthService, messagingService, fileService, jwtManager, logger)
	
	// Test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	router.ServeHTTP(w, req)
	
	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "healthy")
}

func TestReadinessCheck(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	healthService := services.NewHealthService()
	messagingService := services.NewMessagingService(nil, nil, nil)
	fileService := services.NewFileService(nil, nil)
	jwtManager := auth.NewJWTManager("test-secret", "test-issuer", 24)
	logger := logger.NewLogger("debug")
	
	SetupRoutes(router, healthService, messagingService, fileService, jwtManager, logger)
	
	// Test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ready", nil)
	router.ServeHTTP(w, req)
	
	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ready")
}