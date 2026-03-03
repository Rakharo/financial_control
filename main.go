package main

import (
	"financial_control/user"
	"fmt"
)

func main() {
	u := user.User{
		Name: "Andre",
		Age:  30,
	}

	service := user.Service{}
	service.PrintUser(u)

	isAdult := service.IsUserAdult(u)
	if isAdult {
		fmt.Println("O usuário é maior de idade")
	} else {
		fmt.Println("O usuário é menor de idade")
	}

	fmt.Println("Aplicação finalizada")
}
