package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kurogane/gateway/internal/model"
)

func Errors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("error: %v path: %s", err, c.Request.URL.Path)

			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: err.Error(),
				Code:  http.StatusInternalServerError,
			})
		}
	}
}