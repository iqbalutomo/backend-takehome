package main

import (
	"backend-takehome/config"
	"backend-takehome/controllers"
	"backend-takehome/repository"
	"backend-takehome/routers"
	"backend-takehome/services"
	"backend-takehome/utils"
	"os"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// @title Blog Takehome
// @version 1.0.0
// @description Personal Blog powered by Bythen AI :) LFG🚀

// @contact.name Muhlis Iqbal Utomo
// @contact.email muhlisiqbalutomo@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host blog-takehome-243802168923.asia-southeast2.run.app
// @BasePath /
func main() {
	db := config.ConnectDB()
	defer db.Close()

	e := echo.New()
	e.Use(middleware.Logger(), middleware.Recover())
	e.Validator = &utils.CustomValidator{NewValidator: validator.New()}

	redisClient := config.InitRedistClient()
	cachingService := services.NewCachingService(redisClient)

	userRepo := repository.NewUserRepository(db)
	userController := controllers.NewUserController(userRepo)

	postRepo := repository.NewPostRepository(db)
	postController := controllers.NewPostController(postRepo, cachingService)

	commentRepo := repository.NewCommentRepository(db)
	commentController := controllers.NewCommentController(commentRepo, postRepo)

	routers.Echo(e, *userController, *postController, *commentController)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
