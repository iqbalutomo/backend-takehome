package controllers

import (
	"backend-takehome/dto"
	"backend-takehome/helpers"
	"backend-takehome/models"
	"backend-takehome/repository"
	"backend-takehome/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type PostController struct {
	repo repository.Post
}

func NewPostController(repo repository.Post) *PostController {
	return &PostController{repo}
}

func (p *PostController) CreatePost(c echo.Context) error {
	user, err := helpers.GetClaims(c)
	if err != nil {
		return err
	}

	if user.Email == "" {
		return echo.NewHTTPError(utils.ErrUnauthorized.EchoFormatDetails("Access not permitted"))
	}

	data := new(models.Post)
	if err := c.Bind(data); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	if err := c.Validate(data); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	dataPost := &models.Post{
		Title:     data.Title,
		Content:   data.Content,
		AuthorID:  user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := p.repo.Create(dataPost); err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}

	return c.JSON(http.StatusCreated, dto.Response{
		Message: "Post has been created",
		Data:    dataPost,
	})
}
