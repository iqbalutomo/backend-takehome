package models

import "time"

type Comment struct {
	ID         uint      `json:"id"`
	PostID     uint      `json:"post_id"`
	AuthorName string    `json:"author_name"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}
