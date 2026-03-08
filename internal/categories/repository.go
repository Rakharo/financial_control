package category

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAllByUser(userID uint64) ([]Category, error) {

	query := `
	SELECT *
	FROM categories
	WHERE user_id = ? OR user_id IS NULL
	ORDER BY name
	`

	rows, err := r.db.Query(query, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var categories []Category

	for rows.Next() {

		var category Category

		err := rows.Scan(
			&category.ID,
			&category.UserID,
			&category.Name,
			&category.Type,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (r *Repository) GetByID(id uint64) (*Category, error) {

	query := `
	SELECT *
	FROM categories
	WHERE id = ?
	`

	var category Category

	err := r.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&category.Type,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *Repository) Create(category *Category) error {

	query := `
	INSERT INTO categories (user_id, name, type, created_at)
	VALUES (?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		category.UserID,
		category.Name,
		category.Type,
		category.CreatedAt,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	category.ID = uint64(id)

	return nil
}

func (r *Repository) Update(category *Category) error {

	query := `
	UPDATE categories
	SET name = ?, type = ?, updated_at = ?
	WHERE id = ? AND user_id = ?
	`

	_, err := r.db.Exec(
		query,
		category.Name,
		category.Type,
		category.UpdatedAt,
		category.ID,
		category.UserID,
	)

	return err
}

func (r *Repository) Delete(id uint64, userID uint64) error {

	query := `
	DELETE FROM categories
	WHERE id = ? AND user_id = ?
	`

	_, err := r.db.Exec(query, id, userID)

	return err
}

func (r *Repository) GetByNameAndUser(name string, userID uint64) (*Category, error) {

	query := `
	SELECT id, name, user_id
	FROM categories
	WHERE name = ? AND user_id = ?
	LIMIT 1
	`

	row := r.db.QueryRow(query, name, userID)

	var category Category

	err := row.Scan(
		&category.ID,
		&category.Name,
		&category.UserID,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &category, nil
}
