package middleware

import (
	"errors"
	"financial_control/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			if c.Writer.Written() {
				return
			}

			err := c.Errors.Last().Err

			var appErr *utils.AppError

			if errors.As(err, &appErr) {
				c.JSON(appErr.Status, utils.APIError{
					Status:  appErr.Status,
					Message: appErr.Message,
				})
				return
			}

			c.JSON(http.StatusInternalServerError, utils.APIError{
				Status:  http.StatusInternalServerError,
				Message: "Erro interno do servidor",
			})
		}
	}
}
