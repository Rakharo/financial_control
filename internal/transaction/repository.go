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

func (r *Repository) Delete(id uint64) error {

	query := `DELETE FROM transactions WHERE id = ?`

	_, err := r.db.Exec(query, id)

	return err
}
