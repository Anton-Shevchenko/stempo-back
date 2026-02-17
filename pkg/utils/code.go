package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateUniqueCode generates a unique 32-character hex code
func GenerateUniqueCode() (string, error) {
	bytes := make([]byte, 16) // 32 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
