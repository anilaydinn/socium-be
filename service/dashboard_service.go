package service

import "github.com/anilaydinn/socium-be/model"

func (service *Service) GetAdminDashboard() (*model.DashboardInformation, error) {
	userCount, err := service.repository.GetUserCount()
	if err != nil {
		return nil, err
	}

	postCount, err := service.repository.GetPostCount()
	if err != nil {
		return nil, err
	}

	commentCount, err := service.repository.GetCommentCount()
	if err != nil {
		return nil, err
	}

	activatedUserCount, err := service.repository.GetActivatedUserCount()
	if err != nil {
		return nil, err
	}

	return &model.DashboardInformation{
		UserCount:          userCount,
		PostCount:          postCount,
		CommentCount:       commentCount,
		ActivatedUserCount: activatedUserCount,
	}, nil
}
