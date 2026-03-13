package installment

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(installment *Installment) error {
	query := `
	INSERT INTO installment_plans (user_id, total_amount, installments, installment_date, created_at)
	VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		installment.UserID,
		installment.TotalAmount,
		installment.Installments,
		installment.InstallmentDate,
		installment.CreatedAt,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	installment.ID = uint64(id)

	return nil
}
