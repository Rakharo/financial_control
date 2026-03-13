package installment

import "time"

type InstallmentRequest struct {
	UserID          uint64     `json:"user_id" binding:"required"`
	TotalAmount     float64    `json:"total_amount" binding:"required"`
	Installments    int        `json:"installment_number" binding:"required"`
	InstallmentDate *time.Time `json:"installment_date" binding:"required"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type InstallmentResponse struct {
	ID               uint64  `json:"id"`
	TotalAmount      float64 `json:"total_amount"`
	Installments     int     `json:"installments"`
	InstallmentValue float64 `json:"installment_value"`
	InstallmentDate  string  `json:"installment_date"`
	CreatedAt        string  `json:"created_at,omitempty"`
	UpdatedAt        string  `json:"updated_at,omitempty"`
}
