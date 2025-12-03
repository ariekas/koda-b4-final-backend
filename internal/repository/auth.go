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
	defaultPic := "default-avatar.png" // default picture

	err := ctx.BindJSON(&input)
	if err != nil {
		return models.Users{}, fmt.Errorf("error failed type json, %w", err)
	}

	err = pool.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)",
		input.Email,
	).Scan(&checkEmail)

	if err != nil {
		return models.Users{}, fmt.Errorf("failed to check email, %w", err)
	}

	if checkEmail {
		return models.Users{}, fmt.Errorf("email already registered")
	}

	hash, err := argon.HashEncoded([]byte(input.Password))
	if err != nil {
		return models.Users{}, fmt.Errorf("failed to hash password, %w", err)
	}

	tx, err := pool.Begin(context.Background())
	if err != nil {
		return models.Users{}, fmt.Errorf("failed to begin transaction, %w", err)
	}
	defer tx.Rollback(context.Background())

	var userId int
	err = tx.QueryRow(context.Background(),
		"INSERT INTO users (username, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		input.Username, input.Email, hash, now, now,
	).Scan(&userId)

	if err != nil {
		return models.Users{}, fmt.Errorf("failed to insert user, %w", err)
	}

	err = tx.QueryRow(context.Background(),
		"INSERT INTO pic_user (pic, userId, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id",
		defaultPic, userId, now, now,
	).Scan(&input.Id)

	if err != nil {
		return models.Users{}, fmt.Errorf("failed to insert pic_user, %w", err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return models.Users{}, fmt.Errorf("failed to commit transaction, %w", err)
	}

	Users := models.Users{
		Id:        userId,
		Username:  input.Username,
		Email:     input.Email,
		Pic:       defaultPic,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return Users, nil
}

func GetUserWithPic(pool *pgxpool.Pool, userId int) (models.Users, error) {
	var user models.Users

	err := pool.QueryRow(context.Background(),
		`SELECT u.id, u.username, u.email, COALESCE(p.pic, 'default-avatar.png') as pic, u.created_at, u.updated_at
		FROM users u
		LEFT JOIN pic_user p ON u.id = p.userId
		WHERE u.id = $1`,
		userId,
	).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Pic,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return models.Users{}, fmt.Errorf("failed to get user, %w", err)
	}

	return user, nil
}

func UpdateUserPic(pool *pgxpool.Pool, userId int, picPath string) error {
	now := time.Now()

	var exists bool
	err := pool.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM pic_user WHERE userId = $1)",
		userId,
	).Scan(&exists)

	if err != nil {
		return fmt.Errorf("failed to check pic_user, %w", err)
	}

	if exists {
		_, err = pool.Exec(context.Background(),
			"UPDATE pic_user SET pic = $1, updated_at = $2 WHERE userId = $3",
			picPath, now, userId,
		)
	} else {
		_, err = pool.Exec(context.Background(),
			"INSERT INTO pic_user (pic, userId, created_at, updated_at) VALUES ($1, $2, $3, $4)",
			picPath, userId, now, now,
		)
	}

	if err != nil {
		return fmt.Errorf("failed to update pic_user, %w", err)
	}

	return nil
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

