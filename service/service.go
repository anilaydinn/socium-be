package service

import (
	"github.com/anilaydinn/socium-be/repository"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repository *repository.Repository) Service {
	return Service{
		repository: repository,
	}
}
