package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/techagentng/ecommerce-api/errors"
	"net/http"
	"time"
)

const AccessTokenValidity = time.Hour * 24 * 7
const RefreshTokenValidity = time.Hour * 24 * 30

// verifyAccessToken verifies a token
func verifyToken(tokenString string, secret string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func isJWTSecretEmpty(secret string) bool {
	return secret == ""
}

func isAccessTokenEmpty(token string) bool {
	return token == ""
}

func ValidateToken(token string, secret string) (*jwt.Token, error) {
	tk, err := verifyToken(token, secret)
	if err != nil { // TODO: remove
		return nil, fmt.Errorf("invalid token1: %v", err) // TODO: probably need to errors.NEw
	}
	if !tk.Valid {
		return nil, errors.New("invalid token", http.StatusUnauthorized)
	}
	return tk, nil
}

func getClaims(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("could not get claims")
	}
	return claims, claims.Valid()
}

func ValidateAndGetClaims(tokenString string, secret string) (jwt.MapClaims, error) {
	if tokenString == "" {
		return nil, errors.New("invalid token (token is empty)", http.StatusUnauthorized)
	}
	token, err := ValidateToken(tokenString, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}
	claims, err := getClaims(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get claims: %v", err)
	}
	return claims, nil
}

// GenerateToken generates only an access token
func GenerateToken(email string, secret string, isAdmin bool, id uint, roleName string) (string, error) {
	if secret == "" {
		// Return a descriptive error message for missing secret
		return "", errors.New("secret key is required", errors.ErrBadRequest.Status)
	}

	// Generate claims with the role name
	claims := GenerateClaims(email, isAdmin, id, roleName)

	// Create and sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateTokenPair(email string, secret string, isAdmin bool, id uint, roleName string) (accessToken string, refreshToken string, err error) {
	accessToken, err = GenerateToken(email, secret, isAdmin, id, roleName)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = GenerateRefreshToken(email, secret, isAdmin, id, roleName)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func GenerateRefreshToken(email string, secret string, isAdmin bool, id uint, roleName string) (string, error) {
	if secret == "" {
		return "", errors.New("secret key is required", errors.ErrInternalServerError.Status)
	}

	// Create claims with role information if needed
	refreshTokenClaims := jwt.MapClaims{
		"email":    email,
		"exp":      time.Now().Add(RefreshTokenValidity).Unix(),
		"is_admin": isAdmin,
		"id":       id,
		"role":     roleName, // Include roleName if applicable
		"type":     "refresh_token",
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	// Sign and get the complete encoded token as a string using the secret
	refreshTokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return refreshTokenString, nil
}

func GenerateClaims(email string, isAdmin bool, id uint, roleName string) jwt.MapClaims {
	accessClaims := jwt.MapClaims{
		"email":    email,
		"exp":      time.Now().Add(AccessTokenValidity).Unix(),
		"is_admin": isAdmin,
		"id":       id,
		"role":     roleName,
	}
	return accessClaims
}
