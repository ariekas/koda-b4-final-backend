package repository

import (
	"context"
	"fmt"
	"shortlink/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matthewhartstonge/argon2"
)

func CreateUser(ctx *gin.Context, pool *pgxpool.Pool) (models.Users, error) {
	argon := argon2.DefaultConfig()
	var input models.Users
	var checkEmail bool
	now := time.Now()

	err := ctx.BindJSON(&input)
	if err != nil {
		return models.Users{}, fmt.Errorf("error failed type json, %w", err)
	}

	hash, err := argon.HashEncoded([]byte(input.Password))

	if err != nil {
		return models.Users{}, fmt.Errorf("failed to hash password, %w", err)
	}

	if checkEmail {
		return models.Users{}, fmt.Errorf("email already registered")
	}

	err = pool.QueryRow(context.Background(),
		"INSERT INTO users (username, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, username, email,  created_at, updated_at",
		input.Username, input.Email, hash, now, now,
	).Scan(
		&input.Id,
		&input.Username,
		&input.Email,
		&input.CreatedAt,
		&input.UpdatedAt,
	)

	if err != nil {
		return models.Users{}, fmt.Errorf("failed to insert user, %w", err)
	}

	user := models.Users{
		Username:  input.Username,
		Email:     input.Email,
		Password:  string(hash),
		CreatedAt: now,
		UpdatedAt: now,
	}

	return user, nil
}

func FindUserEmail(pool *pgxpool.Pool, email string) (models.Users, error) {
	var user models.Users

	row := pool.QueryRow(context.Background(), `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE users.email = $1
	`, email)

	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return models.Users{}, fmt.Errorf("no user found with this email address, %w", err)
	}
	return user, nil
}

