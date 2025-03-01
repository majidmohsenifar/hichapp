// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repository

import (
	"context"
)

type Querier interface {
	CreateOption(ctx context.Context, db DBTX, arg []CreateOptionParams) (int64, error)
	CreatePoll(ctx context.Context, db DBTX, title string) (Poll, error)
	CreatePollTag(ctx context.Context, db DBTX, arg []CreatePollTagParams) (int64, error)
	CreateVote(ctx context.Context, db DBTX, arg CreateVoteParams) (Vote, error)
	GetLastCreatedPoll(ctx context.Context, db DBTX) (Poll, error)
	GetOptionsByPollID(ctx context.Context, db DBTX, pollID int64) ([]Option, error)
	GetOptionsByPollIDs(ctx context.Context, db DBTX, dollar_1 []int64) ([]Option, error)
	GetOptionsContentAndCountByPollID(ctx context.Context, db DBTX, pollID int64) ([]GetOptionsContentAndCountByPollIDRow, error)
	GetPaginatedPollsByUserID(ctx context.Context, db DBTX, arg GetPaginatedPollsByUserIDParams) ([]Poll, error)
	GetPaginatedPollsByUserIDTagID(ctx context.Context, db DBTX, arg GetPaginatedPollsByUserIDTagIDParams) ([]Poll, error)
	GetPollByID(ctx context.Context, db DBTX, id int64) (Poll, error)
	GetTagByName(ctx context.Context, db DBTX, name string) (Tag, error)
	GetTagsByPollIDs(ctx context.Context, db DBTX, dollar_1 []int64) ([]GetTagsByPollIDsRow, error)
	GetVoteByPollIDAndUserID(ctx context.Context, db DBTX, arg GetVoteByPollIDAndUserIDParams) (Vote, error)
	IncrementOptionVoteCount(ctx context.Context, db DBTX, id int64) error
}

var _ Querier = (*Queries)(nil)
