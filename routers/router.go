package routers

import (
	"backend-takehome/controllers"
	_ "backend-takehome/docs"
	"backend-takehome/middlewares"
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Echo(e *echo.Echo, uc controllers.UserController, pc controllers.PostController, cc controllers.CommentController) {
	apiVersion := "/api/v1"
	api := e.Group(apiVersion)

	e.GET("", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/swagger/index.html")
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// user
	api.POST("/register", uc.Register)
	api.POST("/login", uc.Login)
	api.POST("/logout", uc.Logout, middlewares.ProtectedRoute)

	// post
	posts := api.Group("/posts")
	{
		posts.POST("", pc.CreatePost, middlewares.ProtectedRoute)
		posts.GET("/:id", pc.GetPostDetail)
		posts.GET("", pc.GetPosts)
		posts.PUT("/:id", pc.UpdatePost, middlewares.ProtectedRoute)
		posts.DELETE("/:id", pc.DeletePost, middlewares.ProtectedRoute)
	}

	// comment
	comments := posts.Group("/:id/comments")
	comments.Use(middlewares.ProtectedRoute)
	{
		comments.POST("", cc.CreateComment)
		comments.GET("", cc.GetComments)
	}
}
