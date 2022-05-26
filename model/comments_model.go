package model

import "time"

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
