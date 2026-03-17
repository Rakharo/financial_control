package user

import (
	"database/sql"
	"errors"

	"financial_control/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetAllUsers() ([]User, error)
	GetUserById(userID uint64) (*User, error)
	GetUserWithPasswordById(id uint64) (*User, error)
	GetUserByLogin(login string) (*User, error)
	CreateUser(user *User) error
	UpdateUser(userID uint64, user *User) error
	UpdateUserPassword(userID uint64, password string) error
	DeleteUser(userID uint64) error
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
	user.Login = dto.Login

	err = s.repo.UpdateUser(userID, user)
	if err != nil {
		return nil, err
	}

	response := UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Login: user.Login,
	}

	return &response, nil
}

func (s *Service) UpdateUserPassword(userID uint64, password PasswordRequest) error {
	user, err := s.repo.GetUserWithPasswordById(userID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("usuário não encontrado")
	}

	//validacao da senha atual
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password.CurrentPassword),
	)
	if err != nil {
		return errors.New("senha incorreta")
	}

	//validacao de senha nova e confirmacao de senha nova
	if password.NewPassword != password.ConfirmPassword {
		return errors.New("as senhas não coincidem")
	}

	//validacao de regra da senha
	if err := utils.ValidatePassword(password.NewPassword); err != nil {
		return err
	}

	//nova criptografia da senha
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password.NewPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	return s.repo.UpdateUserPassword(userID, string(hash))
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
