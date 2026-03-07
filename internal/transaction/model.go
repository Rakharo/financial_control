package transaction

import "time"

type Type string
type Frequency string
type Category string

const (
	Income  Type = "income"
	Expense Type = "expense"
)

const (
	Fixed    Frequency = "fixed"
	Variable Frequency = "variable"
)

type Transaction struct {
	ID        uint64
	UserID    uint64
	Title     string
	Amount    float64
	Type      Type
	Category  Category
	Frequency Frequency
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
