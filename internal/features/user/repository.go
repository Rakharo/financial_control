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

	rows, err := r.db.Query("SELECT id, name, email, login FROM users")
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var alb User
		if err := rows.Scan(&alb.ID, &alb.Name, &alb.Email, &alb.Login); err != nil {
			return nil, fmt.Errorf("GetAllUsers: %v", err)
		}
		users = append(users, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllUsers: %v", err)
	}
	return users, nil
}

func (r *Repository) GetUserByLogin(login string) (*User, error) {

	query := `
	SELECT id, name, email, login, password
	FROM users
	WHERE login = ?
	`

	row := r.db.QueryRow(query, login)

	var user User

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Login,
		&user.Password,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetUserById(userID uint64) (*User, error) {
	query := "SELECT id, name, email, login FROM users WHERE id = ?"

	row := r.db.QueryRow(query, userID)

	var u User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // usuário não encontrado
		}
		return nil, err
	}

	return &u, nil
}

func (r *Repository) GetUserWithPasswordById(userID uint64) (*User, error) {
	query := `
	SELECT id, name, email, login, password
	FROM users
	WHERE id = ?
	`

	row := r.db.QueryRow(query, userID)

	var user User

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Login,
		&user.Password,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) CreateUser(user *User) error {

	query := `
	INSERT INTO users (name, email, login, password)
	VALUES (?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query,
		user.Name,
		user.Email,
		user.Login,
		user.Password,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateUser(userID uint64, user *User) error {
	query := "UPDATE users SET name = ?, email = ?, login = ? WHERE id = ?"
	_, err := r.db.Exec(
		query,
		user.Name,
		user.Email,
		user.Login,
		userID,
	)

	return err
}

func (r *Repository) UpdateUserPassword(userID uint64, password string) error {
	query := "UPDATE users SET password = ? WHERE id = ?"

	_, err := r.db.Exec(query, password, userID)
	return err
}

func (r *Repository) DeleteUser(userID uint64) error {
	query := "DELETE FROM users WHERE id = ?"

	_, err := r.db.Exec(query, userID)
	return err
}
