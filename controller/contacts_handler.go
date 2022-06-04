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

func (h *Handler) AdminGetAllContactsHandler(c *fiber.Ctx) error {
	contacts, err := h.service.GetAllContacts()

	switch err {
	case nil:
		c.JSON(contacts)
		c.Status(fiber.StatusOK)
	default:
		c.Status(fiber.StatusInternalServerError)

	}
	return nil
}

func (h *Handler) AdminDeleteContactHandler(c *fiber.Ctx) error {
	contactID := c.Params("contactID")
	if len(contactID) == 0 {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	err := h.service.DeleteContact(contactID)

	switch err {
	case nil:
		c.Status(fiber.StatusNoContent)
	default:
		c.Status(fiber.StatusInternalServerError)

	}
	return nil
}
