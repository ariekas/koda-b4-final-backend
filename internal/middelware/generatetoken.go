package middelware

import (
	"shortlink/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


func GenerateToken(jwtToken string, userId int) (string, error) {
    claims := models.UserClaims{
        UserId: userId,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
            IssuedAt: jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(jwtToken))

    return tokenString, err
}