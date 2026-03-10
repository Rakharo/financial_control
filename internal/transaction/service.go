package transaction

import (
	"errors"
	category "financial_control/internal/categories"
	"time"
)

type TransactionRepository interface {
	GetAllByUser(userID uint64) ([]Transaction, error)
	GetByID(id uint64, userID uint64) (*Transaction, error)
	Create(transaction *Transaction) error
	Update(transaction *Transaction) error
	Delete(id uint64) error
	GetSummaryByUser(userID uint64, month string, year string) (*SummaryDTO, error)
}
type CategoryRepository interface {
	GetByID(id uint64) (*category.Category, error)
}

type Service struct {
	repo         TransactionRepository
	categoryRepo CategoryRepository
}

func NewService(repo TransactionRepository, categoryRepo CategoryRepository) *Service {
	return &Service{
		repo:         repo,
		categoryRepo: categoryRepo,
	}
}

func (s *Service) GetAllTransactions(userID uint64) ([]TransactionResponse, error) {

	transactions, err := s.repo.GetAllByUser(userID)
	if err != nil {
		return nil, err
	}

	var response []TransactionResponse

	for _, t := range transactions {
		response = append(response, ToTransactionResponse(t))
	}

	return response, nil
}

func (s *Service) GetByID(id uint64, userID uint64) (*TransactionResponse, error) {
	transaction, err := s.repo.GetByID(id, userID)

	if err != nil {
		return nil, err
	}

	if transaction == nil {
		return nil, errors.New("transação não encontrada")
	}

	response := ToTransactionResponse(*transaction)

	return &response, nil
}

func (s *Service) Create(userID uint64, dto TransactionRequest) (*TransactionResponse, error) {

	category, err := s.categoryRepo.GetByID(dto.CategoryID)

	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errors.New("categoria não encontrada")
	}

	// categoria de outro usuário
	if category.UserID != nil && *category.UserID != userID {
		return nil, errors.New("categoria inválida")
	}

	now := time.Now()

	transaction := Transaction{
		UserID:    userID,
		Title:     dto.Title,
		Amount:    dto.Amount,
		Type:      dto.Type,
		Category:  category,
		Frequency: dto.Frequency,
		CreatedAt: &now,
	}

	err = s.repo.Create(&transaction)

	if err != nil {
		return nil, err
	}

	response := ToTransactionResponse(transaction)

	return &response, nil
}

func (s *Service) Update(transactionID uint64, userID uint64, dto TransactionRequest) (*TransactionResponse, error) {

	transaction, err := s.repo.GetByID(transactionID, userID)
	if err != nil {
		return nil, err
	}

	if transaction == nil {
		return nil, errors.New("transação não encontrada")
	}

	category, err := s.categoryRepo.GetByID(dto.CategoryID)

	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errors.New("categoria não encontrada")
	}

	if category.UserID != nil && *category.UserID != userID {
		return nil, errors.New("categoria inválida")
	}

	now := time.Now()

	transaction.Title = dto.Title
	transaction.Amount = dto.Amount
	transaction.Type = dto.Type
	transaction.Category = category
	transaction.Frequency = dto.Frequency
	transaction.UpdatedAt = &now

	err = s.repo.Update(transaction)
	if err != nil {
		return nil, err
	}

	response := ToTransactionResponse(*transaction)

	return &response, nil
}

func (s *Service) Delete(id uint64) error {
	return s.repo.Delete(id)
}

func (s *Service) GetSummary(userID uint64, month string, year string) (*SummaryDTO, error) {
	return s.repo.GetSummaryByUser(userID, month, year)
}
