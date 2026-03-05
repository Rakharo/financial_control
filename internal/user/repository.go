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

	rows, err := r.db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var alb User
		if err := rows.Scan(&alb.ID, &alb.Name, &alb.Email); err != nil {
			return nil, fmt.Errorf("GetAllUsers: %v", err)
		}
		users = append(users, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllUsers: %v", err)
	}
	return users, nil
}

func (r *Repository) GetUserById(id int64) (*User, error) {
	query := "SELECT id, name, email FROM users WHERE id = ?"

	row := r.db.QueryRow(query, id)

	var u User
	err := row.Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // usuário não encontrado
		}
		return nil, err
	}

	return &u, nil
}

func (r *Repository) CreateUser(name string, email string) error {
	query := "INSERT INTO users (name, email) VALUES (?, ?)"
	_, err := r.db.Exec(query, name, email)
	if err != nil {
		return fmt.Errorf("createUser %q: %v", name, err)
	}
	return nil
}
