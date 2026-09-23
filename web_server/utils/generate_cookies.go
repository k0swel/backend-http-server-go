package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateSessionCookie генерирует  128-битный session ID
func GenerateSessionCookie() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
