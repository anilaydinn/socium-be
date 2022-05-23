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

type GetFriendPostsDTO struct {
	UserID    string   `json:"userId"`
	FriendIDs []string `json:"friendIds"`
}

type Token struct {
	Token string `json:"token"`
}

type CustomClaims struct {
	UserType string `json:"userType"`
	jwt.StandardClaims
}
