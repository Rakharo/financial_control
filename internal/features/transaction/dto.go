package transaction

import (
	category "financial_control/internal/features/categories"
)

type TransactionRequest struct {
	Title            string    `json:"title" binding:"required"`
	Amount           float64   `json:"amount" binding:"required"`
	Type             Type      `json:"type" binding:"required"`
	CategoryID       uint64    `json:"category_id"`
	Frequency        Frequency `json:"frequency"`
	TransactionDate  string    `json:"transaction_date"`
	InstallmentTotal int       `json:"installment_total"`
}

type TransactionResponse struct {
	ID                uint64                     `json:"id"`
	Title             string                     `json:"title"`
	Amount            float64                    `json:"amount"`
	Type              Type                       `json:"type"`
	Category          *category.CategoryResponse `json:"category"`
	Frequency         Frequency                  `json:"frequency"`
	InstallmentPlanID *uint64                    `json:"installment_plan_id"`
	InstallmentNumber *int                       `json:"installment_number"`
	InstallmentTotal  *int                       `json:"installment_total"`
	InstallmentValue  *float64                   `json:"installment_value"`
	TransactionDate   string                     `json:"transaction_date"`
	CreatedAt         string                     `json:"created_at,omitempty"`
	UpdatedAt         string                     `json:"updated_at,omitempty"`
}
