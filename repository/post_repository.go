package repository

import (
	"backend-takehome/models"
	"database/sql"
	"errors"
)

type Post interface {
	Create(data *models.Post) error
	FindByID(postID uint) (*models.PostDetail, error)
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

func (p *PostRepository) FindByID(postID uint) (*models.PostDetail, error) {
	var postDetail models.PostDetail

	query := `SELECT posts.id, posts.title, posts.content, posts.author_id, posts.created_at, posts.updated_at, users.id, users.name, users.email FROM posts JOIN users ON posts.author_id = users.id WHERE posts.id = ?`
	if err := p.db.QueryRow(query, postID).Scan(&postID, &postDetail.Post.Title, &postDetail.Post.Content, &postDetail.Post.AuthorID, &postDetail.Post.CreatedAt, &postDetail.Post.UpdatedAt, &postDetail.Author.ID, &postDetail.Author.Name, &postDetail.Author.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("post not found")
		}

		return nil, err
	}

	postDetail.Post.ID = postID

	return &postDetail, nil
}
