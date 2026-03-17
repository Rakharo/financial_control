package installment

import "time"

type InstallmentRepository interface {
	Create(installment *Installment) error
	Delete(id uint64) error
}

type Service struct {
	repo InstallmentRepository
}

func NewService(repo InstallmentRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(userID uint64, totalAmount float64, createdAt time.Time, dto *InstallmentRequest) (*InstallmentResponse, error) {

	installmentValue := totalAmount / float64(dto.Installments)

	installment := Installment{
		UserID:           userID,
		TotalAmount:      totalAmount,
		Installments:     dto.Installments,
		InstallmentValue: installmentValue,
		InstallmentDate:  dto.InstallmentDate,
		CreatedAt:        &createdAt,
	}

	err := s.repo.Create(&installment)

	if err != nil {
		return nil, err
	}

	response := ToInstallmentResponse(installment)

	return &response, nil

}

func (s *Service) DeleteByID(id uint64) error {
	return s.repo.Delete(id)
}
