package test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/anilaydinn/socium-be/middleware"
	"github.com/gofiber/fiber/v2"
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

	. "github.com/smartystreets/goconvey/convey"
)

func TestRegisterUser(t *testing.T) {
	Convey("Given a valid user data", t, func() {
		app := fiber.New()
		testRepository := GetCleanTestRepository()
		middleware.SetupMiddleWare(app, *testRepository)
		service := service.NewService(testRepository)
		api := controller.NewAPI(&service)
		api.SetupApp(app)

		Convey("When add user request sent", func() {
			userDTO := model.UserDTO{
				Name:     "John",
				Surname:  "Obama",
				Email:    "john@gmail.com",
				Password: "123123",
			}

			reqBody, err := json.Marshal(userDTO)
			So(err, ShouldBeNil)

			req, _ := http.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(reqBody))
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
				So(actualResult.UserType, ShouldEqual, "user")
				So(actualResult.IsActivated, ShouldBeFalse)
			})
		})
	})
}

func TestAlreadyRegisteredUser(t *testing.T) {
	Convey("Given a registered user and valid user data", t, func() {
		app := fiber.New()
		testRepository := GetCleanTestRepository()
		middleware.SetupMiddleWare(app, *testRepository)
		service := service.NewService(testRepository)
		api := controller.NewAPI(&service)
		api.SetupApp(app)

		user := model.User{
			ID:       utils.GenerateUUID(8),
			Name:     "John",
			Surname:  "Obama",
			Email:    "john@gmail.com",
			Password: "123123",
			UserType: "user",
		}
		testRepository.RegisterUser(user)

		Convey("When add user request sent", func() {
			userDTO := model.UserDTO{
				Name:     "John",
				Surname:  "Obama",
				Email:    "john@gmail.com",
				Password: "123123",
			}

			reqBody, err := json.Marshal(userDTO)
			So(err, ShouldBeNil)

			req, _ := http.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(reqBody))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Set("Content-Length", strconv.Itoa(len(reqBody)))

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 400", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusBadRequest)
			})
		})
	})
}

func TestLoginUser(t *testing.T) {
	Convey("Given already register user", t, func() {
		app := fiber.New()
		testRepository := GetCleanTestRepository()
		middleware.SetupMiddleWare(app, *testRepository)
		service := service.NewService(testRepository)
		api := controller.NewAPI(&service)

		api.SetupApp(app)

		registeredUser := model.User{
			ID:          utils.GenerateUUID(8),
			Email:       "test@gmail.com",
			Name:        "Test Name",
			Surname:     "Test Surname",
			Password:    "$2a$10$WCtghenC3N2Kg6ZjcoN/6O7fEJgTz5UzN65JoCGfxabqfEGJrxdBu",
			UserType:    "user",
			IsActivated: true,
		}
		testRepository.RegisterUser(registeredUser)

		Convey("When login user request sent", func() {
			userCredentialsDTO := model.UserCredentialsDTO{
				Email:    "test@gmail.com",
				Password: "123123",
			}

			reqBody, err := json.Marshal(userCredentialsDTO)
			So(err, ShouldBeNil)

			req, _ := http.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(reqBody))
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

		Convey("When invalid user request sent", func() {
			userCredentialsDTO := model.UserCredentialsDTO{
				Email:    "wrongmail@gmail.com",
				Password: "wrongpassword",
			}

			reqBody, err := json.Marshal(userCredentialsDTO)
			So(err, ShouldBeNil)

			req, _ := http.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(reqBody))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Set("Content-Length", strconv.Itoa(len(reqBody)))

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 400", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusBadRequest)
			})
		})
	})
}

func TestNotActivatedUserLogin(t *testing.T) {
	Convey("Given already register user", t, func() {
		app := fiber.New()
		testRepository := GetCleanTestRepository()
		middleware.SetupMiddleWare(app, *testRepository)
		service := service.NewService(testRepository)
		api := controller.NewAPI(&service)

		api.SetupApp(app)

		registeredUser := model.User{
			ID:          utils.GenerateUUID(8),
			Email:       "test@gmail.com",
			Name:        "Test Name",
			Surname:     "Test Surname",
			Password:    "$2a$10$WCtghenC3N2Kg6ZjcoN/6O7fEJgTz5UzN65JoCGfxabqfEGJrxdBu",
			UserType:    "user",
			IsActivated: false,
		}
		testRepository.RegisterUser(registeredUser)

		Convey("When login user request sent", func() {
			userCredentialsDTO := model.UserCredentialsDTO{
				Email:    "test@gmail.com",
				Password: "123123",
			}

			reqBody, err := json.Marshal(userCredentialsDTO)
			So(err, ShouldBeNil)

			req, _ := http.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(reqBody))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Set("Content-Length", strconv.Itoa(len(reqBody)))

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 401", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusUnauthorized)
			})
		})
	})
}

func TestUserActivation(t *testing.T) {
	Convey("Given already register user", t, func() {
		app := fiber.New()
		testRepository := GetCleanTestRepository()
		middleware.SetupMiddleWare(app, *testRepository)
		service := service.NewService(testRepository)
		api := controller.NewAPI(&service)

		api.SetupApp(app)

		registeredUser := model.User{
			ID:          utils.GenerateUUID(8),
			Email:       "test@gmail.com",
			Name:        "Test Name",
			Surname:     "Test Surname",
			Password:    "$2a$10$WCtghenC3N2Kg6ZjcoN/6O7fEJgTz5UzN65JoCGfxabqfEGJrxdBu",
			UserType:    "user",
			IsActivated: false,
		}
		testRepository.RegisterUser(registeredUser)

		Convey("When activate user request sent", func() {
			req, _ := http.NewRequest(http.MethodGet, "/api/activation/"+registeredUser.ID, nil)
			req.Header.Add("Content-Type", "application/json")

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then user should be activated", func() {
				actualResult := model.User{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult.IsActivated, ShouldBeTrue)
			})
		})
	})
}

func TestForgotPassword(t *testing.T) {
	Convey("Given a registered user", t, func() {
		app := fiber.New()
		testRepository := GetCleanTestRepository()
		middleware.SetupMiddleWare(app, *testRepository)
		service := service.NewService(testRepository)
		api := controller.NewAPI(&service)

		api.SetupApp(app)

		registeredUser := model.User{
			ID:          utils.GenerateUUID(8),
			Email:       "test@gmail.com",
			Name:        "Test Name",
			Surname:     "Test Surname",
			Password:    "$2a$10$WCtghenC3N2Kg6ZjcoN/6O7fEJgTz5UzN65JoCGfxabqfEGJrxdBu",
			UserType:    "user",
			IsActivated: true,
		}
		testRepository.RegisterUser(registeredUser)

		Convey("When forgot password request sent", func() {
			forgotPasswordDTO := model.ForgotPasswordDTO{
				Email: "test@gmail.com",
			}

			reqBody, err := json.Marshal(forgotPasswordDTO)
			So(err, ShouldBeNil)

			req, _ := http.NewRequest(http.MethodPost, "/api/forgotPassword", bytes.NewReader(reqBody))
			req.Header.Add("Content-Type", "application/json")

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})
		})

		Convey("When forgot password request sent with not registered mail", func() {
			forgotPasswordDTO := model.ForgotPasswordDTO{
				Email: "test222@gmail.com",
			}

			reqBody, err := json.Marshal(forgotPasswordDTO)
			So(err, ShouldBeNil)

			req, _ := http.NewRequest(http.MethodPost, "/api/forgotPassword", bytes.NewReader(reqBody))
			req.Header.Add("Content-Type", "application/json")

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 404", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusNotFound)
			})
		})
	})
}

func TestNotActivatedUserForgotPassword(t *testing.T) {
	Convey("Given that activated user", t, func() {
		app := fiber.New()
		testRepository := GetCleanTestRepository()
		middleware.SetupMiddleWare(app, *testRepository)
		service := service.NewService(testRepository)
		api := controller.NewAPI(&service)

		api.SetupApp(app)

		registeredUser := model.User{
			ID:          utils.GenerateUUID(8),
			Email:       "test@gmail.com",
			Name:        "Test Name",
			Surname:     "Test Surname",
			Password:    "$2a$10$WCtghenC3N2Kg6ZjcoN/6O7fEJgTz5UzN65JoCGfxabqfEGJrxdBu",
			UserType:    "user",
			IsActivated: false,
		}
		testRepository.RegisterUser(registeredUser)

		Convey("When forgot password request sent", func() {
			forgotPasswordDTO := model.ForgotPasswordDTO{
				Email: "test@gmail.com",
			}

			reqBody, err := json.Marshal(forgotPasswordDTO)
			So(err, ShouldBeNil)

			req, _ := http.NewRequest(http.MethodPost, "/api/forgotPassword", bytes.NewReader(reqBody))
			req.Header.Add("Content-Type", "application/json")

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 400", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusBadRequest)
			})
		})
	})
}

func TestResetPassword(t *testing.T) {
	Convey("Given that user", t, func() {
		app := fiber.New()
		testRepository := GetCleanTestRepository()
		middleware.SetupMiddleWare(app, *testRepository)
		service := service.NewService(testRepository)
		api := controller.NewAPI(&service)

		api.SetupApp(app)

		registeredUser := model.User{
			ID:          utils.GenerateUUID(8),
			Email:       "test@gmail.com",
			Name:        "Test Name",
			Surname:     "Test Surname",
			Password:    "$2a$10$WCtghenC3N2Kg6ZjcoN/6O7fEJgTz5UzN65JoCGfxabqfEGJrxdBu",
			UserType:    "user",
			IsActivated: false,
		}
		testRepository.RegisterUser(registeredUser)

		Convey("When new user password data sent with user id", func() {

			resetPasswordDTO := model.ResetPasswordDTO{
				Password: "332211",
			}
			reqBody, err := json.Marshal(resetPasswordDTO)
			So(err, ShouldBeNil)

			req, err := http.NewRequest(http.MethodPatch, "/api/resetPassword/"+registeredUser.ID, bytes.NewReader(reqBody))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Set("Content-Length", strconv.Itoa(len(reqBody)))

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})
		})
	})
}

func TestGetUser(t *testing.T) {
	Convey("Given that users", t, func() {
		app := fiber.New()
		testRepository := GetCleanTestRepository()
		middleware.SetupMiddleWare(app, *testRepository)
		service := service.NewService(testRepository)
		api := controller.NewAPI(&service)

		api.SetupApp(app)

		registeredUser1 := model.User{
			ID:          utils.GenerateUUID(8),
			Email:       "test@gmail.com",
			Name:        "Test Name",
			Surname:     "Test Surname",
			Password:    "$2a$10$WCtghenC3N2Kg6ZjcoN/6O7fEJgTz5UzN65JoCGfxabqfEGJrxdBu",
			UserType:    "user",
			IsActivated: true,
		}
		registeredUser2 := model.User{
			ID:          utils.GenerateUUID(8),
			Email:       "test@gmail.com",
			Name:        "Test Name",
			Surname:     "Test Surname",
			Password:    "$2a$10$WCtghenC3N2Kg6ZjcoN/6O7fEJgTz5UzN65JoCGfxabqfEGJrxdBu",
			UserType:    "user",
			IsActivated: true,
		}
		testRepository.RegisterUser(registeredUser1)
		testRepository.RegisterUser(registeredUser2)

		Convey("When get user request sent with user id", func() {
			req, err := http.NewRequest(http.MethodGet, "/api/users/"+registeredUser2.ID, nil)
			req.Header.Add("Content-Type", "application/json")

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then user should retrieved", func() {
				actualResult := model.User{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult.ID, ShouldNotBeNil)
				So(actualResult.ID, ShouldEqual, registeredUser2.ID)
				So(actualResult.Name, ShouldEqual, registeredUser2.Name)
				So(actualResult.Surname, ShouldEqual, registeredUser2.Surname)
				So(actualResult.Password, ShouldEqual, registeredUser2.Password)
				So(actualResult.Email, ShouldEqual, registeredUser2.Email)
				So(actualResult.UserType, ShouldEqual, registeredUser2.UserType)
			})
		})
	})
}

func TestCreatePost(t *testing.T) {
	Convey("Given a authenticated user", t, func() {
		app := fiber.New()
		testRepository := GetCleanTestRepository()
		middleware.SetupMiddleWare(app, *testRepository)
		service := service.NewService(testRepository)
		api := controller.NewAPI(&service)

		api.SetupApp(app)

		registeredUser := model.User{
			ID:          "3c0bbdae",
			Name:        "James",
			Surname:     "Bond",
			Email:       "test@gmail.com",
			Password:    "$2a$10$08qe8bXis2qObLNyEJfzpePCnqSJRyUXIa//ALLJw9l8q5gOTJljq",
			UserType:    "user",
			IsActivated: true,
		}
		testRepository.RegisterUser(registeredUser)

		Convey("When user send create post request", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6InVzZXIiLCJpc3MiOiIzYzBiYmRhZSJ9.F_7cDDzm0THldtJLLNunfdXtoKqLKeMK8BdHG9Dxi-s"

			postDTO := model.PostDTO{
				UserID:      registeredUser.ID,
				Description: "Post description",
				Image:       "cxzcxcxzczxöc",
				IsPrivate:   true,
			}
			reqBody, err := json.Marshal(postDTO)
			So(err, ShouldBeNil)

			req, err := http.NewRequest(http.MethodPost, "/user/posts", bytes.NewReader(reqBody))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)
			req.Header.Set("Content-Length", strconv.Itoa(len(reqBody)))

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 201", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusCreated)
			})

			Convey("Then created post should return", func() {
				actualResult := model.Post{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult.ID, ShouldNotBeNil)
				So(actualResult.UserID, ShouldEqual, postDTO.UserID)
				So(actualResult.Description, ShouldEqual, postDTO.Description)
				So(actualResult.Image, ShouldEqual, postDTO.Image)
				So(actualResult.IsPrivate, ShouldEqual, postDTO.IsPrivate)
				So(actualResult.WhoLikesUserIDs, ShouldBeNil)
				So(actualResult.User, ShouldBeNil)
				So(actualResult.CreatedAt, ShouldEqual, time.Now().UTC().Round(time.Second))
				So(actualResult.UpdatedAt, ShouldEqual, time.Now().UTC().Round(time.Second))
			})
		})
	})
}

func TestGetAllPosts(t *testing.T) {
	Convey("Given posts data", t, func() {
		app := fiber.New()
		testRepository := GetCleanTestRepository()
		middleware.SetupMiddleWare(app, *testRepository)
		service := service.NewService(testRepository)
		api := controller.NewAPI(&service)

		api.SetupApp(app)

		registeredUser1 := model.User{
			ID:          "3c0bbdae",
			Name:        "James",
			Surname:     "Bond",
			Email:       "test@gmail.com",
			Password:    "$2a$10$08qe8bXis2qObLNyEJfzpePCnqSJRyUXIa//ALLJw9l8q5gOTJljq",
			UserType:    "user",
			IsActivated: true,
		}
		registeredUser2 := model.User{
			ID:          utils.GenerateUUID(8),
			Name:        "James",
			Surname:     "Bond",
			Email:       "test@gmail.com",
			Password:    "$2a$10$08qe8bXis2qObLNyEJfzpePCnqSJRyUXIa//ALLJw9l8q5gOTJljq",
			UserType:    "user",
			IsActivated: true,
		}
		testRepository.RegisterUser(registeredUser1)
		testRepository.RegisterUser(registeredUser2)

		testPost1 := model.Post{
			ID:              utils.GenerateUUID(8),
			UserID:          registeredUser2.ID,
			User:            &registeredUser2,
			Description:     "Test Description 1",
			Image:           "zcxçömzcxözcxzzçcmzö 1",
			IsPrivate:       false,
			WhoLikesUserIDs: nil,
			CreatedAt:       time.Now().UTC().Add(-5 * time.Minute).Round(time.Second),
			UpdatedAt:       time.Now().UTC().Add(-5 * time.Minute).Round(time.Second),
		}

		testPost2 := model.Post{
			ID:              utils.GenerateUUID(8),
			UserID:          registeredUser2.ID,
			User:            &registeredUser2,
			Description:     "Test Description 2",
			Image:           "zcxçömzcxözcxzzçcmzö 2",
			IsPrivate:       false,
			WhoLikesUserIDs: nil,
			CreatedAt:       time.Now().UTC().Round(time.Second),
			UpdatedAt:       time.Now().UTC().Round(time.Second),
		}
		testPost3 := model.Post{
			ID:              utils.GenerateUUID(8),
			UserID:          registeredUser2.ID,
			User:            &registeredUser2,
			Description:     "Test Description 3",
			Image:           "zcxçömzcxözcxzzçcmzö 3",
			IsPrivate:       true,
			WhoLikesUserIDs: nil,
			CreatedAt:       time.Now().UTC().Add(-3 * time.Minute).Round(time.Second),
			UpdatedAt:       time.Now().UTC().Add(-3 * time.Minute).Round(time.Second),
		}
		testRepository.CreatePost(testPost1)
		testRepository.CreatePost(testPost2)
		testRepository.CreatePost(testPost3)

		Convey("When user send get posts request", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6InVzZXIiLCJpc3MiOiIzYzBiYmRhZSJ9.F_7cDDzm0THldtJLLNunfdXtoKqLKeMK8BdHG9Dxi-s"

			req, err := http.NewRequest(http.MethodGet, "/user/posts", nil)
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then all public posts should return", func() {
				actualResult := []model.Post{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult, ShouldHaveLength, 3)
				So(actualResult[0].ID, ShouldNotBeNil)
				So(actualResult[0].ID, ShouldEqual, testPost2.ID)
				So(actualResult[0].UserID, ShouldEqual, registeredUser2.ID)
				So(actualResult[0].Description, ShouldEqual, testPost2.Description)
				So(actualResult[0].Image, ShouldEqual, testPost2.Image)
				So(actualResult[0].IsPrivate, ShouldEqual, testPost2.IsPrivate)
				So(actualResult[0].WhoLikesUserIDs, ShouldEqual, testPost2.WhoLikesUserIDs)
				So(actualResult[0].User, ShouldResemble, &registeredUser2)
				So(actualResult[0].IsPrivate, ShouldBeFalse)
				So(actualResult[0].CreatedAt, ShouldEqual, testPost2.CreatedAt)
				So(actualResult[0].UpdatedAt, ShouldEqual, testPost2.UpdatedAt)

				So(actualResult[1].ID, ShouldNotBeNil)
				So(actualResult[1].ID, ShouldEqual, testPost3.ID)
				So(actualResult[1].UserID, ShouldEqual, registeredUser2.ID)
				So(actualResult[1].Description, ShouldEqual, testPost3.Description)
				So(actualResult[1].Image, ShouldEqual, testPost3.Image)
				So(actualResult[1].IsPrivate, ShouldEqual, testPost3.IsPrivate)
				So(actualResult[1].WhoLikesUserIDs, ShouldEqual, testPost3.WhoLikesUserIDs)
				So(actualResult[1].User, ShouldResemble, &registeredUser2)
				So(actualResult[1].IsPrivate, ShouldBeTrue)
				So(actualResult[1].CreatedAt, ShouldEqual, testPost3.CreatedAt)
				So(actualResult[1].UpdatedAt, ShouldEqual, testPost3.UpdatedAt)

				So(actualResult[2].ID, ShouldNotBeNil)
				So(actualResult[2].ID, ShouldEqual, testPost1.ID)
				So(actualResult[2].UserID, ShouldEqual, registeredUser2.ID)
				So(actualResult[2].Description, ShouldEqual, testPost1.Description)
				So(actualResult[2].Image, ShouldEqual, testPost1.Image)
				So(actualResult[2].IsPrivate, ShouldEqual, testPost1.IsPrivate)
				So(actualResult[2].WhoLikesUserIDs, ShouldEqual, testPost1.WhoLikesUserIDs)
				So(actualResult[2].User, ShouldResemble, &registeredUser2)
				So(actualResult[2].IsPrivate, ShouldBeFalse)
				So(actualResult[2].CreatedAt, ShouldEqual, testPost1.CreatedAt)
				So(actualResult[2].UpdatedAt, ShouldEqual, testPost1.UpdatedAt)
			})
		})
	})
}

func TestGetUserPosts(t *testing.T) {
	Convey("Given posts data", t, func() {
		app := fiber.New()
		testRepository := GetCleanTestRepository()
		middleware.SetupMiddleWare(app, *testRepository)
		service := service.NewService(testRepository)
		api := controller.NewAPI(&service)

		api.SetupApp(app)

		registeredUser1 := model.User{
			ID:          "3c0bbdae",
			Name:        "James",
			Surname:     "Bond",
			Email:       "test@gmail.com",
			Password:    "$2a$10$08qe8bXis2qObLNyEJfzpePCnqSJRyUXIa//ALLJw9l8q5gOTJljq",
			UserType:    "user",
			IsActivated: true,
		}
		registeredUser2 := model.User{
			ID:          utils.GenerateUUID(8),
			Name:        "James",
			Surname:     "Bond",
			Email:       "test@gmail.com",
			Password:    "$2a$10$08qe8bXis2qObLNyEJfzpePCnqSJRyUXIa//ALLJw9l8q5gOTJljq",
			UserType:    "user",
			IsActivated: true,
		}
		testRepository.RegisterUser(registeredUser1)
		testRepository.RegisterUser(registeredUser2)

		testPost1 := model.Post{
			ID:              utils.GenerateUUID(8),
			UserID:          registeredUser1.ID,
			User:            &registeredUser1,
			Description:     "Test Description 1",
			Image:           "zcxçömzcxözcxzzçcmzö 1",
			IsPrivate:       false,
			WhoLikesUserIDs: nil,
			CreatedAt:       time.Now().UTC().Add(-5 * time.Minute).Round(time.Second),
			UpdatedAt:       time.Now().UTC().Add(-5 * time.Minute).Round(time.Second),
		}

		testPost2 := model.Post{
			ID:              utils.GenerateUUID(8),
			UserID:          registeredUser2.ID,
			User:            &registeredUser2,
			Description:     "Test Description 2",
			Image:           "zcxçömzcxözcxzzçcmzö 2",
			IsPrivate:       false,
			WhoLikesUserIDs: nil,
			CreatedAt:       time.Now().UTC().Round(time.Second),
			UpdatedAt:       time.Now().UTC().Round(time.Second),
		}
		testPost3 := model.Post{
			ID:              utils.GenerateUUID(8),
			UserID:          registeredUser2.ID,
			User:            &registeredUser2,
			Description:     "Test Description 3",
			Image:           "zcxçömzcxözcxzzçcmzö 3",
			IsPrivate:       true,
			WhoLikesUserIDs: nil,
			CreatedAt:       time.Now().UTC().Add(-3 * time.Minute).Round(time.Second),
			UpdatedAt:       time.Now().UTC().Add(-3 * time.Minute).Round(time.Second),
		}
		testRepository.CreatePost(testPost1)
		testRepository.CreatePost(testPost2)
		testRepository.CreatePost(testPost3)

		Convey("When user send get posts request with userId query", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6InVzZXIiLCJpc3MiOiIzYzBiYmRhZSJ9.F_7cDDzm0THldtJLLNunfdXtoKqLKeMK8BdHG9Dxi-s"

			req, err := http.NewRequest(http.MethodGet, "/user/posts?userId="+registeredUser2.ID, nil)
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then all public posts should return", func() {
				actualResult := []model.Post{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult, ShouldHaveLength, 2)
				So(actualResult[0].ID, ShouldNotBeNil)
				So(actualResult[0].ID, ShouldEqual, testPost2.ID)
				So(actualResult[0].UserID, ShouldEqual, registeredUser2.ID)
				So(actualResult[0].Description, ShouldEqual, testPost2.Description)
				So(actualResult[0].Image, ShouldEqual, testPost2.Image)
				So(actualResult[0].IsPrivate, ShouldEqual, testPost2.IsPrivate)
				So(actualResult[0].WhoLikesUserIDs, ShouldEqual, testPost2.WhoLikesUserIDs)
				So(actualResult[0].User, ShouldResemble, &registeredUser2)
				So(actualResult[0].IsPrivate, ShouldBeFalse)
				So(actualResult[0].CreatedAt, ShouldEqual, testPost2.CreatedAt)
				So(actualResult[0].UpdatedAt, ShouldEqual, testPost2.UpdatedAt)

				So(actualResult[1].ID, ShouldNotBeNil)
				So(actualResult[1].ID, ShouldEqual, testPost3.ID)
				So(actualResult[1].UserID, ShouldEqual, registeredUser2.ID)
				So(actualResult[1].Description, ShouldEqual, testPost3.Description)
				So(actualResult[1].Image, ShouldEqual, testPost3.Image)
				So(actualResult[1].IsPrivate, ShouldEqual, testPost3.IsPrivate)
				So(actualResult[1].WhoLikesUserIDs, ShouldEqual, testPost3.WhoLikesUserIDs)
				So(actualResult[1].User, ShouldResemble, &registeredUser2)
				So(actualResult[1].IsPrivate, ShouldBeTrue)
				So(actualResult[1].CreatedAt, ShouldEqual, testPost3.CreatedAt)
				So(actualResult[1].UpdatedAt, ShouldEqual, testPost3.UpdatedAt)
			})
		})
	})
}

func TestLikePost(t *testing.T) {
	Convey("Given posts data", t, func() {
		app := fiber.New()
		testRepository := GetCleanTestRepository()
		middleware.SetupMiddleWare(app, *testRepository)
		service := service.NewService(testRepository)
		api := controller.NewAPI(&service)

		api.SetupApp(app)

		registeredUser1 := model.User{
			ID:          "3c0bbdae",
			Name:        "James",
			Surname:     "Bond",
			Email:       "test@gmail.com",
			Password:    "$2a$10$08qe8bXis2qObLNyEJfzpePCnqSJRyUXIa//ALLJw9l8q5gOTJljq",
			UserType:    "user",
			IsActivated: true,
		}
		registeredUser2 := model.User{
			ID:          utils.GenerateUUID(8),
			Name:        "James",
			Surname:     "Bond",
			Email:       "test@gmail.com",
			Password:    "$2a$10$08qe8bXis2qObLNyEJfzpePCnqSJRyUXIa//ALLJw9l8q5gOTJljq",
			UserType:    "user",
			IsActivated: true,
		}
		testRepository.RegisterUser(registeredUser1)
		testRepository.RegisterUser(registeredUser2)

		testPost1 := model.Post{
			ID:              utils.GenerateUUID(8),
			UserID:          registeredUser2.ID,
			User:            &registeredUser2,
			Description:     "Test Description 1",
			Image:           "zcxçömzcxözcxzzçcmzö 1",
			IsPrivate:       false,
			WhoLikesUserIDs: nil,
			CreatedAt:       time.Now().UTC().Add(-5 * time.Minute).Round(time.Second),
			UpdatedAt:       time.Now().UTC().Add(-5 * time.Minute).Round(time.Second),
		}
		testRepository.CreatePost(testPost1)

		Convey("When user send like posts request with userId", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6InVzZXIiLCJpc3MiOiIzYzBiYmRhZSJ9.F_7cDDzm0THldtJLLNunfdXtoKqLKeMK8BdHG9Dxi-s"

			likePostDTO := model.LikePostDTO{
				UserID: registeredUser1.ID,
			}
			reqBody, err := json.Marshal(likePostDTO)
			So(err, ShouldBeNil)

			req, err := http.NewRequest(http.MethodPatch, "/user/posts/"+testPost1.ID+"/like", bytes.NewReader(reqBody))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then all public posts should return", func() {
				actualResult := model.Post{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult.WhoLikesUserIDs, ShouldContain, registeredUser1.ID)
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
