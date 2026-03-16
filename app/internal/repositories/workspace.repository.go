package repositories

import (
	"ponial/internal/models"

	"gorm.io/gorm"
)

type WorkspaceRepository interface {
	Create(workspace *models.Workspace) error
	FindByID(id int64) (*models.Workspace, error)
	FindByUserHeadID(userHeadID int64) ([]models.Workspace, error)
	Update(workspace *models.Workspace) error
	Delete(id int64) error
	SoftDelete(id int64) error
	List(limit, offset int, filters ...WorkspaceFilter) ([]models.Workspace, int64, error)
	Count(filters ...WorkspaceFilter) (int64, error)
	AddUserToWorkspace(workspaceID, userID int64) error
	RemoveUserFromWorkspace(workspaceID, userID int64) error
	CheckUserInWorkspace(workspaceID, userID int64) (bool, error)
}

type WorkspaceFilter func(db *gorm.DB) *gorm.DB

type workspaceRepository struct {
	db *gorm.DB
}

func NewWorkspaceRepository(db *gorm.DB) WorkspaceRepository {
	return &workspaceRepository{db: db}
}

func (r *workspaceRepository) Create(workspace *models.Workspace) error {
	return r.db.Create(workspace).Error
}

func (r *workspaceRepository) FindByID(id int64) (*models.Workspace, error) {
	var workspace models.Workspace
	err := r.db.First(&workspace, id).Error
	if err != nil {
		return nil, err
	}
	return &workspace, nil
}

func (r *workspaceRepository) FindByUserHeadID(userHeadID int64) ([]models.Workspace, error) {
	var workspaces []models.Workspace
	err := r.db.Where("user_head_id = ?", userHeadID).Find(&workspaces).Error
	if err != nil {
		return nil, err
	}
	return workspaces, nil
}

func (r *workspaceRepository) Update(workspace *models.Workspace) error {
	return r.db.Save(workspace).Error
}

func (r *workspaceRepository) Delete(id int64) error {
	return r.db.Unscoped().Delete(&models.Workspace{}, id).Error
}

func (r *workspaceRepository) SoftDelete(id int64) error {
	return r.db.Delete(&models.Workspace{}, id).Error
}

func (r *workspaceRepository) List(limit, offset int, filters ...WorkspaceFilter) ([]models.Workspace, int64, error) {
	var workspaces []models.Workspace
	var total int64

	db := r.db.Model(&models.Workspace{})

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

	err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&workspaces).Error
	if err != nil {
		return nil, 0, err
	}

	return workspaces, total, nil
}

func (r *workspaceRepository) Count(filters ...WorkspaceFilter) (int64, error) {
	var count int64

	db := r.db.Model(&models.Workspace{})
	for _, filter := range filters {
		db = filter(db)
	}

	err := db.Count(&count).Error
	return count, err
}

func (r *workspaceRepository) AddUserToWorkspace(workspaceID, userID int64) error {
	userIDJSON := models.Int64Array{userID}
	return r.db.Exec(`
		UPDATE workspaces 
		SET includes_users_id = 
			CASE 
				WHEN includes_users_id IS NULL OR includes_users_id = 'null'::jsonb 
				THEN ?::jsonb
				ELSE (includes_users_id || ?::jsonb)
			END
		WHERE id = ? AND (
			includes_users_id IS NULL 
			OR includes_users_id = 'null'::jsonb 
			OR NOT EXISTS (
				SELECT 1 
				FROM jsonb_array_elements(includes_users_id) elem
				WHERE elem::text::bigint = ?
			)
		)
	`, userIDJSON, userIDJSON, workspaceID, userID).Error
}

func (r *workspaceRepository) RemoveUserFromWorkspace(workspaceID, userID int64) error {
	return r.db.Exec(`
		UPDATE workspaces 
		SET includes_users_id = (
			SELECT COALESCE(
				jsonb_agg(elem),
				'[]'::jsonb
			)
			FROM jsonb_array_elements(includes_users_id) elem
			WHERE elem::text::bigint != ?
		)
		WHERE id = ? AND EXISTS (
			SELECT 1 
			FROM jsonb_array_elements(includes_users_id) elem
			WHERE elem::text::bigint = ?
		)
	`, userID, workspaceID, userID).Error
}

func (r *workspaceRepository) CheckUserInWorkspace(workspaceID, userID int64) (bool, error) {
	var count int64

	err := r.db.Raw(`
		SELECT COUNT(*) FROM workspaces w
		WHERE w.id = ? AND EXISTS (
			SELECT 1 
			FROM jsonb_array_elements(w.includes_users_id) elem
			WHERE elem::text::bigint = ?
		)
	`, workspaceID, userID).Scan(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func WorkspaceWithNameLike(name string) WorkspaceFilter {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name ILIKE ?", "%"+name+"%")
	}
}

func WithUserHeadID(userHeadID int64) WorkspaceFilter {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_head_id = ?", userHeadID)
	}
}

func WithIncludesUserID(userID int64) WorkspaceFilter {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(`
			EXISTS (
				SELECT 1 
				FROM jsonb_array_elements(includes_users_id) elem
				WHERE elem::text::bigint = ?
			)
		`, userID)
	}
}

func WorkspaceActiveOnly() WorkspaceFilter {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("deleted_at IS NULL")
	}
}

func WithID(id int64) WorkspaceFilter {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}
