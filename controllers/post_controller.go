package controllers

import (
	"backend-takehome/dto"
	"backend-takehome/helpers"
	"backend-takehome/models"
	"backend-takehome/repository"
	"backend-takehome/utils"
	"fmt"
	"net/http"
	"strconv"
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

func (p *PostController) GetPostDetail(c echo.Context) error {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	postDetail, err := p.repo.FindByID(uint(postID))
	if err != nil {
		if err.Error() == "post not found" {
			return echo.NewHTTPError(utils.ErrNotFound.EchoFormatDetails("Post not found"))
		}

		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}

	return c.JSON(http.StatusOK, dto.Response{
		Message: fmt.Sprintf("Get post detail with id %d successfully", postID),
		Data:    postDetail,
	})
}

func (p *PostController) GetPosts(c echo.Context) error {
	posts, err := p.repo.GetAll()
	if err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}

	return c.JSON(http.StatusOK, dto.Response{
		Message: "Get all posts successfully",
		Data:    posts,
	})
}
