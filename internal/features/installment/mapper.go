package installment

import (
	"time"
)

func ToInstallmentResponse(i Installment) InstallmentResponse {

	var updatedAt string
	if i.UpdatedAt != nil {
		updatedAt = i.UpdatedAt.Format(time.RFC3339)
	}

	var createdAt string

	if i.CreatedAt != nil {
		createdAt = i.CreatedAt.Format(time.RFC3339)
	}

	var installmentDate string

	if i.InstallmentDate != nil {
		installmentDate = i.InstallmentDate.Format(time.RFC3339)
	}

	return InstallmentResponse{
		ID:              i.ID,
		TotalAmount:     i.TotalAmount,
		Installments:    i.Installments,
		InstallmentDate: installmentDate,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
}
