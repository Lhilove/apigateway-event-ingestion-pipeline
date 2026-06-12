package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Temporary in memory store. In production, this would be a database or secrets manager.
var validAPIKeys = map[string]string{
	"sk_live_abc123": "linux-agent-01",
	"sk_live_def456": "windows-agent-01",
}

// APIKeyAuth is a gin middleware that checks for a valid API key in the Authorization header.
func APIKeyAuth() gin.HandlerFunc {
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
