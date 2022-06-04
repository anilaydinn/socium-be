package test

import (
	"bytes"
	"encoding/json"
	"github.com/anilaydinn/socium-be/controller"
	"github.com/anilaydinn/socium-be/middleware"
	"github.com/anilaydinn/socium-be/model"
	"github.com/anilaydinn/socium-be/service"
	"github.com/anilaydinn/socium-be/utils"
	"github.com/gofiber/fiber/v2"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestCreateContact(t *testing.T) {
	Convey("Given guest user", t, func() {
		app := fiber.New()
		testRepository := GetCleanTestRepository()
		middleware.SetupMiddleWare(app, *testRepository)
		service := service.NewService(testRepository)
		api := controller.NewAPI(&service)

		api.SetupApp(app)

		Convey("When user send contact information request", func() {
			contactDTO := model.ContactDTO{
				Name:    "Test",
				Surname: "Surname",
				Email:   "testsurname@yopmail.com",
				Message: "Hello",
			}
			reqBody, err := json.Marshal(contactDTO)
			So(err, ShouldBeNil)

			req, err := http.NewRequest(http.MethodPost, "/api/contacts", bytes.NewReader(reqBody))
			req.Header.Add("Content-Type", "application/json")

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 201", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusCreated)
			})

			Convey("Then newly contact information should return", func() {
				actualResult := model.Contact{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult.ID, ShouldNotBeNil)
				So(actualResult.Name, ShouldEqual, contactDTO.Name)
				So(actualResult.Surname, ShouldEqual, contactDTO.Surname)
				So(actualResult.Email, ShouldEqual, contactDTO.Email)
				So(actualResult.Message, ShouldEqual, contactDTO.Message)
			})
		})
	})
}

func TestGetAllContacts(t *testing.T) {
	Convey("Given admin user", t, func() {
		app := fiber.New()
		testRepository := GetCleanTestRepository()
		middleware.SetupMiddleWare(app, *testRepository)
		service := service.NewService(testRepository)
		api := controller.NewAPI(&service)

		api.SetupApp(app)

		testContact1 := model.Contact{
			ID:      utils.GenerateUUID(8),
			Name:    "James",
			Surname: "Bond",
			Email:   "jamesbond@gmail.com",
			Message: "This is very nice social media web site.",
		}
		testContact2 := model.Contact{
			ID:      utils.GenerateUUID(8),
			Name:    "Natalie",
			Surname: "Jason",
			Email:   "nataliejason@gmail.com",
			Message: "This is very nice social media web site. It is amazing.",
		}
		testRepository.CreateContact(testContact1)
		testRepository.CreateContact(testContact2)

		registeredUser := model.User{
			ID:                   "3c0bbdae",
			Name:                 "James",
			Surname:              "Bond",
			Email:                "test@gmail.com",
			Password:             "$2a$10$08qe8bXis2qObLNyEJfzpePCnqSJRyUXIa//ALLJw9l8q5gOTJljq",
			FriendRequestUserIDs: []string{},
			FriendIDs:            []string{"123123"},
			UserType:             "admin",
			IsActivated:          true,
		}
		testRepository.RegisterUser(registeredUser)

		Convey("When admin user send get all contacts request", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6ImFkbWluIiwiaXNzIjoiM2MwYmJkYWUifQ.aYf3WQryPbYoexgG18Q9iWYbnLtnH2ueE_rgTFdqBx4"

			req, err := http.NewRequest(http.MethodGet, "/admin/contacts", nil)
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then all contact informations should return", func() {
				actualResult := []model.Contact{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult[0].ID, ShouldNotBeNil)
				So(actualResult[0].ID, ShouldEqual, testContact1.ID)
				So(actualResult[0].Name, ShouldEqual, testContact1.Name)
				So(actualResult[0].Surname, ShouldEqual, testContact1.Surname)
				So(actualResult[0].Message, ShouldEqual, testContact1.Message)
				So(actualResult[0].Email, ShouldEqual, testContact1.Email)

				So(actualResult[1].ID, ShouldNotBeNil)
				So(actualResult[1].ID, ShouldEqual, testContact2.ID)
				So(actualResult[1].Name, ShouldEqual, testContact2.Name)
				So(actualResult[1].Surname, ShouldEqual, testContact2.Surname)
				So(actualResult[1].Message, ShouldEqual, testContact2.Message)
				So(actualResult[1].Email, ShouldEqual, testContact2.Email)
			})
		})
	})
}

func TestAdminDeleteContact(t *testing.T) {
	Convey("Given admin and contacts data", t, func() {
		app := fiber.New()
		testRepository := GetCleanTestRepository()
		middleware.SetupMiddleWare(app, *testRepository)
		service := service.NewService(testRepository)
		api := controller.NewAPI(&service)

		api.SetupApp(app)

		registeredUser1 := model.User{
			ID:                   "3c0bbdae",
			Name:                 "James",
			Surname:              "Bond",
			Email:                "test@gmail.com",
			Password:             "$2a$10$08qe8bXis2qObLNyEJfzpePCnqSJRyUXIa//ALLJw9l8q5gOTJljq",
			FriendRequestUserIDs: []string{},
			FriendIDs:            []string{"123123"},
			UserType:             "admin",
			IsActivated:          true,
		}
		testRepository.RegisterUser(registeredUser1)

		contact1 := model.Contact{
			ID:      utils.GenerateUUID(8),
			Name:    "Anıl",
			Surname: "Aydın",
			Email:   "test@gmail.com",
			Message: "Test message.",
		}
		testRepository.CreateContact(contact1)

		Convey("When admin user send delete user post request with id params", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6ImFkbWluIiwiaXNzIjoiM2MwYmJkYWUifQ.aYf3WQryPbYoexgG18Q9iWYbnLtnH2ueE_rgTFdqBx4"

			req, err := http.NewRequest(http.MethodDelete, "/admin/contacts/"+contact1.ID, nil)
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 204", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusNoContent)
			})

			Convey("Then contact should deleted", func() {
				contact, err := testRepository.GetContact(contact1.ID)
				So(contact, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
