package utils

import "net/http"

type ErrResponse struct {
	Status int    `json:"status"`
	Type   string `json:"type"`
	Detail string `json:"detail"`
}

var (
	ErrBadRequest = ErrResponse{
		Status: http.StatusBadRequest,
		Type:   "Bad Request",
	}

	ErrInternalServer = ErrResponse{
		Status: http.StatusInternalServerError,
		Type:   "Internal Server",
	}

	ErrNotFound = ErrResponse{
		Status: http.StatusNotFound,
		Type:   "Not Found",
	}

	ErrUnauthorized = ErrResponse{
		Status: http.StatusUnauthorized,
		Type:   "Unauthorized",
	}
)

func (er *ErrResponse) EchoFormatDetails(detail string) (int, *ErrResponse) {
	er.Detail = detail

	return er.Status, er
}
