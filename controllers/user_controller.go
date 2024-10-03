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

func (u *UserController) Login(c echo.Context) error {
	data := new(dto.LoginRequest)
	if err := c.Bind(data); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	if err := c.Validate(data); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	userDataTmp, err := u.repo.FindByEmail(data.Email)
	if err != nil {
		if err.Error() == "user not found" {
			return echo.NewHTTPError(utils.ErrUnauthorized.EchoFormatDetails("Invalid email/password"))
		}

		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}

	userData := models.User{
		ID:           userDataTmp.ID,
		Name:         userDataTmp.Name,
		Email:        userDataTmp.Email,
		PasswordHash: userDataTmp.PasswordHash,
		CreatedAt:    userDataTmp.CreatedAt,
		UpdateAt:     userDataTmp.UpdateAt,
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userData.PasswordHash), []byte(data.Password)); err != nil {
		return echo.NewHTTPError(utils.ErrUnauthorized.EchoFormatDetails("Invalid email/password"))
	}

	if err := helpers.SignNewJWT(c, userData); err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}

	return c.JSON(http.StatusOK, dto.Response{
		Message: "Login successfully",
		Data:    "Authorization is stored in cookie",
	})
}
