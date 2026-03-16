package analytics

type DashboardDTO struct {
	Income        float64            `json:"income"`
	Expenses      float64            `json:"expenses"`
	Balance       float64            `json:"balance"`
	TopCategories []CategoryUsageDTO `json:"top_categories"`
	DailyExpenses []DailyExpenseDTO  `json:"daily_expenses"`
}

type CategoryUsageDTO struct {
	Category string
	Total    float64
}

type DailyExpenseDTO struct {
	Day   int
	Total float64
}
