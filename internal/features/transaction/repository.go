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

type installmentRow struct {
	installmentID          sql.NullInt64
	installmentTotalAmount sql.NullFloat64
	installments           sql.NullInt64
	installmentDate        sql.NullTime
	installmentCreatedAt   sql.NullTime
	installmentUpdatedAt   sql.NullTime
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
	t.installment_plan_id,
	t.installment_number,
	t.installment_total,
	t.installment_value,
	t.transaction_date,
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
	AND MONTH(t.transaction_date) = ?
	AND YEAR(t.transaction_date) = ?
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
			&transaction.InstallmentPlanID,
			&transaction.InstallmentNumber,
			&transaction.InstallmentTotal,
			&transaction.InstallmentValue,
			&transaction.TransactionDate,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
			&category.ID,
			&category.Name,
			&category.Type,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		transaction.Category = &category
		transactions = append(transactions, transaction)
	}

	var total int

	err = r.db.QueryRow("SELECT COUNT(*) FROM transactions WHERE user_id = ? AND MONTH(transaction_date) = ? AND YEAR(transaction_date) = ?", userID, month, year).Scan(&total)

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
	t.installment_plan_id,
	t.installment_number,
	t.installment_total,
	t.installment_value,
	t.transaction_date,
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
		&transaction.InstallmentPlanID,
		&transaction.InstallmentNumber,
		&transaction.InstallmentTotal,
		&transaction.InstallmentValue,
		&transaction.TransactionDate,
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
	(user_id, title, amount, type, category_id, frequency,
	installment_plan_id, installment_number, installment_total,
	installment_value, transaction_date, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
		transaction.InstallmentPlanID,
		transaction.InstallmentNumber,
		transaction.InstallmentTotal,
		transaction.InstallmentValue,
		transaction.TransactionDate,
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
