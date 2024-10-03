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
