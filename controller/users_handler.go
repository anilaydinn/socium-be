package controller

import (
	"github.com/anilaydinn/socium-be/errors"
	"github.com/anilaydinn/socium-be/model"
	"github.com/anilaydinn/socium-be/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

func (h *Handler) RegisterUserHandler(c *fiber.Ctx) error {
	userDTO := model.UserDTO{}

	err := c.BodyParser(&userDTO)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	user, err := h.service.RegisterUser(userDTO)

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

func (h *Handler) LoginUserHandler(c *fiber.Ctx) error {
	userCredentialsDTO := model.UserCredentialsDTO{}

	err := c.BodyParser(&userCredentialsDTO)

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	token, cookie, err := h.service.LoginUser(userCredentialsDTO)

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

func (h *Handler) ActivationHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")

	user, err := h.service.Activation(userID)

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

func (h *Handler) ForgotPasswordHandler(c *fiber.Ctx) error {
	forgotPasswordDTO := model.ForgotPasswordDTO{}
	err := c.BodyParser(&forgotPasswordDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	err = h.service.ForgotPassword(forgotPasswordDTO)

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

func (h *Handler) ResetPasswordHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	resetPasswordDTO := model.ResetPasswordDTO{}
	err := c.BodyParser(&resetPasswordDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	err = h.service.ResetPassword(userID, resetPasswordDTO)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (h *Handler) GetUserHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")

	if len(userID) == 0 {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	user, err := h.service.GetUser(userID)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(user)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (h *Handler) UpdateUserHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	updateUserDTO := model.UpdateUserDTO{}
	err := c.BodyParser(&updateUserDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	updatedUser, err := h.service.UpdateUser(userID, updateUserDTO)
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

func (h *Handler) SendFriendRequestHandler(c *fiber.Ctx) error {
	targetUserID := c.Params("targetUserID")
	friendRequestDTO := model.FriendRequestDTO{}
	err := c.BodyParser(&friendRequestDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	updatedUser, err := h.service.SendFriendRequest(targetUserID, friendRequestDTO)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(updatedUser)
	default:
		c.Status(fiber.StatusInternalServerError)

	}
	return nil
}

func (h *Handler) GetUserFriendRequestsHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	users, err := h.service.GetUserFriendRequests(userID)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(users)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (h *Handler) AcceptOrDeclineUserFriendRequestHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	targetID := c.Params("targetID")
	acceptOrDeclineFriendRequestDTO := model.AcceptOrDeclineFriendRequestDTO{}
	err := c.BodyParser(&acceptOrDeclineFriendRequestDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	user, err := h.service.AcceptOrDeclineUserFriendRequest(userID, targetID, acceptOrDeclineFriendRequestDTO)
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

func (h *Handler) GetUserFriendsHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")

	friends, err := h.service.GetUserFriends(userID)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(friends)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (h *Handler) GetUsersWithFilterHandler(c *fiber.Ctx) error {
	filter := c.Query("filter")
	var filterArr []string
	if strings.Contains(filter, " ") {
		filterArr = strings.Split(filter, " ")
	} else {
		filterArr = append(filterArr, filter)
	}

	users, err := h.service.GetUsersWithFilter(filterArr)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(users)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (h *Handler) GetAllUsersHandler(c *fiber.Ctx) error {
	filter := c.Query("filter")
	var filterArr []string
	if strings.Contains(filter, " ") {
		filterArr = strings.Split(filter, " ")
	} else {
		filterArr = append(filterArr, filter)
	}

	pageStr := c.Query("page")
	page := 0
	if len(pageStr) != 0 {
		var err error
		page, err = strconv.Atoi(pageStr)
		if page < 0 || err != nil {
			c.Status(fiber.StatusBadRequest)
			return err
		}
	}

	sizeStr := c.Query("size")
	size := utils.MaxInt
	if len(sizeStr) != 0 {
		var err error
		size, err = strconv.Atoi(sizeStr)
		if size < 0 || err != nil {
			c.Status(fiber.StatusBadRequest)
			return err
		}
	}

	users, err := h.service.GetAllUsers(page, size, filterArr)
	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(users)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (h *Handler) AdminGetUserHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")

	if len(userID) == 0 {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	user, err := h.service.AdminGetUser(userID)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(user)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (h *Handler) AdminGetUserPosts(c *fiber.Ctx) error {
	userID := c.Params("userID")
	if len(userID) == 0 {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	posts, err := h.service.GetUserPosts(userID)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(posts)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (h *Handler) GetNearUsersHandler(c *fiber.Ctx) error {
	userID := c.Params("userID")
	if len(userID) == 0 {
		c.Status(fiber.StatusBadRequest)
		return nil
	}
	getNearUsersDTO := model.GetNearUsersDTO{}
	err := c.BodyParser(&getNearUsersDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
	}
	users, err := h.service.GetNearUsers(userID, getNearUsersDTO)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(users)
	default:
		c.Status(fiber.StatusInternalServerError)

	}
	return nil
}
