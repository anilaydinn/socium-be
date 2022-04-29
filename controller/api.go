package controller

import (
	"github.com/anilaydinn/socium-be/errors"
	"github.com/anilaydinn/socium-be/model"
	"github.com/anilaydinn/socium-be/service"
	"github.com/gofiber/fiber/v2"
)

type API struct {
	service *service.Service
}

func (api *API) SetupApp(app *fiber.App) {
	app.Post("/register", api.RegisterUserHandler)
	app.Post("/login", api.LoginUserHandler)
	app.Get("/activation/:userID", api.ActivationHandler)
	app.Post("/forgotPassword", api.ForgotPasswordHandler)
	app.Patch("/resetPassword/:userID", api.ResetPasswordHandler)
}

func NewAPI(service *service.Service) API {
	return API{
		service: service,
	}
}

func (api *API) RegisterUserHandler(c *fiber.Ctx) error {
	userDTO := model.UserDTO{}

	err := c.BodyParser(&userDTO)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	user, err := api.service.RegisterUser(userDTO)

	switch err {
	case nil:
		c.JSON(user)
		c.Status(fiber.StatusCreated)
	case errors.UserAlreadyRegistered:
		c.Status(fiber.StatusBadRequest)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) LoginUserHandler(c *fiber.Ctx) error {
	userCredentialsDTO := model.UserCredentialsDTO{}

	err := c.BodyParser(&userCredentialsDTO)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	token, cookie, err := api.service.LoginUser(userCredentialsDTO)

	switch err {
	case nil:
		c.JSON(token)
		c.Cookie(cookie)
		c.Status(fiber.StatusOK)
	case errors.UserNotFound:
		c.Status(fiber.StatusBadRequest)
	case errors.Unauthorized:
		c.Status(fiber.StatusUnauthorized)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) ActivationHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")

	user, err := api.service.Activation(userID)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(user)
	case errors.UserNotFound:
		c.Status(fiber.StatusNotFound)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) ForgotPasswordHandler(c *fiber.Ctx) error {
	forgotPasswordDTO := model.ForgotPasswordDTO{}
	err := c.BodyParser(&forgotPasswordDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	err = api.service.ForgotPassword(forgotPasswordDTO)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
	case errors.UserNotFound:
		c.Status(fiber.StatusNotFound)
	case errors.UserNotActivated:
		c.Status(fiber.StatusBadRequest)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) ResetPasswordHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	resetPasswordDTO := model.ResetPasswordDTO{}
	err := c.BodyParser(&resetPasswordDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	err = api.service.ResetPassword(userID, resetPasswordDTO)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}
