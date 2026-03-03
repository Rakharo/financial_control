package main

import (
	// "financial_control/user"
	"financial_control/payment"
	"fmt"
)

func main() {
	// u := user.User{
	// 	Name: "Andre",
	// 	Age:  30,
	// }

	// service := user.Service{}
	// service.PrintUser(u)

	// isAdult := service.IsUserAdult(u)
	// if isAdult {
	// 	fmt.Println("O usuário é maior de idade")
	// } else {
	// 	fmt.Println("O usuário é menor de idade")
	// }

	// fmt.Println("Aplicação finalizada")

	payment := payment.Payment{
		TotalAmount:      4000,
		FixedExpenses:    1000,
		VariableExpenses: 500,
		FixedIncome:      2000,
		VariableIncome:   800,
	}
	total := payment.CalculateFinalAmount()
	totalExpenses := payment.TotalExpenses()
	totalIncome := payment.TotalIncome()

	fmt.Println("Total Amount:", payment.TotalAmount)
	fmt.Println("Total Expenses:", totalExpenses)
	fmt.Println("Total Income:", totalIncome)
	fmt.Println("Result:", total)
}
