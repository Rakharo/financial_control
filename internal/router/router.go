package router

import (
	"financial_control/internal/user"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *user.Handler) *gin.Engine {
	r := gin.Default()

	users := r.Group("/users")
	{
		users.GET("", userHandler.GetUsers)
		users.GET("/:id", userHandler.GetUserByID)
		users.POST("", userHandler.CreateUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}

	// transactions := r.Group("/transactions")
	// {
	// 	// transactions.GET("", transactionHandler.GetTransactions)
	// 	// transactions.GET("/:id", transactionHandler.GetTransactionByID)
	// 	// transactions.POST("", transactionHandler.CreateTransaction)
	// }

	return r
}
