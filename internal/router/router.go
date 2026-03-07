package router

import (
	"financial_control/internal/middleware"
	"financial_control/internal/transaction"
	"financial_control/internal/user"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "financial_control/docs"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *user.Handler, transactionHandler *transaction.Handler) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Logger())

	// SWAGGER
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//ROTAS PÚBLICAS
	public := r.Group("/")
	{
		public.POST("/login", userHandler.Login)
		public.POST("/register", userHandler.CreateUser)
	}

	//ROTAS PROTEGIDAS
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		users := protected.Group("/user")
		{
			users.GET("", userHandler.GetUsers)          //GET Users List
			users.GET("/:id", userHandler.GetUserByID)   // GET User by ID
			users.PUT("/:id", userHandler.UpdateUser)    // PUT Update User
			users.DELETE("/:id", userHandler.DeleteUser) // DELETE User
		}

		transactions := protected.Group("/transaction")
		{
			transactions.GET("", transactionHandler.GetTransactions)
			transactions.GET("/:id", transactionHandler.GetTransactionByID)
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
			transactions.PUT("/:id", transactionHandler.UpdateTransaction)
		}
	}

	return r
}
