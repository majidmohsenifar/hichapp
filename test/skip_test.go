package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/majidmohsenifar/hichapp/handler/api"
	"github.com/majidmohsenifar/hichapp/repository"
	"github.com/stretchr/testify/assert"
)

func TestSkip_InvalidInputs(t *testing.T) {
	assert := assert.New(t)
	app := spawn_app()
	defer app.close()

	tests := []api.VoteReq{
		{
			UserID: 0,
		},
	}

	for _, test := range tests {
		params, err := json.Marshal(test)
		assert.Nil(err)
		body := bytes.NewReader(params)
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/polls/1/skip", app.address), body)
		assert.Nil(err)
		res, err := http.DefaultClient.Do(req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, res.StatusCode)
	}
}

func TestSkip_PollNotFound(t *testing.T) {
	assert := assert.New(t)
	app := spawn_app()
	defer app.close()
	params := api.VoteReq{
		UserID: 1,
	}

	paramsBytes, err := json.Marshal(params)
	assert.Nil(err)
	body := bytes.NewReader(paramsBytes)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/polls/1/skip", app.address), body)
	assert.Nil(err)
	res, err := http.DefaultClient.Do(req)
	assert.Nil(err)
	assert.Equal(http.StatusNotFound, res.StatusCode)
	resBody, err := io.ReadAll(res.Body)
	assert.Nil(err)
	apiErr := &api.ResponseFailure{}
	err = json.Unmarshal(resBody, apiErr)
	assert.Nil(err)
	assert.Equal("poll not found", apiErr.Error.Message)
}

func TestSkip_AlreadyVoted(t *testing.T) {
	assert := assert.New(t)
	app := spawn_app()
	defer app.close()
	err := app.redis.FlushAll(t.Context()).Err()
	assert.Nil(err)

	poll := app.CreatePollWithOptionsAndTags(t.Context(), "poll1", []string{"op1", "op2"}, []string{"tag1", "tag2"})
	options, err := app.repo.GetOptionsByPollID(t.Context(), app.db, poll.ID)
	assert.Nil(err)

	_, err = app.repo.CreateVote(t.Context(), app.db, repository.CreateVoteParams{
		PollID:    poll.ID,
		OptionID:  pgtype.Int8{Int64: int64(options[0].ID), Valid: true},
		UserID:    1,
		IsSkipped: false,
	})
	assert.Nil(err)

	params := api.VoteReq{
		UserID: 1,
	}

	paramsBytes, err := json.Marshal(params)
	assert.Nil(err)
	body := bytes.NewReader(paramsBytes)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/polls/%d/skip", app.address, poll.ID), body)
	assert.Nil(err)
	res, err := http.DefaultClient.Do(req)
	assert.Nil(err)
	assert.Equal(http.StatusUnprocessableEntity, res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	assert.Nil(err)
	apiErr := &api.ResponseFailure{}
	err = json.Unmarshal(resBody, apiErr)
	assert.Nil(err)
	assert.Equal("user already voted or skipped", apiErr.Error.Message)
}

func TestSkip_AlreadySkipped(t *testing.T) {
	assert := assert.New(t)
	app := spawn_app()
	defer app.close()
	err := app.redis.FlushAll(t.Context()).Err()
	assert.Nil(err)

	poll := app.CreatePollWithOptionsAndTags(t.Context(), "poll1", []string{"op1", "op2"}, []string{"tag1", "tag2"})
	_, err = app.repo.GetOptionsByPollID(t.Context(), app.db, poll.ID)
	assert.Nil(err)

	_, err = app.repo.CreateVote(t.Context(), app.db, repository.CreateVoteParams{
		PollID:    poll.ID,
		OptionID:  pgtype.Int8{},
		UserID:    1,
		IsSkipped: true,
	})
	assert.Nil(err)

	params := api.VoteReq{
		UserID: 1,
	}

	paramsBytes, err := json.Marshal(params)
	assert.Nil(err)
	body := bytes.NewReader(paramsBytes)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/polls/%d/skip", app.address, poll.ID), body)
	assert.Nil(err)
	res, err := http.DefaultClient.Do(req)
	assert.Nil(err)
	assert.Equal(http.StatusUnprocessableEntity, res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	assert.Nil(err)
	apiErr := &api.ResponseFailure{}
	err = json.Unmarshal(resBody, apiErr)
	assert.Nil(err)
	assert.Equal("user already voted or skipped", apiErr.Error.Message)
}

func TestSkip_Successful(t *testing.T) {
	assert := assert.New(t)
	app := spawn_app()
	defer app.close()
	err := app.redis.FlushAll(t.Context()).Err()
	assert.Nil(err)

	poll := app.CreatePollWithOptionsAndTags(t.Context(), "poll1", []string{"op1", "op2"}, []string{"tag1", "tag2"})
	_, err = app.repo.GetOptionsByPollID(t.Context(), app.db, poll.ID)
	assert.Nil(err)

	params := api.VoteReq{
		UserID: 1,
	}

	paramsBytes, err := json.Marshal(params)
	assert.Nil(err)
	body := bytes.NewReader(paramsBytes)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/polls/%d/skip", app.address, poll.ID), body)
	assert.Nil(err)
	res, err := http.DefaultClient.Do(req)
	assert.Nil(err)
	assert.Equal(http.StatusCreated, res.StatusCode)

	//check the vote table for the created vote
	vote, err := app.repo.GetVoteByPollIDAndUserID(t.Context(), app.db, repository.GetVoteByPollIDAndUserIDParams{
		PollID: poll.ID,
		UserID: 1,
	})
	assert.Nil(err)
	assert.False(vote.OptionID.Valid)
	assert.Equal(true, vote.IsSkipped)
}
