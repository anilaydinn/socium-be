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
	ID       string `bson:"id"`
	Name     string `bson:"name"`
	Surname  string `bson:"surname"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
	UserType string `bson:"userType"`
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

func convertUserModelToUserEntity(user model.User) UserEntity {
	return UserEntity{
		ID:       user.ID,
		Name:     user.Name,
		Surname:  user.Surname,
		Email:    user.Email,
		Password: user.Password,
		UserType: user.UserType,
	}
}

func convertUserEntityToUserModel(userEntity UserEntity) model.User {
	return model.User{
		ID:       userEntity.ID,
		Name:     userEntity.Name,
		Surname:  userEntity.Surname,
		Email:    userEntity.Email,
		Password: userEntity.Password,
		UserType: userEntity.UserType,
	}
}
