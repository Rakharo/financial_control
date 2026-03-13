package transaction

import (
	"errors"
	category "financial_control/internal/features/categories"
	"financial_control/internal/features/installment"
	"financial_control/internal/utils"
	"math"
	"time"
)

type TransactionRepository interface {
	GetAllByUser(userID uint64, limit int, offset int, month int, year int) ([]Transaction, int, error)
	GetByID(id uint64, userID uint64) (*Transaction, error)
	Create(transaction *Transaction) error
	Update(transaction *Transaction) error
	Delete(id uint64) error
	GetSummaryByUser(userID uint64, month string, year string) (*SummaryDTO, error)
}
type CategoryRepository interface {
	GetByID(id uint64) (*category.Category, error)
}

type InstallmentService interface {
	Create(userID uint64, totalAmount float64, createdAt time.Time, dto *installment.InstallmentRequest) (*installment.InstallmentResponse, error)
}

type Service struct {
	repo               TransactionRepository
	categoryRepo       CategoryRepository
	installmentService InstallmentService
}

func NewService(repo TransactionRepository, categoryRepo CategoryRepository, installmentService InstallmentService) *Service {
	return &Service{
		repo:               repo,
		categoryRepo:       categoryRepo,
		installmentService: installmentService,
	}
}

func (s *Service) GetAllTransactions(userID uint64, page int, limit int, month int, year int) ([]TransactionResponse, int, error) {

	offset := (page - 1) * limit

	transactions, total, err := s.repo.GetAllByUser(userID, limit, offset, month, year)
	if err != nil {
		return nil, 0, err
	}

	var response []TransactionResponse

	for _, t := range transactions {
		response = append(response, ToTransactionResponse(t))
	}

	return response, total, nil
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
	date, err := utils.ParseDate(dto.TransactionDate)
	if err != nil {
		return nil, err
	}

	//Criação da transaction com parcelamento
	if dto.InstallmentTotal > 1 {
		installmentDTO := installment.InstallmentRequest{
			UserID:          userID,
			TotalAmount:     dto.Amount,
			Installments:    dto.InstallmentTotal,
			InstallmentDate: &date,
		}

		installmentPlan, err := s.installmentService.Create(
			userID,
			dto.Amount,
			now,
			&installmentDTO,
		)

		if err != nil {
			return nil, err
		}

		installmentValue := math.Floor((dto.Amount/float64(dto.InstallmentTotal))*100) / 100

		remaining := dto.Amount

		for i := 1; i <= dto.InstallmentTotal; i++ {
			installmentDate := utils.AddMonthsSafe(date, i-1)

			value := installmentValue

			if i == dto.InstallmentTotal {
				value = remaining
			}

			remaining -= value

			number := i
			total := dto.InstallmentTotal
			planID := installmentPlan.ID

			transaction := Transaction{
				UserID:            userID,
				Title:             dto.Title,
				Amount:            dto.Amount,
				Type:              dto.Type,
				Category:          category,
				Frequency:         dto.Frequency,
				InstallmentPlanID: &planID,
				InstallmentNumber: &number,
				InstallmentTotal:  &total,
				InstallmentValue:  &value,
				TransactionDate:   &installmentDate,
				CreatedAt:         &now,
			}

			err := s.repo.Create(&transaction)
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	}

	transaction := Transaction{
		UserID:          userID,
		Title:           dto.Title,
		Amount:          dto.Amount,
		Type:            dto.Type,
		Category:        category,
		Frequency:       dto.Frequency,
		TransactionDate: &date,
		CreatedAt:       &now,
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
