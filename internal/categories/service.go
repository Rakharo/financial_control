package category

import (
	"errors"
	"strings"
	"time"
)

type CategoryRepository interface {
	GetAllByUser(userID uint64) ([]Category, error)
	GetByID(id uint64) (*Category, error)
	Create(category *Category) error
	Update(category *Category) error
	Delete(id uint64, userID uint64) error
	GetByNameAndUser(name string, userID uint64) (*Category, error)
}

type Service struct {
	repo CategoryRepository
}

func NewService(repo CategoryRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll(userID uint64) ([]Category, error) {
	return s.repo.GetAllByUser(userID)
}

func (s *Service) GetByID(id uint64) (*Category, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Create(userID uint64, dto CategoryRequest) (*Category, error) {

	dto.Name = strings.ToLower(strings.TrimSpace(dto.Name))

	existing, err := s.repo.GetByNameAndUser(dto.Name, userID)

	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, errors.New("categoria já existe")
	}

	category := Category{
		Name:   dto.Name,
		UserID: &userID,
	}

	err = s.repo.Create(&category)

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (s *Service) Update(id uint64, userID uint64, dto CategoryRequest) error {

	dto.Name = strings.ToLower(strings.TrimSpace(dto.Name))

	existing, err := s.repo.GetByNameAndUser(dto.Name, userID)

	if err != nil {
		return err
	}

	if existing != nil && existing.ID != id {
		return errors.New("categoria já existe")
	}

	category, err := s.repo.GetByID(id)

	if err != nil {
		return err
	}

	if category == nil {
		return errors.New("categoria não encontrada")
	}

	if category.UserID == nil {
		return errors.New("categoria do sistema não pode ser editada")
	}

	if *category.UserID != userID {
		return errors.New("não autorizado")
	}

	now := time.Now()

	category.Name = dto.Name
	category.Type = dto.Type
	category.UpdatedAt = &now

	return s.repo.Update(category)
}

func (s *Service) Delete(id uint64, userID uint64) error {

	category, err := s.repo.GetByID(id)

	if err != nil {
		return err
	}

	if category.UserID == nil {
		return errors.New("categoria do sistema não pode ser deletada")
	}

	if *category.UserID != userID {
		return errors.New("não autorizado")
	}

	return s.repo.Delete(id, userID)
}
