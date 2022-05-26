package repository

import (
	"context"
	"github.com/anilaydinn/socium-be/errors"
	"github.com/anilaydinn/socium-be/model"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

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
