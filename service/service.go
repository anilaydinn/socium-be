package service

import (
	"github.com/anilaydinn/socium-be/email"
	"github.com/anilaydinn/socium-be/errors"
	"github.com/anilaydinn/socium-be/model"
	"github.com/anilaydinn/socium-be/repository"
	"github.com/anilaydinn/socium-be/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repository *repository.Repository) Service {
	return Service{
		repository: repository,
	}
}

func (service *Service) RegisterUser(userDTO model.UserDTO) (*model.User, error) {
	alreadyRegisteredUser, err := service.repository.GetUserByEmail(userDTO.Email)
	if alreadyRegisteredUser != nil {
		return nil, errors.UserAlreadyRegistered
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := model.User{
		ID:          utils.GenerateUUID(8),
		Name:        userDTO.Name,
		Surname:     userDTO.Surname,
		Email:       userDTO.Email,
		Password:    string(hashedPassword),
		UserType:    "user",
		IsActivated: false,
	}

	newUser, err := service.repository.RegisterUser(user)

	if err != nil {
		return nil, err
	}

	err = email.SendMail(newUser.Email, "Complete Registration", "Please click "+os.Getenv("PROD_HOSTNAME")+"/activation/"+user.ID)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (service *Service) LoginUser(userCredentialsDTO model.UserCredentialsDTO) (*model.Token, *fiber.Cookie, error) {
	user, err := service.repository.GetUserByEmail(userCredentialsDTO.Email)

	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		return nil, nil, errors.UserNotFound
	}

	if !user.IsActivated {
		return nil, nil, errors.Unauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userCredentialsDTO.Password)); err != nil {
		return nil, nil, errors.WrongPassword
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, model.CustomClaims{
		UserType: user.UserType,
		StandardClaims: jwt.StandardClaims{
			Issuer: user.ID,
		},
	})

	token, err := claims.SignedString([]byte(""))

	if err != nil {
		return nil, nil, err
	}

	cookie := fiber.Cookie{
		Name:    "user-token",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 24),
	}
	return &model.Token{
		Token: token,
	}, &cookie, nil
}

func (service *Service) Activation(userID string) (*model.User, error) {
	user, err := service.repository.GetUser(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.UserNotFound
	}

	if !user.IsActivated {
		user.IsActivated = true
	} else {
		return nil, errors.UserAlreadyActivated
	}

	return service.repository.UpdateUser(userID, *user)
}

func (service *Service) ForgotPassword(forgotPasswordDTO model.ForgotPasswordDTO) error {
	registeredUser, _ := service.repository.GetUserByEmail(forgotPasswordDTO.Email)
	if registeredUser == nil {
		return errors.UserNotFound
	}

	if !registeredUser.IsActivated {
		return errors.UserNotActivated
	}

	err := email.SendMail(forgotPasswordDTO.Email, "Reset Password", "You can reset your password click "+os.Getenv("REACT_HOSTNAME")+"/reset-password/"+registeredUser.ID+" here.")
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) ResetPassword(userID string, resetPasswordDTO model.ResetPasswordDTO) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(resetPasswordDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user, err := service.repository.GetUser(userID)
	if err != nil {
		return errors.UserNotFound
	}

	user.Password = string(hashedPassword)

	_, err = service.repository.UpdateUser(userID, *user)
	if err != nil {
		return err
	}
	return nil
}

func (service *Service) GetUser(userID string) (*model.User, error) {
	return service.repository.GetUser(userID)
}

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

func (service *Service) GetPosts(userID string) ([]model.Post, error) {
	posts, err := service.repository.GetPosts(userID)
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
