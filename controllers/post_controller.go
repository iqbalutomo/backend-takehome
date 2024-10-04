package controllers

import (
	"backend-takehome/dto"
	"backend-takehome/helpers"
	"backend-takehome/models"
	"backend-takehome/repository"
	"backend-takehome/services"
	"backend-takehome/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type PostController struct {
	repo    repository.Post
	caching services.CachingService
}

func NewPostController(repo repository.Post, caching services.CachingService) *PostController {
	return &PostController{repo, caching}
}

// @Summary Create Post
// @Description Create post on Blog Takehome
// @Tags Post
// @Accept json
// @Produce json
// @Param request body dto.CreatePostRequest true "Create post details"
// @Success 200 {object} dto.Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 401 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/posts [post]
func (p *PostController) CreatePost(c echo.Context) error {
	user, err := helpers.GetClaims(c)
	if err != nil {
		return err
	}

	if user.Email == "" {
		return echo.NewHTTPError(utils.ErrUnauthorized.EchoFormatDetails("Access not permitted"))
	}

	data := new(dto.CreatePostRequest)
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

	_, err = p.CachePostDetailed(dataPost.ID)
	if err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}

	return c.JSON(http.StatusCreated, dto.Response{
		Message: "Post has been created",
		Data:    dataPost,
	})
}

// @Summary Get Post Detail
// @Description Get post detail with Post ID (this data get from cache by Redis)
// @Tags Post
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} dto.Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/posts/{id} [get]
func (p *PostController) GetPostDetail(c echo.Context) error {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	// from cache
	postDetail, err := p.caching.GetPostDetailed(uint(postID))
	if err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}

	// fallback
	if postDetail == nil {
		postDetail, err = p.repo.FindByID(uint(postID))
		if err != nil {
			if err.Error() == "post not found" {
				return echo.NewHTTPError(utils.ErrNotFound.EchoFormatDetails("Post not found"))
			}

			return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
		}

		if err := p.caching.SetPostDetailed(uint(postID), postDetail); err != nil {
			return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
		}
	}

	return c.JSON(http.StatusOK, dto.Response{
		Message: fmt.Sprintf("Get post detail with id %d successfully", postID),
		Data:    postDetail,
	})
}

// @Summary Get All Posts
// @Description Get all posts with pagination and sorting options
// @Tags Post
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of posts per page" default(10)
// @Param sort query string false "Sort by 'newest' ord 'oldest'" default(newest)
// @Success 200 {object} dto.Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/posts [get]
func (p *PostController) GetPosts(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = utils.POSTS_PAGE_DEFAULT
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 {
		limit = utils.POSTS_LIMIT_DEFAULT
	}

	sort := c.QueryParam("sort") // newest || oldest (default: newest)

	offset := (page - 1) * limit

	posts, err := p.repo.GetAll(limit, offset, sort)
	if err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}

	return c.JSON(http.StatusOK, dto.Response{
		Message: "Get all posts successfully",
		Data:    posts,
	})
}

// @Summary Update Post
// @Description Update post with param post id
// @Tags Post
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param request body dto.UpdatePostRequest true "Update post details"
// @Success 200 {object} dto.Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 401 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/posts/{id} [put]
func (p *PostController) UpdatePost(c echo.Context) error {
	user, err := helpers.GetClaims(c)
	if err != nil {
		return err
	}

	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	dataTmp, err := p.repo.FindByID(uint(postID))
	if err != nil {
		if err.Error() == "post not found" {
			return echo.NewHTTPError(utils.ErrNotFound.EchoFormatDetails("Post not found"))
		}

		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}

	if user.ID != dataTmp.Author.ID {
		return echo.NewHTTPError(utils.ErrUnauthorized.EchoFormatDetails("Access not permitted"))
	}

	dataReq := new(dto.UpdatePostRequest)

	if err := c.Bind(dataReq); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	if err := c.Validate(dataReq); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	dataPostUpdate := &models.Post{
		ID:       dataTmp.Post.ID,
		Title:    dataReq.Title,
		Content:  dataReq.Content,
		AuthorID: dataTmp.Author.ID,
	}

	if err := p.repo.Update(dataPostUpdate); err != nil {
		if err.Error() == "post not found" {
			return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
		}

		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}

	_, err = p.CachePostDetailed(uint(postID))
	if err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}

	return c.JSON(http.StatusOK, dto.Response{
		Message: fmt.Sprintf("Update post with id %d successfully", postID),
		Data:    dataPostUpdate,
	})
}

// @Summary Delete Post
// @Description Delete post with param post id
// @Tags Post
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} dto.Response
// @Failure 400 {object} utils.ErrResponse
// @Failure 401 {object} utils.ErrResponse
// @Failure 404 {object} utils.ErrResponse
// @Failure 500 {object} utils.ErrResponse
// @Router /api/v1/posts/{id} [delete]
func (p *PostController) DeletePost(c echo.Context) error {
	user, err := helpers.GetClaims(c)
	if err != nil {
		return err
	}

	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	dataTmp, err := p.repo.FindByID(uint(postID))
	if err != nil {
		if err.Error() == "post not found" {
			return echo.NewHTTPError(utils.ErrNotFound.EchoFormatDetails("Post not found"))
		}

		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}

	if user.ID != dataTmp.Author.ID {
		return echo.NewHTTPError(utils.ErrUnauthorized.EchoFormatDetails("Access not permitted"))
	}

	if err := p.repo.Delete(uint(postID)); err != nil {
		if err.Error() == "post not found" {
			return echo.NewHTTPError(utils.ErrNotFound.EchoFormatDetails("Post not found"))
		}

		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}

	return c.JSON(http.StatusOK, dto.Response{
		Message: "Post has been deleted",
		Data:    fmt.Sprintf("Post ID %d", postID),
	})
}

func (p *PostController) CachePostDetailed(postID uint) (*models.PostDetail, error) {
	data, err := p.repo.FindByID(postID)
	if err != nil {
		return nil, err
	}

	if err := p.caching.SetPostDetailed(postID, data); err != nil {
		return nil, err
	}

	return data, nil
}
