package main

import (
	"github.com/anilaydinn/socium-be/controller"
	"github.com/anilaydinn/socium-be/middleware"
	"github.com/anilaydinn/socium-be/repository"
	"github.com/anilaydinn/socium-be/service"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	repository := repository.NewRepository("mongodb+srv://sociumtest:Se6iRf8elvL6Fcn6@cluster0.g1hlq.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	service := service.NewService(repository)
	api := controller.NewAPI(&service)

	app := controller.Handler(&api)
	middleware.SetupMiddleWare(app, *repository)
	app.Use(cors.New())
	app.Use(logger.New())

	app.Listen(":8080")
}
