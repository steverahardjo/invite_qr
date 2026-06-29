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

// Admin represents an admin user with a username and bcrypt password hash.
type Admin struct {
	Username     string
	PasswordHash string
}

// JwtService provides JWT token creation and admin credential validation.
type JwtService struct {
	SecretKey      []byte
	ExpirationTime time.Duration
	Issuer         string
}

// NewJwtService creates a JwtService with the given secret, token expiry duration, and issuer name.
func NewJwtService(secret string, expiry time.Duration, issuer string) *JwtService {
	return &JwtService{
		SecretKey:      []byte(secret),
		ExpirationTime: expiry,
		Issuer:         issuer,
	}
}

// SetPasswordHashEnv hashes the plaintext password with bcrypt and stores it in the
// HASHED_ADMIN_PASSWORD environment variable for subsequent login comparisons.
func (j *JwtService) SetPasswordHashEnv(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	os.Setenv("HASHED_ADMIN_PASSWORD", string(hashed))
	return nil
}

// ComparePasswordHash reads the hashed admin password from the HASHED_ADMIN_PASSWORD
// environment variable, compares it against the provided plaintext password using
// bcrypt, and returns whether they match. Internal errors are logged but not leaked
// to the caller.
func ComparePasswordHash(password string, log *zap.Logger) (bool, error) {
	hashedAdmin := os.Getenv("HASHED_ADMIN_PASSWORD")
	if hashedAdmin == "" {
		log.Error("no hashed admin password set")
		return false, errors.New("server configuration error")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedAdmin), []byte(password))
	if err != nil {
		log.Warn("failed password comparison attempt", zap.Error(err))
		return false, nil
	}

	return true, nil
}

// LoginAdmin validates the provided username and password against the stored
// hash and returns a signed JWT (HS256) containing the username, admin flag,
// issuer, expiration, and issued-at claims on success.
func (j *JwtService) LoginAdmin(ctx context.Context, username, password string) (string, error) {
	logger := server.LoggerFromContext(ctx)

	isValid, err := ComparePasswordHash(password, logger)
	if err != nil {
		return "", err
	}

	if !isValid {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"sub":      username,
		"is_admin": true,
		"iss":      j.Issuer,
		"exp":      time.Now().Add(j.ExpirationTime).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(j.SecretKey)
	if err != nil {
		logger.Error("failed to sign token", zap.Error(err))
		return "", err
	}

	return signedToken, nil
}
