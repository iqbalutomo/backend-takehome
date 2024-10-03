package routers

import (
	"backend-takehome/controllers"
	"backend-takehome/middlewares"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func Echo(e *echo.Echo, uc controllers.UserController, pc controllers.PostController) {
	apiVersion := "/api/v1"
	api := e.Group(apiVersion)

	e.GET("", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, apiVersion)
	})

	api.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"message": fmt.Sprintf("server running on port %s", os.Getenv("PORT")),
		})
	})

	// user
	api.POST("/register", uc.Register)
	api.POST("/login", uc.Login)

	// post
	posts := api.Group("/posts")
	{
		posts.POST("", pc.CreatePost, middlewares.ProtectedRoute)
		posts.GET("/:id", pc.GetPostDetail)
		posts.GET("", pc.GetPosts)
		posts.PUT("/:id", pc.UpdatePost, middlewares.ProtectedRoute)
	}
}
