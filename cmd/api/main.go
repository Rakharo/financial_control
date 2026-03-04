package main

import (
	connection "financial_control/internal/database"
	"financial_control/internal/user"

	"github.com/gin-gonic/gin"
)

func main() {
	db := connection.Connect()

	userRepo := user.NewRepository(db)
	userHandler := user.NewHandler(userRepo)

	r := gin.Default()

	r.GET("/users", userHandler.GetUsers)
	r.GET("/users/:id", userHandler.GetUserByID)
	r.POST("/users", userHandler.CreateUser)

	r.Run(":8080")
}
