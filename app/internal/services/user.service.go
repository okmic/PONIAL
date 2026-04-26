package services

import (
	"errors"
	"ponial/internal/models"
	"ponial/internal/repositories"
	"ponial/pkg/jwt"

	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(req *models.UserCreateRequest, creatorRole models.Role, creatorID int64) (*models.UserResponse, error)
	GetUser(id int64) (*models.UserResponse, error)
	UpdateUser(id int64, req *models.UserUpdateRequest) (*models.UserResponse, error)
	DeleteUser(id int64) error
	ListUsers(limit, offset int, role string, name string, currentUserRole models.Role, currentUserID int64) (*UserListResponse, error)
	Signup(req *models.UserSignupRequest) (*models.AuthResponse, error)
	Signin(req *models.UserSigninRequest) (*models.AuthResponse, error)
}

type UserListResponse struct {
	Users      []models.UserResponse `json:"users"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalPages int                   `json:"total_pages"`
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Signup(req *models.UserSignupRequest) (*models.AuthResponse, error) {
	existing, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existing != nil {
		return nil, errors.New("user with this email already exists")
	}
	user, err := req.ToUser()
	if err != nil {
		return nil, err
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	token, err := jwt.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	response := &models.AuthResponse{
		Token: token,
		User:  s.userToResponse(user),
	}

	return response, nil
}

func (s *userService) Signin(req *models.UserSigninRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid email or password")
	}
	token, err := jwt.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	response := &models.AuthResponse{
		Token: token,
		User:  s.userToResponse(user),
	}
	return response, nil
}

func (s *userService) CreateUser(req *models.UserCreateRequest, creatorRole models.Role, creatorID int64) (*models.UserResponse, error) {
	if creatorRole != models.RoleRoot && creatorRole != models.RoleHead {
		return nil, errors.New("only root and head can create users")
	}

	if creatorRole == models.RoleRoot && req.Role != models.RoleHead {
		return nil, errors.New("root can only create head users")
	}

	if creatorRole == models.RoleHead && req.Role == models.RoleRoot {
		return nil, errors.New("head cannot create root users")
	}

	if creatorRole == models.RoleHead {
		req.UserHeadID = &creatorID
	}

	existing, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existing != nil {
		return nil, errors.New("user with this email already exists")
	}

	user, err := req.ToUser()
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	if creatorRole == models.RoleHead && req.Role == models.RoleUser {
		if err != nil {
			return s.userToResponsePtr(user), nil
		}
	}

	return s.userToResponsePtr(user), nil
}

func (s *userService) GetUser(id int64) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.userToResponsePtr(user), nil
}

func (s *userService) UpdateUser(id int64, req *models.UserUpdateRequest) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if err := user.UpdateFromRequest(req); err != nil {
		return nil, err
	}
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return s.userToResponsePtr(user), nil
}

func (s *userService) DeleteUser(id int64) error {
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}
	return s.userRepo.SoftDelete(id)
}

func (s *userService) ListUsers(limit, offset int, role, name string, currentUserRole models.Role, currentUserID int64) (*UserListResponse, error) {
	var filters []repositories.UserFilter

	if currentUserRole == models.RoleRoot {
	} else if currentUserRole == models.RoleHead {
		filters = append(filters, func(db *gorm.DB) *gorm.DB {
			return db.Where("user_head_id = ? OR id = ?", currentUserID, currentUserID)
		})
	} else {
		filters = append(filters, func(db *gorm.DB) *gorm.DB {
			return db.Where("id = ?", currentUserID)
		})
	}

	if role != "" {
		roleEnum := models.Role(role)
		if roleEnum.IsValid() {
			filters = append(filters, repositories.WithRole(roleEnum))
		}
	}
	if name != "" {
		filters = append(filters, repositories.WithNameLike(name))
	}
	filters = append(filters, repositories.ActiveOnly())
	users, total, err := s.userRepo.List(limit, offset, filters...)
	if err != nil {
		return nil, err
	}

	userResponses := make([]models.UserResponse, len(users))

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}
	page := (offset / limit) + 1
	return &UserListResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *userService) userToResponse(user *models.User) models.UserResponse {
	return models.UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Role:       user.Role,
		UserHeadID: user.UserHeadID,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

func (s *userService) userToResponsePtr(user *models.User) *models.UserResponse {
	response := s.userToResponse(user)
	return &response
}
