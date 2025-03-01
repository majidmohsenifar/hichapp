package limiter

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	UserVoteLimiterPrefix = "user_vote_limiter:"
	MaxVotesPerDay        = 100
)

type UserVoteLimiter struct {
	redisClient redis.UniversalClient
}

func (l *UserVoteLimiter) IsUserAllowedToVote(ctx context.Context, userID int64) (bool, error) {
	voteCount, err := l.redisClient.Get(ctx, fmt.Sprintf("%s%d", UserVoteLimiterPrefix, userID)).Int64()
	if err != nil && err != redis.Nil {
		slog.Error("failed to get user vote count from cache", "err", err)
		return false, err
	}
	if errors.Is(err, redis.Nil) {
		return true, nil
	}

	if voteCount < MaxVotesPerDay {
		return true, nil
	}
	return false, nil
}

func (l *UserVoteLimiter) IncreaseUserVoteCount(ctx context.Context, userID int64) error {
	count, err := l.redisClient.Incr(ctx, fmt.Sprintf("%s%d", UserVoteLimiterPrefix, userID)).Result()
	if err != nil {
		slog.Error("failed to increase user vote count", "err", err)
		return err
	}
	if count == 1 {
		now := time.Now()
		endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
		expireTime := endOfDay.Sub(now)
		_, err := l.redisClient.Expire(ctx, fmt.Sprintf("%s%d", UserVoteLimiterPrefix, userID), expireTime).Result()
		if err != nil {
			slog.Error("failed to set expiration for user vote count", "err", err)
			return err
		}
	}
	return nil
}

func NewUserVoteLimiter(
	redisClient redis.UniversalClient,
) *UserVoteLimiter {
	return &UserVoteLimiter{
		redisClient: redisClient,
	}
}
