package services

import (
	"errors"
	"log"

	"github.com/techagentng/ecommerce-api/config"
	"github.com/techagentng/ecommerce-api/db"
	"github.com/techagentng/ecommerce-api/models"
	"golang.org/x/crypto/bcrypt"
	apiError "github.com/techagentng/ecommerce-api/errors"
)

// AuthService interface
type AuthService interface {
	GetRoleByName(name string) (*models.Role, error)
	SignupUser(request *models.User) (*models.User, error)
}

// authService struct
type authService struct {
	Config   *config.Config
	authRepo db.AuthRepository
}

func NewAuthService(authRepo db.AuthRepository, conf *config.Config) AuthService {
	return &authService{
		Config:   conf,
		authRepo: authRepo,
	}
}

func (a *authService) GetRoleByName(name string) (*models.Role, error) {
    // Call the repository method to fetch the role
    role, err := a.authRepo.FindRoleByName(name)
    if err != nil {
        return nil, err
    }
    return role, nil
}

func (s *authService) SignupUser(user *models.User) (*models.User, error) {
	if user == nil {
		log.Println("SignupUser error: user is nil")
		return nil, errors.New("user is nil")
	}

	if user.Email == "" {
		log.Println("SignupUser error: email is empty")
		return nil, errors.New("email is empty")
	}

	// Check if the email already exists
	err := s.authRepo.IsEmailExist(user.Email)
	if err != nil {
		log.Printf("SignupUser error: %v", err)
		return nil, apiError.GetUniqueContraintError(err)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("SignupUser error hashing password: %v", err)
		return nil, apiError.ErrInternalServerError
	}
	user.HashedPassword = string(hashedPassword)
	user.Password = "" // Clear the plain password

	// Create the user in the database
	user, err = s.authRepo.CreateUser(user)
	if err != nil {
		log.Printf("SignupUser error creating user: %v", err)
		return nil, apiError.ErrInternalServerError
	}

	// Fetch the created user
	createdUser, err := s.authRepo.FindUserByEmail(user.Email)
	if err != nil {
		log.Printf("SignupUser error fetching created user: %v", err)
		return nil, apiError.ErrInternalServerError
	}

	return createdUser, nil
}