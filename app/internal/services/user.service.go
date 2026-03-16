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
	userRepo      repositories.UserRepository
	workspaceRepo repositories.WorkspaceRepository
}

func NewUserService(userRepo repositories.UserRepository, workspaceRepo repositories.WorkspaceRepository) UserService {
	return &userService{
		userRepo:      userRepo,
		workspaceRepo: workspaceRepo,
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

	var workspace *models.Workspace
	if user.Role == models.RoleHead {
		workspaces, _ := s.workspaceRepo.FindByUserHeadID(user.ID)
		if len(workspaces) > 0 {
			workspace = &workspaces[0]
		}
	}

	response := &models.AuthResponse{
		Token: token,
		User:  s.userToResponse(user, workspace),
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

	var workspace *models.Workspace
	if user.Role == models.RoleHead {
		workspaces, _ := s.workspaceRepo.FindByUserHeadID(user.ID)
		if len(workspaces) > 0 {
			workspace = &workspaces[0]
		}
	} else if user.Role == models.RoleUser || user.Role == models.RoleManager {
		var userHeadID int64
		if user.UserHeadID != nil {
			userHeadID = *user.UserHeadID
		} else {
			userHeadID = user.ID
		}
		workspaces, _ := s.workspaceRepo.FindByUserHeadID(userHeadID)
		if len(workspaces) > 0 {
			workspace = &workspaces[0]
		}
	}

	response := &models.AuthResponse{
		Token: token,
		User:  s.userToResponse(user, workspace),
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

	var workspace *models.Workspace
	var createdWorkspace *models.Workspace

	if creatorRole == models.RoleRoot && req.Role == models.RoleHead {
		workspaceReq := &models.WorkspaceCreateRequest{
			Name:            user.Name + " Workspace",
			UserHeadId:      user.ID,
			TgBotToken:      "",
			IncludesUsersId: []int64{user.ID},
		}

		createdWorkspace = workspaceReq.ToWorkspace()
		if err := s.workspaceRepo.Create(createdWorkspace); err != nil {
			return nil, err
		}
		workspace = createdWorkspace
	}

	if creatorRole == models.RoleHead && req.Role == models.RoleUser {
		workspaces, err := s.workspaceRepo.FindByUserHeadID(creatorID)
		if err != nil {
			return s.userToResponsePtr(user, nil), nil
		}

		if len(workspaces) > 0 {
			workspaceItem := workspaces[0]
			if err := s.workspaceRepo.AddUserToWorkspace(workspaceItem.ID, user.ID); err != nil {
				return s.userToResponsePtr(user, nil), nil
			}
			workspace = &workspaceItem
		}
	}

	return s.userToResponsePtr(user, workspace), nil
}

func (s *userService) GetUser(id int64) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	var workspace *models.Workspace
	switch user.Role {
	case models.RoleHead:
		workspaces, _ := s.workspaceRepo.FindByUserHeadID(user.ID)
		if len(workspaces) > 0 {
			workspace = &workspaces[0]
		}
	case models.RoleUser, models.RoleManager:
		var userHeadID int64
		if user.UserHeadID != nil {
			userHeadID = *user.UserHeadID
		} else {
			userHeadID = user.ID
		}
		workspaces, _ := s.workspaceRepo.FindByUserHeadID(userHeadID)
		if len(workspaces) > 0 {
			workspace = &workspaces[0]
		}
	}

	return s.userToResponsePtr(user, workspace), nil
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

	var workspace *models.Workspace
	if user.Role == models.RoleHead {
		workspaces, _ := s.workspaceRepo.FindByUserHeadID(user.ID)
		if len(workspaces) > 0 {
			workspace = &workspaces[0]
		}
	}

	return s.userToResponsePtr(user, workspace), nil
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
	for i, user := range users {
		var workspace *models.Workspace

		if user.Role == models.RoleHead {
			workspaces, _ := s.workspaceRepo.FindByUserHeadID(user.ID)
			if len(workspaces) > 0 {
				workspace = &workspaces[0]
			}
		} else if user.Role == models.RoleUser || user.Role == models.RoleManager {
			var userHeadID int64
			if user.UserHeadID != nil {
				userHeadID = *user.UserHeadID
			} else {
				userHeadID = user.ID
			}
			workspaces, _ := s.workspaceRepo.FindByUserHeadID(userHeadID)
			if len(workspaces) > 0 {
				workspace = &workspaces[0]
			}
		}

		userResponses[i] = s.userToResponse(&user, workspace)
	}

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

func (s *userService) userToResponse(user *models.User, workspace *models.Workspace) models.UserResponse {
	var safeWorkspace *models.WorkspaceResponse
	if workspace != nil {
		workspaceResp := workspace.ToSafeResponse()
		safeWorkspace = &workspaceResp
	}

	return models.UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Role:       user.Role,
		UserHeadID: user.UserHeadID,
		Workspace:  safeWorkspace,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

func (s *userService) userToResponsePtr(user *models.User, workspace *models.Workspace) *models.UserResponse {
	response := s.userToResponse(user, workspace)
	return &response
}
