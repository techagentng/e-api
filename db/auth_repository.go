package db

import (
	"errors"
	"fmt"
	"log"
	"strings"

	ew "github.com/pkg/errors"
	"github.com/techagentng/ecommerce-api/models"
	"gorm.io/gorm"
)

type AuthRepository interface {
	FindUserByID(id uint) (*models.User, error)
	IsTokenInBlacklist(token string) bool
	FindRoleByName(name string) (*models.Role, error)
	CreateUser(user *models.User) (*models.User, error)
	IsEmailExist(email string) error
	FindUserByEmail(email string) (*models.User, error)
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

func (a *authRepo) FindRoleByName(name string) (*models.Role, error) {
    var role models.Role
    if err := a.DB.Where("name = ?", name).First(&role).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Role not foundx-:", name)
            return nil, errors.New("role not found--x")
        }
        return nil, err
    }
    return &role, nil
}

func (a *authRepo) CreateUser(user *models.User) (*models.User, error) {
	if user == nil {
		log.Println("CreateUser error: user is nil")
		return nil, errors.New("user is nil")
	}

	// Create the user in the database
	err := a.DB.Create(user).Error
	if err != nil {
		log.Printf("CreateUser error: %v", err)
		return nil, err
	}

	return user, nil
}

func (a *authRepo) IsEmailExist(email string) error {
	var count int64
	err := a.DB.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No user found with this email, return nil
			return nil
		}
		// Return wrapped error for other errors
		return ew.Wrap(err, "gorm count error")
	}
	if count > 0 {
		// Email already exists, return specific error
		return errors.New("email already in use")
	}
	return nil
}

func (a *authRepo) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := a.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}
	return &user, nil
}