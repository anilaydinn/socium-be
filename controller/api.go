package controller

import (
	"github.com/anilaydinn/socium-be/errors"
	"github.com/anilaydinn/socium-be/model"
	"github.com/anilaydinn/socium-be/service"
	"github.com/gofiber/fiber/v2"
)

type API struct {
	service *service.Service
}

func (api *API) SetupApp(app *fiber.App) {
	app.Post("/api/register", api.RegisterUserHandler)
	app.Post("/api/login", api.LoginUserHandler)
	app.Get("/api/activation/:userID", api.ActivationHandler)
	app.Post("/api/forgotPassword", api.ForgotPasswordHandler)
	app.Patch("/api/resetPassword/:userID", api.ResetPasswordHandler)
	app.Get("/api/users/:userID", api.GetUserHandler)
	app.Post("/user/posts", api.CreatePostHandler)
	app.Get("/user/posts", api.GetPostsHandler)
	app.Patch("/user/posts/:postID/like", api.LikePostHandler)
	app.Post("/user/posts/:postID/comments", api.AddPostCommentHandler)
	app.Patch("/user/users/:userID", api.UpdateUserHandler)
	app.Post("/user/users/:targetUserID/friendRequests", api.SendFriendRequestHandler)
	app.Get("/user/users/:userID/friendRequests", api.GetUserFriendRequestsHandler)
	app.Post("/user/users/:userID/friendRequests/:targetID", api.AcceptOrDeclineUserFriendRequestHandler)
	app.Get("/user/users/:userID/friends", api.GetUserFriendsHandler)
}

func NewAPI(service *service.Service) API {
	return API{
		service: service,
	}
}

func (api *API) RegisterUserHandler(c *fiber.Ctx) error {
	userDTO := model.UserDTO{}

	err := c.BodyParser(&userDTO)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	user, err := api.service.RegisterUser(userDTO)

	switch err {
	case nil:
		c.JSON(user)
		c.Status(fiber.StatusCreated)
	case errors.UserAlreadyRegistered:
		c.Status(fiber.StatusBadRequest)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) LoginUserHandler(c *fiber.Ctx) error {
	userCredentialsDTO := model.UserCredentialsDTO{}

	err := c.BodyParser(&userCredentialsDTO)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	token, cookie, err := api.service.LoginUser(userCredentialsDTO)

	switch err {
	case nil:
		c.JSON(token)
		c.Cookie(cookie)
		c.Status(fiber.StatusOK)
	case errors.UserNotFound:
		c.Status(fiber.StatusBadRequest)
	case errors.Unauthorized:
		c.Status(fiber.StatusUnauthorized)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) ActivationHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")

	user, err := api.service.Activation(userID)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(user)
	case errors.UserNotFound:
		c.Status(fiber.StatusNotFound)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) ForgotPasswordHandler(c *fiber.Ctx) error {
	forgotPasswordDTO := model.ForgotPasswordDTO{}
	err := c.BodyParser(&forgotPasswordDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	err = api.service.ForgotPassword(forgotPasswordDTO)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
	case errors.UserNotFound:
		c.Status(fiber.StatusNotFound)
	case errors.UserNotActivated:
		c.Status(fiber.StatusBadRequest)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) ResetPasswordHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	resetPasswordDTO := model.ResetPasswordDTO{}
	err := c.BodyParser(&resetPasswordDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	err = api.service.ResetPassword(userID, resetPasswordDTO)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) GetUserHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")

	if len(userID) == 0 {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	user, err := api.service.GetUser(userID)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(user)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) CreatePostHandler(c *fiber.Ctx) error {
	postDTO := model.PostDTO{}
	err := c.BodyParser(&postDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	post, err := api.service.CreatePost(postDTO)

	switch err {
	case nil:
		c.Status(fiber.StatusCreated)
		c.JSON(post)
	case errors.PostNotFound:
		c.Status(fiber.StatusNotFound)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) GetPostsHandler(c *fiber.Ctx) error {
	userID := c.Query("userId")
	getFriendPostsDTO := model.GetFriendPostsDTO{}
	_ = c.BodyParser(&getFriendPostsDTO)

	posts, err := api.service.GetPosts(userID, getFriendPostsDTO.FriendIDs)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(posts)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) LikePostHandler(c *fiber.Ctx) error {
	postID := c.Params("postID")
	likePostDTO := model.LikePostDTO{}
	err := c.BodyParser(&likePostDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	post, err := api.service.LikePost(postID, likePostDTO)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(post)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) AddPostCommentHandler(c *fiber.Ctx) error {
	postID := c.Params("postID")
	commentDTO := model.CommentDTO{}
	err := c.BodyParser(&commentDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	post, err := api.service.AddPostComment(postID, commentDTO)

	switch err {
	case nil:
		c.Status(fiber.StatusCreated)
		c.JSON(post)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) UpdateUserHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	updateUserDTO := model.UpdateUserDTO{}
	err := c.BodyParser(&updateUserDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	updatedUser, err := api.service.UpdateUser(userID, updateUserDTO)
	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(updatedUser)
	case errors.UserNotFound:
		c.Status(fiber.StatusNotFound)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) SendFriendRequestHandler(c *fiber.Ctx) error {
	targetUserID := c.Params("targetUserID")
	friendRequestDTO := model.FriendRequestDTO{}
	err := c.BodyParser(&friendRequestDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	updatedUser, err := api.service.SendFriendRequest(targetUserID, friendRequestDTO)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(updatedUser)
	default:
		c.Status(fiber.StatusInternalServerError)

	}
	return nil
}

func (api *API) GetUserFriendRequestsHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	users, err := api.service.GetUserFriendRequests(userID)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(users)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) AcceptOrDeclineUserFriendRequestHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	targetID := c.Params("targetID")
	acceptOrDeclineFriendRequestDTO := model.AcceptOrDeclineFriendRequestDTO{}
	err := c.BodyParser(&acceptOrDeclineFriendRequestDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	user, err := api.service.AcceptOrDeclineUserFriendRequest(userID, targetID, acceptOrDeclineFriendRequestDTO)
	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(user)
	case errors.UserNotFound:
		c.Status(fiber.StatusNotFound)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (api *API) GetUserFriendsHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")

	friends, err := api.service.GetUserFriends(userID)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(friends)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}
