package models

import "time"

type Session struct {
	Id int `json:"id"`
	UserId int `json:"userId"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}