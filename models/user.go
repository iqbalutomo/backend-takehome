package models

import "time"

type User struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name" validate:"required"`
	Email        string    `json:"email" validate:"required,email"`
	PasswordHash string    `json:"password_hash" validate:"required"`
	CreatedAt    time.Time `json:"created_at"`
	UpdateAt     time.Time `json:"updated_at"`
}
