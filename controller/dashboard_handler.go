package controller

import "github.com/gofiber/fiber/v2"

func (h *Handler) GetAdminDashboard(c *fiber.Ctx) error {
	adminDashboard, err := h.service.GetAdminDashboard()

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(adminDashboard)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}
