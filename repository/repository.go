package repository

import (
	"context"
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
	ID          string `bson:"id"`
	Name        string `bson:"name"`
	Surname     string `bson:"surname"`
	Email       string `bson:"email"`
	Password    string `bson:"password"`
	UserType    string `bson:"userType"`
	IsActivated bool   `bson:"isActivated"`
}

type PostEntity struct {
	ID              string    `bson:"id"`
	UserID          string    `bson:"userId"`
	Description     string    `bson:"description"`
	Image           string    `bson:"image"`
	IsPrivate       bool      `bson:"isPrivate"`
	WhoLikesUserIDs []string  `bson:"whoLikesUserIds"`
	CreatedAt       time.Time `bson:"createdAt"`
	UpdatedAt       time.Time `bson:"updatedAt"`
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

func (repository *Repository) GetPosts(userID string) ([]model.Post, error) {
	collection := repository.MongoClient.Database("socium").Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	options := options.Find()
	options.SetSort(bson.M{"createdAt": -1})

	var filter bson.M
	if len(userID) == 0 {
		filter = bson.M{}
	} else {
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

func convertUserModelToUserEntity(user model.User) UserEntity {
	return UserEntity{
		ID:          user.ID,
		Name:        user.Name,
		Surname:     user.Surname,
		Email:       user.Email,
		Password:    user.Password,
		UserType:    user.UserType,
		IsActivated: user.IsActivated,
	}
}

func convertUserEntityToUserModel(userEntity UserEntity) model.User {
	return model.User{
		ID:          userEntity.ID,
		Name:        userEntity.Name,
		Surname:     userEntity.Surname,
		Email:       userEntity.Email,
		Password:    userEntity.Password,
		UserType:    userEntity.UserType,
		IsActivated: userEntity.IsActivated,
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
		CreatedAt:       postEntity.CreatedAt,
		UpdatedAt:       postEntity.UpdatedAt,
	}
}
