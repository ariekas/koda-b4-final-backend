package database

import (
	"context"
	"shortlink/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Database() *pgxpool.Pool {
	dbUrl := config.GetDatabase()

	pool, err := pgxpool.New(context.Background(), dbUrl)

	if err != nil {
		panic("rrror : failed to connect database")
	}
	
	return pool
}