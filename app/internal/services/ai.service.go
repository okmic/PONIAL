package services

import (
	"ponial/internal/models"
	"ponial/internal/repositories"
)

type AIService interface {
	getAIMsg(req *models.UserCreateRequest, creatorRole models.Role, creatorID int64) (*models.UserResponse, error)
}

func getAIMsg(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}
