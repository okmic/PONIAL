package services

import (
	"errors"
	"ponial/internal/models"
	"ponial/internal/repositories"
)

type YandexAIService interface {
	CreateYandexAI(req *models.YandexAICreateRequest) (*models.YandexAIResponse, error)
	GetYandexAI(id int64) (*models.YandexAIResponse, error)
	GetYandexAIByToken(token string) (*models.YandexAIResponse, error)
	UpdateYandexAI(id int64, req *models.YandexAIUpdateRequest) (*models.YandexAIResponse, error)
	DeleteYandexAI(id int64) error
	ListYandexAIs(limit, offset int, aiType, token, folderID string, minPrice, maxPrice int64) (*YandexAIListResponse, error)
}

type YandexAIListResponse struct {
	YandexAIs  []models.YandexAIResponse `json:"yandexAIs"`
	Total      int64                     `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"totalPages"`
}

type yandexAIService struct {
	yandexAIRepo repositories.YandexAIRepository
}

func NewYandexAIService(yandexAIRepo repositories.YandexAIRepository) YandexAIService {
	return &yandexAIService{yandexAIRepo: yandexAIRepo}
}

func (s *yandexAIService) CreateYandexAI(req *models.YandexAICreateRequest) (*models.YandexAIResponse, error) {
	existing, err := s.yandexAIRepo.FindByToken(req.OAuthToken)
	if err == nil && existing != nil {
		return nil, errors.New("yandex AI with this token already exists")
	}

	existingByFolder, err := s.yandexAIRepo.FindByFolderID(req.FolderID)
	if err == nil && existingByFolder != nil {
		return nil, errors.New("yandex AI with this folder ID already exists")
	}

	yandexAI := req.ToYandexAI()
	if err := s.yandexAIRepo.Create(yandexAI); err != nil {
		return nil, err
	}

	response := yandexAI.ToResponse()
	return &response, nil
}

func (s *yandexAIService) GetYandexAI(id int64) (*models.YandexAIResponse, error) {
	yandexAI, err := s.yandexAIRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("yandex AI not found")
	}

	response := yandexAI.ToResponse()
	return &response, nil
}

func (s *yandexAIService) GetYandexAIByToken(token string) (*models.YandexAIResponse, error) {
	yandexAI, err := s.yandexAIRepo.FindByToken(token)
	if err != nil {
		return nil, errors.New("yandex AI not found")
	}

	response := yandexAI.ToResponse()
	return &response, nil
}

func (s *yandexAIService) UpdateYandexAI(id int64, req *models.YandexAIUpdateRequest) (*models.YandexAIResponse, error) {
	yandexAI, err := s.yandexAIRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("yandex AI not found")
	}

	if req.OAuthToken != "" {
		existing, err := s.yandexAIRepo.FindByToken(req.OAuthToken)
		if err == nil && existing != nil && existing.ID != id {
			return nil, errors.New("yandex AI with this token already exists")
		}
	}

	if req.FolderID != "" {
		existingByFolder, err := s.yandexAIRepo.FindByFolderID(req.FolderID)
		if err == nil && existingByFolder != nil && existingByFolder.ID != id {
			return nil, errors.New("yandex AI with this folder ID already exists")
		}
	}

	yandexAI.UpdateFromRequest(req)
	if err := s.yandexAIRepo.Update(yandexAI); err != nil {
		return nil, err
	}

	response := yandexAI.ToResponse()
	return &response, nil
}

func (s *yandexAIService) DeleteYandexAI(id int64) error {
	_, err := s.yandexAIRepo.FindByID(id)
	if err != nil {
		return errors.New("yandex AI not found")
	}

	return s.yandexAIRepo.SoftDelete(id)
}

func (s *yandexAIService) ListYandexAIs(limit, offset int, aiType, token, folderID string, minPrice, maxPrice int64) (*YandexAIListResponse, error) {
	var filters []repositories.YandexAIFilter

	if aiType != "" {
		typeEnum := models.YandexAIType(aiType)
		if typeEnum.IsValid() {
			filters = append(filters, repositories.WithType(typeEnum))
		}
	}

	if token != "" {
		filters = append(filters, repositories.WithTokenLike(token))
	}

	if folderID != "" {
		filters = append(filters, repositories.WithFolderIDLike(folderID))
	}

	if minPrice > 0 {
		filters = append(filters, repositories.WithPriceGreaterThan(minPrice))
	}

	if maxPrice > 0 {
		filters = append(filters, repositories.WithPriceLessThan(maxPrice))
	}

	yandexAIs, total, err := s.yandexAIRepo.List(limit, offset, filters...)
	if err != nil {
		return nil, err
	}

	yandexAIResponses := make([]models.YandexAIResponse, len(yandexAIs))
	for i, yandexAI := range yandexAIs {
		yandexAIResponses[i] = yandexAI.ToResponse()
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	page := (offset / limit) + 1
	if page < 1 {
		page = 1
	}

	return &YandexAIListResponse{
		YandexAIs:  yandexAIResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
