package middleware

import (
	"github.com/anilaydinn/socium-be/auth"
	"github.com/anilaydinn/socium-be/repository"
	"github.com/gofiber/fiber/v2"
)

func SetupMiddleWare(app *fiber.App, userRepository repository.Repository) {
	authService := auth.NewService(userRepository)
	authHandler := auth.NewHandler(authService)
	app.Use("/user", authHandler.AuthUserHandler)
	app.Use("/admin", authHandler.AuthAdminHandler)
}
