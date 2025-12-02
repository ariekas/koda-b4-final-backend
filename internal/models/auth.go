package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Register struct {
	Id int `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Role string `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Login struct {
	UserId int `json:"userId"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type InputLogin struct{
	Email string `json:"email"`
	Password string `json:"password"`
}	