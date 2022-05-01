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
