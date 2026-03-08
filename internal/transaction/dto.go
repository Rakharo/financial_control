package transaction

import category "financial_control/internal/categories"

type TransactionRequest struct {
	Title      string    `json:"title" binding:"required"`
	Amount     float64   `json:"amount" binding:"required"`
	Type       Type      `json:"type" binding:"required"`
	CategoryID uint64    `json:"category_id"`
	Frequency  Frequency `json:"frequency"`
}

type TransactionResponse struct {
	ID        uint64                     `json:"id"`
	Title     string                     `json:"title"`
	Amount    float64                    `json:"amount"`
	Type      Type                       `json:"type"`
	Category  *category.CategoryResponse `json:"category"`
	Frequency Frequency                  `json:"frequency"`
	CreatedAt string                     `json:"created_at,omitempty"`
	UpdatedAt string                     `json:"updated_at,omitempty"`
}

type SummaryDTO struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`
}
