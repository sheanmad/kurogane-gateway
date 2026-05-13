package middleware

import "github.com/gin-gonic/gin"

// Auth is a placeholder middleware. For MVP, it passes through.
// TODO: Implement API key / JWT validation.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}