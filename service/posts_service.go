package service

import (
	"github.com/anilaydinn/socium-be/errors"
	"github.com/anilaydinn/socium-be/model"
	"github.com/anilaydinn/socium-be/utils"
	"time"
)

func (service *Service) CreatePost(postDTO model.PostDTO) (*model.Post, error) {
	post := model.Post{
		ID:          utils.GenerateUUID(8),
		UserID:      postDTO.UserID,
		Description: postDTO.Description,
		Image:       postDTO.Image,
		IsPrivate:   postDTO.IsPrivate,
		CreatedAt:   time.Now().UTC().Round(time.Second),
		UpdatedAt:   time.Now().UTC().Round(time.Second),
	}

	return service.repository.CreatePost(post)
}

func (service *Service) GetPosts(userID string, isHomePage bool, friendIDList []string) ([]model.Post, error) {
	if isHomePage {
		friendIDList = append(friendIDList, userID)
	}

	posts, err := service.repository.GetPosts(userID, isHomePage, friendIDList)
	if err != nil {
		return nil, err
	}

	for i, post := range posts {
		postUser, err := service.repository.GetUser(post.UserID)
		if err != nil {
			return nil, err
		}
		posts[i].User = postUser

		comments, err := service.repository.GetCommentsByIDList(post.CommentIDs)
		if err != nil {
			return nil, err
		}

		for j, comment := range comments {
			commentUser, err := service.repository.GetUser(comment.UserID)
			if err != nil {
				return nil, err
			}
			comments[j].User = commentUser
		}
		posts[i].Comments = comments
	}

	return posts, nil
}

func (service *Service) GetPost(postID string) (*model.Post, error) {
	post, err := service.repository.GetPost(postID)
	if err != nil {
		return nil, err
	}

	comments, err := service.repository.GetCommentsByIDList(post.CommentIDs)
	if err != nil {
		return nil, err
	}

	post.Comments = comments

	return post, nil
}

func (service *Service) LikePost(postID string, likePostDTO model.LikePostDTO) (*model.Post, error) {
	post, err := service.repository.GetPost(postID)
	if err != nil {
		return nil, errors.PostNotFound
	}

	if utils.Contains(post.WhoLikesUserIDs, likePostDTO.UserID) {
		post.WhoLikesUserIDs = utils.RemoveElement(post.WhoLikesUserIDs, likePostDTO.UserID)
	} else {
		post.WhoLikesUserIDs = append(post.WhoLikesUserIDs, likePostDTO.UserID)
	}

	updatedPost, err := service.repository.UpdatePost(postID, *post)
	if err != nil {
		return nil, err
	}

	return updatedPost, nil
}

func (service *Service) AddPostComment(postID string, commentDTO model.CommentDTO) (*model.Post, error) {
	comment := model.Comment{
		ID:        utils.GenerateUUID(8),
		UserID:    commentDTO.UserID,
		PostID:    postID,
		User:      nil,
		Content:   commentDTO.Content,
		CreatedAt: time.Now().UTC().Round(time.Second),
		UpdatedAt: time.Now().UTC().Round(time.Second),
	}

	newComment, err := service.repository.AddComment(comment)
	if err != nil {
		return nil, err
	}

	post, err := service.repository.GetPost(postID)
	if err != nil {
		return nil, err
	}

	post.CommentIDs = append(post.CommentIDs, newComment.ID)

	_, err = service.repository.UpdatePost(postID, *post)
	if err != nil {
		return nil, err
	}

	return service.GetPost(postID)
}

func (service *Service) DeleteAdminUserPost(postID, userID string) error {
	return service.repository.DeleteUserPost(postID, userID)
}
