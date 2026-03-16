package analytics

import (
	"strconv"
	"time"
)

type AnalyticsRepository interface {
	GetDashboardSummary(userID uint64, startDate time.Time, endDate time.Time) (float64, float64, error)
	GetTopCategories(userID uint64, startDate time.Time, endDate time.Time) ([]CategoryUsageDTO, error)
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

	m, _ := strconv.Atoi(month)
	y, _ := strconv.Atoi(year)

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
		return nil, err
	}

	yearIncome, yearExpense, err := s.repo.GetYearSummary(userID, y)
	if err != nil {
		return nil, err
	}

	lifetimeIncome, lifetimeExpense, err := s.repo.GetLifetimeBalance(userID)
	if err != nil {
		return nil, err
	}

	monthlyInstallments, err := s.repo.GetMonthlyInstallments(userID, startMonth, endMonth)
	if err != nil {
		return nil, err
	}

	futureInstallments, err := s.repo.GetFutureInstallments(userID)
	if err != nil {
		return nil, err
	}

	biggestInstallment, remainingInstallments, err := s.repo.GetBiggestActiveInstallment(userID)
	if err != nil {
		return nil, err
	}

	topCategories, _ := s.repo.GetTopCategories(userID, startMonth, endMonth)
	dailyExpenses, _ := s.repo.GetDailyExpenses(userID, startMonth, endMonth)

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
