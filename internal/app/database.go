package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

var wireDB = func(ctx context.Context, dbURL string) (*Database, error) {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return &Database{pool: pool}, nil
}

func (db *Database) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}
