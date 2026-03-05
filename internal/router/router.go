package router

import (
	"financial_control/internal/user"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *user.Handler) *gin.Engine {
	r := gin.Default()

	users := r.Group("/users")
	{
		users.GET("", userHandler.GetUsers)          //GET Users List
		users.GET("/:id", userHandler.GetUserByID)   // GET User by ID
		users.POST("", userHandler.CreateUser)       // POST Create User
		users.PUT("/:id", userHandler.UpdateUser)    // PUT Update User
		users.DELETE("/:id", userHandler.DeleteUser) // DELETE User
	}

	// transactions := r.Group("/transactions")
	// {
	// 	// transactions.GET("", transactionHandler.GetTransactions)
	// 	// transactions.GET("/:id", transactionHandler.GetTransactionByID)
	// 	// transactions.POST("", transactionHandler.CreateTransaction)
	// }

	return r
}
