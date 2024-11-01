package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Model
	ID             uint      `gorm:"primaryKey"`
	Name           string    `gorm:"size:255"`
	Fullname       string    `json:"fullname" binding:"required,min=2"`
	Username       string    `json:"username" binding:"required,min=2"`
	Telephone      string    `json:"telephone" gorm:"unique;default:null" binding:"required"`
	Email          string    `gorm:"unique;not null"`
	Password       string    `json:"password,omitempty" gorm:"-"`
	IsEmailActive  bool      `json:"-"`
	HashedPassword string    `json:"-"`
	AdminStatus    bool      `json:"is_admin" gorm:"foreignKey:Status"`
	ThumbNailURL   string    `json:"thumbnail_url,omitempty"`
	RoleID         uuid.UUID `gorm:"type:uuid" json:"role_id"`
	Role           Role      `gorm:"foreignKey:RoleID" json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Fullname  string `json:"fullname"`
	Username  string `json:"username"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
	RoleName      string             `json:"role_name"`
}

type LoginResponse struct {
	UserResponse
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (u *User) VerifyPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	if err != nil {
		return err // Passwords do not match
	}
	return nil // Passwords match
}