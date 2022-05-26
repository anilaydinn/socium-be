package test

import (
	"context"
	"time"

	"github.com/anilaydinn/socium-be/repository"
)

func GetCleanTestRepository() *repository.Repository {
	repository := repository.NewRepository("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	testDB := repository.MongoClient.Database("socium")
	testDB.Drop(ctx)

	return repository
}
