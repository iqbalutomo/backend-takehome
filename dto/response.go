package dto

import "time"

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type RegisterResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
}
