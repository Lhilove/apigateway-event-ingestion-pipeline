package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// loadAPIKeys loads API keys from environment variables and returns a map of valid keys to agent names.
func loadAPIKeys() map[string]string {
	return map[string]string{
		os.Getenv("API_KEY_LINUX_AGENT"):   "linux-agent-01",
		os.Getenv("API_KEY_WINDOWS_AGENT"): "windows-agent-01",
	}
}

// APIKeyAuth is a gin middleware that checks for a valid API key in the Authorization header.
func APIKeyAuth() gin.HandlerFunc {
	validAPIKeys := loadAPIKeys()
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing or malformed Authorization header",
			})
			c.Abort()
			return
		}

		key := strings.TrimPrefix(authHeader, "Bearer ")

		agentName, ok := validAPIKeys[key]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid API key",
			})
			c.Abort()
			return
		}

		// Store agent identity in context for downstream handlers/logging
		c.Set("agent_name", agentName)
		c.Next()
	}
}
