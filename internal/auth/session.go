package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"strings"
)

func GenerateSessionToken() (string, error) {
	bytes := make([]byte, 20)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	base32Encoded := base32.StdEncoding.EncodeToString(bytes)
	base32LowerCaseNoPadding := strings.ToLower(strings.TrimRight(base32Encoded, "="))

	return base32LowerCaseNoPadding, nil
}

func CreateSessionId(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	t := h.Sum(nil)

	session := hex.EncodeToString(t)

	return session
}

// func ValidateSessionToken(ctx context.Context, token string) (Session, error) {

// }

// func InvalidateSession(ctx context.Context, sessionId string) error {

// }
