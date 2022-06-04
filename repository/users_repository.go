package repository

import (
	"context"
	"github.com/anilaydinn/socium-be/errors"
	"github.com/anilaydinn/socium-be/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

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

func (repository *Repository) GetAllUsers(page, size int, filterArr []string) ([]model.User, int, error) {
	collection := repository.MongoClient.Database("socium").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	options := options.Find()
	if size != 0 {
		options.SetSkip(int64(page * size))
		options.SetLimit(int64(size))
	}

	var filter bson.D
	if len(filterArr) > 1 {
		filter = bson.D{{"name", primitive.Regex{Pattern: filterArr[0], Options: "i"}}, {"surname", primitive.Regex{Pattern: filterArr[1], Options: "i"}}}
	} else if len(filterArr) == 1 {
		filter = bson.D{{"name", primitive.Regex{Pattern: filterArr[0], Options: "i"}}}
	}

	cur, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, 0, err
	}

	var users []model.User
	for cur.Next(ctx) {
		userEntity := UserEntity{}
		err := cur.Decode(&userEntity)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, convertUserEntityToUserModel(userEntity))
	}

	totalElements, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, int(totalElements), nil
}

func (repository *Repository) GetUserCount() (int, error) {
	collection := repository.MongoClient.Database("socium").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userCount, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return int(userCount), nil
}

func (repository *Repository) GetActivatedUserCount() (int, error) {
	collection := repository.MongoClient.Database("socium").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"isActivated": true}

	activatedUserCount, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(activatedUserCount), nil
}
