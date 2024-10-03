package dto

type PostCommentRequest struct {
	Content string `json:"content" validate:"required"`
}
