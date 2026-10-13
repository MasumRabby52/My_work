package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	var ctx context.Context = context.Background()

	var config *pgxpool.Config
	var err error
	config, err = pgxpool.ParseConfig(databaseURL)

	if err != nil {
		log.Println("unable to parse DATABASE_URL", err)
		return nil, err
	}

	var pool *pgxpool.Pool
	pool, err = pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		log.Println("Unable to create connection pool", err)
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		log.Println("Unable to ping database: ", err)
		pool.Close()
		return nil, err
	}
	return pool, nil
}
