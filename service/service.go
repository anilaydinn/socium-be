package service

import (
	"github.com/anilaydinn/socium-be/errors"
	"github.com/anilaydinn/socium-be/model"
	"github.com/anilaydinn/socium-be/repository"
	"github.com/anilaydinn/socium-be/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	"golang.org/x/crypto/bcrypt"
	"time"
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
		UserType: "user",
	}

	newUser, err := service.repository.RegisterUser(user)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (service *Service) LoginUser(userCredentialsDTO model.UserCredentialsDTO) (*model.Token, *fiber.Cookie, error) {
	user, err := service.repository.GetUserByEmail(userCredentialsDTO.Email)

	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		return nil, nil, errors.UserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userCredentialsDTO.Password)); err != nil {
		return nil, nil, errors.WrongPassword
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, model.CustomClaims{
		UserType: user.UserType,
		StandardClaims: jwt.StandardClaims{
			Issuer: user.ID,
		},
	})

	token, err := claims.SignedString([]byte("fe7999d6-47fa-11ec-81d3-0242ac130003"))

	if err != nil {
		return nil, nil, err
	}

	cookie := fiber.Cookie{
		Name:    "user-token",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 24),
	}
	return &model.Token{
		Token: token,
	}, &cookie, nil
}
