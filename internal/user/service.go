package user

import "errors"

type UserRepository interface {
	GetAllUsers() ([]User, error)
	GetUserById(id int64) (*User, error)
	CreateUser(name string, email string) error
}

type Service struct {
	repo UserRepository
}

func NewService(repo UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllUsers() ([]User, error) {
	return s.repo.GetAllUsers()
}

func (s *Service) GetUserById(id int64) (*User, error) {
	return s.repo.GetUserById(id)
}

func (s *Service) CreateUser(name string, email string) error {
	if name == "" || email == "" {
		return errors.New("Nome e email são obrigatórios")
	}
	return s.repo.CreateUser(name, email)
}
