package models

import "github.com/google/uuid"

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