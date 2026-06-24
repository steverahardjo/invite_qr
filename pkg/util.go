package pkg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

type IDEncryptor struct {
	salt string
}

func NewIDEncryptor(salt string) *IDEncryptor {
	return &IDEncryptor{
		salt: salt,
	}
}

func (e *IDEncryptor) Encode(id int32, wa string, email string) string {
	payload := fmt.Sprintf("%d|%s|%s", id, wa, email)

	mac := hmac.New(sha256.New, []byte(e.salt))
	mac.Write([]byte(payload))

	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	token := base64.RawURLEncoding.EncodeToString([]byte(payload))

	return token + "." + signature
}

func (e *IDEncryptor) Decode(token string) (int32, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid token")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return 0, err
	}

	payload := string(payloadBytes)

	mac := hmac.New(sha256.New, []byte(e.salt))
	mac.Write([]byte(payload))

	expected := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expected), []byte(parts[1])) {
		return 0, fmt.Errorf("invalid signature")
	}

	fields := strings.Split(payload, "|")
	if len(fields) < 1 {
		return 0, fmt.Errorf("invalid payload")
	}

	id, err := strconv.ParseInt(fields[0], 10, 32)
	if err != nil {
		return 0, err
	}

	return int32(id), nil
}
