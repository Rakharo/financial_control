package installment

import "time"

type Installment struct {
	ID               uint64
	UserID           uint64
	TotalAmount      float64
	Installments     int
	InstallmentValue float64
	InstallmentDate  *time.Time
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
}
