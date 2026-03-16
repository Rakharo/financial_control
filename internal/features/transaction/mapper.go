package transaction

import (
	category "financial_control/internal/features/categories"
	"time"
)

func ToTransactionResponse(t Transaction) TransactionResponse {

	var categoryResponse *category.CategoryResponse

	if t.Category != nil {
		c := category.ToCategoryResponse(*t.Category)
		categoryResponse = &c
	}

	var updatedAt string
	if t.UpdatedAt != nil {
		updatedAt = t.UpdatedAt.Format(time.RFC3339)
	}

	var createdAt string

	if t.CreatedAt != nil {
		createdAt = t.CreatedAt.Format(time.RFC3339)
	}

	var transactionDate string

	if t.TransactionDate != nil {
		transactionDate = t.TransactionDate.Format(time.RFC3339)
	}

	return TransactionResponse{
		ID:                t.ID,
		Title:             t.Title,
		Amount:            t.Amount,
		Type:              t.Type,
		Category:          categoryResponse,
		Frequency:         t.Frequency,
		InstallmentPlanID: t.InstallmentPlanID,
		InstallmentNumber: t.InstallmentNumber,
		InstallmentTotal:  t.InstallmentTotal,
		InstallmentValue:  t.InstallmentValue,
		TransactionDate:   transactionDate,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
}
