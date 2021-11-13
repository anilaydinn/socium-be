package main

import (
	"github.com/anilaydinn/socium-be/controller"
	"github.com/anilaydinn/socium-be/repository"
	"github.com/anilaydinn/socium-be/service"
	"github.com/gofiber/fiber/v2"
)

func main() {
	repository := repository.NewRepository("mongodb+srv://sociumtest:Se6iRf8elvL6Fcn6@cluster0.g1hlq.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	service := service.NewService(repository)
	api := controller.NewAPI(&service)

	app := SetupApp(&api)

	app.Listen(":8080")
}

func SetupApp(api *controller.API) *fiber.App {

	app := fiber.New()

	return app
}
