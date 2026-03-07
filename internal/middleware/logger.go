package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("Request:", c.Request.Method, c.Request.URL.Path)

		c.Next()

		fmt.Println("Status:", c.Writer.Status())
	}
}
