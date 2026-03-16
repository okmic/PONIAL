package repositories

import (
	"ponial/internal/models"

	"gorm.io/gorm"
)

type YandexAIRepository interface {
	Create(yandexAI *models.YandexAI) error
	FindByID(id int64) (*models.YandexAI, error)
	FindByToken(token string) (*models.YandexAI, error)
	FindByFolderID(folderID string) (*models.YandexAI, error)
	Update(yandexAI *models.YandexAI) error
	Delete(id int64) error
	SoftDelete(id int64) error
	List(limit, offset int, filters ...YandexAIFilter) ([]models.YandexAI, int64, error)
	Count(filters ...YandexAIFilter) (int64, error)
}

type YandexAIFilter func(db *gorm.DB) *gorm.DB

type yandexAIRepository struct {
	db *gorm.DB
}

func NewYandexAIRepository(db *gorm.DB) YandexAIRepository {
	return &yandexAIRepository{db: db}
}

func (r *yandexAIRepository) Create(yandexAI *models.YandexAI) error {
	return r.db.Create(yandexAI).Error
}

func (r *yandexAIRepository) FindByID(id int64) (*models.YandexAI, error) {
	var yandexAI models.YandexAI
	err := r.db.First(&yandexAI, id).Error
	if err != nil {
		return nil, err
	}
	return &yandexAI, nil
}

func (r *yandexAIRepository) FindByToken(token string) (*models.YandexAI, error) {
	var yandexAI models.YandexAI
	err := r.db.Where("o_auth_token = ?", token).First(&yandexAI).Error
	if err != nil {
		return nil, err
	}
	return &yandexAI, nil
}

func (r *yandexAIRepository) FindByFolderID(folderID string) (*models.YandexAI, error) {
	var yandexAI models.YandexAI
	err := r.db.Where("folder_id = ?", folderID).First(&yandexAI).Error
	if err != nil {
		return nil, err
	}
	return &yandexAI, nil
}

func (r *yandexAIRepository) Update(yandexAI *models.YandexAI) error {
	return r.db.Save(yandexAI).Error
}

func (r *yandexAIRepository) Delete(id int64) error {
	return r.db.Unscoped().Delete(&models.YandexAI{}, id).Error
}

func (r *yandexAIRepository) SoftDelete(id int64) error {
	return r.db.Delete(&models.YandexAI{}, id).Error
}

func (r *yandexAIRepository) List(limit, offset int, filters ...YandexAIFilter) ([]models.YandexAI, int64, error) {
	var yandexAIs []models.YandexAI
	var total int64

	db := r.db.Model(&models.YandexAI{})

	for _, filter := range filters {
		db = filter(db)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = r.db
	for _, filter := range filters {
		db = filter(db)
	}

	err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&yandexAIs).Error
	if err != nil {
		return nil, 0, err
	}

	return yandexAIs, total, nil
}

func (r *yandexAIRepository) Count(filters ...YandexAIFilter) (int64, error) {
	var count int64

	db := r.db.Model(&models.YandexAI{})
	for _, filter := range filters {
		db = filter(db)
	}

	err := db.Count(&count).Error
	return count, err
}

func WithType(aiType models.YandexAIType) YandexAIFilter {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("type = ?", aiType)
	}
}

func WithTokenLike(token string) YandexAIFilter {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("o_auth_token ILIKE ?", "%"+token+"%")
	}
}

func WithFolderIDLike(folderID string) YandexAIFilter {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("folder_id ILIKE ?", "%"+folderID+"%")
	}
}

func WithPriceGreaterThan(price int64) YandexAIFilter {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("price_per_1000_tokens_in_rub > ?", price)
	}
}

func WithPriceLessThan(price int64) YandexAIFilter {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("price_per_1000_tokens_in_rub < ?", price)
	}
}
