package transaction

import (
	"errors"
	"time"
)

type TransactionRepository interface {
	GetAllByUser(userID uint64) ([]Transaction, error)
	GetByID(id uint64, userID uint64) (*Transaction, error)
	Create(transaction *Transaction) error
	Update(transaction *Transaction) error
	Delete(id uint64) error
}

type Service struct {
	repo TransactionRepository
}

func NewService(repo TransactionRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll(userID uint64) ([]Transaction, error) {
	return s.repo.GetAllByUser(userID)
}

func (s *Service) GetByID(id uint64, userID uint64) (*Transaction, error) {
	return s.repo.GetByID(id, userID)
}

func (s *Service) Create(userID uint64, dto TransactionRequest) (*Transaction, error) {

	now := time.Now()

	transaction := Transaction{
		UserID:    userID,
		Title:     dto.Title,
		Amount:    dto.Amount,
		Type:      dto.Type,
		Category:  dto.Category,
		Frequency: dto.Frequency,
		CreatedAt: &now,
	}

	err := s.repo.Create(&transaction)

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (s *Service) Update(transactionID uint64, userID uint64, dto TransactionRequest) error {
	transaction, err := s.repo.GetByID(transactionID, userID)

	if err != nil {
		return err
	}

	if transaction == nil {
		return errors.New("Transação não encontrada")
	}

	if userID == 0 {
		return errors.New("Usuário não encontrado")
	}

	if userID != transaction.UserID {
		return errors.New("Usuário não autorizado")
	}

	if dto.Title == "" || dto.Amount <= 0 || dto.Type == "" || dto.Category == "" || dto.Frequency == "" {
		return errors.New("Campos obrigatórios não preenchidos")
	}

	now := time.Now()
	transaction.Title = dto.Title
	transaction.Amount = dto.Amount
	transaction.Type = dto.Type
	transaction.Category = dto.Category
	transaction.Frequency = dto.Frequency
	transaction.UpdatedAt = &now

	return s.repo.Update(transaction)

}

func (s *Service) Delete(id uint64) error {
	return s.repo.Delete(id)
}
