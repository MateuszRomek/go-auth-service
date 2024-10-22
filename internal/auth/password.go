package auth

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

type argon2Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var params = &argon2Params{
	memory:      64 * 1024,
	iterations:  3,
	parallelism: 2,
	saltLength:  16,
	keyLength:   32,
}

func generateSalt(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func HashPassword(password string) ([]string, error) {
	salt, err := generateSalt(params.saltLength)
	if err != nil {
		return nil, err
	}

	hash := argon2.IDKey([]byte(password), salt, params.iterations, params.memory, params.parallelism, params.keyLength)

	encodedHash := base64.RawStdEncoding.EncodeToString(hash)
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)

	return []string{encodedSalt, encodedHash}, nil
}

func VerifyPassword(password, encodedHash, encodedSalt string) (bool, error) {
	decodedSalt, err := base64.RawStdEncoding.DecodeString(encodedSalt)
	if err != nil {
		return false, err
	}

	hash := argon2.IDKey([]byte(password), decodedSalt, params.iterations, params.memory, params.parallelism, params.keyLength)

	hashToCompare := base64.RawStdEncoding.EncodeToString(hash)

	return hashToCompare == encodedHash, nil

}
