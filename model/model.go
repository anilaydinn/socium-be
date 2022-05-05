package model

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
	Email        string    `json:"email"`
	BirthDate    time.Time `json:"birthDate"`
	Description  string    `json:"description"`
	ProfileImage string    `json:"profileImage"`
	Password     string    `json:"password"`
	UserType     string    `json:"userType"`
	IsActivated  bool      `json:"isActivated"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type UserDTO struct {
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Email     string    `json:"email"`
	BirthDate time.Time `json:"birthDate"`
	Password  string    `json:"password"`
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
	CommentIDs      []string  `json:"commentIds"`
	Comments        []Comment `json:"comments"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type Comment struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	PostID    string    `json:"postId"`
	User      *User     `json:"user"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CommentDTO struct {
	UserID  string `json:"userId"`
	Content string `json:"content"`
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
