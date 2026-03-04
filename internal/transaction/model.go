package transaction

type Type string
type Category string

const (
	Income  Type = "income"
	Expense Type = "expense"

	Fixed    Category = "fixed"
	Variable Category = "variable"
)

type Transaction struct {
	ID       uint64
	UserID   uint64
	Type     Type
	Category Category
	Amount   float32
}
