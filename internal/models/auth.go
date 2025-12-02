package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Users struct {
	Id int `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserClaims struct {
	UserId int `json:"userId"`
	jwt.RegisteredClaims
}

type InputLogin struct{
	Email string `json:"email"`
	Password string `json:"password"`
}	
