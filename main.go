package main

import (
	"github.com/anilaydinn/socium-be/controller"
	"github.com/anilaydinn/socium-be/repository"
	"github.com/anilaydinn/socium-be/service"
)

func main() {
	repository := repository.NewRepository("mongodb+srv://sociumtest:Se6iRf8elvL6Fcn6@cluster0.g1hlq.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	service := service.NewService(repository)
	api := controller.NewAPI(&service)

	app := controller.Handler(&api)

	app.Listen(":8080")
}
