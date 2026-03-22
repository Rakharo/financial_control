package main

// @title Financial Control API
// @version 1.0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description API para controle financeiro
// @host localhost:8080
// @BasePath /

import (
	connection "financial_control/internal/database"
	"financial_control/internal/features/analytics"
	"financial_control/internal/features/auth"
	category "financial_control/internal/features/categories"
	"financial_control/internal/features/installment"
	"financial_control/internal/features/transaction"
	"financial_control/internal/features/user"
	"financial_control/internal/router"
)

func main() {
	db := connection.Connect()

	err := connection.SeedCategories(db)
	if err != nil {
		panic(err)
	}

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	authRepo := auth.NewRepository(db)
	authService := auth.NewService(*authRepo, *userRepo)
	authHandler := auth.NewHandler(authService, userService)

	installmentRepo := installment.NewRepository(db)
	installmentService := installment.NewService(installmentRepo)

	categoryRepo := category.NewRepository(db)
	categoryService := category.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categoryService)

	transactionRepo := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepo, categoryRepo, installmentService)
	transactionHandler := transaction.NewHandler(transactionService)

	analyticsRepo := analytics.NewRepository(db)
	analyticsService := analytics.NewService(analyticsRepo)
	analyticsHandler := analytics.NewHandler(analyticsService)

	r := router.SetupRouter(userHandler, transactionHandler, authHandler, categoryHandler, analyticsHandler)

	r.Run(":8080")
}
