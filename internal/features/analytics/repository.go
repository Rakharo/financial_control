package analytics

import (
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetDashboardSummary(userID uint64, startDate time.Time, endDate time.Time) (float64, float64, error) {

	query := `
		SELECT
			COALESCE(SUM(
				CASE 
					WHEN type = 'income' 
					THEN CASE 
						WHEN installment_number IS NOT NULL THEN installment_value
						ELSE amount
					END
				END
			),0),

			COALESCE(SUM(
				CASE 
					WHEN type = 'expense' 
					THEN CASE 
						WHEN installment_number IS NOT NULL THEN installment_value
						ELSE amount
					END
				END
			),0)
		FROM transactions
		WHERE user_id = ?
		AND transaction_date >= ?
		AND transaction_date < ?
	`

	var income float64
	var expense float64

	err := r.db.QueryRow(query, userID, startDate, endDate).Scan(
		&income,
		&expense,
	)

	if err != nil {
		return 0, 0, err
	}

	return income, expense, nil
}

func (r *Repository) GetTopCategories(userID uint64, startDate time.Time, endDate time.Time) ([]CategoryUsageDTO, error) {

	query := `
		SELECT 
			c.name,
			SUM(
				CASE 
					WHEN t.installment_number IS NOT NULL THEN t.installment_value
					ELSE t.amount
				END
			) as total
		FROM transactions t
		JOIN categories c ON c.id = t.category_id
		WHERE t.user_id = ?
		AND t.type = 'expense'
		AND t.transaction_date >= ?
		AND t.transaction_date < ?
		GROUP BY c.id
		ORDER BY total DESC
		LIMIT 5
	`

	rows, err := r.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []CategoryUsageDTO

	for rows.Next() {

		var c CategoryUsageDTO

		err := rows.Scan(
			&c.Category,
			&c.Total,
		)

		if err != nil {
			return nil, err
		}

		categories = append(categories, c)
	}

	return categories, nil
}

func (r *Repository) GetDailyExpenses(userID uint64, startDate time.Time, endDate time.Time) ([]DailyExpenseDTO, error) {

	query := `
		SELECT 
			DAY(transaction_date) as day,
			SUM(
				CASE 
					WHEN installment_number IS NOT NULL THEN installment_value
					ELSE amount
				END
			) as total
		FROM transactions
		WHERE user_id = ?
		AND type = 'expense'
		AND transaction_date >= ?
		AND transaction_date < ?
		GROUP BY day
		ORDER BY day
	`

	rows, err := r.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []DailyExpenseDTO

	for rows.Next() {

		var d DailyExpenseDTO

		err := rows.Scan(
			&d.Day,
			&d.Total,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, d)
	}

	return result, nil
}

func (r *Repository) GetYearSummary(userID uint64, year int) (float64, float64, error) {

	query := `
	SELECT
	COALESCE(SUM(
	 CASE WHEN type='income'
	 THEN CASE
	  WHEN installment_number IS NOT NULL THEN installment_value
	  ELSE amount
	 END
	 END),0),

	COALESCE(SUM(
	 CASE WHEN type='expense'
	 THEN CASE
	  WHEN installment_number IS NOT NULL THEN installment_value
	  ELSE amount
	 END
	 END),0)

	FROM transactions
	WHERE user_id = ?
	AND YEAR(transaction_date) = ?
	`

	var income float64
	var expense float64

	err := r.db.QueryRow(query, userID, year).Scan(&income, &expense)

	if err != nil {
		return 0, 0, err
	}

	return income, expense, nil
}

func (r *Repository) GetLifetimeBalance(userID uint64) (float64, float64, error) {

	query := `
	SELECT
	COALESCE(SUM(
	 CASE WHEN type='income'
	 THEN CASE
	  WHEN installment_number IS NOT NULL THEN installment_value
	  ELSE amount
	 END
	 END),0),

	COALESCE(SUM(
	 CASE WHEN type='expense'
	 THEN CASE
	  WHEN installment_number IS NOT NULL THEN installment_value
	  ELSE amount
	 END
	 END),0)

	FROM transactions
	WHERE user_id = ?
	`

	var income float64
	var expense float64

	err := r.db.QueryRow(query, userID).Scan(&income, &expense)

	if err != nil {
		return 0, 0, err
	}

	return income, expense, nil
}

func (r *Repository) GetMonthlyInstallments(userID uint64, startDate time.Time, endDate time.Time) (float64, error) {

	query := `
	SELECT COALESCE(SUM(installment_value),0)
	FROM transactions
	WHERE user_id = ?
	AND installment_number IS NOT NULL
	AND type = 'expense'
	AND transaction_date >= ?
	AND transaction_date < ?
	`

	var total float64

	err := r.db.QueryRow(query, userID, startDate, endDate).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *Repository) GetFutureInstallments(userID uint64) (float64, error) {

	query := `
	SELECT COALESCE(SUM(installment_value),0)
	FROM transactions
	WHERE user_id = ?
	AND installment_number IS NOT NULL
	AND type = 'expense'
	AND transaction_date > CURDATE()
	`

	var total float64

	err := r.db.QueryRow(query, userID).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *Repository) GetBiggestActiveInstallment(userID uint64) (float64, int, error) {

	query := `
	SELECT 
	installment_value,
	COUNT(*) as remaining
	FROM transactions
	WHERE user_id = ?
	AND installment_number IS NOT NULL
	AND type = 'expense'
	AND transaction_date >= CURDATE()
	GROUP BY installment_value
	ORDER BY installment_value DESC
	LIMIT 1
	`

	var value float64
	var remaining int

	err := r.db.QueryRow(query, userID).Scan(&value, &remaining)

	if err != nil {
		return 0, 0, err
	}

	return value, remaining, nil
}
