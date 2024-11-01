package services

import (
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	gofrsUUID "github.com/gofrs/uuid"
	"github.com/techagentng/ecommerce-api/config"
	"github.com/techagentng/ecommerce-api/db"
	apiError "github.com/techagentng/ecommerce-api/errors"
	"github.com/techagentng/ecommerce-api/models"
	"github.com/techagentng/ecommerce-api/services/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService interface
type AuthService interface {
	GetRoleByName(name string) (*models.Role, error)
	SignupUser(request *models.User) (*models.User, error)
	LoginUser(request *models.LoginRequest) (*models.LoginResponse, *apiError.Error)
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

// LoginUser logs in a user and returns the login response
func (a *authService) LoginUser(loginRequest *models.LoginRequest) (*models.LoginResponse, *apiError.Error) {
    // Find the user by email
    foundUser, err := a.authRepo.FindUserByEmail(loginRequest.Email)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, apiError.New("invalid email or password", http.StatusUnprocessableEntity)
        }
        log.Printf("Error finding user by email: %v", err)
        return nil, apiError.New("unable to find user", http.StatusInternalServerError)
    }

    // Verify user password
    if err := foundUser.VerifyPassword(loginRequest.Password); err != nil {
        log.Printf("Invalid password for user %s", foundUser.Email)
        return nil, apiError.ErrInvalidPassword
    }

    if foundUser.RoleID == uuid.Nil {
        log.Printf("User %s does not have a role assigned", foundUser.Email)
        return nil, apiError.New("user role not assigned", http.StatusInternalServerError)
    }

	convertedRoleID, err := gofrsUUID.FromString(foundUser.RoleID.String())
	if err != nil {
		log.Printf("Error converting RoleID for user %s: %v", foundUser.Email, err)
		return nil, apiError.New("unable to convert role ID", http.StatusInternalServerError)
	}

    // Fetch the user's role
    log.Printf("Fetching role with ID: %s for user %s", foundUser.RoleID.String(), foundUser.Email)
    role, err := a.authRepo.FindRoleByID(convertedRoleID)
    if err != nil {
        log.Printf("Error fetching role for user %s: %v", foundUser.Email, err)
        return nil, apiError.New("unable to fetch role", http.StatusInternalServerError)
    }
    
    roleName := role.Name
    log.Printf("Generating token pair for user %s with role %s", foundUser.Email, roleName)
    accessToken, refreshToken, err := jwt.GenerateTokenPair(foundUser.Email, a.Config.JWTSecret, foundUser.AdminStatus, foundUser.ID, roleName)
    if err != nil {
        log.Printf("Error generating token pair for user %s: %v", foundUser.Email, err)
        return nil, apiError.ErrInternalServerError
    }

    return &models.LoginResponse{
        UserResponse: models.UserResponse{
            ID:        foundUser.ID,
            Fullname:  foundUser.Fullname,
            Username:  foundUser.Username,
            Telephone: foundUser.Telephone,
            Email:     foundUser.Email,
            RoleName:  roleName, 
        },
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }, nil
}