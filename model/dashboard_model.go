package model

type DashboardInformation struct {
	UserCount    int `json:"userCount"`
	PostCount    int `json:"postCount"`
	CommentCount int `json:"commentCount"`
}
