package category

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAllByUser(userID uint64, limit int, offset int) ([]Category, int, error) {

	query := `
	SELECT id, user_id, name, type, created_at, updated_at
	FROM categories
	WHERE user_id = ? OR user_id IS NULL
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, userID, limit, offset)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	categories := []Category{}

	for rows.Next() {

		var c Category
		var userID sql.NullInt64

		err := rows.Scan(
			&c.ID,
			&userID,
			&c.Name,
			&c.Type,
			&c.CreatedAt,
			&c.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		if userID.Valid {
			uid := uint64(userID.Int64)
			c.UserID = &uid
		}

		categories = append(categories, c)
	}

	var total int

	err = r.db.QueryRow("SELECT COUNT(*) FROM categories WHERE user_id = ?", userID).Scan(&total)

	if err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

func (r *Repository) GetByID(id uint64) (*Category, error) {

	query := `
	SELECT *
	FROM categories
	WHERE id = ?
	`

	var category Category
	var userID sql.NullInt64

	err := r.db.QueryRow(query, id).Scan(
		&category.ID,
		&userID,
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
