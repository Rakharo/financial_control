package payment

type Payment struct {
	TotalAmount      float64
	FixedExpenses    float64
	VariableExpenses float64
	FixedIncome      float64
	VariableIncome   float64
}

func (p Payment) TotalIncome() float64 {
	return p.FixedIncome + p.VariableIncome
}

func (p Payment) TotalExpenses() float64 {
	return p.FixedExpenses + p.VariableExpenses
}

func (p Payment) CalculateFinalAmount() float64 {
	return p.TotalAmount + p.TotalIncome() - p.TotalExpenses()
}
