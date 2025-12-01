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

func CreateUser(ctx *gin.Context, pool *pgxpool.Pool) (models.Register, error) {
	argon := argon2.DefaultConfig()
	var input models.Register
	var checkEmail bool
	now := time.Now()

	if err := ctx.ShouldBindJSON(&input); err != nil {
		return models.Register{}, fmt.Errorf("invalid request body: %w", err)
	}
	
	hash, err := argon.HashEncoded([]byte(input.Password))

	if err != nil {
		return models.Register{}, fmt.Errorf("failed to hash password, %w", err)
	}

	if checkEmail {
		return models.Register{}, fmt.Errorf("email already registered")
	}

	_, err = pool.Exec(context.Background(), "INSERT INTO users (username, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", input.Username, input.Email, hash, input.CreatedAt, input.UpdatedAt)

	if err != nil {
		return models.Register{}, fmt.Errorf("failed to insert user, %w", err)
	}

	user := models.Register{
		Username: input.Username,
		Email: input.Email,
		Password: string(hash),
		CreatedAt: now,
		UpdatedAt: now,
	}

	return user, nil
}