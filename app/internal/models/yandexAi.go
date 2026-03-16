package models

import (
	"time"

	"gorm.io/gorm"
)

type YandexAIType string

const (
	TypeVoice YandexAIType = "voice"
	TypeText  YandexAIType = "txt"
)

type YandexAI struct {
	ID                      int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Selected                bool           `gorm:"type:bool;default:false;not null" json:"selected"`
	FolderID                string         `gorm:"type:varchar(100);not null" json:"folderId"`
	OAuthToken              string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"oAuthToken"`
	PricePer1000TokensInRub int64          `gorm:"type:int;default:5;not null" json:"pricePer1000TokensInRub"`
	Type                    YandexAIType   `gorm:"type:varchar(20);default:'txt';not null" json:"type"`
	CreatedAt               time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt               time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt               gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

type YandexAICreateRequest struct {
	FolderID                string       `json:"folderId" binding:"required,min=2,max=100"`
	OAuthToken              string       `json:"oAuthToken" binding:"required,min=2,max=255"`
	PricePer1000TokensInRub int64        `json:"pricePer1000TokensInRub" binding:"required,min=1,max=1000"`
	Type                    YandexAIType `json:"type" binding:"required,oneof=voice txt"`
	Selected                bool         `json:"selected"` // Добавлено
}

type YandexAIUpdateRequest struct {
	FolderID                string       `json:"folderId" binding:"omitempty,min=2,max=100"`
	OAuthToken              string       `json:"oAuthToken" binding:"omitempty,min=2,max=255"`
	PricePer1000TokensInRub int64        `json:"pricePer1000TokensInRub" binding:"omitempty,min=1,max=1000"`
	Type                    YandexAIType `json:"type" binding:"omitempty,oneof=voice txt"`
	Selected                *bool        `json:"selected"`
}

type YandexAIResponse struct {
	ID                      int64        `json:"id"`
	Selected                bool         `json:"selected"`
	FolderID                string       `json:"folderId"`
	OAuthToken              string       `json:"oAuthToken"`
	PricePer1000TokensInRub int64        `json:"pricePer1000TokensInRub"`
	Type                    YandexAIType `json:"type"`
	CreatedAt               time.Time    `json:"createdAt"`
	UpdatedAt               time.Time    `json:"updatedAt"`
}

func (YandexAI) TableName() string {
	return "yandex_ai"
}

func (y *YandexAI) ToResponse() YandexAIResponse {
	return YandexAIResponse{
		ID:                      y.ID,
		Selected:                y.Selected,
		FolderID:                y.FolderID,
		OAuthToken:              y.OAuthToken,
		PricePer1000TokensInRub: y.PricePer1000TokensInRub,
		Type:                    y.Type,
		CreatedAt:               y.CreatedAt,
		UpdatedAt:               y.UpdatedAt,
	}
}

func (req *YandexAICreateRequest) ToYandexAI() *YandexAI {
	return &YandexAI{
		Selected:                req.Selected,
		FolderID:                req.FolderID,
		OAuthToken:              req.OAuthToken,
		PricePer1000TokensInRub: req.PricePer1000TokensInRub,
		Type:                    req.Type,
	}
}

func (y *YandexAI) UpdateFromRequest(req *YandexAIUpdateRequest) {
	if req.FolderID != "" {
		y.FolderID = req.FolderID
	}
	if req.OAuthToken != "" {
		y.OAuthToken = req.OAuthToken
	}
	if req.PricePer1000TokensInRub != 0 {
		y.PricePer1000TokensInRub = req.PricePer1000TokensInRub
	}
	if req.Type != "" {
		y.Type = req.Type
	}
	if req.Selected != nil { // Проверяем указатель
		y.Selected = *req.Selected
	}
}

func (t YandexAIType) IsValid() bool {
	switch t {
	case TypeVoice, TypeText:
		return true
	default:
		return false
	}
}
