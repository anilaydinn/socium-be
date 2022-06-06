package controller

import (
	"github.com/anilaydinn/socium-be/service"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *service.Service
}

func NewAPI(service *service.Service) Handler {
	return Handler{
		service: service,
	}
}

func (h *Handler) SetupApp(app *fiber.App) {
	app.Post("/api/register", h.RegisterUserHandler)
	app.Post("/api/login", h.LoginUserHandler)
	app.Get("/api/activation/:userID", h.ActivationHandler)
	app.Post("/api/forgotPassword", h.ForgotPasswordHandler)
	app.Patch("/api/resetPassword/:userID", h.ResetPasswordHandler)
	app.Get("/api/users/:userID", h.GetUserHandler)
	app.Post("/user/posts", h.CreatePostHandler)
	app.Get("/user/posts", h.GetPostsHandler)
	app.Patch("/user/posts/:postID/like", h.LikePostHandler)
	app.Post("/user/posts/:postID/comments", h.AddPostCommentHandler)
	app.Patch("/user/users/:userID", h.UpdateUserHandler)
	app.Post("/user/users/:targetUserID/friendRequests", h.SendFriendRequestHandler)
	app.Get("/user/users/:userID/friendRequests", h.GetUserFriendRequestsHandler)
	app.Post("/user/users/:userID/friendRequests/:targetID", h.AcceptOrDeclineUserFriendRequestHandler)
	app.Get("/user/users/:userID/friends", h.GetUserFriendsHandler)
	app.Post("/api/contacts", h.CreateContactHandler)
	app.Get("/user/users", h.GetUsersWithFilterHandler)
	app.Get("/admin/users", h.GetAllUsersHandler)
	app.Get("/admin/users/:userID", h.AdminGetUserHandler)
	app.Get("/admin/users/:userID/posts", h.AdminGetUserPosts)
	app.Get("/admin/dashboard", h.GetAdminDashboard)
	app.Delete("/admin/users/:userID/posts/:postID", h.DeleteAdminUserPostHandler)
	app.Get("/admin/contacts", h.AdminGetAllContactsHandler)
	app.Delete("/admin/contacts/:contactID", h.AdminDeleteContactHandler)
	app.Post("/user/users/:userID/near", h.GetNearUsersHandler)
}
