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
	"financial_control/internal/auth"
	category "financial_control/internal/categories"
	connection "financial_control/internal/database"
	"financial_control/internal/router"
	"financial_control/internal/transaction"
	"financial_control/internal/user"
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

	authHandler := auth.NewHandler(userService)

	categoryRepo := category.NewRepository(db)
	categoryService := category.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categoryService)

	transactionRepo := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepo, categoryRepo)
	transactionHandler := transaction.NewHandler(transactionService)

	r := router.SetupRouter(userHandler, transactionHandler, authHandler, categoryHandler)

	r.Run(":8080")
}
