package middleware

import (
	"net/http"
	"strings"

	"task-manager-api/internal/services"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware выполняет проверку JWT токена.
type AuthMiddleware struct {
	tokenService *services.TokenService
}

// NewAuthMiddleware создает middleware авторизации.
func NewAuthMiddleware(
	tokenService *services.TokenService,
) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
	}
}

// Handle выполняет обработку JWT токена.
func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "требуется токен авторизации",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "некорректный формат токена",
			})
			c.Abort()
			return
		}

		claims, err := m.tokenService.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "невалидный токен",
			})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
