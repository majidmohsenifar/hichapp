package poll

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/majidmohsenifar/hichapp/repository"
	"github.com/majidmohsenifar/hichapp/service/limiter"
	"github.com/majidmohsenifar/hichapp/service/tag"
)

var (
	ErrPollNotFound              = errors.New("poll not found")
	ErrUserAlreadyVotedOrSkipped = errors.New("user already voted or skipped")
	ErrInvalidOptionIndex        = errors.New("invalid option index")
	ErrUserNotAllowedToVote      = errors.New("user not allowed to vote")
)

type Service struct {
	db              *pgxpool.Pool
	repository      *repository.Queries
	tagService      *tag.Service
	userVoteLimiter *limiter.UserVoteLimiter
}

type CreatePollParams struct {
	Title   string
	Options []string
	Tags    []string
}

type VoteOrSkipParams struct {
	PollID      int64
	OptionIndex int8
	UserID      int64
	IsSkipped   bool
}

type SinglePollList struct {
	ID      int64
	Title   string
	Options []string
	Tags    []string
}

type PollListFilters struct {
	Page     int64
	PageSize int16
	Tag      string
	UserID   int64
}

func (s *Service) CreatePoll(ctx context.Context, params CreatePollParams) (int64, error) {
	dbTx, err := s.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to begin transaction", "err", err)
		return 0, errors.New("something went wrong")
	}

	poll, err := s.repository.CreatePoll(ctx, dbTx, params.Title)
	if err != nil {
		dbTx.Rollback(ctx)
		slog.Error("failed to create poll", "err", err)
		return 0, errors.New("something went wrong creating poll")
	}

	createOptionParams := make([]repository.CreateOptionParams, len(params.Options))
	for i, option := range params.Options {
		createOptionParams[i] = repository.CreateOptionParams{
			PollID:  poll.ID,
			Content: option,
		}
	}

	_, err = s.repository.CreateOption(ctx, dbTx, createOptionParams)
	if err != nil {
		dbTx.Rollback(ctx)
		slog.Error("failed to create option", "err", err)
		return 0, errors.New("something went wrong creating option")
	}

	tagIDs, err := s.repository.CreateTags(ctx, dbTx, params.Tags)
	if err != nil {
		dbTx.Rollback(ctx)
		slog.Error("failed to create tag", "err", err)
		return 0, errors.New("something went wrong creating tag")
	}

	createPollTagsParams := make([]repository.CreatePollTagParams, len(params.Tags))
	for i, t := range tagIDs {
		createPollTagsParams[i] = repository.CreatePollTagParams{
			PollID: poll.ID,
			TagID:  t,
		}
	}
	_, err = s.repository.CreatePollTag(ctx, dbTx, createPollTagsParams)
	if err != nil {
		dbTx.Rollback(ctx)
		slog.Error("failed to create tag", "err", err)
		return 0, errors.New("something went wrong creating tag")
	}

	err = dbTx.Commit(ctx)
	if err != nil {
		dbTx.Rollback(ctx)
		slog.Error("failed to commit transaction", "err", err)
		return 0, errors.New("something went wrong")
	}
	return poll.ID, nil
}

func (s *Service) GetPollsList(ctx context.Context, params PollListFilters) ([]SinglePollList, error) {
	var polls []repository.Poll
	var err error
	if params.Tag != "" {
		tagID, err := s.tagService.GetTagIDByName(ctx, params.Tag)
		if err != nil && !errors.Is(err, tag.ErrTagNotFound) {
			slog.Error("failed to get tagID by name", "err", err)
			return nil, errors.New("something went wrong getting tag")
		}
		if errors.Is(err, tag.ErrTagNotFound) {
			return nil, tag.ErrTagNotFound
		}
		polls, err = s.repository.GetPaginatedPollsByUserIDTagID(ctx, s.db, repository.GetPaginatedPollsByUserIDTagIDParams{
			UserID:  params.UserID,
			TagID:   tagID,
			Limit:   int32(params.PageSize),
			Column4: params.Page * int64(params.PageSize),
		})
		if err != nil {
			slog.Error("failed to get paginated polls by user id and tag id", "err", err)
			return nil, errors.New("something went wrong getting polls")
		}
	} else {
		polls, err = s.repository.GetPaginatedPollsByUserID(ctx, s.db, repository.GetPaginatedPollsByUserIDParams{
			UserID:  params.UserID,
			Limit:   int32(params.PageSize),
			Column3: params.Page * int64(params.PageSize),
		})
		if err != nil {
			slog.Error("failed to get paginated polls by user id", "err", err)
			return nil, errors.New("something went wrong getting polls")
		}
	}

	res := make([]SinglePollList, 0, params.PageSize)
	for _, poll := range polls {
		res = append(res, SinglePollList{
			ID:    poll.ID,
			Title: poll.Title,
		})
	}
	//getting options and tags for each poll
	pollIDs := make([]int64, len(res))
	for i, poll := range res {
		pollIDs[i] = poll.ID
	}
	options, err := s.repository.GetOptionsByPollIDs(ctx, s.db, pollIDs)
	if err != nil {
		slog.Error("failed to get options by poll ids", "err", err)
		return nil, errors.New("something went wrong getting options")
	}

	tags, err := s.repository.GetTagsByPollIDs(ctx, s.db, pollIDs)
	if err != nil {
		slog.Error("failed to get tags by poll ids", "err", err)
		return nil, errors.New("something went wrong getting tags")
	}
	for i, poll := range res {
		for _, option := range options {
			if option.PollID == poll.ID {
				res[i].Options = append(res[i].Options, option.Content)
			}
		}
		for _, tag := range tags {
			if tag.PollID == poll.ID {
				res[i].Tags = append(res[i].Tags, tag.Name)
			}
		}
	}
	return res, nil
}

func (s *Service) VoteOrSkip(ctx context.Context, params VoteOrSkipParams) error {
	poll, err := s.repository.GetPollByID(ctx, s.db, params.PollID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("failed to get poll", "err", err)
		return errors.New("something went wrong getting poll")
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrPollNotFound
	}

	vote, err := s.repository.GetVoteByPollIDAndUserID(ctx, s.db, repository.GetVoteByPollIDAndUserIDParams{
		PollID: poll.ID,
		UserID: params.UserID,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("failed to get vote", "err", err)
		return errors.New("something went wrong checking if user has already voted")
	}
	if err == nil || vote.ID > 0 {
		return ErrUserAlreadyVotedOrSkipped
	}

	optionID := pgtype.Int8{}
	if !params.IsSkipped {
		allowed, err := s.userVoteLimiter.IsUserAllowedToVote(ctx, params.UserID)
		if err != nil {
			slog.Error("failed to check if user is allowed to vote", "err", err)
			//we do not return error as it is not critical
		}
		if !allowed {
			fmt.Println("here we goooooooooooooooooooo d")
			return ErrUserNotAllowedToVote
		}

		options, err := s.repository.GetOptionsByPollID(ctx, s.db, poll.ID)
		if err != nil {
			slog.Error("failed to get options", "err", err)
			return errors.New("something went wrong getting options")
		}
		if int(params.OptionIndex) >= len(options) {
			return ErrInvalidOptionIndex
		}
		optionID.Int64 = int64(options[params.OptionIndex].ID)
		optionID.Valid = true
	}

	dbTx, err := s.db.Begin(ctx)
	if err != nil {
		slog.Error("failed to begin transaction", "err", err)
		return errors.New("something went wrong")
	}

	_, err = s.repository.CreateVote(ctx, dbTx, repository.CreateVoteParams{
		PollID:    poll.ID,
		OptionID:  optionID,
		UserID:    params.UserID,
		IsSkipped: params.IsSkipped,
	})
	if err != nil {
		dbTx.Rollback(ctx)
		slog.Error("failed to create vote", "err", err)
		return errors.New("something went wrong creating vote")
	}

	if !params.IsSkipped {
		err = s.repository.IncrementOptionVoteCount(ctx, dbTx, optionID.Int64)
		if err != nil {
			dbTx.Rollback(ctx)
			slog.Error("failed to increment option vote count", "err", err)
			return errors.New("something went wrong incrementing option vote count")
		}
		err = s.userVoteLimiter.IncreaseUserVoteCount(ctx, params.UserID)
		if err != nil {
			dbTx.Rollback(ctx)
			slog.Error("failed to increment user vote count", "err", err)
			return errors.New("something went wrong incrementing option vote count")
		}
	}
	err = dbTx.Commit(ctx)
	if err != nil {
		dbTx.Rollback(ctx)
		slog.Error("failed to commit transaction", "err", err)
		return errors.New("something went wrong")
	}
	return nil
}

func New(
	db *pgxpool.Pool,
	repository *repository.Queries,
	tagService *tag.Service,
	userVoteLimiter *limiter.UserVoteLimiter,
) *Service {
	return &Service{
		db:              db,
		repository:      repository,
		tagService:      tagService,
		userVoteLimiter: userVoteLimiter,
	}
}
