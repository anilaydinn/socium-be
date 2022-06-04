package service

import (
	"github.com/anilaydinn/socium-be/model"
	"github.com/anilaydinn/socium-be/utils"
)

func (service *Service) CreateContact(contactDTO model.ContactDTO) (*model.Contact, error) {
	contact := model.Contact{
		ID:      utils.GenerateUUID(8),
		Name:    contactDTO.Name,
		Surname: contactDTO.Surname,
		Email:   contactDTO.Email,
		Message: contactDTO.Message,
	}

	newContact, err := service.repository.CreateContact(contact)
	if err != nil {
		return nil, err
	}
	return newContact, nil
}

func (service *Service) GetAllContacts() ([]model.Contact, error) {
	return service.repository.GetAllContacts()
}
