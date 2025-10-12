package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Database struct encapsulates the database connection pool
type Database struct {
	pool *pgxpool.Pool
}

func wireDatabaseFn(ctx context.Context, dbURL string) (*Database, error) {
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

// Close closes the database connection pool.
func (db *Database) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}
