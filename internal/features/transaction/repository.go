package transaction

import (
	"database/sql"
	category "financial_control/internal/features/categories"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAllByUser(userID uint64, limit int, offset int, month int, year int) ([]Transaction, int, error) {

	query := `
	SELECT
	t.id,
	t.user_id,
	t.title,
	t.amount,
	t.type,
	t.frequency,
	t.created_at,
	t.updated_at,
	c.id,
	c.name,
	c.type,
	c.created_at,
	c.updated_at
	FROM transactions t
	LEFT JOIN categories c ON c.id = t.category_id
	WHERE t.user_id = ?
	AND MONTH(t.created_at) = ?
	AND YEAR(t.created_at) = ?
	ORDER BY t.created_at DESC
	LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, userID, month, year, limit, offset)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	transactions := []Transaction{}

	for rows.Next() {
		var transaction Transaction
		var category category.Category

		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.Title,
			&transaction.Amount,
			&transaction.Type,
			&transaction.Frequency,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
			&category.ID,
			&category.Name,
			&category.Type,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		transaction.Category = &category

		if err != nil {
			return nil, 0, err
		}

		transaction.Category = &category
		transactions = append(transactions, transaction)
	}

	var total int

	err = r.db.QueryRow("SELECT COUNT(*) FROM transactions WHERE user_id = ? AND MONTH(created_at) = ? AND YEAR(created_at) = ?", userID, month, year).Scan(&total)

	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *Repository) GetByID(id uint64, userID uint64) (*Transaction, error) {

	query := `
	SELECT
	t.id,
	t.user_id,
	t.title,
	t.amount,
	t.type,
	t.frequency,
	t.created_at,
	t.updated_at,
	c.id,
	c.name,
	c.type,
	c.created_at,
	c.updated_at
	FROM transactions t
	LEFT JOIN categories c ON c.id = t.category_id
	WHERE t.id = ? AND t.user_id = ?
	`

	var transaction Transaction
	var category category.Category

	err := r.db.QueryRow(query, id, userID).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.Title,
		&transaction.Amount,
		&transaction.Type,
		&transaction.Frequency,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&category.ID,
		&category.Name,
		&category.Type,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	transaction.Category = &category

	return &transaction, nil
}

func (r *Repository) Create(transaction *Transaction) error {

	query := `
	INSERT INTO transactions 
	(user_id, title, amount, type, category_id, frequency, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	var categoryID interface{}

	if transaction.Category != nil {
		categoryID = transaction.Category.ID
	}

	result, err := r.db.Exec(
		query,
		transaction.UserID,
		transaction.Title,
		transaction.Amount,
		transaction.Type,
		categoryID,
		transaction.Frequency,
		transaction.CreatedAt,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	transaction.ID = uint64(id)

	return nil
}

func (r *Repository) Update(transaction *Transaction) error {

	query := `
	UPDATE transactions 
	SET title = ?, amount = ?, type = ?, category_id = ?, frequency = ?, updated_at = ?
	WHERE id = ? AND user_id = ?
	`

	var categoryID interface{}

	if transaction.Category != nil {
		categoryID = transaction.Category.ID
	}

	_, err := r.db.Exec(
		query,
		transaction.Title,
		transaction.Amount,
		transaction.Type,
		categoryID,
		transaction.Frequency,
		transaction.UpdatedAt,
		transaction.ID,
		transaction.UserID,
	)

	return err
}

func (r *Repository) Delete(id uint64) error {

	query := `DELETE FROM transactions WHERE id = ?`

	_, err := r.db.Exec(query, id)

	return err
}

func (r *Repository) GetSummaryByUser(userID uint64, month string, year string) (*SummaryDTO, error) {

	query := `
		SELECT
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END),0),
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END),0)
		FROM transactions
		WHERE user_id = ?
	`

	args := []interface{}{userID}

	if month != "" && year != "" {
		query += " AND MONTH(created_at) = ? AND YEAR(created_at) = ?"
		args = append(args, month, year)
	}

	var summary SummaryDTO

	err := r.db.QueryRow(query, args...).Scan(
		&summary.TotalIncome,
		&summary.TotalExpense,
	)

	if err != nil {
		return nil, err
	}

	summary.Balance = summary.TotalIncome - summary.TotalExpense

	return &summary, nil
}
