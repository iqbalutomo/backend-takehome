package controllers

import (
	"backend-takehome/dto"
	"backend-takehome/models"
	"backend-takehome/repository"
	"backend-takehome/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	repo repository.User
}

func NewUserController(repo repository.User) *UserController {
	return &UserController{repo}
}

func (u *UserController) Register(c echo.Context) error {
	data := new(models.User)
	if err := c.Bind(data); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	if err := c.Validate(data); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}

	dataUser := &models.User{
		Name:         data.Name,
		Email:        data.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdateAt:     time.Now(),
	}

	if err := u.repo.Create(dataUser); err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}

	dataResp := dto.RegisterResponse{
		ID:        dataUser.ID,
		Name:      dataUser.Name,
		Email:     dataUser.Email,
		CreatedAt: dataUser.CreatedAt,
		UpdateAt:  dataUser.UpdateAt,
	}

	return c.JSON(http.StatusCreated, dto.Response{
		Message: "Registered successfully",
		Data:    dataResp,
	})
}
