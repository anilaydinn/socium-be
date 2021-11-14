package service

import (
	"github.com/anilaydinn/socium-be/model"
	"github.com/anilaydinn/socium-be/repository"
	"github.com/anilaydinn/socium-be/utils"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repository *repository.Repository) Service {
	return Service{
		repository: repository,
	}
}

func (service *Service) RegisterUser(userDTO model.UserDTO) (*model.User, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := model.User{
		ID:       utils.GenerateUUID(8),
		Name:     userDTO.Name,
		Surname:  userDTO.Surname,
		Email:    userDTO.Email,
		Password: string(hashedPassword),
	}

	newUser, err := service.repository.RegisterUser(user)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}
