package repository

import "time"

type UserEntity struct {
	ID                   string    `bson:"id"`
	Name                 string    `bson:"name"`
	Surname              string    `bson:"surname"`
	Email                string    `bson:"email"`
	BirthDate            time.Time `bson:"birthDate"`
	Description          string    `bson:"description"`
	ProfileImage         string    `bson:"profileImage"`
	FriendRequestUserIDs []string  `bson:"friendRequestUserIDs"`
	FriendIDs            []string  `json:"friendIds"`
	Password             string    `bson:"password"`
	UserType             string    `bson:"userType"`
	IsActivated          bool      `bson:"isActivated"`
	CreatedAt            time.Time `bson:"createdAt"`
	UpdatedAt            time.Time `bson:"updatedAt"`
	Latitude             float64   `bson:"latitude"`
	Longitude            float64   `bson:"longitude"`
}

type PostEntity struct {
	ID              string    `bson:"id"`
	UserID          string    `bson:"userId"`
	Description     string    `bson:"description"`
	Image           string    `bson:"image"`
	IsPrivate       bool      `bson:"isPrivate"`
	WhoLikesUserIDs []string  `bson:"whoLikesUserIds"`
	CommentIDs      []string  `bson:"commentIds"`
	CreatedAt       time.Time `bson:"createdAt"`
	UpdatedAt       time.Time `bson:"updatedAt"`
}

type CommentEntity struct {
	ID        string    `bson:"id"`
	UserID    string    `bson:"userId"`
	PostID    string    `bson:"postId"`
	Content   string    `bson:"content"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

type ContactEntity struct {
	ID      string `bson:"id"`
	Name    string `bson:"name"`
	Surname string `bson:"surname"`
	Email   string `bson:"email"`
	Message string `bson:"message"`
}
