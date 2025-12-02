package middelware

import (
    "net/http"
    "shortlink/internal/config"
    "shortlink/internal/models"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

func VerifToken() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "message": "missing authorization header",
            })
            c.Abort()
            return
        }

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

        secret := config.GetJwtToken()

        claims := &models.UserClaims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{
                "message": "invalid token",
            })
            c.Abort()
            return
        }

        c.Set("userId", claims.UserId)
        c.Set("role", claims.Role)

        c.Next()
    }
}
