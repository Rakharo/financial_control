package analytics

type DashboardDTO struct {
	Month           BalanceDTO             `json:"month"`
	Year            BalanceDTO             `json:"year"`
	LifetimeBalance float64                `json:"lifetime_balance"`
	Installments    InstallmentInsightsDTO `json:"installments"`
	TopCategories   []CategoryUsageDTO     `json:"top_categories"`
	DailyExpenses   []DailyExpenseDTO      `json:"daily_expenses"`
}

type CategoryUsageDTO struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
	Color    string  `json:"color"`
	UserID   *uint64 `json:"user_id"`
}

type DailyExpenseDTO struct {
	Day   int     `json:"day"`
	Total float64 `json:"total"`
}

type BalanceDTO struct {
	Balance float64 `json:"balance"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

type InstallmentInsightsDTO struct {
	MonthlyInstallments   float64 `json:"monthly_installments"`
	FutureInstallments    float64 `json:"future_installments"`
	BiggestInstallment    float64 `json:"biggest_installment"`
	RemainingInstallments int     `json:"remaining_installments"`
}
