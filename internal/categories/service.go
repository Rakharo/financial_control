package category

import (
	"errors"
	"strings"
	"time"
)

type CategoryRepository interface {
	GetAllByUser(userID uint64, limit int, offset int) ([]Category, int, error)
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

func (s *Service) GetAll(userID uint64, page int, limit int) ([]CategoryResponse, int, error) {
	offset := (page - 1) * limit

	categories, total, err := s.repo.GetAllByUser(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var response []CategoryResponse

	for _, c := range categories {
		response = append(response, ToCategoryResponse(c))
	}

	return response, total, nil
}

func (s *Service) GetByID(id uint64) (*CategoryResponse, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errors.New("categoria não encontrada")
	}

	response := ToCategoryResponse(*category)

	return &response, nil
}

func (s *Service) Create(userID uint64, dto CategoryRequest) (*CategoryResponse, error) {

	dto.Name = strings.ToLower(strings.TrimSpace(dto.Name))

	existing, err := s.repo.GetByNameAndUser(dto.Name, userID)

	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, errors.New("categoria já existe")
	}

	now := time.Now()

	category := Category{
		Name:      dto.Name,
		Type:      dto.Type,
		UserID:    &userID,
		CreatedAt: &now,
	}

	err = s.repo.Create(&category)

	if err != nil {
		return nil, err
	}

	response := ToCategoryResponse(category)

	return &response, nil
}

func (s *Service) Update(id uint64, userID uint64, dto CategoryRequest) (*CategoryResponse, error) {

	dto.Name = strings.ToLower(strings.TrimSpace(dto.Name))

	existing, err := s.repo.GetByNameAndUser(dto.Name, userID)

	if err != nil {
		return nil, err
	}

	if existing != nil && existing.ID != id {
		return nil, errors.New("categoria já existe")
	}

	category, err := s.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errors.New("categoria não encontrada")
	}

	if category.UserID == nil {
		return nil, errors.New("categoria do sistema não pode ser editada")
	}

	if *category.UserID != userID {
		return nil, errors.New("não autorizado")
	}

	now := time.Now()

	category.Name = dto.Name
	category.Type = dto.Type
	category.UpdatedAt = &now

	err = s.repo.Update(category)
	if err != nil {
		return nil, err
	}

	response := ToCategoryResponse(*category)

	return &response, nil
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
