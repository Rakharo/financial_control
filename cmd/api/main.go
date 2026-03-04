package main

import (
	"log"
	"net/http"
	"strings"

	connection "financial_control/internal/database"
	"financial_control/internal/user"
)

func main() {
	db := connection.Connect()

	userRepo := user.NewRepository(db)
	userHandler := user.NewHandler(userRepo)

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users" {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			userHandler.GetUsers(w, r)
		case http.MethodPost:
			userHandler.CreateUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		// Só aceita GET por enquanto
		if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/users/") {
			userHandler.GetUserByID(w, r)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	log.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
