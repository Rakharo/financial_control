package router

import (
	"financial_control/internal/auth"
	category "financial_control/internal/categories"
	"financial_control/internal/middleware"
	"financial_control/internal/transaction"
	"financial_control/internal/user"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "financial_control/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *user.Handler, transactionHandler *transaction.Handler, authHandler *auth.Handler, categoryHandler *category.Handler) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Logger())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// SWAGGER
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//ROTAS PÚBLICAS
	public := r.Group("/")
	{
		r.POST("/auth/login", authHandler.Login)
		public.POST("/auth/register", userHandler.CreateUser)
		public.POST("/auth/refresh", authHandler.RefreshToken)
	}

	//ROTAS PROTEGIDAS
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		users := protected.Group("/user")
		{
			users.GET("/me", userHandler.GetMe)
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

			transactions.GET("/summary", transactionHandler.GetSummary)
		}

		categories := protected.Group("/category")
		{
			categories.GET("", categoryHandler.GetAll)
			categories.GET("/:id", categoryHandler.GetByID)
			categories.POST("", categoryHandler.Create)
			categories.PUT("/:id", categoryHandler.Update)
			categories.DELETE("/:id", categoryHandler.Delete)
		}
	}

	return r
}
