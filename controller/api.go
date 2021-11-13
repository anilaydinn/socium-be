package controller

import "github.com/anilaydinn/socium-be/service"

type API struct {
	service *service.Service
}

func NewAPI(service *service.Service) API {
	return API{
		service: service,
	}
}
