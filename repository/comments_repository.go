package repository

import (
	"context"
	"github.com/anilaydinn/socium-be/errors"
	"github.com/anilaydinn/socium-be/model"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

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
