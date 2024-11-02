package models

import "github.com/google/uuid"

type Role struct {
	ID   uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name string    `json:"name"`
	UserID      uint   `json:"user_id"`
}

const (
	RoleUser  = "User"
	RoleAdmin = "Admin"
)
