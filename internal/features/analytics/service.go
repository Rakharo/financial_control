package analytics

import (
	"financial_control/internal/utils"
	"strconv"
	"time"
)

type AnalyticsRepository interface {
	GetDashboardSummary(userID uint64, startDate time.Time, endDate time.Time) (float64, float64, error)
	GetTopCategories(userID uint64, startDate time.Time, endDate time.Time) ([]CategoryUsageDTO, error)
	GetTransactionByCategory(userID uint64, categoryID uint64, startDate time.Time, endDate time.Time) ([]TransactionDTO, error)
	GetDailyExpenses(userID uint64, startDate time.Time, endDate time.Time) ([]DailyExpenseDTO, error)
	GetYearSummary(userID uint64, year int) (float64, float64, error)
	GetLifetimeBalance(userID uint64) (float64, float64, error)
	GetMonthlyInstallments(userID uint64, startDate time.Time, endDate time.Time) (float64, error)
	GetFutureInstallments(userID uint64) (float64, error)
	GetBiggestActiveInstallment(userID uint64) (float64, int, error)
}

type Service struct {
	repo AnalyticsRepository
}

func NewService(repo AnalyticsRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetDashboard(userID uint64, month string, year string) (*DashboardDTO, error) {

	now := time.Now()

	m, err := strconv.Atoi(month)
	if err != nil {
		return nil, utils.NewBadRequest("Mês inválido", err)
	}

	y, err := strconv.Atoi(year)
	if err != nil {
		return nil, utils.NewBadRequest("Ano inválido", err)
	}

	if m == 0 {
		m = int(now.Month())
	}

	if y == 0 {
		y = now.Year()
	}

	startMonth := time.Date(y, time.Month(m), 1, 0, 0, 0, 0, now.Location())
	endMonth := startMonth.AddDate(0, 1, 0)

	monthIncome, monthExpense, err := s.repo.GetDashboardSummary(userID, startMonth, endMonth)
	if err != nil {
		return nil, utils.NewInternal(err)
	}

	yearIncome, yearExpense, err := s.repo.GetYearSummary(userID, y)
	if err != nil {
		return nil, utils.NewInternal(err)
	}

	lifetimeIncome, lifetimeExpense, err := s.repo.GetLifetimeBalance(userID)
	if err != nil {
		return nil, utils.NewInternal(err)
	}

	monthlyInstallments, err := s.repo.GetMonthlyInstallments(userID, startMonth, endMonth)
	if err != nil {
		return nil, utils.NewInternal(err)
	}

	futureInstallments, err := s.repo.GetFutureInstallments(userID)
	if err != nil {
		return nil, utils.NewInternal(err)
	}

	biggestInstallment, remainingInstallments, err := s.repo.GetBiggestActiveInstallment(userID)
	if err != nil {
		return nil, utils.NewInternal(err)
	}

	topCategories, err := s.repo.GetTopCategories(userID, startMonth, endMonth)
	if err != nil {
		return nil, utils.NewInternal(err)
	}

	for i, category := range topCategories {

		transactions, err := s.repo.GetTransactionByCategory(
			userID,
			category.CategoryID,
			startMonth,
			endMonth,
		)
		if err != nil {
			return nil, utils.NewInternal(err)
		}

		topCategories[i].Transactions = transactions
	}

	dailyExpenses, err := s.repo.GetDailyExpenses(userID, startMonth, endMonth)
	if err != nil {
		return nil, utils.NewInternal(err)
	}

	return &DashboardDTO{
		Month: BalanceDTO{
			Income:  monthIncome,
			Expense: monthExpense,
			Balance: monthIncome - monthExpense,
		},

		Year: BalanceDTO{
			Income:  yearIncome,
			Expense: yearExpense,
			Balance: yearIncome - yearExpense,
		},

		LifetimeBalance: lifetimeIncome - lifetimeExpense,

		Installments: InstallmentInsightsDTO{
			MonthlyInstallments:   monthlyInstallments,
			FutureInstallments:    futureInstallments,
			BiggestInstallment:    biggestInstallment,
			RemainingInstallments: remainingInstallments,
		},

		TopCategories: topCategories,
		DailyExpenses: dailyExpenses,
	}, nil
}
