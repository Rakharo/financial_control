package user

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetAllUsers() ([]User, error)
	GetUserById(id int64) (*User, error)
	GetUserByLogin(login string) (*User, error)
	CreateUser(user *User) error
	UpdateUser(id int64, name string, email string) error
	DeleteUser(id int64) error
}

type Service struct {
	repo UserRepository
}

func NewService(repo UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Login(login string, password string) (*User, error) {
	user, err := s.repo.GetUserByLogin(login)

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *Service) GetAllUsers() ([]UserResponse, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var response []UserResponse
	for _, u := range users {
		response = append(response, UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
			Login: u.Login,
		})
	}

	return response, nil
}

func (s *Service) GetUserById(id int64) (*User, error) {
	if id <= 0 {
		return nil, errors.New("ID inválido")
	}
	return s.repo.GetUserById(id)
}

func (s *Service) CreateUser(req CreateUserRequest) error {

	existingUser, err := s.repo.GetUserByLogin(req.Login)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if existingUser != nil {
		return errors.New("usuário já existe")
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	user := User{
		Name:     req.Name,
		Email:    req.Email,
		Login:    req.Login,
		Password: string(hash),
	}

	return s.repo.CreateUser(&user)
}

func (s *Service) UpdateUser(id int64, name string, email string) error {
	if id <= 0 {
		return errors.New("ID inválido")
	}
	if name == "" || email == "" {
		return errors.New("Nome e email são obrigatórios")
	}
	return s.repo.UpdateUser(id, name, email)
}

func (s *Service) DeleteUser(id int64) error {
	if id <= 0 {
		return errors.New("ID inválido")
	}
	return s.repo.DeleteUser(id)
}
