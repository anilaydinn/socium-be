package model

import "time"

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

type LikePostDTO struct {
	UserID string `json:"userId"`
}

type GetPostsQuery struct {
	UserID       string   `query:"userId"`
	Homepage     string   `query:"homepage"`
	FriendIDList []string `query:"friendIdList"`
}

type WhoLikesQuery struct {
	WhoLikesUserIDs []string `query:"whoLikesUserIds"`
}
