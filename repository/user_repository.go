package repository

import (
	"backend-takehome/models"
	"database/sql"
	"errors"
)

type User interface {
	Create(data *models.User) error
	FindByEmail(email string) (*models.User, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) User {
	return &UserRepository{db}
}

func (u *UserRepository) Create(data *models.User) error {
	query := `INSERT INTO users(name, email, password_hash, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	result, err := u.db.Exec(query, data.Name, data.Email, data.PasswordHash, data.CreatedAt, data.UpdateAt)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows affected")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	data.ID = uint(id)

	return nil
}

func (u *UserRepository) FindByEmail(email string) (*models.User, error) {
	var data models.User

	query := `SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE email = ?`
	if err := u.db.QueryRow(query, email).Scan(&data.ID, &data.Name, &data.Email, &data.PasswordHash, &data.CreatedAt, &data.UpdateAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	return &data, nil
}
