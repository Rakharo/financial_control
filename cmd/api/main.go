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
	"financial_control/internal/router"
	"financial_control/internal/transaction"
	"financial_control/internal/user"
)

func main() {
	db := connection.Connect()

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	transactionRepo := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepo)
	transactionHandler := transaction.NewHandler(transactionService)

	r := router.SetupRouter(userHandler, transactionHandler)

	r.Run(":8080")
}
