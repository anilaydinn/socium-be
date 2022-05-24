package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"

	"github.com/anilaydinn/socium-be/errors"
	"github.com/anilaydinn/socium-be/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	MongoClient *mongo.Client
}

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

func NewRepository(uri string) *Repository {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	defer cancel()
	client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	return &Repository{client}
}

func (repository *Repository) RegisterUser(user model.User) (*model.User, error) {
	collection := repository.MongoClient.Database("socium").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userEntity := convertUserModelToUserEntity(user)

	_, err := collection.InsertOne(ctx, userEntity)

	if err != nil {
		return nil, err
	}

	return repository.GetUser(userEntity.ID)
}

func (repository *Repository) GetUser(userID string) (*model.User, error) {
	collection := repository.MongoClient.Database("socium").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": userID}

	cur := collection.FindOne(ctx, filter)

	if cur.Err() != nil {
		return nil, cur.Err()
	}

	if cur == nil {
		return nil, errors.UserNotFound
	}

	userEntity := UserEntity{}
	err := cur.Decode(&userEntity)

	if err != nil {
		return nil, err
	}

	user := convertUserEntityToUserModel(userEntity)

	return &user, nil
}

func (repository *Repository) GetUserByEmail(email string) (*model.User, error) {
	collection := repository.MongoClient.Database("socium").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"email": email}

	cur := collection.FindOne(ctx, filter)

	if cur.Err() != nil {
		return nil, errors.UserNotFound
	}

	if cur == nil {
		return nil, errors.UserNotFound
	}

	userEntity := UserEntity{}
	err := cur.Decode(&userEntity)

	if err != nil {
		return nil, err
	}

	user := convertUserEntityToUserModel(userEntity)

	return &user, nil
}

func (repository *Repository) UpdateUser(userID string, user model.User) (*model.User, error) {
	collection := repository.MongoClient.Database("socium").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": userID}

	userEntity := convertUserModelToUserEntity(user)

	cur := collection.FindOneAndReplace(ctx, filter, userEntity)

	if cur.Err() != nil {
		return nil, cur.Err()
	}

	if cur == nil {
		return nil, errors.UserNotFound
	}

	return repository.GetUser(userID)

}

func (repository *Repository) CreatePost(post model.Post) (*model.Post, error) {
	collection := repository.MongoClient.Database("socium").Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	postEntity := convertPostModelToPostEntity(post)

	_, err := collection.InsertOne(ctx, postEntity)

	if err != nil {
		return nil, err
	}

	return repository.GetPost(postEntity.ID)
}

func (repository *Repository) GetPost(postID string) (*model.Post, error) {
	collection := repository.MongoClient.Database("socium").Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": postID}

	cur := collection.FindOne(ctx, filter)

	if cur.Err() != nil {
		return nil, cur.Err()
	}

	if cur == nil {
		return nil, errors.PostNotFound
	}

	postEntity := PostEntity{}
	err := cur.Decode(&postEntity)

	if err != nil {
		return nil, err
	}

	post := convertPostEntityToPostModel(postEntity)

	return &post, nil
}

func (repository *Repository) GetPosts(userID string, isHomePage bool, friendIDs []string) ([]model.Post, error) {
	collection := repository.MongoClient.Database("socium").Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	options := options.Find()
	options.SetSort(bson.M{"createdAt": -1})

	var filter bson.M
	if len(friendIDs) > 0 && isHomePage {
		filter = bson.M{"userId": bson.M{"$in": friendIDs}}
	} else if !isHomePage {
		filter = bson.M{"userId": userID}
	}

	cur, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}

	var posts []model.Post
	for cur.Next(ctx) {
		postEntity := PostEntity{}
		err := cur.Decode(&postEntity)
		if err != nil {
			return nil, err
		}
		posts = append(posts, convertPostEntityToPostModel(postEntity))
	}

	return posts, nil
}

func (repository *Repository) UpdatePost(postID string, post model.Post) (*model.Post, error) {
	collection := repository.MongoClient.Database("socium").Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": postID}

	postEntity := convertPostModelToPostEntity(post)

	cur := collection.FindOneAndReplace(ctx, filter, postEntity)

	if cur.Err() != nil {
		return nil, cur.Err()
	}

	if cur == nil {
		return nil, errors.PostNotFound
	}

	return repository.GetPost(postID)
}

func (repository *Repository) AddComment(comment model.Comment) (*model.Comment, error) {
	collection := repository.MongoClient.Database("socium").Collection("comments")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	commentEntity := convertCommentModelToCommentEntity(comment)

	_, err := collection.InsertOne(ctx, commentEntity)

	if err != nil {
		return nil, err
	}

	return repository.GetComment(commentEntity.ID)
}

func (repository *Repository) GetComment(commentID string) (*model.Comment, error) {
	collection := repository.MongoClient.Database("socium").Collection("comments")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": commentID}

	cur := collection.FindOne(ctx, filter)

	if cur.Err() != nil {
		return nil, cur.Err()
	}

	if cur == nil {
		return nil, errors.PostNotFound
	}

	commentEntity := CommentEntity{}
	err := cur.Decode(&commentEntity)

	if err != nil {
		return nil, err
	}

	comment := convertCommentEntityToCommentModel(commentEntity)

	return &comment, nil
}

func (repository *Repository) GetCommentsByIDList(commentIDs []string) ([]model.Comment, error) {
	collection := repository.MongoClient.Database("socium").Collection("comments")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var filter bson.M
	if len(commentIDs) == 0 {
		filter = bson.M{"id": bson.M{"$in": []string{}}}
	} else {
		filter = bson.M{"id": bson.M{"$in": commentIDs}}
	}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var comments []model.Comment
	for cur.Next(ctx) {
		commentEntity := CommentEntity{}
		err := cur.Decode(&commentEntity)
		if err != nil {
			return nil, err
		}
		comments = append(comments, convertCommentEntityToCommentModel(commentEntity))
	}
	return comments, nil
}

func (repository *Repository) GetUsersByIDList(userIDs []string) ([]model.User, error) {
	collection := repository.MongoClient.Database("socium").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var filter bson.M
	if len(userIDs) == 0 {
		filter = bson.M{"id": bson.M{"$in": []string{}}}
	} else {
		filter = bson.M{"id": bson.M{"$in": userIDs}}
	}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var users []model.User
	for cur.Next(ctx) {
		userEntity := UserEntity{}
		err := cur.Decode(&userEntity)
		if err != nil {
			return nil, err
		}
		users = append(users, convertUserEntityToUserModel(userEntity))
	}
	return users, nil
}

func (repository *Repository) CreateContact(contact model.Contact) (*model.Contact, error) {
	collection := repository.MongoClient.Database("socium").Collection("contacts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	contactEntity := convertContactModelToContactEntity(contact)

	_, err := collection.InsertOne(ctx, contactEntity)

	if err != nil {
		return nil, err
	}

	return repository.GetContact(contactEntity.ID)
}

func (repository *Repository) GetContact(contactID string) (*model.Contact, error) {
	collection := repository.MongoClient.Database("socium").Collection("contacts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": contactID}

	cur := collection.FindOne(ctx, filter)

	if cur.Err() != nil {
		return nil, cur.Err()
	}

	if cur == nil {
		return nil, errors.ContactNotFound
	}

	contactEntity := ContactEntity{}
	err := cur.Decode(&contactEntity)

	if err != nil {
		return nil, err
	}

	contact := convertContactEntityToContactModel(contactEntity)

	return &contact, nil
}

func (repository *Repository) GetUsersWithFilter(filterArr []string) ([]model.User, error) {
	collection := repository.MongoClient.Database("socium").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var filter bson.D
	if len(filterArr) > 1 {
		filter = bson.D{{"name", primitive.Regex{Pattern: filterArr[0], Options: "i"}}, {"surname", primitive.Regex{Pattern: filterArr[1], Options: "i"}}}
	} else if len(filterArr) == 1 {
		filter = bson.D{{"name", primitive.Regex{Pattern: filterArr[0], Options: "i"}}}
	}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var users []model.User
	for cur.Next(ctx) {
		userEntity := UserEntity{}
		err := cur.Decode(&userEntity)
		if err != nil {
			return nil, err
		}
		users = append(users, convertUserEntityToUserModel(userEntity))
	}

	return users, nil
}

func convertUserModelToUserEntity(user model.User) UserEntity {
	return UserEntity{
		ID:                   user.ID,
		Name:                 user.Name,
		Surname:              user.Surname,
		Email:                user.Email,
		BirthDate:            user.BirthDate,
		Description:          user.Description,
		ProfileImage:         user.ProfileImage,
		FriendRequestUserIDs: user.FriendRequestUserIDs,
		FriendIDs:            user.FriendIDs,
		Password:             user.Password,
		UserType:             user.UserType,
		IsActivated:          user.IsActivated,
		CreatedAt:            user.CreatedAt,
		UpdatedAt:            user.UpdatedAt,
		Latitude:             user.Latitude,
		Longitude:            user.Longitude,
	}
}

func convertUserEntityToUserModel(userEntity UserEntity) model.User {
	return model.User{
		ID:                   userEntity.ID,
		Name:                 userEntity.Name,
		Surname:              userEntity.Surname,
		Email:                userEntity.Email,
		BirthDate:            userEntity.BirthDate,
		Description:          userEntity.Description,
		ProfileImage:         userEntity.ProfileImage,
		FriendRequestUserIDs: userEntity.FriendRequestUserIDs,
		FriendIDs:            userEntity.FriendIDs,
		Password:             userEntity.Password,
		UserType:             userEntity.UserType,
		IsActivated:          userEntity.IsActivated,
		CreatedAt:            userEntity.CreatedAt,
		UpdatedAt:            userEntity.UpdatedAt,
		Latitude:             userEntity.Latitude,
		Longitude:            userEntity.Longitude,
	}
}

func convertPostModelToPostEntity(post model.Post) PostEntity {
	return PostEntity{
		ID:              post.ID,
		UserID:          post.UserID,
		Description:     post.Description,
		Image:           post.Image,
		IsPrivate:       post.IsPrivate,
		WhoLikesUserIDs: post.WhoLikesUserIDs,
		CommentIDs:      post.CommentIDs,
		CreatedAt:       post.CreatedAt,
		UpdatedAt:       post.UpdatedAt,
	}
}

func convertPostEntityToPostModel(postEntity PostEntity) model.Post {
	return model.Post{
		ID:              postEntity.ID,
		UserID:          postEntity.UserID,
		Description:     postEntity.Description,
		Image:           postEntity.Image,
		IsPrivate:       postEntity.IsPrivate,
		WhoLikesUserIDs: postEntity.WhoLikesUserIDs,
		CommentIDs:      postEntity.CommentIDs,
		CreatedAt:       postEntity.CreatedAt,
		UpdatedAt:       postEntity.UpdatedAt,
	}
}

func convertCommentModelToCommentEntity(comment model.Comment) CommentEntity {
	return CommentEntity{
		ID:        comment.ID,
		UserID:    comment.UserID,
		PostID:    comment.PostID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
}

func convertCommentEntityToCommentModel(commentEntity CommentEntity) model.Comment {
	return model.Comment{
		ID:        commentEntity.ID,
		UserID:    commentEntity.UserID,
		PostID:    commentEntity.PostID,
		Content:   commentEntity.Content,
		CreatedAt: commentEntity.CreatedAt,
		UpdatedAt: commentEntity.UpdatedAt,
	}
}

func convertContactModelToContactEntity(contact model.Contact) ContactEntity {
	return ContactEntity{
		ID:      contact.ID,
		Name:    contact.Name,
		Surname: contact.Surname,
		Email:   contact.Email,
		Message: contact.Message,
	}
}

func convertContactEntityToContactModel(contactEntity ContactEntity) model.Contact {
	return model.Contact{
		ID:      contactEntity.ID,
		Name:    contactEntity.Name,
		Surname: contactEntity.Surname,
		Email:   contactEntity.Email,
		Message: contactEntity.Message,
	}
}
