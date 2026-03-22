package user

import (
	"database/sql"
	"errors"
	"time"

	"financial_control/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetAllUsers() ([]User, error)
	GetUserById(userID uint64) (*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateUser(user *User) error
	UpdateUser(userID uint64, user *User) error
	DeleteUser(userID uint64) error
}

type Service struct {
	repo UserRepository
}

func NewService(repo UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllUsers() ([]UserResponse, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var response []UserResponse
	for _, u := range users {
		response = append(response, UserResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Phone:     u.Phone,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	return response, nil
}

func (s *Service) GetUserById(id uint64) (*User, error) {
	if id <= 0 {
		return nil, errors.New("ID inválido")
	}
	return s.repo.GetUserById(id)
}

func (s *Service) CreateUser(req CreateUserRequest) error {

	if err := utils.ValidatePassword(req.Password); err != nil {
		return err
	}

	existingUser, err := s.repo.GetUserByEmail(req.Email)

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
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  string(hash),
		CreatedAt: time.Now(),
	}

	return s.repo.CreateUser(&user)
}

func (s *Service) UpdateUser(userID uint64, dto UpdateUserRequest) (*UserResponse, error) {
	user, err := s.repo.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("usuário não encontrado")
	}

	user.Name = dto.Name
	user.Email = dto.Email
	user.Phone = dto.Phone

	err = s.repo.UpdateUser(userID, user)
	if err != nil {
		return nil, err
	}

	response := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		UpdatedAt: time.Now(),
	}

	return &response, nil
}

func (s *Service) DeleteUser(userID uint64) error {
	user, err := s.repo.GetUserById(userID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("usuário não encontrado")
	}
	return s.repo.DeleteUser(userID)
}
