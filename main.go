package main

import (
	"fmt"
	"github.com/anilaydinn/socium-be/controller"
	"github.com/anilaydinn/socium-be/middleware"
	"github.com/anilaydinn/socium-be/repository"
	"github.com/anilaydinn/socium-be/service"
	"github.com/anilaydinn/socium-be/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	dbURL := utils.GetDBUrl()

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	repository := repository.NewRepository(dbURL)
	middleware.SetupMiddleWare(app, *repository)
	service := service.NewService(repository)
	api := controller.NewAPI(&service)

	api.SetupApp(app)

	port := utils.SetPort()

	if err := app.Listen(":" + port); err != nil {
		fmt.Println(err)
		return
	}
}
