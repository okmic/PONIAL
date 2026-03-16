package repositories

import (
	"ponial/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id int64) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id int64) error
	SoftDelete(id int64) error
	List(limit, offset int, filters ...UserFilter) ([]models.User, int64, error)
	Count(filters ...UserFilter) (int64, error)
}

type UserFilter func(db *gorm.DB) *gorm.DB

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id int64) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id int64) error {
	return r.db.Unscoped().Delete(&models.User{}, id).Error
}

func (r *userRepository) SoftDelete(id int64) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *userRepository) List(limit, offset int, filters ...UserFilter) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	db := r.db.Model(&models.User{})

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

	err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Count(filters ...UserFilter) (int64, error) {
	var count int64

	db := r.db.Model(&models.User{})
	for _, filter := range filters {
		db = filter(db)
	}

	err := db.Count(&count).Error
	return count, err
}

func WithRole(role models.Role) UserFilter {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("role = ?", role)
	}
}

func WithNameLike(name string) UserFilter {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name ILIKE ?", "%"+name+"%")
	}
}

func ActiveOnly() UserFilter {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("deleted_at IS NULL")
	}
}
