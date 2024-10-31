package db

import (
	"errors"
	"strings"

	"github.com/techagentng/ecommerce-api/models"
	"gorm.io/gorm"
)

type AuthRepository interface {
	FindUserByID(id uint) (*models.User, error)
	IsTokenInBlacklist(token string) bool
}

type authRepo struct {
	DB *gorm.DB
}

func NewAuthRepo(db *GormDB) AuthRepository {
	return &authRepo{db.DB}
}

func normalizeToken(token string) string {
	// Trim leading and trailing white spaces
	return strings.TrimSpace(token)
}

func (a *authRepo) FindUserByID(id uint) (*models.User, error) {
	var user models.User
	err := a.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (a *authRepo) IsTokenInBlacklist(token string) bool {
	// Normalize the token
	normalizedToken := normalizeToken(token)

	var count int64
	// Assuming you have a Blacklist model with a Token field
	a.DB.Model(&models.Blacklist{}).Where("token = ?", normalizedToken).Count(&count)
	return count > 0
}