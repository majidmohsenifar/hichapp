package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/majidmohsenifar/hichapp/handler/api"
	"github.com/stretchr/testify/assert"
)

func TestCreatePoll_InvalidInputs(t *testing.T) {
	assert := assert.New(t)
	app := spawn_app()
	defer app.close()

	tests := []api.CreatePollReq{
		{
			Title:   "",
			Options: []string{},
			Tags:    []string{},
		},
		{
			Title:   "title1",
			Options: []string{},
			Tags:    []string{},
		},
	}

	for _, test := range tests {
		params, err := json.Marshal(test)
		assert.Nil(err)
		body := bytes.NewReader(params)
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/polls", app.address), body)
		assert.Nil(err)
		res, err := http.DefaultClient.Do(req)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, res.StatusCode)
	}
}

func TestCreatePoll_Successful(t *testing.T) {
	assert := assert.New(t)
	app := spawn_app()
	defer app.close()
	params := api.CreatePollReq{
		Title:   "poll1",
		Options: []string{"op1", "op2"},
		Tags:    []string{"tag1", "tag2"},
	}
	paramsBytes, err := json.Marshal(params)
	assert.Nil(err)
	body := bytes.NewReader(paramsBytes)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/polls", app.address), body)
	assert.Nil(err)
	res, err := http.DefaultClient.Do(req)
	assert.Nil(err)
	assert.Equal(http.StatusCreated, res.StatusCode)

	//check the db for the created poll
	poll, err := app.repo.GetLastCreatedPoll(context.Background(), app.db)
	assert.Nil(err)
	assert.Equal(params.Title, poll.Title)

	//check the db for the created options
	options, err := app.repo.GetOptionsByPollID(context.Background(), app.db, poll.ID)
	assert.Nil(err)
	assert.Len(options, len(params.Options))
	assert.Equal(params.Options[0], options[0].Content)
	assert.Equal(params.Options[1], options[1].Content)

	//check the db for the created tags
	tags, err := app.repo.GetTagsByPollIDs(t.Context(), app.db, []int64{poll.ID})
	assert.Nil(err)
	assert.Len(tags, len(params.Tags))
	assert.Equal(params.Tags[0], tags[0].Name)
	assert.Equal(params.Tags[1], tags[1].Name)
}
