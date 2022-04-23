package auth

import (
	"context"
	"github.com/anilaydinn/socium-be/model"
	"github.com/anilaydinn/socium-be/repository"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAuthUser(t *testing.T) {
	Convey("Given that bearer token", t, func() {
		repository := GetCleanTestRepository()
		authService := NewService(*repository)
		bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6InVzZXIiLCJpc3MiOiJjMzJmNDMxMCJ9.mlqFqv-z43Skfh6qbFJCitYptQXd1yVGZyYEDwgA82U"

		user := model.User{
			ID:       "c32f4310",
			UserType: "user",
		}
		repository.RegisterUser(user)

		Convey("When bearer token sent verify token method", func() {
			isAuth := authService.VerifyToken(bearerToken, "user")

			Convey("Then isAuth should be true", func() {
				So(isAuth, ShouldBeTrue)
			})
		})
	})
}

func TestAuthAdmin(t *testing.T) {
	Convey("Given that bearer token", t, func() {
		repository := GetCleanTestRepository()
		authService := NewService(*repository)
		bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6ImFkbWluIiwiaXNzIjoiYzMyZjQzMTAifQ.WY_8-y3Wh8EoogG2FKtmxiFLBW3JE3d-j6pxaP21BQA"

		user := model.User{
			ID:       "c32f4310",
			UserType: "admin",
		}
		repository.RegisterUser(user)

		Convey("When bearer token sent verify token method", func() {
			isAuth := authService.VerifyToken(bearerToken, "admin")

			Convey("Then isAuth should be true", func() {
				So(isAuth, ShouldBeTrue)
			})
		})
	})
}

func TestAuthUserWrongToken(t *testing.T) {
	Convey("Given that wrong bearer token", t, func() {
		repository := GetCleanTestRepository()
		authService := NewService(*repository)
		bearerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6InVzZXIiLCJpc3MiOiJjMzJmNDMxMCJ9.mlqFqv-z43Skfh6qbFJCitYptQXd1yVGZyYEDwgA82U"

		user := model.User{
			ID:       "c32f4310",
			UserType: "user",
		}
		repository.RegisterUser(user)

		Convey("When wrong bearer token sent verify token method", func() {
			isAuth := authService.VerifyToken(bearerToken, "user")

			Convey("Then isAuth should be true", func() {
				So(isAuth, ShouldBeFalse)
			})
		})
	})
}

func GetCleanTestRepository() *repository.Repository {
	repository := repository.NewRepository("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	testDB := repository.MongoClient.Database("socium")
	testDB.Drop(ctx)

	return repository
}
