package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/majidmohsenifar/hichapp/infra/metrics"
)

func NewDBClient(ctx context.Context, dbDSN string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dbDSN)
	if err != nil {
		return nil, err
	}

	cfg.MinConns = 2
	cfg.MaxConns = 20
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.MaxConnLifetime = 1 * time.Hour

	cfg.ConnConfig.Tracer = &metrics.DBTracer{}

	dbPool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if err := dbPool.Ping(ctx); err != nil {
		return nil, err
	}
	return dbPool, nil
}
