package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role string

const (
	RoleUser    Role = "user"
	RoleManager Role = "manager"
	RoleHead    Role = "head"
	RoleRoot    Role = "root"
)

type User struct {
	ID         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string         `gorm:"type:varchar(100);not null" json:"name"`
	Email      string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password   string         `gorm:"type:varchar(255);not null" json:"-"`
	Vin        string         `gorm:"type:varchar(255);not null" json:"-"`
	Role       Role           `gorm:"type:varchar(20);default:'user';not null" json:"role"`
	UserHeadID *int64         `gorm:"type:int;default:NULL" json:"user_head_id"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type UserSignupRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	Vin         string `json:"vin" binding:"required,min=10"`
	AdminSecret string `json:"adminSecret" binding:"required,min=2,max=100"`
}

type UserSigninRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserCreateRequest struct {
	Name       string `json:"name" binding:"required,min=2,max=100"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	Role       Role   `json:"role" binding:"required,oneof=user manager head root"`
	Vin        string `json:"vin" binding:"required,min=10"`
	UserHeadID *int64 `json:"user_head_id"`
}

type UserUpdateRequest struct {
	Name       string `json:"name" binding:"omitempty,min=2,max=100"`
	Email      string `json:"email" binding:"omitempty,email"`
	Password   string `json:"password" binding:"omitempty,min=6"`
	Role       Role   `json:"role" binding:"omitempty,oneof=user manager head root"`
	Vin        string `json:"vin" binding:"required,min=10"`
	UserHeadID *int64 `json:"user_head_id"`
}

type UserResponse struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Role       Role      `json:"role"`
	Vin        string    `json:"vin"`
	UserHeadID *int64    `json:"user_head_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		Role:       u.Role,
		Vin:        u.Vin,
		UserHeadID: u.UserHeadID,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

func (req *UserCreateRequest) ToUser() (*User, error) {
	user := &User{
		Name:       req.Name,
		Email:      req.Email,
		Role:       req.Role,
		UserHeadID: req.UserHeadID,
	}

	if err := user.HashPassword(req.Password); err != nil {
		return nil, err
	}

	return user, nil
}

func (req *UserSignupRequest) ToUser() (*User, error) {
	user := &User{
		Name:  req.Name,
		Email: req.Email,
		Role:  RoleUser,
	}

	if err := user.HashPassword(req.Password); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) UpdateFromRequest(req *UserUpdateRequest) error {
	if req.Name != "" {
		u.Name = req.Name
	}
	if req.Email != "" {
		u.Email = req.Email
	}
	if req.Password != "" {
		if err := u.HashPassword(req.Password); err != nil {
			return err
		}
	}
	if req.Role != "" {
		u.Role = req.Role
	}
	if req.UserHeadID != nil {
		u.UserHeadID = req.UserHeadID
	}
	return nil
}

func (r Role) IsValid() bool {
	switch r {
	case RoleUser, RoleManager, RoleHead, RoleRoot:
		return true
	default:
		return false
	}
}
