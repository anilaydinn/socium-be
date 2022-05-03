package model

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	UserType    string `json:"userType"`
	IsActivated bool   `json:"isActivated"`
}

type UserDTO struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserCredentialsDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgotPasswordDTO struct {
	Email string `json:"email"`
}

type ResetPasswordDTO struct {
	Password string `json:"password"`
}

type PostDTO struct {
	UserID      string `json:"userId"`
	Description string `json:"description"`
	Image       string `json:"image"`
	IsPrivate   bool   `json:"isPrivate"`
}

type Post struct {
	ID              string    `json:"id"`
	UserID          string    `json:"userId"`
	User            *User     `json:"user"`
	Description     string    `json:"description"`
	Image           string    `json:"image"`
	IsPrivate       bool      `json:"isPrivate"`
	WhoLikesUserIDs []string  `json:"whoLikesUserIds"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type LikePostDTO struct {
	UserID string `json:"userId"`
}

type Token struct {
	Token string `json:"token"`
}

type CustomClaims struct {
	UserType string `json:"userType"`
	jwt.StandardClaims
}
