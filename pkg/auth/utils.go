package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

func GenerateSessionToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(token), nil
}
func getSessionID(token string) string {
	return base64.RawStdEncoding.EncodeToString(sha256.New().Sum([]byte(token)))
}
