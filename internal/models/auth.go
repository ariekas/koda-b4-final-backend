package models

import "time"

type Register struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}