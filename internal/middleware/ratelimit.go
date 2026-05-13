package middleware

import "github.com/gin-gonic/gin"

// RateLimit is a placeholder middleware. For MVP, it passes through.
// TODO: Implement token bucket rate limiting.
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}