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
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

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
				Name:      "John",
				Surname:   "Obama",
				Email:     "john@gmail.com",
				BirthDate: time.Date(1998, 3, 16, 0, 0, 0, 0, time.Local),
				Password:  "123123",
				Latitude:  41.3843848,
				Longitude: 26.8328428,
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
				So(actualResult.BirthDate, ShouldEqual, userDTO.BirthDate)
				So(actualResult.UserType, ShouldEqual, "user")
				So(actualResult.IsActivated, ShouldBeFalse)
				So(actualResult.CreatedAt, ShouldEqual, time.Now().UTC().Round(time.Minute))
				So(actualResult.UpdatedAt, ShouldEqual, time.Now().UTC().Round(time.Minute))
				So(actualResult.Latitude, ShouldEqual, userDTO.Latitude)
				So(actualResult.Longitude, ShouldEqual, userDTO.Longitude)
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

func TestUpdateUser(t *testing.T) {
	Convey("Given registered user", t, func() {
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

		Convey("When user send update user request with userId and valid data", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6InVzZXIiLCJpc3MiOiIzYzBiYmRhZSJ9.F_7cDDzm0THldtJLLNunfdXtoKqLKeMK8BdHG9Dxi-s"

			updateUserDTO := model.UpdateUserDTO{
				Description:  "Test description",
				ProfileImage: "asdmasdlmkdsaads",
			}
			reqBody, err := json.Marshal(updateUserDTO)
			So(err, ShouldBeNil)

			req, err := http.NewRequest(http.MethodPatch, "/user/users/"+registeredUser.ID, bytes.NewReader(reqBody))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then user should be updated", func() {
				actualResult := model.User{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult.Description, ShouldEqual, updateUserDTO.Description)
				So(actualResult.ProfileImage, ShouldEqual, updateUserDTO.ProfileImage)
			})
		})
	})
}

func TestSendFriendRequest(t *testing.T) {
	Convey("Given registered users", t, func() {
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

		Convey("When user send friend request to target user with userId", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6InVzZXIiLCJpc3MiOiIzYzBiYmRhZSJ9.F_7cDDzm0THldtJLLNunfdXtoKqLKeMK8BdHG9Dxi-s"

			friendRequestDTO := model.FriendRequestDTO{
				UserID: registeredUser1.ID,
			}
			reqBody, err := json.Marshal(friendRequestDTO)
			So(err, ShouldBeNil)

			req, err := http.NewRequest(http.MethodPost, "/user/users/"+registeredUser2.ID+"/friendRequests", bytes.NewReader(reqBody))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then friend request should be sended", func() {
				actualResult := model.User{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult.FriendRequestUserIDs, ShouldNotBeNil)
				So(actualResult.FriendRequestUserIDs, ShouldHaveLength, 1)
				So(actualResult.FriendRequestUserIDs, ShouldContain, friendRequestDTO.UserID)

			})
		})
	})
}

func TestGetUserFriendRequests(t *testing.T) {
	Convey("Given registered users", t, func() {
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
			FriendRequestUserIDs: []string{"123123"},
			UserType:             "user",
			IsActivated:          true,
		}

		registeredUser2 := model.User{
			ID:          "123123",
			Name:        "James",
			Surname:     "Bond",
			Email:       "test@gmail.com",
			Password:    "$2a$10$08qe8bXis2qObLNyEJfzpePCnqSJRyUXIa//ALLJw9l8q5gOTJljq",
			UserType:    "user",
			IsActivated: true,
		}
		testRepository.RegisterUser(registeredUser1)
		testRepository.RegisterUser(registeredUser2)

		Convey("When user send friend request ids", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6InVzZXIiLCJpc3MiOiIzYzBiYmRhZSJ9.F_7cDDzm0THldtJLLNunfdXtoKqLKeMK8BdHG9Dxi-s"

			req, err := http.NewRequest(http.MethodGet, "/user/users/"+registeredUser1.ID+"/friendRequests", nil)
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then friend request should be sended", func() {
				actualResult := []model.User{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult, ShouldHaveLength, 1)
				So(actualResult[0].ID, ShouldEqual, registeredUser2.ID)
			})
		})
	})
}

func TestAcceptFriendRequest(t *testing.T) {
	Convey("Given registered users", t, func() {
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
			FriendRequestUserIDs: []string{"123123"},
			FriendIDs:            []string{},
			UserType:             "user",
			IsActivated:          true,
		}

		registeredUser2 := model.User{
			ID:          "123123",
			Name:        "James",
			Surname:     "Bond",
			Email:       "test@gmail.com",
			Password:    "$2a$10$08qe8bXis2qObLNyEJfzpePCnqSJRyUXIa//ALLJw9l8q5gOTJljq",
			UserType:    "user",
			IsActivated: true,
		}
		testRepository.RegisterUser(registeredUser1)
		testRepository.RegisterUser(registeredUser2)

		Convey("When user send accept friend request with user id", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6InVzZXIiLCJpc3MiOiIzYzBiYmRhZSJ9.F_7cDDzm0THldtJLLNunfdXtoKqLKeMK8BdHG9Dxi-s"

			acceptFriendRequestDTO := model.AcceptOrDeclineFriendRequestDTO{
				Accept: true,
			}
			reqBody, err := json.Marshal(acceptFriendRequestDTO)
			So(err, ShouldBeNil)

			req, err := http.NewRequest(http.MethodPost, "/user/users/"+registeredUser1.ID+"/friendRequests/123123", bytes.NewReader(reqBody))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then new friend should be added to user friends array", func() {
				actualResult := model.User{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult.ID, ShouldEqual, registeredUser1.ID)
				So(actualResult.FriendIDs, ShouldHaveLength, 1)
				So(actualResult.FriendIDs[0], ShouldEqual, "123123")
				So(actualResult.FriendRequestUserIDs, ShouldHaveLength, 0)
			})
		})
	})
}

func TestDeclineFriendRequest(t *testing.T) {
	Convey("Given registered users", t, func() {
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
			FriendRequestUserIDs: []string{"123123"},
			FriendIDs:            []string{},
			UserType:             "user",
			IsActivated:          true,
		}

		registeredUser2 := model.User{
			ID:          "123123",
			Name:        "James",
			Surname:     "Bond",
			Email:       "test@gmail.com",
			Password:    "$2a$10$08qe8bXis2qObLNyEJfzpePCnqSJRyUXIa//ALLJw9l8q5gOTJljq",
			UserType:    "user",
			IsActivated: true,
		}
		testRepository.RegisterUser(registeredUser1)
		testRepository.RegisterUser(registeredUser2)

		Convey("When user send accept friend request with user id", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6InVzZXIiLCJpc3MiOiIzYzBiYmRhZSJ9.F_7cDDzm0THldtJLLNunfdXtoKqLKeMK8BdHG9Dxi-s"

			acceptFriendRequestDTO := model.AcceptOrDeclineFriendRequestDTO{
				Accept: false,
			}
			reqBody, err := json.Marshal(acceptFriendRequestDTO)
			So(err, ShouldBeNil)

			req, err := http.NewRequest(http.MethodPost, "/user/users/"+registeredUser1.ID+"/friendRequests/123123", bytes.NewReader(reqBody))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then new friend should be added to user friends array", func() {
				actualResult := model.User{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult.ID, ShouldEqual, registeredUser1.ID)
				So(actualResult.FriendRequestUserIDs, ShouldHaveLength, 0)
			})
		})
	})
}

func TestGetUserFriends(t *testing.T) {
	Convey("Given registered users", t, func() {
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
			UserType:             "user",
			IsActivated:          true,
		}

		registeredUser2 := model.User{
			ID:          "123123",
			Name:        "James",
			Surname:     "Bond",
			Email:       "test@gmail.com",
			Password:    "$2a$10$08qe8bXis2qObLNyEJfzpePCnqSJRyUXIa//ALLJw9l8q5gOTJljq",
			UserType:    "user",
			IsActivated: true,
		}
		testRepository.RegisterUser(registeredUser1)
		testRepository.RegisterUser(registeredUser2)

		Convey("When user send get user friends request with user id", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6InVzZXIiLCJpc3MiOiIzYzBiYmRhZSJ9.F_7cDDzm0THldtJLLNunfdXtoKqLKeMK8BdHG9Dxi-s"

			req, err := http.NewRequest(http.MethodGet, "/user/users/"+registeredUser1.ID+"/friends", nil)
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then user friends should return", func() {
				actualResult := []model.User{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult, ShouldHaveLength, 1)
				So(actualResult[0].ID, ShouldEqual, registeredUser2.ID)
				So(actualResult[0].Name, ShouldEqual, registeredUser2.Name)
				So(actualResult[0].Surname, ShouldEqual, registeredUser2.Surname)
				So(actualResult[0].Email, ShouldEqual, registeredUser2.Email)
				So(actualResult[0].Password, ShouldEqual, registeredUser2.Password)
				So(actualResult[0].UserType, ShouldEqual, registeredUser2.UserType)
				So(actualResult[0].IsActivated, ShouldBeTrue)
			})
		})
	})
}

func TestSearchUser(t *testing.T) {
	Convey("Given registered users", t, func() {
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
			UserType:             "user",
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
		registeredUser3 := model.User{
			ID:          "321321",
			Name:        "Ahmet",
			Surname:     "Bond",
			Email:       "test2@gmail.com",
			Password:    "$2a$10$08qe8bXis2qObLNyEJfzpePCnqSJRyUXIa//ALLJw9l8q5gOTJljq",
			UserType:    "user",
			IsActivated: true,
		}
		testRepository.RegisterUser(registeredUser1)
		testRepository.RegisterUser(registeredUser2)
		testRepository.RegisterUser(registeredUser3)

		Convey("When user send search user request with name query", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6InVzZXIiLCJpc3MiOiIzYzBiYmRhZSJ9.F_7cDDzm0THldtJLLNunfdXtoKqLKeMK8BdHG9Dxi-s"

			req, err := http.NewRequest(http.MethodGet, "/user/users?filter=Ahmet", nil)
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then user friends should return", func() {
				actualResult := []model.User{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult, ShouldHaveLength, 1)
				So(actualResult[0].ID, ShouldEqual, registeredUser3.ID)
				So(actualResult[0].Name, ShouldEqual, registeredUser3.Name)
				So(actualResult[0].Surname, ShouldEqual, registeredUser3.Surname)
				So(actualResult[0].Email, ShouldEqual, registeredUser3.Email)
				So(actualResult[0].Password, ShouldEqual, registeredUser3.Password)
				So(actualResult[0].UserType, ShouldEqual, registeredUser3.UserType)
				So(actualResult[0].IsActivated, ShouldBeTrue)
			})
		})
	})
}

func TestAdminGetAllUsers(t *testing.T) {
	Convey("Given admin user", t, func() {
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
		registeredUser3 := model.User{
			ID:          "321321",
			Name:        "Ahmet",
			Surname:     "Bond",
			Email:       "test2@gmail.com",
			Password:    "$2a$10$08qe8bXis2qObLNyEJfzpePCnqSJRyUXIa//ALLJw9l8q5gOTJljq",
			UserType:    "user",
			IsActivated: true,
		}
		registeredUser4 := model.User{
			ID:          "32132131",
			Name:        "Sinan",
			Surname:     "Bond",
			Email:       "test3@gmail.com",
			Password:    "$2a$10$08qe8bXis2qObLNyEJfzpePCnqSJRyUXIa//ALLJw9l8q5gOTJljq",
			UserType:    "user",
			IsActivated: true,
		}
		testRepository.RegisterUser(registeredUser1)
		testRepository.RegisterUser(registeredUser2)
		testRepository.RegisterUser(registeredUser3)
		testRepository.RegisterUser(registeredUser4)

		Convey("When admin send get all users request", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6ImFkbWluIiwiaXNzIjoiM2MwYmJkYWUifQ.aYf3WQryPbYoexgG18Q9iWYbnLtnH2ueE_rgTFdqBx4"

			req, err := http.NewRequest(http.MethodGet, "/admin/users?page=0&size=2", nil)
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then all users should return", func() {
				actualResult := model.UsersPageableResponse{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult.Users, ShouldNotBeNil)
				So(actualResult.Users, ShouldHaveLength, 2)
				So(actualResult.Users[0].ID, ShouldEqual, registeredUser1.ID)
				So(actualResult.Users[1].ID, ShouldEqual, registeredUser2.ID)
				So(actualResult.Page.TotalPages, ShouldEqual, 2)
				So(actualResult.Page.Size, ShouldEqual, 2)
				So(actualResult.Page.Number, ShouldEqual, 0)
				So(actualResult.Page.TotalElements, ShouldEqual, 4)

			})
		})
	})
}

func TestAdminSearchUser(t *testing.T) {
	Convey("Given admin and registered users", t, func() {
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
		registeredUser3 := model.User{
			ID:          "321321",
			Name:        "Ahmet",
			Surname:     "Bond",
			Email:       "test2@gmail.com",
			Password:    "$2a$10$08qe8bXis2qObLNyEJfzpePCnqSJRyUXIa//ALLJw9l8q5gOTJljq",
			UserType:    "user",
			IsActivated: true,
		}
		testRepository.RegisterUser(registeredUser1)
		testRepository.RegisterUser(registeredUser2)
		testRepository.RegisterUser(registeredUser3)

		Convey("When admin user send search user request with name query", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6ImFkbWluIiwiaXNzIjoiM2MwYmJkYWUifQ.aYf3WQryPbYoexgG18Q9iWYbnLtnH2ueE_rgTFdqBx4"

			req, err := http.NewRequest(http.MethodGet, "/admin/users?filter=Ahmet", nil)
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then searched users should return", func() {
				actualResult := model.UsersPageableResponse{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult.Users, ShouldHaveLength, 1)
				So(actualResult.Users[0].ID, ShouldEqual, registeredUser3.ID)
				So(actualResult.Users[0].Name, ShouldEqual, registeredUser3.Name)
				So(actualResult.Users[0].Surname, ShouldEqual, registeredUser3.Surname)
				So(actualResult.Users[0].Email, ShouldEqual, registeredUser3.Email)
				So(actualResult.Users[0].Password, ShouldEqual, registeredUser3.Password)
				So(actualResult.Users[0].UserType, ShouldEqual, registeredUser3.UserType)
				So(actualResult.Users[0].IsActivated, ShouldBeTrue)
				So(actualResult.Page.TotalPages, ShouldEqual, 1)
				So(actualResult.Page.Size, ShouldEqual, utils.MaxInt)
				So(actualResult.Page.Number, ShouldEqual, 0)
				So(actualResult.Page.TotalElements, ShouldEqual, 1)
			})
		})
	})
}

func TestAdminGetUser(t *testing.T) {
	Convey("Given admin and registered users", t, func() {
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

		Convey("When admin user send get user request with id params", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6ImFkbWluIiwiaXNzIjoiM2MwYmJkYWUifQ.aYf3WQryPbYoexgG18Q9iWYbnLtnH2ueE_rgTFdqBx4"

			req, err := http.NewRequest(http.MethodGet, "/admin/users/"+registeredUser2.ID, nil)
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then searched users should return", func() {
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

func TestAdminGetUserPosts(t *testing.T) {
	Convey("Given admin and registered user", t, func() {
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

			req, err := http.NewRequest(http.MethodGet, "/admin/users/"+registeredUser2.ID+"/posts", nil)
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then user posts should return", func() {
				actualResult := []model.Post{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult[0].ID, ShouldNotBeNil)
				So(actualResult[0].ID, ShouldEqual, post1.ID)
				So(actualResult[0].User, ShouldResemble, post1.User)
				So(actualResult[0].UserID, ShouldEqual, registeredUser2.ID)
				So(actualResult[0].Description, ShouldEqual, post1.Description)
				So(actualResult[0].UpdatedAt, ShouldEqual, post1.UpdatedAt)
				So(actualResult[0].CreatedAt, ShouldEqual, post1.CreatedAt)
				So(actualResult[0].Image, ShouldEqual, post1.Image)
			})
		})
	})
}

func TestGetNearUsers(t *testing.T) {
	Convey("Given registered user", t, func() {
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
			UserType:             "user",
			IsActivated:          true,
			Longitude:            26.679363372044076,
			Latitude:             41.262426815780856,
		}

		registeredUser2 := model.User{
			ID:          "123123",
			Name:        "James",
			Surname:     "Bond",
			Email:       "test@gmail.com",
			Password:    "$2a$10$08qe8bXis2qObLNyEJfzpePCnqSJRyUXIa//ALLJw9l8q5gOTJljq",
			UserType:    "user",
			IsActivated: true,
			Longitude:   26.653,
			Latitude:    41.25,
		}
		registeredUser3 := model.User{
			ID:          "3123123",
			Name:        "James",
			Surname:     "Bond",
			Email:       "test@gmail.com",
			Password:    "$2a$10$08qe8bXis2qObLNyEJfzpePCnqSJRyUXIa//ALLJw9l8q5gOTJljq",
			UserType:    "user",
			IsActivated: true,
			Longitude:   24.253,
			Latitude:    39.2,
		}
		testRepository.RegisterUser(registeredUser1)
		testRepository.RegisterUser(registeredUser2)
		testRepository.RegisterUser(registeredUser3)

		Convey("When user send get near user request", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6InVzZXIiLCJpc3MiOiIzYzBiYmRhZSJ9.F_7cDDzm0THldtJLLNunfdXtoKqLKeMK8BdHG9Dxi-s"

			getNearUsersDTO := model.GetNearUsersDTO{
				Longitude: registeredUser1.Longitude,
				Latitude:  registeredUser1.Latitude,
			}
			reqBody, err := json.Marshal(getNearUsersDTO)
			So(err, ShouldBeNil)

			req, err := http.NewRequest(http.MethodPost, "/user/users/"+registeredUser1.ID+"/near", bytes.NewReader(reqBody))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then near user should return", func() {
				actualResult := []model.User{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult, ShouldHaveLength, 1)
				So(actualResult[0].ID, ShouldEqual, registeredUser2.ID)
				So(actualResult[0].Name, ShouldEqual, registeredUser2.Name)
				So(actualResult[0].Surname, ShouldEqual, registeredUser2.Surname)
				So(actualResult[0].Email, ShouldEqual, registeredUser2.Email)
				So(actualResult[0].Password, ShouldEqual, registeredUser2.Password)
				So(actualResult[0].UserType, ShouldEqual, registeredUser2.UserType)
				So(actualResult[0].IsActivated, ShouldBeTrue)
			})
		})
	})
}

func TestDeleteUserFriend(t *testing.T) {
	Convey("Given registered user", t, func() {
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
			UserType:             "user",
			IsActivated:          true,
			Longitude:            26.679363372044076,
			Latitude:             41.262426815780856,
		}

		registeredUser2 := model.User{
			ID:          "123123",
			Name:        "James",
			Surname:     "Bond",
			Email:       "test@gmail.com",
			FriendIDs:   []string{"3c0bbdae"},
			Password:    "$2a$10$08qe8bXis2qObLNyEJfzpePCnqSJRyUXIa//ALLJw9l8q5gOTJljq",
			UserType:    "user",
			IsActivated: true,
			Longitude:   26.653,
			Latitude:    41.25,
		}
		testRepository.RegisterUser(registeredUser1)
		testRepository.RegisterUser(registeredUser2)

		Convey("When delete user friend request send", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6InVzZXIiLCJpc3MiOiIzYzBiYmRhZSJ9.F_7cDDzm0THldtJLLNunfdXtoKqLKeMK8BdHG9Dxi-s"

			req, err := http.NewRequest(http.MethodGet, "/user/users/"+registeredUser1.ID+"/friends/"+registeredUser2.ID, nil)
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then near user should return", func() {
				actualResult := model.User{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult.ID, ShouldEqual, registeredUser1.ID)
				So(actualResult.Name, ShouldEqual, registeredUser1.Name)
				So(actualResult.Surname, ShouldEqual, registeredUser1.Surname)
				So(actualResult.Email, ShouldEqual, registeredUser1.Email)
				So(actualResult.Password, ShouldEqual, registeredUser1.Password)
				So(actualResult.UserType, ShouldEqual, registeredUser1.UserType)
				So(actualResult.FriendIDs, ShouldHaveLength, 0)
				So(actualResult.IsActivated, ShouldBeTrue)
			})
		})
	})
}
