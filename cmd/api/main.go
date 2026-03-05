package main

import (
	connection "financial_control/internal/database"
	"financial_control/internal/router"
	"financial_control/internal/user"
)

func main() {
	db := connection.Connect()

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	r := router.SetupRouter(userHandler)

	r.Run(":8080")
}
