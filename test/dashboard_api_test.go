package test

import (
	"encoding/json"
	"github.com/anilaydinn/socium-be/controller"
	"github.com/anilaydinn/socium-be/middleware"
	"github.com/anilaydinn/socium-be/model"
	"github.com/anilaydinn/socium-be/service"
	"github.com/anilaydinn/socium-be/utils"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAdminGetDashboard(t *testing.T) {
	Convey("Given admin", t, func() {
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

		registeredUser2 := model.User{
			ID:          "123123",
			Name:        "Mehmet",
			Surname:     "Bond",
			Email:       "test1@gmail.com",
			Password:    "$2a$10$08qe8bXis2qObLNyEJfzpePCnqSJRyUXIa//ALLJw9l8q5gOTJljq",
			UserType:    "user",
			IsActivated: true,
		}
		testRepository.RegisterUser(registeredUser1)
		testRepository.RegisterUser(registeredUser2)

		post1 := model.Post{
			ID:              utils.GenerateUUID(8),
			UserID:          registeredUser2.ID,
			User:            &registeredUser2,
			Description:     "Test Post Description",
			Image:           "asdşasdöls",
			IsPrivate:       false,
			WhoLikesUserIDs: nil,
			CommentIDs:      nil,
			Comments:        nil,
			CreatedAt:       time.Now().UTC().Round(time.Minute),
			UpdatedAt:       time.Now().UTC().Round(time.Minute),
		}
		testRepository.CreatePost(post1)

		Convey("When admin user send get user posts request with id params", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6ImFkbWluIiwiaXNzIjoiM2MwYmJkYWUifQ.aYf3WQryPbYoexgG18Q9iWYbnLtnH2ueE_rgTFdqBx4"

			req, err := http.NewRequest(http.MethodGet, "/admin/dashboard", nil)
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then dashboard informations should return", func() {
				actualResult := model.DashboardInformation{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult.CommentCount, ShouldEqual, 0)
				So(actualResult.UserCount, ShouldEqual, 2)
				So(actualResult.PostCount, ShouldEqual, 1)
				So(actualResult.ActivatedUserCount, ShouldEqual, 2)
			})
		})
	})
}
