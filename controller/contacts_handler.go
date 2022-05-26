package controller

import (
	"github.com/anilaydinn/socium-be/model"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateContactHandler(c *fiber.Ctx) error {
	contactDTO := model.ContactDTO{}
	err := c.BodyParser(&contactDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	contact, err := h.service.CreateContact(contactDTO)

	switch err {
	case nil:
		c.Status(fiber.StatusCreated)
		c.JSON(contact)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}
