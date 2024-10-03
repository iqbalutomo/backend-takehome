package repository

import (
	"backend-takehome/models"
	"database/sql"
	"errors"
)

type Comment interface {
	Create(data *models.Comment) error
	GetAllByPostID(postID uint, limit, offset int) ([]models.Comment, error)
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

func (c *CommentRepository) GetAllByPostID(postID uint, limit, offset int) ([]models.Comment, error) {
	var datas []models.Comment

	query := `SELECT c.id, c.post_id, c.author_name, c.content, c.created_at FROM comments c JOIN posts p ON c.post_id = p.id WHERE p.id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := c.db.Query(query, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var data models.Comment
		if err := rows.Scan(&data.ID, &data.PostID, &data.AuthorName, &data.Content, &data.CreatedAt); err != nil {
			return nil, err
		}

		datas = append(datas, data)
	}

	return datas, nil
}
