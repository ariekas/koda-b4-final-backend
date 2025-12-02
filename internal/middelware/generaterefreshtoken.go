package middelware

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func GenerateRefreshToken() (string, string, error) {
	refreshToken := make([]byte, 32)
	
	_, err := rand.Read(refreshToken)
	if err != nil {
		fmt.Printf("error read refresh token, %s", err)
	}

	token := base64.StdEncoding.EncodeToString(refreshToken)

	hash := sha256.Sum256([]byte(token))
	hashToken := base64.StdEncoding.EncodeToString(hash[:])

	return token, hashToken, nil
}