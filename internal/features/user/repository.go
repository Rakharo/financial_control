package user

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAllUsers() ([]User, error) {
	var users []User

	rows, err := r.db.Query("SELECT id, name, email, phone, created_at, updated_at FROM users")
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var alb User
		if err := rows.Scan(&alb.ID, &alb.Name, &alb.Email, &alb.Phone, &alb.CreatedAt, &alb.UpdatedAt); err != nil {
			return nil, fmt.Errorf("GetAllUsers: %v", err)
		}
		users = append(users, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllUsers: %v", err)
	}
	return users, nil
}

func (r *Repository) GetUserByEmail(email string) (*User, error) {

	query := `
	SELECT id, name, email, phone, password, created_at, updated_at
	FROM users
	WHERE email = ?
	`

	row := r.db.QueryRow(query, email)

	var user User

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetUserById(userID uint64) (*User, error) {
	query := "SELECT id, name, email, phone, created_at, updated_at FROM users WHERE id = ?"

	row := r.db.QueryRow(query, userID)

	var u User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // usuário não encontrado
		}
		return nil, err
	}

	return &u, nil
}

func (r *Repository) CreateUser(user *User) error {
	query := `
	INSERT INTO users (name, email, password, phone)
	VALUES (?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		user.Name,
		user.Email,
		user.Password,
		user.Phone,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = uint64(id)

	return nil
}

func (r *Repository) UpdateUser(userID uint64, user *User) error {
	query := "UPDATE users SET name = ?, email = ?, phone = ? WHERE id = ?"
	_, err := r.db.Exec(
		query,
		user.Name,
		user.Email,
		user.Phone,
		userID,
	)

	return err
}

func (r *Repository) DeleteUser(userID uint64) error {
	query := "DELETE FROM users WHERE id = ?"

	_, err := r.db.Exec(query, userID)
	return err
}
