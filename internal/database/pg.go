package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresDatabase will connect to a postgres database
func NewPostgresDatabase(connectionString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return pool, nil
}

type Database struct {
	db *pgxpool.Pool
}

func NewDatabase(db *pgxpool.Pool) Database {
	return Database{db}
}
