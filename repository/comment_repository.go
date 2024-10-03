package repository

import (
	"backend-takehome/models"
	"database/sql"
	"errors"
)

type Comment interface {
	Create(data *models.Comment) error
}

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) Comment {
	return &CommentRepository{db}
}

func (c *CommentRepository) Create(data *models.Comment) error {
	query := `INSERT INTO comments(post_id, author_name, content, created_at) VALUES(?, ?, ?, ?)`
	result, err := c.db.Exec(query, data.PostID, data.AuthorName, data.Content, data.CreatedAt)
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
