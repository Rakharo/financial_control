package router

import (
	"financial_control/internal/auth"
	"financial_control/internal/features/analytics"
	category "financial_control/internal/features/categories"
	"financial_control/internal/features/transaction"
	"financial_control/internal/features/user"
	"financial_control/internal/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "financial_control/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userHandler *user.Handler,
	transactionHandler *transaction.Handler,
	authHandler *auth.Handler,
	categoryHandler *category.Handler,
	analyticsHandler *analytics.Handler) *gin.Engine {

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
			users.GET("", userHandler.GetUsers)                    //GET Users List
			users.GET("/:id", userHandler.GetUserByID)             // GET User by ID
			users.PUT("", userHandler.UpdateUser)                  // PUT Update User
			users.PUT("/password", userHandler.UpdateUserPassword) // PUT Update Password
			users.DELETE("", userHandler.DeleteUser)               // DELETE User
		}

		transactions := protected.Group("/transaction")
		{
			transactions.GET("", transactionHandler.GetTransactions)
			transactions.GET("/:id", transactionHandler.GetTransactionByID)
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
			transactions.PUT("/:id", transactionHandler.UpdateTransaction)
		}

		categories := protected.Group("/category")
		{
			categories.GET("", categoryHandler.GetAll)
			categories.GET("/:id", categoryHandler.GetByID)
			categories.POST("", categoryHandler.Create)
			categories.PUT("/:id", categoryHandler.Update)
			categories.DELETE("/:id", categoryHandler.Delete)
		}

		insights := protected.Group("/analytics")
		{
			insights.GET("/dashboard", analyticsHandler.GetDashboard)
		}
	}

	return r
}
