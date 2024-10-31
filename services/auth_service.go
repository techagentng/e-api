package services

import (
	"github.com/techagentng/ecommerce-api/config"
	"github.com/techagentng/ecommerce-api/db"
)

// AuthService interface
type AuthService interface {

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
