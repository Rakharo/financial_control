package transaction

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAllByUser(userID uint64) ([]Transaction, error) {

	query := `SELECT * FROM transactions WHERE user_id = ?`

	rows, err := r.db.Query(query, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var transactions []Transaction

	for rows.Next() {
		var transaction Transaction

		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.Title,
			&transaction.Amount,
			&transaction.Type,
			&transaction.Category,
			&transaction.Frequency,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *Repository) GetByID(id uint64, userID uint64) (*Transaction, error) {

	query := `SELECT * FROM transactions WHERE id = ? AND user_id = ?`

	var transaction Transaction

	err := r.db.QueryRow(query, id, userID).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.Title,
		&transaction.Amount,
		&transaction.Type,
		&transaction.Category,
		&transaction.Frequency,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *Repository) Create(transaction *Transaction) error {

	query := `
		INSERT INTO transactions 
		(user_id, title, amount, type, category, frequency, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		transaction.UserID,
		transaction.Title,
		transaction.Amount,
		transaction.Type,
		transaction.Category,
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
	query := `UPDATE transactions SET title = ?, amount = ?, type = ?, category = ?, frequency = ?, updated_at = ? WHERE id = ? AND user_id = ?`

	_, err := r.db.Exec(
		query,
		transaction.Title,
		transaction.Amount,
		transaction.Type,
		transaction.Category,
		transaction.Frequency,
		transaction.UpdatedAt,
		transaction.ID,
		transaction.UserID,
	)

	if err != nil {
		return err
	}

	return nil
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
