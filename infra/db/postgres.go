package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDBClient(ctx context.Context, dbDSN string) (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(ctx, dbDSN)
	if err != nil {
		return nil, err
	}
	if err := dbPool.Ping(ctx); err != nil {
		return nil, err
	}
	return dbPool, nil
}
