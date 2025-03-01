package tag

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/majidmohsenifar/hichapp/repository"
	"github.com/redis/go-redis/v9"
)

const (
	TagCachePrefix = "tag:"
)

var (
	ErrTagNotFound = errors.New("tag not found")
)

type Service struct {
	db          *pgxpool.Pool
	repository  *repository.Queries
	redisClient redis.UniversalClient
}

func (s *Service) GetTagIDByName(ctx context.Context, name string) (int64, error) {
	tagID, err := s.redisClient.Get(ctx, fmt.Sprintf("%s%s", TagCachePrefix, name)).Int64()
	if err != nil && err != redis.Nil {
		slog.Error("failed to get tag from cache", "err", err)
		return 0, err
	}
	if err == nil {
		return tagID, nil
	}
	tag, err := s.repository.GetTagByName(ctx, s.db, name)
	if err != nil && err != pgx.ErrNoRows {
		return 0, err
	}
	if err == pgx.ErrNoRows {
		return 0, ErrTagNotFound
	}
	err = s.redisClient.Set(ctx, fmt.Sprintf("%s%s", TagCachePrefix, name), tag.ID, 0).Err()
	if err != nil {
		slog.Error("failed to set tag in cache", "err", err)
	}
	return tag.ID, nil
}

func New(
	db *pgxpool.Pool,
	repository *repository.Queries,
	redisClient redis.UniversalClient,
) *Service {
	return &Service{
		db:          db,
		repository:  repository,
		redisClient: redisClient,
	}
}
