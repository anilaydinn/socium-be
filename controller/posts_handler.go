package controller

import (
	"github.com/anilaydinn/socium-be/errors"
	"github.com/anilaydinn/socium-be/model"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreatePostHandler(c *fiber.Ctx) error {
	postDTO := model.PostDTO{}
	err := c.BodyParser(&postDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	post, err := h.service.CreatePost(postDTO)

	switch err {
	case nil:
		c.Status(fiber.StatusCreated)
		c.JSON(post)
	case errors.PostNotFound:
		c.Status(fiber.StatusNotFound)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (h *Handler) GetPostsHandler(c *fiber.Ctx) error {
	q := new(model.GetPostsQuery)

	if err := c.QueryParser(q); err != nil {
		return err
	}

	var isHomepage bool
	if q.Homepage == "true" {
		isHomepage = true
	} else {
		isHomepage = false
	}

	posts, err := h.service.GetPosts(q.UserID, isHomepage, q.FriendIDList)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(posts)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (h *Handler) LikePostHandler(c *fiber.Ctx) error {
	postID := c.Params("postID")
	likePostDTO := model.LikePostDTO{}
	err := c.BodyParser(&likePostDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	post, err := h.service.LikePost(postID, likePostDTO)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(post)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (h *Handler) AddPostCommentHandler(c *fiber.Ctx) error {
	postID := c.Params("postID")
	commentDTO := model.CommentDTO{}
	err := c.BodyParser(&commentDTO)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	post, err := h.service.AddPostComment(postID, commentDTO)

	switch err {
	case nil:
		c.Status(fiber.StatusCreated)
		c.JSON(post)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (h *Handler) DeleteAdminUserPostHandler(c *fiber.Ctx) error {
	postID := c.Params("postID")
	userID := c.Params("userID")
	if len(postID) == 0 || len(userID) == 0 {
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	err := h.service.DeleteAdminUserPost(postID, userID)

	switch err {
	case nil:
		c.Status(fiber.StatusNoContent)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}

func (h *Handler) GetWholikesPostHandler(c *fiber.Ctx) error {
	postID := c.Params("postID")
	if len(postID) == 0 {
		c.Status(fiber.StatusBadRequest)
		return nil
	}
	q := new(model.WhoLikesQuery)

	if err := c.QueryParser(q); err != nil {
		return err
	}

	users, err := h.service.GetWhoLikesPost(postID, q.WhoLikesUserIDs)

	switch err {
	case nil:
		c.Status(fiber.StatusOK)
		c.JSON(users)
	case errors.WhoLikesArrayNotEqual:
		c.Status(fiber.StatusBadRequest)
	default:
		c.Status(fiber.StatusInternalServerError)
	}
	return nil
}
