package repository

import (
	"backend-takehome/models"
	"database/sql"
	"errors"
)

type Post interface {
	Create(data *models.Post) error
}

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) Post {
	return &PostRepository{db}
}

func (p *PostRepository) Create(data *models.Post) error {
	query := `INSERT INTO posts(title, content, author_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	result, err := p.db.Exec(query, data.Title, data.Content, data.AuthorID, data.CreatedAt, data.UpdatedAt)
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
