package repository

import (
	"context"
	"github.com/anilaydinn/socium-be/errors"
	"github.com/anilaydinn/socium-be/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

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
