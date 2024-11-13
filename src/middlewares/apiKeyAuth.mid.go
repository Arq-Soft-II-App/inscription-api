// middlewares/api_key_auth.go
package middlewares

import (
	"inscription-api/src/config/envs"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func APIKeyAuthMiddleware(logger *zap.Logger) gin.HandlerFunc {
	KEY := envs.LoadEnvs(".env").Get("INSCRIPTION_API_KEY")
	return func(c *gin.Context) {
		apiKey := c.GetHeader("Authorization")

		if apiKey != KEY {
			logger.Warn("API Key inválida", zap.String("ip", c.ClientIP()))
			ErrorResponse(c, 401, "Invalid API Key")
			return
		}

		c.Next()
	}
}

// ErrorResponse sets CORS headers and aborts the request with a JSON error response
func ErrorResponse(c *gin.Context, status int, message string) {
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Header("Access-Control-Allow-Credentials", "true")
	c.AbortWithStatusJSON(status, gin.H{"error": message})
}
