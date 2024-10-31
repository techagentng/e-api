package models

type User struct {
    Model
    Email    string `gorm:"unique" json:"email" binding:"required,email"`
    Password string `json:"password,omitempty" binding:"required"`
    Role     string `json:"role" gorm:"default:'user'"` 
}