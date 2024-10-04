package controllers

import (
	"backend-takehome/dto"
	"backend-takehome/helpers"
	"backend-takehome/models"
	"backend-takehome/repository"
	"backend-takehome/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type CommentController struct {
	commRepo repository.Comment
	postRepo repository.Post
}

func NewCommentController(commRepo repository.Comment, postRepo repository.Post) *CommentController {
	return &CommentController{commRepo, postRepo}
}

// @Summary Create Comment
// @Description Create comment on post by post id
// @Tags Comment
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param request body dto.PostCommentRequest true "Create comment details"
// @Success 200 {object} dto.Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 401 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/posts/{id}/comments [post]
func (cc *CommentController) CreateComment(c echo.Context) error {
	user, err := helpers.GetClaims(c)
	if err != nil {
		return err
	}

	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	_, err = cc.postRepo.FindByID(uint(postID))
	if err != nil {
		if err.Error() == "post not found" {
			return echo.NewHTTPError(utils.ErrNotFound.EchoFormatDetails("Post not found"))
		}

		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}

	data := new(dto.PostCommentRequest)
	if err := c.Bind(data); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	if err := c.Validate(data); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	dataComment := &models.Comment{
		PostID:     uint(postID),
		AuthorName: user.Name,
		Content:    data.Content,
		CreatedAt:  time.Now(),
	}

	if err := cc.commRepo.Create(dataComment); err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}

	return c.JSON(http.StatusCreated, dto.Response{
		Message: "Comment has been created",
		Data:    dataComment,
	})
}

// @Summary Get All Comments
// @Description Get all comments on post by post id
// @Tags Comment
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of posts per page" default(10)
// @Success 200 {object} dto.Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 401 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/posts/{id}/comments [get]
func (cc *CommentController) GetComments(c echo.Context) error {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = utils.COMMENTS_PAGE_DEFAULT
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 {
		limit = utils.COMMENTS_LIMIT_DEFAULT
	}

	offset := (page - 1) * limit

	_, err = cc.postRepo.FindByID(uint(postID))
	if err != nil {
		if err.Error() == "post not found" {
			return echo.NewHTTPError(utils.ErrNotFound.EchoFormatDetails("Post not found"))
		}

		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}

	dataComments, err := cc.commRepo.GetAllByPostID(uint(postID), limit, offset)
	if err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}

	return c.JSON(http.StatusOK, dto.Response{
		Message: "get list all comments successfully",
		Data:    dataComments,
	})
}
