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
		"INSERT INTO users (username, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, username, email, role, created_at, updated_at",
		input.Username, input.Email, hash, now, now,
	).Scan(
		&input.Id,
		&input.Username,
		&input.Email,
		&input.Role,
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
		Role:      input.Role,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return user, nil
}

func FindUserEmail(pool *pgxpool.Pool, email string) (models.Users, error) {
	var user models.Users

	row := pool.QueryRow(context.Background(), `
		SELECT id, username, email, password,role, created_at, updated_at
		FROM users
		WHERE users.email = $1
	`, email)

	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return models.Users{}, fmt.Errorf("no user found with this email address, %w", err)
	}
	return user, nil
}

func SaveSession(pool *pgxpool.Pool, s models.Session) error {
	_, err := pool.Exec(context.Background(),
		"INSERT INTO sessions (userid, refreshtoken, revoked, created_at, expires_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		s.UserId, s.RefreshToken, s.Revoked, s.CreatedAt, s.ExpiresAt, s.UpdatedAt,
	)
	return err
}


func FindSesstionByToken(pool *pgxpool.Pool, hash string) (models.Session, error) {
	var session models.Session
	now := time.Now()

	err := pool.QueryRow(context.Background(), "SELECT id, userid, refreshtoken,revoked, created_at, expires_at, updated_at  FROM sessions WHERE refreshtoken=$1", hash).Scan(&session.Id, &session.UserId, &session.RefreshToken, &session.Revoked, &session.CreatedAt, &session.ExpiresAt, &session.UpdatedAt)

	if err != nil {
		return session, fmt.Errorf("invalid to get refresh token, %w", err)
	}

	if session.Revoked || now.After(session.ExpiresAt) {
		return session, fmt.Errorf("refresh token expired or revoked, %w", err)
	}

	return session, nil
}

func UpdateSesstion(pool *pgxpool.Pool, id int, newToken string) error {
	_, err := pool.Exec(context.Background(), "UPDATE sessions SET refreshtoken=$1, updated_at=$2 WHERE id=$3", newToken, time.Now(), id)

	return err
}
