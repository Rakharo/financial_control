package auth

import (
	"database/sql"
	"financial_control/internal/features/user"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) UpdateUserPassword(userID uint64, password string) error {
	query := "UPDATE users SET password = ? WHERE id = ?"

	_, err := r.db.Exec(query, password, userID)
	return err
}

func (r *Repository) GetUserWithPasswordById(userID uint64) (*user.User, error) {
	query := `
	SELECT id, name, email, login, password
	FROM users
	WHERE id = ?
	`

	row := r.db.QueryRow(query, userID)

	var user user.User

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Password,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetUserProvider(providerName, providerUserID string) (*UserProvider, error) {
	query := `
	SELECT id, user_id, provider_name, provider_user_id, created_at
	FROM user_providers
	WHERE provider_name = ? AND provider_user_id = ?
	`
	row := r.db.QueryRow(query, providerName, providerUserID)

	var up UserProvider
	err := row.Scan(&up.ID, &up.UserID, &up.ProviderName, &up.ProviderUserID, &up.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &up, nil
}

func (r *Repository) GetProvidersByUserID(userID uint64) ([]string, error) {
	query := `
	SELECT provider_name
	FROM user_providers
	WHERE user_id = ?
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []string

	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		providers = append(providers, p)
	}

	return providers, nil
}

func (r *Repository) CreateUserProvider(up *UserProvider) error {
	query := `
	INSERT INTO user_providers (user_id, provider_name, provider_user_id, created_at)
	VALUES (?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, up.UserID, up.ProviderName, up.ProviderUserID, up.CreatedAt)
	return err
}

func (r *Repository) CreateRefreshToken(rt *RefreshToken) error {
	query := `
	INSERT INTO refresh_tokens (user_id, token, expires_at, created_at)
	VALUES (?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, rt.UserID, rt.Token, rt.ExpiresAt, rt.CreatedAt)
	return err
}

func (r *Repository) GetRefreshToken(token string) (*RefreshToken, error) {
	query := `
	SELECT id, user_id, token, expires_at, created_at
	FROM refresh_tokens
	WHERE token = ?
	`
	row := r.db.QueryRow(query, token)

	var rt RefreshToken
	err := row.Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Verifica se já expirou
	if time.Now().After(rt.ExpiresAt) {
		return nil, nil
	}

	return &rt, nil
}

func (r *Repository) DeleteRefreshToken(token string) error {
	query := "DELETE FROM refresh_tokens WHERE token = ?"
	_, err := r.db.Exec(query, token)
	return err
}
