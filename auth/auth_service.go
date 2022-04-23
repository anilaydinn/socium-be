package auth

import (
	"encoding/json"
	"github.com/anilaydinn/socium-be/model"
	"github.com/anilaydinn/socium-be/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"

	"strings"
)

type Service struct {
	Repository repository.Repository
}

func NewService(repository repository.Repository) Service {
	return Service{
		Repository: repository,
	}
}

func (s *Service) VerifyToken(bearerToken, role string) bool {

	splitToken := strings.Split(strings.Replace(bearerToken, "\"", "", -1), "Bearer ")
	if len(splitToken) != 2 {
		return false
	}

	token := splitToken[1]

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, &fiber.Error{Code: 400, Message: "Invalid Token"}
		}
		return []byte(""), nil
	})

	if err != nil {
		return false
	}

	jwtClaims, isVerified := parsedToken.Claims.(jwt.MapClaims)

	if !isVerified && !parsedToken.Valid {
		return false
	}

	jsonString, _ := json.Marshal(jwtClaims)
	parsedClaims := model.CustomClaims{}
	json.Unmarshal(jsonString, &parsedClaims)

	user, _ := s.Repository.GetUser(parsedClaims.Issuer)

	if user != nil || parsedClaims.Issuer == user.ID {
		if parsedClaims.UserType == user.UserType && parsedClaims.UserType == role {
			return true
		}
	}

	return false
}
