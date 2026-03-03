package user

import (
	"fmt"
)

type Service struct{}

func (s Service) PrintUser(u User) {
	fmt.Println("Nome:", u.Name)
	fmt.Println("Idade:", u.Age)
}

func (s Service) IsUserAdult(u User) bool {
	return u.IsAdult()
}
