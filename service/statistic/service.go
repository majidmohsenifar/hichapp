package statistic

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/majidmohsenifar/hichapp/repository"
)

type Service struct {
	db         *pgxpool.Pool
	repository *repository.Queries
}

type StatsVote struct {
	Option string
	Count  int32
}

type StatsResult struct {
	PollID int64
	Votes  []StatsVote
}

func (s *Service) GetStatsForPoll(ctx context.Context, pollID int64) (StatsResult, error) {
	options, err := s.repository.GetOptionsContentAndCountByPollID(ctx, s.db, pollID)
	if err != nil {
		slog.Error("failed to get options by poll id", "err", err)
		return StatsResult{}, errors.New("failed to get options by poll id")
	}
	votes := make([]StatsVote, len(options))
	for i, option := range options {
		votes[i] = StatsVote{
			Option: option.Content,
			Count:  option.Counts,
		}

	}
	stats := StatsResult{
		PollID: pollID,
		Votes:  votes,
	}
	return stats, nil
}

func New(
	db *pgxpool.Pool,
	repository *repository.Queries,
) *Service {
	return &Service{
		db:         db,
		repository: repository,
	}
}
