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
	Pic       string    `json:"pic"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PicUser struct {
	Id        int       `json:"id"`
	Pic       string    `json:"pic"`
	UserId    int       `json:"userId"`
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
