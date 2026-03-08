package connection

import "database/sql"

func SeedCategories(db *sql.DB) error {

	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM categories WHERE user_id IS NULL").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	query := `
	INSERT INTO categories (user_id, name, type, created_at)
	VALUES
	(NULL, 'salário', 'income', NOW()),
	(NULL, 'freelance', 'income', NOW()),
	(NULL, 'investimentos', 'income', NOW()),
	(NULL, 'alimentação', 'expense', NOW()),
	(NULL, 'transporte', 'expense', NOW()),
	(NULL, 'moradia', 'expense', NOW()),
	(NULL, 'saúde', 'expense', NOW()),
	(NULL, 'educação', 'expense', NOW()),
	(NULL, 'entretenimento', 'expense', NOW()),
	(NULL, 'outros', 'expense', NOW())
	`

	_, err = db.Exec(query)

	return err
}
