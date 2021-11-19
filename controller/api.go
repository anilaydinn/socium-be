package controller

import (
	"github.com/anilaydinn/socium-be/errors"
	"github.com/anilaydinn/socium-be/model"
	"github.com/anilaydinn/socium-be/service"
	"github.com/gofiber/fiber"
)

type API struct {
	service *service.Service
}

func Handler(api *API) *fiber.App {

	app := fiber.New()

	app.Post("/users", api.RegisterUserHandler)
	app.Post("/users/login", api.LoginUserHandler)

	return app
}

func NewAPI(service *service.Service) API {
	return API{
		service: service,
	}
}

func (api *API) RegisterUserHandler(c *fiber.Ctx) {
	userDTO := model.UserDTO{}

	err := c.BodyParser(&userDTO)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
	}

	user, err := api.service.RegisterUser(userDTO)

	switch err {
	case nil:
		c.JSON(user)
		c.Status(fiber.StatusCreated)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
}

func (api *API) LoginUserHandler(c *fiber.Ctx) {
	userCredentialsDTO := model.UserCredentialsDTO{}

	err := c.BodyParser(&userCredentialsDTO)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
	}

	token, cookie, err := api.service.LoginUser(userCredentialsDTO)

	switch err {
	case nil:
		c.JSON(token)
		c.Cookie(cookie)
		c.Status(fiber.StatusOK)
	case errors.UserNotFound:
		c.Status(fiber.StatusBadRequest)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
}
