package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/majidmohsenifar/hichapp/repository"
	"github.com/majidmohsenifar/hichapp/service/poll"
	"github.com/stretchr/testify/assert"
)

func TestPollList_Successful(t *testing.T) {
	//we have 6 polls
	//the user has voted poll1, and skipped poll3
	assert := assert.New(t)
	app := spawn_app()
	defer app.close()

	poll1 := app.CreatePollWithOptionsAndTags(t.Context(), "poll1", []string{"op1", "op2"}, []string{"tag1", "tag2"})
	poll1Options, err := app.repo.GetOptionsByPollID(t.Context(), app.db, poll1.ID)
	assert.Nil(err)
	poll2 := app.CreatePollWithOptionsAndTags(t.Context(), "poll2", []string{"op3", "op4"}, []string{"tag3", "tag4"})
	poll3 := app.CreatePollWithOptionsAndTags(t.Context(), "poll3", []string{"op5", "op6"}, []string{"tag5", "tag6"})
	poll4 := app.CreatePollWithOptionsAndTags(t.Context(), "poll4", []string{"op7", "op8"}, []string{"tag7", "tag8"})
	poll5 := app.CreatePollWithOptionsAndTags(t.Context(), "poll5", []string{"op9", "op10"}, []string{"tag9", "tag10"})
	poll6 := app.CreatePollWithOptionsAndTags(t.Context(), "poll6", []string{"op11", "op12"}, []string{"tag11", "tag12"})

	_, err = app.repo.CreateVote(t.Context(), app.db, repository.CreateVoteParams{
		PollID:    poll1.ID,
		OptionID:  pgtype.Int8{Int64: int64(poll1Options[0].ID), Valid: true},
		UserID:    1,
		IsSkipped: false,
	})
	assert.Nil(err)

	_, err = app.repo.CreateVote(t.Context(), app.db, repository.CreateVoteParams{
		PollID:    poll3.ID,
		OptionID:  pgtype.Int8{},
		UserID:    1,
		IsSkipped: true,
	})
	assert.Nil(err)

	//page 1
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/polls?user_id=%d&page=%d&page_size=%d", app.address, 1, 0, 2), nil)
	assert.Nil(err)
	res, err := http.DefaultClient.Do(req)
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	assert.Nil(err)
	apiResponse := struct {
		Success bool                  `json:"success" example:"true"`
		Message string                `json:"message,omitempty"`
		Data    []poll.SinglePollList `json:"data,omitempty"`
	}{}

	err = json.Unmarshal(resBody, &apiResponse)
	assert.Nil(err)
	assert.Equal(2, len(apiResponse.Data))

	//the first one should be poll2
	poll2Res := apiResponse.Data[0]
	assert.Equal(poll2Res.ID, poll2.ID)
	assert.Equal(poll2Res.Title, poll2.Title)

	slices.Contains(poll2Res.Options, "op3")
	slices.Contains(poll2Res.Options, "op4")

	slices.Contains(poll2Res.Tags, "tag3")
	slices.Contains(poll2Res.Tags, "tag4")

	//the second one should be poll4
	poll4Res := apiResponse.Data[1]
	assert.Equal(poll4Res.ID, poll4.ID)
	assert.Equal(poll4Res.Title, poll4.Title)

	slices.Contains(poll4Res.Options, "op7")
	slices.Contains(poll4Res.Options, "op8")

	slices.Contains(poll4Res.Tags, "tag7")
	slices.Contains(poll4Res.Tags, "tag8")

	//page 2
	req, err = http.NewRequest("GET", fmt.Sprintf("%s/api/v1/polls?user_id=%d&page=%d&page_size=%d", app.address, 1, 1, 2), nil)
	assert.Nil(err)
	res, err = http.DefaultClient.Do(req)
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	resBody, err = io.ReadAll(res.Body)
	assert.Nil(err)
	apiResponse = struct {
		Success bool                  `json:"success" example:"true"`
		Message string                `json:"message,omitempty"`
		Data    []poll.SinglePollList `json:"data,omitempty"`
	}{}

	err = json.Unmarshal(resBody, &apiResponse)
	assert.Nil(err)
	assert.Equal(2, len(apiResponse.Data))

	//the first one should be poll5
	poll5Res := apiResponse.Data[0]
	assert.Equal(poll5Res.ID, poll5.ID)
	assert.Equal(poll5Res.Title, poll5.Title)

	slices.Contains(poll5Res.Options, "op9")
	slices.Contains(poll5Res.Options, "op10")

	slices.Contains(poll5Res.Tags, "tag9")
	slices.Contains(poll5Res.Tags, "tag10")

	//the second one should be poll6
	poll6Res := apiResponse.Data[1]
	assert.Equal(poll6Res.ID, poll6.ID)
	assert.Equal(poll6Res.Title, poll6.Title)

	slices.Contains(poll6Res.Options, "op11")
	slices.Contains(poll6Res.Options, "op12")

	slices.Contains(poll6Res.Tags, "tag11")
	slices.Contains(poll6Res.Tags, "tag12")

	//page 3 must be empty
	req, err = http.NewRequest("GET", fmt.Sprintf("%s/api/v1/polls?user_id=%d&page=%d&page_size=%d", app.address, 1, 2, 2), nil)
	assert.Nil(err)
	res, err = http.DefaultClient.Do(req)
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	resBody, err = io.ReadAll(res.Body)
	assert.Nil(err)
	apiResponse = struct {
		Success bool                  `json:"success" example:"true"`
		Message string                `json:"message,omitempty"`
		Data    []poll.SinglePollList `json:"data,omitempty"`
	}{}
	err = json.Unmarshal(resBody, &apiResponse)
	assert.Nil(err)
	assert.Equal(0, len(apiResponse.Data))

	//page 1 filter by tag
	req, err = http.NewRequest("GET", fmt.Sprintf("%s/api/v1/polls?user_id=%d&page=%d&page_size=%d&tag=%s", app.address, 1, 0, 2, "tag12"), nil)
	assert.Nil(err)
	res, err = http.DefaultClient.Do(req)
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	resBody, err = io.ReadAll(res.Body)
	assert.Nil(err)
	apiResponse = struct {
		Success bool                  `json:"success" example:"true"`
		Message string                `json:"message,omitempty"`
		Data    []poll.SinglePollList `json:"data,omitempty"`
	}{}

	err = json.Unmarshal(resBody, &apiResponse)
	assert.Nil(err)
	assert.Equal(1, len(apiResponse.Data))
	poll6Res = apiResponse.Data[0]
	assert.Equal(poll6Res.ID, poll6.ID)
	assert.Equal(poll6Res.Title, poll6.Title)

	slices.Contains(poll6Res.Options, "op11")
	slices.Contains(poll6Res.Options, "op12")

	slices.Contains(poll6Res.Tags, "tag11")
	slices.Contains(poll6Res.Tags, "tag12")
}
