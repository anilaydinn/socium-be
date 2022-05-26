package test

import (
	"bytes"
	"encoding/json"
	"github.com/anilaydinn/socium-be/controller"
	"github.com/anilaydinn/socium-be/middleware"
	"github.com/anilaydinn/socium-be/model"
	"github.com/anilaydinn/socium-be/service"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
