package transaction

import "time"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(userID uint64, dto CreateTransaction) (*Transaction, error) {

	transaction := Transaction{
		UserID:    userID,
		Title:     dto.Title,
		Amount:    dto.Amount,
		Type:      dto.Type,
		Category:  dto.Category,
		Frequency: dto.Frequency,
		CreatedAt: time.Now(),
	}

	err := s.repo.Create(&transaction)

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (s *Service) GetAll(userID uint64) ([]Transaction, error) {
	return s.repo.GetAllByUser(userID)
}

func (s *Service) GetByID(id uint64, userID uint64) (*Transaction, error) {
	return s.repo.GetByID(id, userID)
}

func (s *Service) Delete(id uint64) error {
	return s.repo.Delete(id)
}
