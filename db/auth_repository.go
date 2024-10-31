package db

import (
	"gorm.io/gorm"
)

type AuthRepository interface {
	
}

type authRepo struct {
	DB *gorm.DB
}

func NewAuthRepo(db *GormDB) AuthRepository {
	return &authRepo{db.DB}
}