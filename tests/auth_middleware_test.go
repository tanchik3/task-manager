package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"task-manager-api/internal/config"
	"task-manager-api/internal/middleware"
	"task-manager-api/internal/services"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		JWTSecret:      "test-secret",
		JWTExpireHours: 24,
	}

	tokenService := services.NewTokenService(cfg)

	validToken, err := tokenService.GenerateToken(1)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "valid jwt",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid jwt",
			authHeader:     "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()

			authMiddleware := middleware.NewAuthMiddleware(
				tokenService,
			)

			router.Use(authMiddleware.Handle())

			router.GET("/protected", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(
				http.MethodGet,
				"/protected",
				nil,
			)

			req.Header.Set(
				"Authorization",
				tt.authHeader,
			)

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf(
					"expected status %d, got %d",
					tt.expectedStatus,
					recorder.Code,
				)
			}
		})
	}
}
