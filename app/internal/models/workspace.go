package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Workspace struct {
	ID              int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string     `json:"name" binding:"required,min=2,max=100"`
	UserHeadId      int64      `gorm:"type:int;default:5;not null" json:"user_head_id"`
	TgBotToken      string     `json:"tg_bot_token" binding:"required,min=10,max=300"`
	IncludesUsersId Int64Array `gorm:"type:jsonb;default:'[]'" json:"includes_users_id"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type Int64Array []int64

func (a *Int64Array) Scan(value interface{}) error {
	if value == nil {
		*a = []int64{}
		return nil
	}
	return json.Unmarshal(value.([]byte), a)
}

func (a Int64Array) Value() (driver.Value, error) {
	if a == nil {
		return json.Marshal([]int64{})
	}
	return json.Marshal(a)
}

func (a *Int64Array) GormDataType() string {
	return "jsonb"
}

type WorkspaceCreateRequest struct {
	Name            string  `json:"name" binding:"required,min=2,max=100"`
	UserHeadId      int64   `json:"user_head_id" binding:"required,min=1"`
	TgBotToken      string  `json:"tg_bot_token" binding:"required,min=10,max=300"`
	IncludesUsersId []int64 `json:"includes_users_id"`
}

type WorkspaceUpdateRequest struct {
	Name            string  `json:"name" binding:"omitempty,min=2,max=100"`
	UserHeadId      int64   `json:"user_head_id" binding:"omitempty,min=1"`
	TgBotToken      string  `json:"tg_bot_token" binding:"omitempty,min=10,max=300"`
	IncludesUsersId []int64 `json:"includes_users_id"`
}

type WorkspaceResponse struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	UserHeadId      int64     `json:"user_head_id"`
	TgBotToken      string    `json:"tg_bot_token"`
	IncludesUsersId []int64   `json:"includes_users_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (req *WorkspaceCreateRequest) ToWorkspace() *Workspace {
	includesUsers := Int64Array{}
	if req.IncludesUsersId != nil {
		includesUsers = req.IncludesUsersId
	}

	return &Workspace{
		Name:            req.Name,
		UserHeadId:      req.UserHeadId,
		TgBotToken:      req.TgBotToken,
		IncludesUsersId: includesUsers,
	}
}

func (req *WorkspaceUpdateRequest) UpdateWorkspace(workspace *Workspace) {
	if req.Name != "" {
		workspace.Name = req.Name
	}
	if req.UserHeadId > 0 {
		workspace.UserHeadId = req.UserHeadId
	}
	if req.TgBotToken != "" {
		workspace.TgBotToken = req.TgBotToken
	}
	if req.IncludesUsersId != nil {
		workspace.IncludesUsersId = req.IncludesUsersId
	}
}

func (w *Workspace) ToResponse() WorkspaceResponse {
	return WorkspaceResponse{
		ID:              w.ID,
		Name:            w.Name,
		UserHeadId:      w.UserHeadId,
		TgBotToken:      w.TgBotToken,
		IncludesUsersId: []int64(w.IncludesUsersId),
		CreatedAt:       w.CreatedAt,
		UpdatedAt:       w.UpdatedAt,
	}
}

func (w *Workspace) ToSafeResponse() WorkspaceResponse {

	return WorkspaceResponse{
		ID:              w.ID,
		Name:            w.Name,
		UserHeadId:      w.UserHeadId,
		TgBotToken:      "",
		IncludesUsersId: []int64(w.IncludesUsersId),
		CreatedAt:       w.CreatedAt,
		UpdatedAt:       w.UpdatedAt,
	}
}
