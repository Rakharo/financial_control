package analytics

import (
	"fmt"
	"strconv"
	"time"
)

type AnalyticsRepository interface {
	GetDashboardSummary(userID uint64, startDate time.Time, endDate time.Time) (float64, float64, error)
	GetTopCategories(userID uint64, startDate time.Time, endDate time.Time) ([]CategoryUsageDTO, error)
	GetDailyExpenses(userID uint64, startDate time.Time, endDate time.Time) ([]DailyExpenseDTO, error)
}

type Service struct {
	repo AnalyticsRepository
}

func NewService(repo AnalyticsRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetDashboard(userID uint64, month string, year string) (*DashboardDTO, error) {

	now := time.Now()

	m := int(now.Month())
	y := now.Year()

	if month != "" {
		parsedMonth, err := strconv.Atoi(month)
		if err != nil {
			return nil, fmt.Errorf("mês inválido")
		}
		m = parsedMonth
	}

	if year != "" {
		parsedYear, err := strconv.Atoi(year)
		if err != nil {
			return nil, fmt.Errorf("ano inválido")
		}
		y = parsedYear
	}

	startDate := time.Date(y, time.Month(m), 1, 0, 0, 0, 0, now.Location())
	endDate := startDate.AddDate(0, 1, 0)

	income, expense, err := s.repo.GetDashboardSummary(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	categories, err := s.repo.GetTopCategories(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	dailyExpenses, err := s.repo.GetDailyExpenses(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return &DashboardDTO{
		Income:        income,
		Expenses:      expense,
		Balance:       income - expense,
		TopCategories: categories,
		DailyExpenses: dailyExpenses,
	}, nil
}
