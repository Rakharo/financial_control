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

type AuthRepository interface {
	GetProvidersByUserID(userID uint64) ([]string, error)
}

type Service struct {
	userRepo UserRepository
	authRepo AuthRepository
}

func NewService(userRepo UserRepository, authRepo AuthRepository) *Service {
	return &Service{userRepo: userRepo, authRepo: authRepo}
}

func (s *Service) GetAllUsers() ([]UserResponse, error) {
	users, err := s.userRepo.GetAllUsers()
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

func (s *Service) GetUserById(userID uint64) (*UserResponse, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("usuário não encontrado")
	}

	providers, err := s.authRepo.GetProvidersByUserID(userID)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Providers: &providers,
	}, nil

	// return s.userRepo.GetUserById(userID)
}

func (s *Service) CreateUser(req CreateUserRequest) error {

	if err := utils.ValidatePassword(req.Password); err != nil {
		return err
	}

	existingUser, err := s.userRepo.GetUserByEmail(req.Email)

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

	return s.userRepo.CreateUser(&user)
}

func (s *Service) UpdateUser(userID uint64, dto UpdateUserRequest) (*UserResponse, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("usuário não encontrado")
	}

	user.Name = dto.Name
	user.Email = dto.Email
	user.Phone = dto.Phone

	err = s.userRepo.UpdateUser(userID, user)
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
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("usuário não encontrado")
	}
	return s.userRepo.DeleteUser(userID)
}
