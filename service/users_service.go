package service

import (
	"github.com/anilaydinn/socium-be/email"
	"github.com/anilaydinn/socium-be/errors"
	"github.com/anilaydinn/socium-be/model"
	"github.com/anilaydinn/socium-be/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"math"
	"os"
	"time"
)

func (service *Service) RegisterUser(userDTO model.UserDTO) (*model.User, error) {
	alreadyRegisteredUser, err := service.repository.GetUserByEmail(userDTO.Email)
	if alreadyRegisteredUser != nil {
		return nil, errors.UserAlreadyRegistered
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := model.User{
		ID:          utils.GenerateUUID(8),
		Name:        userDTO.Name,
		Surname:     userDTO.Surname,
		Email:       userDTO.Email,
		BirthDate:   userDTO.BirthDate,
		Password:    string(hashedPassword),
		UserType:    "user",
		IsActivated: false,
		CreatedAt:   time.Now().UTC().Round(time.Minute),
		UpdatedAt:   time.Now().UTC().Round(time.Minute),
		Latitude:    userDTO.Latitude,
		Longitude:   userDTO.Longitude,
	}

	newUser, err := service.repository.RegisterUser(user)

	if err != nil {
		return nil, err
	}

	err = email.SendMail(newUser.Email, "Complete Registration", "Please click "+os.Getenv("REACT_HOSTNAME")+"/activation/"+user.ID)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (service *Service) LoginUser(userCredentialsDTO model.UserCredentialsDTO) (*model.Token, *fiber.Cookie, error) {
	user, err := service.repository.GetUserByEmail(userCredentialsDTO.Email)

	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		return nil, nil, errors.UserNotFound
	}

	if !user.IsActivated {
		return nil, nil, errors.Unauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userCredentialsDTO.Password)); err != nil {
		return nil, nil, errors.WrongPassword
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, model.CustomClaims{
		UserType: user.UserType,
		StandardClaims: jwt.StandardClaims{
			Issuer: user.ID,
		},
	})

	token, err := claims.SignedString([]byte(""))

	if err != nil {
		return nil, nil, err
	}

	cookie := fiber.Cookie{
		Name:    "user-token",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 24),
	}
	return &model.Token{
		Token: token,
	}, &cookie, nil
}

func (service *Service) Activation(userID string) (*model.User, error) {
	user, err := service.repository.GetUser(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.UserNotFound
	}

	if !user.IsActivated {
		user.IsActivated = true
	} else {
		return nil, errors.UserAlreadyActivated
	}

	return service.repository.UpdateUser(userID, *user)
}

func (service *Service) ForgotPassword(forgotPasswordDTO model.ForgotPasswordDTO) error {
	registeredUser, _ := service.repository.GetUserByEmail(forgotPasswordDTO.Email)
	if registeredUser == nil {
		return errors.UserNotFound
	}

	if !registeredUser.IsActivated {
		return errors.UserNotActivated
	}

	err := email.SendMail(forgotPasswordDTO.Email, "Reset Password", "You can reset your password click "+os.Getenv("REACT_HOSTNAME")+"/reset-password/"+registeredUser.ID+" here.")
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) ResetPassword(userID string, resetPasswordDTO model.ResetPasswordDTO) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(resetPasswordDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user, err := service.repository.GetUser(userID)
	if err != nil {
		return errors.UserNotFound
	}

	user.Password = string(hashedPassword)

	_, err = service.repository.UpdateUser(userID, *user)
	if err != nil {
		return err
	}
	return nil
}

func (service *Service) GetUser(userID string) (*model.User, error) {
	return service.repository.GetUser(userID)
}

func (service *Service) UpdateUser(userID string, updateUserDTO model.UpdateUserDTO) (*model.User, error) {
	user, err := service.repository.GetUser(userID)
	if err != nil {
		return nil, errors.UserNotFound
	}
	user.Description = updateUserDTO.Description
	user.ProfileImage = updateUserDTO.ProfileImage

	updatedUser, err := service.repository.UpdateUser(userID, *user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (service *Service) SendFriendRequest(targetUserID string, friendRequestDTO model.FriendRequestDTO) (*model.User, error) {
	user, err := service.GetUser(targetUserID)
	if err != nil {
		return nil, errors.UserNotFound
	}
	if !utils.Contains(user.FriendRequestUserIDs, friendRequestDTO.UserID) {
		user.FriendRequestUserIDs = append(user.FriendRequestUserIDs, friendRequestDTO.UserID)
	}

	updatedUser, err := service.repository.UpdateUser(targetUserID, *user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (service *Service) GetUserFriendRequests(userID string) ([]model.User, error) {
	user, err := service.repository.GetUser(userID)
	if err != nil {
		return nil, errors.UserNotFound
	}

	friendRequestUsers, err := service.repository.GetUsersByIDList(user.FriendRequestUserIDs)
	if err != nil {
		return nil, err
	}

	return friendRequestUsers, nil
}

func (service *Service) AcceptOrDeclineUserFriendRequest(userID, targetID string, acceptOrDeclineFriendRequestDTO model.AcceptOrDeclineFriendRequestDTO) (*model.User, error) {
	user, err := service.repository.GetUser(userID)
	if err != nil {
		return nil, errors.UserNotFound
	}

	targetUser, err := service.repository.GetUser(targetID)
	if err != nil {
		return nil, errors.UserNotFound
	}

	if acceptOrDeclineFriendRequestDTO.Accept {
		user.FriendIDs = append(user.FriendIDs, targetID)
		targetUser.FriendIDs = append(targetUser.FriendIDs, userID)
		user.FriendRequestUserIDs = utils.RemoveElement(user.FriendRequestUserIDs, targetID)
	} else {
		user.FriendRequestUserIDs = utils.RemoveElement(user.FriendRequestUserIDs, targetID)
	}

	updatedUser, err := service.repository.UpdateUser(userID, *user)
	if err != nil {
		return nil, err
	}

	_, err = service.repository.UpdateUser(targetID, *targetUser)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (service *Service) GetUserFriends(userID string) ([]model.User, error) {
	user, err := service.repository.GetUser(userID)
	if err != nil {
		return nil, errors.UserNotFound
	}

	friends, err := service.repository.GetUsersByIDList(user.FriendIDs)
	if err != nil {
		return nil, err
	}

	return friends, nil
}

func (service *Service) GetUsersWithFilter(filterArr []string) ([]model.User, error) {
	users, err := service.repository.GetUsersWithFilter(filterArr)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (service *Service) GetAllUsers(pageNumber, size int, filterArr []string) (*model.UsersPageableResponse, error) {
	users, totalElements, err := service.repository.GetAllUsers(pageNumber, size, filterArr)
	if err != nil {
		return nil, err
	}
	page := model.Page{
		Number:        pageNumber,
		Size:          size,
		TotalElements: totalElements,
		TotalPages:    int(math.Ceil(float64(totalElements) / float64(size))),
	}

	return &model.UsersPageableResponse{
		Users: users,
		Page:  page,
	}, nil
}

func (service *Service) AdminGetUser(userID string) (*model.User, error) {
	return service.repository.GetUser(userID)
}

func (service *Service) GetUserPosts(userID string) ([]model.Post, error) {
	user, err := service.repository.GetUser(userID)
	if err != nil {
		return nil, errors.UserNotFound
	}

	posts, err := service.repository.GetUserPosts(userID)
	if err != nil {
		return nil, err
	}

	var postResults []model.Post
	for _, post := range posts {
		post.User = user
		postResults = append(postResults, post)
	}

	return postResults, nil
}

func (service *Service) GetNearUsers(userID string, getNearUsersDTO model.GetNearUsersDTO) ([]model.User, error) {
	users, _, err := service.repository.GetAllUsers(0, 0, []string{})
	if err != nil {
		return nil, err
	}

	nearUsers := []model.User{}
	for _, user := range users {
		if utils.CalculateDistanceKM(user.Latitude, user.Longitude, getNearUsersDTO.Latitude, getNearUsersDTO.Longitude, "K") <= 20 && user.ID != userID && user.Longitude > 0 && user.Latitude > 0 {
			nearUsers = append(nearUsers, user)
		}
	}

	return nearUsers, nil
}

func (service *Service) DeleteUserFriend(userID, friendID string) (*model.User, error) {
	user, err := service.repository.GetUser(userID)
	if err != nil {
		return nil, errors.UserNotFound
	}

	friend, err := service.repository.GetUser(friendID)
	if err != nil {
		return nil, errors.UserNotFound
	}

	user.FriendIDs = utils.RemoveElement(user.FriendIDs, friend.ID)
	updatedUser, err := service.repository.UpdateUser(user.ID, *user)
	if err != nil {
		return nil, err
	}

	friend.FriendIDs = utils.RemoveElement(friend.FriendIDs, user.ID)
	_, err = service.repository.UpdateUser(friend.ID, *friend)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
