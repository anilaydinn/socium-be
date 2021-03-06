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
			FriendIDs:   []string{"2dbbds32"},
			UserType:    "user",
			IsActivated: true,
		}
		registeredUser2 := model.User{
			ID:          "2dbbds32",
			Name:        "James",
			Surname:     "Bond",
			Email:       "test@gmail.com",
			Password:    "$2a$10$08qe8bXis2qObLNyEJfzpePCnqSJRyUXIa//ALLJw9l8q5gOTJljq",
			FriendIDs:   []string{"3c0bbdae"},
			UserType:    "user",
			IsActivated: true,
		}
		registeredUser3 := model.User{
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
		testRepository.RegisterUser(registeredUser3)

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
			IsPrivate:       false,
			WhoLikesUserIDs: nil,
			CreatedAt:       time.Now().UTC().Add(-3 * time.Minute).Round(time.Second),
			UpdatedAt:       time.Now().UTC().Add(-3 * time.Minute).Round(time.Second),
		}
		testPost4 := model.Post{
			ID:              utils.GenerateUUID(8),
			UserID:          registeredUser3.ID,
			User:            &registeredUser3,
			Description:     "Test Description 4",
			Image:           "zcxçömzcxözcxzzçcmzö 4",
			IsPrivate:       false,
			WhoLikesUserIDs: nil,
			CreatedAt:       time.Now().UTC().Add(-3 * time.Minute).Round(time.Second),
			UpdatedAt:       time.Now().UTC().Add(-3 * time.Minute).Round(time.Second),
		}
		testPost5 := model.Post{
			ID:              utils.GenerateUUID(8),
			UserID:          registeredUser1.ID,
			User:            &registeredUser1,
			Description:     "Test Description 5",
			Image:           "zcxçömzcxözcxzzçcmzö 5",
			IsPrivate:       false,
			WhoLikesUserIDs: nil,
			CreatedAt:       time.Now().UTC().Add(-2 * time.Minute).Round(time.Second),
			UpdatedAt:       time.Now().UTC().Add(-2 * time.Minute).Round(time.Second),
		}
		testRepository.CreatePost(testPost1)
		testRepository.CreatePost(testPost2)
		testRepository.CreatePost(testPost3)
		testRepository.CreatePost(testPost4)
		testRepository.CreatePost(testPost5)

		Convey("When user send get posts request", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6InVzZXIiLCJpc3MiOiIzYzBiYmRhZSJ9.F_7cDDzm0THldtJLLNunfdXtoKqLKeMK8BdHG9Dxi-s"

			req, err := http.NewRequest(http.MethodGet, "/user/posts?userId=3c0bbdae&homepage=true&friendIdList=2dbbds32", nil)
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

				So(actualResult, ShouldHaveLength, 4)
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
				So(actualResult[1].ID, ShouldEqual, testPost5.ID)
				So(actualResult[1].UserID, ShouldEqual, registeredUser1.ID)
				So(actualResult[1].Description, ShouldEqual, testPost5.Description)
				So(actualResult[1].Image, ShouldEqual, testPost5.Image)
				So(actualResult[1].IsPrivate, ShouldEqual, testPost5.IsPrivate)
				So(actualResult[1].WhoLikesUserIDs, ShouldEqual, testPost5.WhoLikesUserIDs)
				So(actualResult[1].User, ShouldResemble, &registeredUser1)
				So(actualResult[1].IsPrivate, ShouldBeFalse)
				So(actualResult[1].CreatedAt, ShouldEqual, testPost5.CreatedAt)
				So(actualResult[1].UpdatedAt, ShouldEqual, testPost5.UpdatedAt)

				So(actualResult[2].ID, ShouldNotBeNil)
				So(actualResult[2].ID, ShouldEqual, testPost3.ID)
				So(actualResult[2].UserID, ShouldEqual, registeredUser2.ID)
				So(actualResult[2].Description, ShouldEqual, testPost3.Description)
				So(actualResult[2].Image, ShouldEqual, testPost3.Image)
				So(actualResult[2].IsPrivate, ShouldEqual, testPost3.IsPrivate)
				So(actualResult[2].WhoLikesUserIDs, ShouldEqual, testPost3.WhoLikesUserIDs)
				So(actualResult[2].User, ShouldResemble, &registeredUser2)
				So(actualResult[2].IsPrivate, ShouldBeFalse)
				So(actualResult[2].CreatedAt, ShouldEqual, testPost3.CreatedAt)
				So(actualResult[2].UpdatedAt, ShouldEqual, testPost3.UpdatedAt)

				So(actualResult[3].ID, ShouldNotBeNil)
				So(actualResult[3].ID, ShouldEqual, testPost1.ID)
				So(actualResult[3].UserID, ShouldEqual, registeredUser2.ID)
				So(actualResult[3].Description, ShouldEqual, testPost1.Description)
				So(actualResult[3].Image, ShouldEqual, testPost1.Image)
				So(actualResult[3].IsPrivate, ShouldEqual, testPost1.IsPrivate)
				So(actualResult[3].WhoLikesUserIDs, ShouldEqual, testPost1.WhoLikesUserIDs)
				So(actualResult[3].User, ShouldResemble, &registeredUser2)
				So(actualResult[3].IsPrivate, ShouldBeFalse)
				So(actualResult[3].CreatedAt, ShouldEqual, testPost1.CreatedAt)
				So(actualResult[3].UpdatedAt, ShouldEqual, testPost1.UpdatedAt)
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

		commentID := utils.GenerateUUID(8)

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
			CommentIDs: []string{
				commentID,
			},
			CreatedAt: time.Now().UTC().Round(time.Second),
			UpdatedAt: time.Now().UTC().Round(time.Second),
		}
		testPost3 := model.Post{
			ID:              utils.GenerateUUID(8),
			UserID:          registeredUser2.ID,
			User:            &registeredUser2,
			Description:     "Test Description 3",
			Image:           "zcxçömzcxözcxzzçcmzö 3",
			IsPrivate:       true,
			WhoLikesUserIDs: nil,
			CommentIDs: []string{
				commentID,
			},
			CreatedAt: time.Now().UTC().Add(-3 * time.Minute).Round(time.Second),
			UpdatedAt: time.Now().UTC().Add(-3 * time.Minute).Round(time.Second),
		}
		testRepository.CreatePost(testPost1)
		testRepository.CreatePost(testPost2)
		testRepository.CreatePost(testPost3)

		testComment1 := model.Comment{
			ID:        commentID,
			UserID:    registeredUser2.ID,
			PostID:    testPost1.ID,
			Content:   "Comment",
			CreatedAt: time.Now().UTC().Round(time.Second),
			UpdatedAt: time.Now().UTC().Round(time.Second),
		}
		testRepository.AddComment(testComment1)

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
				So(actualResult[0].Comments[0].Content, ShouldEqual, testComment1.Content)
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
				So(actualResult[1].Comments[0].Content, ShouldEqual, testComment1.Content)
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

			Convey("Then user id should be inside the whoLikes array", func() {
				actualResult := model.Post{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult.WhoLikesUserIDs, ShouldContain, registeredUser1.ID)
			})
		})
	})
}

func TestLikeAlreadyLikedPost(t *testing.T) {
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
			ID:          utils.GenerateUUID(8),
			UserID:      registeredUser2.ID,
			User:        &registeredUser2,
			Description: "Test Description 1",
			Image:       "zcxçömzcxözcxzzçcmzö 1",
			IsPrivate:   false,
			WhoLikesUserIDs: []string{
				"3c0bbdae",
			},
			CreatedAt: time.Now().UTC().Add(-5 * time.Minute).Round(time.Second),
			UpdatedAt: time.Now().UTC().Add(-5 * time.Minute).Round(time.Second),
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

			Convey("Then user id should not be inside the whoLikes array", func() {
				actualResult := model.Post{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult.WhoLikesUserIDs, ShouldHaveLength, 0)
				So(actualResult.WhoLikesUserIDs, ShouldNotContain, registeredUser1.ID)
			})
		})
	})
}

func TestAddCommentToPost(t *testing.T) {
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
			ID:          utils.GenerateUUID(8),
			UserID:      registeredUser2.ID,
			User:        &registeredUser2,
			Description: "Test Description 1",
			Image:       "zcxçömzcxözcxzzçcmzö 1",
			IsPrivate:   false,
			WhoLikesUserIDs: []string{
				"3c0bbdae",
			},
			CommentIDs: []string{},
			CreatedAt:  time.Now().UTC().Add(-5 * time.Minute).Round(time.Second),
			UpdatedAt:  time.Now().UTC().Add(-5 * time.Minute).Round(time.Second),
		}
		testRepository.CreatePost(testPost1)

		Convey("When user send add comment request with userId and postId", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6InVzZXIiLCJpc3MiOiIzYzBiYmRhZSJ9.F_7cDDzm0THldtJLLNunfdXtoKqLKeMK8BdHG9Dxi-s"

			commentDTO := model.CommentDTO{
				UserID:  registeredUser1.ID,
				Content: "New comment",
			}
			reqBody, err := json.Marshal(commentDTO)
			So(err, ShouldBeNil)

			req, err := http.NewRequest(http.MethodPost, "/user/posts/"+testPost1.ID+"/comments", bytes.NewReader(reqBody))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 201", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusCreated)
			})

			Convey("Then comment should be added to post comments", func() {
				actualResult := model.Post{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult.CommentIDs, ShouldHaveLength, 1)
				So(actualResult.Comments, ShouldHaveLength, 1)
				So(actualResult.Comments[0].Content, ShouldEqual, commentDTO.Content)
			})
		})
	})
}

func TestAdminDeleteUserPost(t *testing.T) {
	Convey("Given admin, registered user and post data", t, func() {
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

		Convey("When admin user send delete user post request with id params", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6ImFkbWluIiwiaXNzIjoiM2MwYmJkYWUifQ.aYf3WQryPbYoexgG18Q9iWYbnLtnH2ueE_rgTFdqBx4"

			req, err := http.NewRequest(http.MethodDelete, "/admin/users/"+registeredUser2.ID+"/posts/"+post1.ID, nil)
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 204", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusNoContent)
			})

			Convey("Then user post should deleted", func() {
				post, err := testRepository.GetPost(post1.ID)
				So(post, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestGetWhoLikesPost(t *testing.T) {
	Convey("Given that register user", t, func() {
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
			UserID:          utils.GenerateUUID(8),
			User:            nil,
			Description:     "Hello",
			Image:           "adsknasnjkadsnjka",
			IsPrivate:       false,
			WhoLikesUserIDs: []string{registeredUser2.ID},
			CommentIDs:      nil,
			Comments:        nil,
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
		}
		testRepository.CreatePost(post1)

		Convey("When user send get who likes request", func() {
			bearerToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyVHlwZSI6InVzZXIiLCJpc3MiOiIzYzBiYmRhZSJ9.F_7cDDzm0THldtJLLNunfdXtoKqLKeMK8BdHG9Dxi-s"

			req, err := http.NewRequest(http.MethodGet, "/user/posts/"+post1.ID+"/likes?whoLikesUserIds="+registeredUser2.ID, nil)
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)

			res, err := app.Test(req, 30000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 200", func() {
				So(res.StatusCode, ShouldEqual, fiber.StatusOK)
			})

			Convey("Then all likes users should return", func() {
				actualResult := []model.User{}
				httpResponseBody, _ := ioutil.ReadAll(res.Body)
				err := json.Unmarshal(httpResponseBody, &actualResult)
				So(err, ShouldBeNil)

				So(actualResult[0].ID, ShouldEqual, registeredUser2.ID)
				So(actualResult[0].Name, ShouldEqual, registeredUser2.Name)
			})
		})

	})
}
