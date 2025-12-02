package repository

import (
	"context"
	"fmt"
	"shortlink/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)


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

func RevokedSesstion(pgx *pgxpool.Pool, token string) error {
	_, err := pgx.Exec(context.Background(), "DELETE FROM sessions WHERE refreshtoken=$1", token)

	if err != nil {
		return fmt.Errorf("failed to delete sesstions, %w", err)
	}

	return nil
}