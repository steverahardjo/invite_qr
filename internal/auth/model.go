package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
)

// function to hash token before storing it in the database
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func DBOneTimeToken(ctx context.Context, UserID int, token string) error {
	hashedToken := hashToken(token)

}
