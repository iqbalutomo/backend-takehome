package repository

import (
	"backend-takehome/models"
	"database/sql"
	"errors"
)

type Post interface {
	Create(data *models.Post) error
	FindByID(postID uint) (*models.PostDetail, error)
	GetAll() ([]models.PostDetail, error)
	Update(data *models.Post) error
	Delete(postID uint) error
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

	query := `SELECT p.id, p.title, p.content, p.author_id, p.created_at, p.updated_at, u.id, u.name, u.email FROM posts p JOIN users u ON p.author_id = u.id WHERE p.id = ?`
	if err := p.db.QueryRow(query, postID).Scan(&postID, &postDetail.Post.Title, &postDetail.Post.Content, &postDetail.Post.AuthorID, &postDetail.Post.CreatedAt, &postDetail.Post.UpdatedAt, &postDetail.Author.ID, &postDetail.Author.Name, &postDetail.Author.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("post not found")
		}

		return nil, err
	}

	postDetail.Post.ID = postID

	return &postDetail, nil
}

func (p *PostRepository) GetAll() ([]models.PostDetail, error) {
	var datas []models.PostDetail

	query := `SELECT p.id, p.title, p.content, p.author_id, p.created_at, p.updated_at, u.id, u.name, u.email FROM posts p JOIN users u ON p.author_id = u.id`
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var data models.PostDetail
		if err := rows.Scan(&data.Post.ID, &data.Post.Title, &data.Post.Content, &data.Post.AuthorID, &data.Post.CreatedAt, &data.Post.UpdatedAt, &data.Author.ID, &data.Author.Name, &data.Author.Email); err != nil {
			return nil, err
		}

		datas = append(datas, data)
	}

	return datas, nil
}

func (p *PostRepository) Update(data *models.Post) error {
	query := `UPDATE posts SET title = ?, content = ? WHERE id = ?`
	result, err := p.db.Exec(query, data.Title, data.Content, data.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("post not found")
	}

	return nil
}

func (p *PostRepository) Delete(postID uint) error {
	query := `DELETE FROM posts WHERE id = ?`
	result, err := p.db.Exec(query, postID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("post not found")
	}

	return nil
}
