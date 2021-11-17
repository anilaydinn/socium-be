package test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/anilaydinn/socium-be/utils"

	"github.com/anilaydinn/socium-be/controller"
	"github.com/anilaydinn/socium-be/model"
	"github.com/anilaydinn/socium-be/repository"
	"github.com/anilaydinn/socium-be/service"
	"github.com/gofiber/fiber"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRegisterUser(t *testing.T) {
	Convey("Given a valid user data", t, func() {
		testRepository := GetCleanTestRepository()
		service := service.NewService(testRepository)
		api := controller.NewAPI(&service)

		app := controller.Handler(&api)

		Convey("When add user request sent", func() {
			userDTO := model.UserDTO{
				Name:     "John",
				Surname:  "Obama",
				Email:    "john@gmail.com",
				Password: "123123",
			}

			reqBody, err := json.Marshal(userDTO)
			So(err, ShouldBeNil)

			req, _ := http.NewRequest("POST", "/users", bytes.NewReader(reqBody))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Set("Content-Length", strconv.Itoa(len(reqBody)))

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 201", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusCreated)
			})

			Convey("Then newly created user should return as a response", func() {
				actualResult := model.User{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult.ID, ShouldNotBeEmpty)
				So(actualResult.Name, ShouldEqual, userDTO.Name)
				So(actualResult.Surname, ShouldEqual, userDTO.Surname)
				So(actualResult.Email, ShouldEqual, userDTO.Email)
			})
		})
	})
}

func TestLoginUser(t *testing.T) {
	Convey("Given already register user", t, func() {
		testRepository := GetCleanTestRepository()
		service := service.NewService(testRepository)
		api := controller.NewAPI(&service)

		app := controller.Handler(&api)

		registeredUser := model.User{
			ID:       utils.GenerateUUID(8),
			Email:    "test@gmail.com",
			Name:     "Test Name",
			Surname:  "Test Surname",
			Password: "$2a$10$WCtghenC3N2Kg6ZjcoN/6O7fEJgTz5UzN65JoCGfxabqfEGJrxdBu",
		}
		testRepository.RegisterUser(registeredUser)

		Convey("When login user request sent", func() {
			userCredentialsDTO := model.UserCredentialsDTO{
				Email:    "test@gmail.com",
				Password: "123123",
			}

			reqBody, err := json.Marshal(userCredentialsDTO)
			So(err, ShouldBeNil)

			req, _ := http.NewRequest("POST", "/users/login", bytes.NewReader(reqBody))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Set("Content-Length", strconv.Itoa(len(reqBody)))

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then user token should returned", func() {
				actualResult := model.Token{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult.Token, ShouldNotBeNil)
				So(len(actualResult.Token), ShouldBeGreaterThan, 0)
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
