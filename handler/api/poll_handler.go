package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/majidmohsenifar/hichapp/service/poll"
	"github.com/majidmohsenifar/hichapp/service/tag"
)

type PollHandler struct {
	pollService *poll.Service
	validator   *validator.Validate
}

type CreatePollReq struct {
	Title   string   `json:"title" validate:"required"`
	Options []string `json:"options" validate:"required,min=2"`
	Tags    []string `json:"tags"`
}

type PollListReq struct {
	Page     int64  `form:"page" validate:"gte=0"`
	PageSize int16  `form:"page_size" validate:"gte=2,lte=100000"`
	Tag      string `form:"tag"`
	UserID   int64  `form:"user_id" validate:"required,gt=0"`
}

type VoteReq struct {
	UserID      int64 `json:"user_id" validate:"required,gt=0"`
	OptionIndex int8  `json:"option_index"`
}

//	@Summary		create poll
//	@Description    create poll
//	@Tags			Poll
//	@ID				CreatePoll
//	@Produce		json
//
// @Param		request body CreatePollReq true "Create-Poll-Params"
// @Success		200
// @Failure		500		{object}	ResponseFailure
// @Router			/api/v1/polls [post]
func (h *PollHandler) Create(c *gin.Context) {
	params := CreatePollReq{}
	err := c.ShouldBindJSON(&params)
	if err != nil {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			err,
		)
		return
	}
	err = h.validator.Struct(params)
	if err != nil {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			err)
		return
	}
	_, err = h.pollService.CreatePoll(c, poll.CreatePollParams(params))
	if err != nil {
		MakeErrorResponseWithCode(c.Writer, http.StatusInternalServerError, err)
		return
	}
	MakeSuccessResponseWithoutBody(c.Writer, http.StatusCreated)
}

//	@Summary		list of polls
//	@Description    list of polls
//	@Tags			Poll
//	@ID				PollList
//	@Produce		json
//
// @Param        page    query uint64 false  "Page"
// @Param        page_size    query uint64 false  "PageSize"
// @Param        tag query string false  "Tag"
// @Param        user_id query uint64 true "User-ID"
//
// @Success		 200	{object}	ResponseSuccess{data=[]poll.SinglePollList}
// @Failure		500		{object}	ResponseFailure
// @Router			/api/v1/polls [get]
func (h *PollHandler) List(c *gin.Context) {
	params := PollListReq{}
	err := c.ShouldBindQuery(&params)
	if err != nil {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			err,
		)
		return
	}
	err = h.validator.Struct(params)
	if err != nil {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			err)
		return
	}

	res, err := h.pollService.GetPollsList(c, poll.PollListFilters(params))
	if errors.Is(err, tag.ErrTagNotFound) {
		MakeErrorResponseWithCode(c.Writer, http.StatusNotFound, err)
		return
	}
	if err != nil {
		MakeErrorResponseWithCode(c.Writer, http.StatusInternalServerError, err)
		return
	}
	MakeSuccessResponse(c.Writer, res, "")

}

//	@Summary		vote to poll
//	@Description    vote poll
//	@Tags			Poll
//	@ID				Vote
//	@Produce		json
//
// @Param		request body VoteReq true "Vote-Params"
// @Param		id path int true "Vote ID"
// @Success		200
// @Failure		500		{object}	ResponseFailure
// @Router			/api/v1/polls/{id}/vote [post]
func (h *PollHandler) Vote(c *gin.Context) {
	params := VoteReq{}
	err := c.ShouldBindJSON(&params)
	if err != nil {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			err,
		)
		return
	}
	err = h.validator.Struct(params)
	if err != nil {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			err)
		return
	}
	if params.OptionIndex < 0 {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			errors.New("option index should be non-negative"))
		return
	}
	pollIDStr := c.Param("id")
	pollID, err := strconv.ParseInt(pollIDStr, 10, 64)
	if err != nil {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			err)
		return
	}

	srvParams := poll.VoteOrSkipParams{
		PollID:      pollID,
		OptionIndex: params.OptionIndex,
		UserID:      params.UserID,
		IsSkipped:   false,
	}
	err = h.pollService.VoteOrSkip(c, srvParams)
	if errors.Is(err, poll.ErrPollNotFound) {
		MakeErrorResponseWithCode(c.Writer, http.StatusNotFound, err)
		return
	}
	if errors.Is(err, poll.ErrUserAlreadyVotedOrSkipped) {
		MakeErrorResponseWithCode(c.Writer, http.StatusUnprocessableEntity, err)
		return
	}
	if errors.Is(err, poll.ErrUserNotAllowedToVote) {
		MakeErrorResponseWithCode(c.Writer, http.StatusUnprocessableEntity, err)
		return
	}
	if errors.Is(err, poll.ErrInvalidOptionIndex) {
		MakeErrorResponseWithCode(c.Writer, http.StatusBadRequest, err)
		return
	}
	if err != nil {
		MakeErrorResponseWithCode(c.Writer, http.StatusInternalServerError, err)
		return
	}
	MakeSuccessResponseWithoutBody(c.Writer, http.StatusCreated)
}

//	@Summary		skip poll
//	@Description    skip poll
//	@Tags			Poll
//	@ID				Skip
//	@Produce		json
//
// @Param		request body VoteReq true "Vote-Params"
// @Param		id path int true "Vote ID"
// @Success		200
// @Failure		500		{object}	ResponseFailure
// @Router			/api/v1/polls/{id}/skip [post]
func (h *PollHandler) Skip(c *gin.Context) {
	params := VoteReq{}
	err := c.ShouldBindJSON(&params)
	if err != nil {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			err,
		)
		return
	}
	err = h.validator.Struct(params)
	if err != nil {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			err)
		return
	}
	pollIDStr := c.Param("id")
	pollID, err := strconv.ParseInt(pollIDStr, 10, 64)
	if err != nil {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			err)
		return
	}

	srvParams := poll.VoteOrSkipParams{
		PollID:      pollID,
		OptionIndex: params.OptionIndex,
		UserID:      params.UserID,
		IsSkipped:   true,
	}
	err = h.pollService.VoteOrSkip(c, srvParams)
	if errors.Is(err, poll.ErrPollNotFound) {
		MakeErrorResponseWithCode(c.Writer, http.StatusNotFound, err)
		return
	}
	if errors.Is(err, poll.ErrUserAlreadyVotedOrSkipped) {
		MakeErrorResponseWithCode(c.Writer, http.StatusUnprocessableEntity, err)
		return
	}
	if err != nil {
		MakeErrorResponseWithCode(c.Writer, http.StatusInternalServerError, err)
		return
	}
	MakeSuccessResponseWithoutBody(c.Writer, http.StatusCreated)

}

func NewPollHandler(
	pollService *poll.Service,
	validator *validator.Validate,
) *PollHandler {
	return &PollHandler{
		pollService: pollService,
		validator:   validator,
	}
}
