package services

import (
	"errors"
	"ponial/internal/models"
	"ponial/internal/repositories"
)

type WorkspaceService interface {
	CreateWorkspace(req *models.WorkspaceCreateRequest) (*models.WorkspaceResponse, error)
	GetWorkspace(id int64) (*models.WorkspaceResponse, error)
	UpdateWorkspace(id int64, req *models.WorkspaceUpdateRequest, currentUserID int64, currentUserRole models.Role) (*models.WorkspaceResponse, error)
	DeleteWorkspace(id int64) error
	ListWorkspaces(limit, offset int, name string, userHeadID int64) (*WorkspaceListResponse, error)
	GetWorkspacesByUserHead(userHeadID int64) ([]models.WorkspaceResponse, error)
	GetWorkspacesByUser(userID int64) ([]models.WorkspaceResponse, error)
	AddUserToWorkspace(workspaceID, userID int64) error
	RemoveUserFromWorkspace(workspaceID, userID int64) error
	CheckUserInWorkspace(workspaceID, userID int64) (bool, error)
	GetWorkspaceSafe(id int64) (*models.WorkspaceResponse, error)
	GetMyWorkspace(currentUserID int64, currentUserRole models.Role) (*models.WorkspaceResponse, error)
	UpdateMyWorkspace(req *models.WorkspaceUpdateRequest, currentUserID int64, currentUserRole models.Role) (*models.WorkspaceResponse, error)
}

type WorkspaceListResponse struct {
	Workspaces []models.WorkspaceResponse `json:"workspaces"`
	Total      int64                      `json:"total"`
	Page       int                        `json:"page"`
	Limit      int                        `json:"limit"`
	TotalPages int                        `json:"total_pages"`
}

type workspaceService struct {
	workspaceRepo repositories.WorkspaceRepository
	userRepo      repositories.UserRepository
}

func NewWorkspaceService(workspaceRepo repositories.WorkspaceRepository, userRepo repositories.UserRepository) WorkspaceService {
	return &workspaceService{
		workspaceRepo: workspaceRepo,
		userRepo:      userRepo,
	}
}

func (s *workspaceService) CreateWorkspace(req *models.WorkspaceCreateRequest) (*models.WorkspaceResponse, error) {
	workspace := req.ToWorkspace()

	if err := s.workspaceRepo.Create(workspace); err != nil {
		return nil, err
	}

	response := workspace.ToResponse()
	return &response, nil
}

func (s *workspaceService) GetWorkspace(id int64) (*models.WorkspaceResponse, error) {
	workspace, err := s.workspaceRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("workspace not found")
	}

	response := workspace.ToResponse()
	return &response, nil
}

func (s *workspaceService) GetWorkspaceSafe(id int64) (*models.WorkspaceResponse, error) {
	workspace, err := s.workspaceRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("workspace not found")
	}

	response := workspace.ToSafeResponse()
	return &response, nil
}

func (s *workspaceService) UpdateWorkspace(id int64, req *models.WorkspaceUpdateRequest, currentUserID int64, currentUserRole models.Role) (*models.WorkspaceResponse, error) {
	workspace, err := s.workspaceRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("workspace not found")
	}

	if currentUserRole != models.RoleRoot {
		if currentUserRole == models.RoleHead && workspace.UserHeadId != currentUserID {
			return nil, errors.New("you can only update your own workspace")
		}
		if currentUserRole == models.RoleManager || currentUserRole == models.RoleUser {
			inWorkspace, err := s.workspaceRepo.CheckUserInWorkspace(workspace.ID, currentUserID)
			if err != nil {
				return nil, err
			}
			if !inWorkspace {
				return nil, errors.New("you don't have access to this workspace")
			}
		}
	}

	req.UpdateWorkspace(workspace)

	if err := s.workspaceRepo.Update(workspace); err != nil {
		return nil, err
	}

	response := workspace.ToResponse()
	return &response, nil
}

func (s *workspaceService) DeleteWorkspace(id int64) error {
	_, err := s.workspaceRepo.FindByID(id)
	if err != nil {
		return errors.New("workspace not found")
	}

	return s.workspaceRepo.SoftDelete(id)
}

func (s *workspaceService) ListWorkspaces(limit, offset int, name string, userHeadID int64) (*WorkspaceListResponse, error) {
	var filters []repositories.WorkspaceFilter

	if name != "" {
		filters = append(filters, repositories.WorkspaceWithNameLike(name))
	}

	if userHeadID > 0 {
		filters = append(filters, repositories.WithUserHeadID(userHeadID))
	}

	filters = append(filters, repositories.WorkspaceActiveOnly())

	workspaces, total, err := s.workspaceRepo.List(limit, offset, filters...)
	if err != nil {
		return nil, err
	}

	workspaceResponses := make([]models.WorkspaceResponse, len(workspaces))
	for i, workspace := range workspaces {
		workspaceResponses[i] = workspace.ToSafeResponse()
	}

	totalPages := 0
	if limit > 0 {
		totalPages = int(total) / limit
		if int(total)%limit > 0 {
			totalPages++
		}
	}

	page := 1
	if limit > 0 {
		page = (offset / limit) + 1
	}

	return &WorkspaceListResponse{
		Workspaces: workspaceResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *workspaceService) GetWorkspacesByUserHead(userHeadID int64) ([]models.WorkspaceResponse, error) {
	workspaces, err := s.workspaceRepo.FindByUserHeadID(userHeadID)
	if err != nil {
		return nil, err
	}

	responses := make([]models.WorkspaceResponse, len(workspaces))
	for i, workspace := range workspaces {
		responses[i] = workspace.ToSafeResponse()
	}

	return responses, nil
}

func (s *workspaceService) GetWorkspacesByUser(userID int64) ([]models.WorkspaceResponse, error) {

	workspaces, err := s.workspaceRepo.FindByUserHeadID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]models.WorkspaceResponse, len(workspaces))
	for i, workspace := range workspaces {
		responses[i] = workspace.ToSafeResponse()
	}

	return responses, nil
}

func (s *workspaceService) AddUserToWorkspace(workspaceID, userID int64) error {
	exists, err := s.workspaceRepo.CheckUserInWorkspace(workspaceID, userID)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("user already in workspace")
	}

	return s.workspaceRepo.AddUserToWorkspace(workspaceID, userID)
}

func (s *workspaceService) RemoveUserFromWorkspace(workspaceID, userID int64) error {
	exists, err := s.workspaceRepo.CheckUserInWorkspace(workspaceID, userID)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("user not found in workspace")
	}

	return s.workspaceRepo.RemoveUserFromWorkspace(workspaceID, userID)
}

func (s *workspaceService) CheckUserInWorkspace(workspaceID, userID int64) (bool, error) {
	return s.workspaceRepo.CheckUserInWorkspace(workspaceID, userID)
}

func (s *workspaceService) GetMyWorkspace(currentUserID int64, currentUserRole models.Role) (*models.WorkspaceResponse, error) {
	if currentUserRole == models.RoleRoot {
		return nil, errors.New("root user does not belong to a workspace")
	}
	user, userErr := s.userRepo.FindByID(currentUserID)
	if userErr != nil {
		return nil, userErr
	}
	var workspaceHeadID int64
	if currentUserRole == models.RoleHead {
		workspaceHeadID = user.ID
	} else {
		if user.UserHeadID == nil {
			return nil, errors.New("user head not found")
		}
		workspaceHeadID = *user.UserHeadID
	}
	workspaces, err := s.workspaceRepo.FindByUserHeadID(workspaceHeadID)
	if err != nil {
		return nil, err
	}

	if len(workspaces) == 0 {
		return nil, errors.New("workspace not found")
	}

	response := workspaces[0].ToResponse()
	return &response, nil
}

func (s *workspaceService) UpdateMyWorkspace(req *models.WorkspaceUpdateRequest, currentUserID int64, currentUserRole models.Role) (*models.WorkspaceResponse, error) {
	var workspace *models.Workspace

	if currentUserRole == models.RoleHead {
		workspaces, err := s.workspaceRepo.FindByUserHeadID(currentUserID)
		if err != nil {
			return nil, err
		}
		if len(workspaces) == 0 {
			return nil, errors.New("workspace not found")
		}
		workspace = &workspaces[0]
	} else if currentUserRole == models.RoleManager {
		user, errUser := s.userRepo.FindByID(currentUserID)
		if errUser != nil {
			return nil, errUser
		}
		workspaces, err := s.workspaceRepo.FindByUserHeadID(*user.UserHeadID)
		if err != nil {
			return nil, err
		}
		if len(workspaces) == 0 {
			return nil, errors.New("workspace not found")
		}
		workspace = &workspaces[0]
	} else {
		return nil, errors.New("only head, manager and user can update their workspace")
	}

	req.UpdateWorkspace(workspace)

	if err := s.workspaceRepo.Update(workspace); err != nil {
		return nil, err
	}

	response := workspace.ToResponse()
	return &response, nil
}
