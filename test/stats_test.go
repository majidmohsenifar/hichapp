package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/majidmohsenifar/hichapp/service/statistic"
	"github.com/stretchr/testify/assert"
)

func TestStats_Successful(t *testing.T) {
	assert := assert.New(t)
	app := spawn_app()
	defer app.close()

	poll := app.CreatePollWithOptionsAndTags(t.Context(), "poll1", []string{"op1", "op2"}, []string{"tag1", "tag2"})
	options, err := app.repo.GetOptionsByPollID(t.Context(), app.db, poll.ID)
	assert.Nil(err)
	for i := 0; i < 10; i++ {
		app.repo.IncrementOptionVoteCount(t.Context(), app.db, options[0].ID)
	}

	app.repo.IncrementOptionVoteCount(t.Context(), app.db, options[1].ID)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/polls/%d/stats", app.address, poll.ID), nil)
	assert.Nil(err)
	res, err := http.DefaultClient.Do(req)
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	assert.Nil(err)
	apiResponse := struct {
		Success bool                  `json:"success" example:"true"`
		Message string                `json:"message,omitempty"`
		Data    statistic.StatsResult `json:"data,omitempty"`
	}{}

	err = json.Unmarshal(resBody, &apiResponse)
	assert.Nil(err)
	assert.Equal(2, len(apiResponse.Data.Votes))

	for _, d := range apiResponse.Data.Votes {
		switch d.Option {
		case "op1":
			assert.Equal(int32(10), d.Count)
		case "op2":
			assert.Equal(int32(1), d.Count)
		default:
			assert.FailNow("unexpected option")

		}
	}
}
