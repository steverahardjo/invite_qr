package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"invite_qr/internal/server"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	Username     string
	PasswordHash string
}

type JwtService struct {
	SecretKey      []byte
	ExpirationTime time.Duration
	Issuer         string
}

// NewJwtService is a constructor to safely initialize your service
func NewJwtService(secret string, expiry time.Duration, issuer string) *JwtService {
	return &JwtService{
		SecretKey:      []byte(secret),
		ExpirationTime: expiry,
		Issuer:         issuer,
	}
}

func (j *JwtService) SetPasswordHashEnv(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	os.Setenv("HASHED_ADMIN_PASSWORD", string(hashed))
	return nil
}

// ComparePasswordHash gets the password hash from the environment and compares it with the user input
func ComparePasswordHash(password string, log *zap.Logger) (bool, error) {
	// Get the hash directly from the environment variables
	hashedAdmin := os.Getenv("HASHED_ADMIN_PASSWORD")
	if hashedAdmin == "" {
		log.Error("no hashed admin password set")
		return false, errors.New("server configuration error")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedAdmin), []byte(password))
	if err != nil {
		// Log the error internally, but don't crash or leak data
		log.Warn("failed password comparison attempt", zap.Error(err))
		return false, nil
	}

	return true, nil
}

// LoginAdmin validates the credentials and returns a signed JWT if successful
func (j *JwtService) LoginAdmin(ctx context.Context, username, password string) (string, error) {
	logger := server.LoggerFromContext(ctx)

	// Correctly handle the two return values from ComparePasswordHash
	isValid, err := ComparePasswordHash(password, logger)
	if err != nil {
		return "", err
	}

	if !isValid {
		return "", errors.New("invalid credentials")
	}

	// Generate claims if password is valid
	claims := jwt.MapClaims{
		"sub":      username,
		"is_admin": true,
		"iss":      j.Issuer, // Using the struct's config instead of hardcoded "admin"
		"exp":      time.Now().Add(j.ExpirationTime).Unix(),
		"iat":      time.Now().Unix(),
	}

	// Create and sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(j.SecretKey)
	if err != nil {
		logger.Error("failed to sign token", zap.Error(err))
		return "", err
	}

	return signedToken, nil
}
