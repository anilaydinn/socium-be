package model

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type User struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	Surname              string    `json:"surname"`
	Email                string    `json:"email"`
	BirthDate            time.Time `json:"birthDate"`
	Description          string    `json:"description"`
	ProfileImage         string    `json:"profileImage"`
	FriendRequestUserIDs []string  `json:"friendRequestUserIDs"`
	FriendIDs            []string  `json:"friendIds"`
	Password             string    `json:"password"`
	UserType             string    `json:"userType"`
	IsActivated          bool      `json:"isActivated"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	Latitude             float64   `json:"latitude"`
	Longitude            float64   `json:"longitude"`
}

type UserDTO struct {
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Email     string    `json:"email"`
	BirthDate time.Time `json:"birthDate"`
	Password  string    `json:"password"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}

type UsersPageableResponse struct {
	Users []User `json:"users"`
	Page  Page   `json:"page"`
}

type Page struct {
	Number        int `json:"number"`
	Size          int `json:"size",omitempty`
	TotalElements int `json:"totalElements",omitempty"`
	TotalPages    int `json:"totalPages",omitempty"`
}

type UserCredentialsDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserDTO struct {
	Description  string `json:"description"`
	ProfileImage string `json:"profileImage"`
}

type ForgotPasswordDTO struct {
	Email string `json:"email"`
}

type ResetPasswordDTO struct {
	Password string `json:"password"`
}

type FriendRequestDTO struct {
	UserID string `json:"userId"`
}

type FriendRequestIDsDTO struct {
	UserIDs []string `json:"userIds"`
}

type AcceptOrDeclineFriendRequestDTO struct {
	Accept bool `json:"accept"`
}

type GetFriendPostsDTO struct {
	FriendIDs []string `json:"friendIds"`
}

type GetNearUsersDTO struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Token struct {
	Token string `json:"token"`
}

type CustomClaims struct {
	UserType string `json:"userType"`
	jwt.StandardClaims
}
