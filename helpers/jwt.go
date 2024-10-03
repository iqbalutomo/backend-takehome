package helpers

import (
	"backend-takehome/models"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func SignNewJWT(c echo.Context, user models.User) error {
	claims := jwt.MapClaims{
		"exp":   time.Now().Add(2 * time.Hour).Unix(),
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.Value = tokenString
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Expires = time.Now().Add(2 * time.Hour)

	c.SetCookie(cookie)

	return nil
}
