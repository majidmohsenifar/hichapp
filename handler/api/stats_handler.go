package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/majidmohsenifar/hichapp/service/statistic"
)

type StatsHandler struct {
	statsService *statistic.Service
	validator    *validator.Validate
}

//	@Summary		list of poll stats
//	@Description    list of poll stats
//	@Tags			PollStats
//	@ID				PollStats
//	@Produce		json
//
// @Param		id path int true "Vote ID"
//
// @Success		 200	{object}	ResponseSuccess{data=statistic.StatsResult}
// @Failure		500		{object}	ResponseFailure
// @Router			/api/v1/polls/{id}/stats [get]
func (h *StatsHandler) Stats(c *gin.Context) {
	pollIDStr := c.Param("id")
	//convert pollIDStr to int64 and return badRequest error
	pollID, err := strconv.ParseInt(pollIDStr, 10, 64)
	if err != nil {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			err)
		return
	}
	res, err := h.statsService.GetStatsForPoll(c, pollID)
	if err != nil {
		MakeErrorResponseWithCode(c.Writer, http.StatusInternalServerError, err)
		return
	}
	MakeSuccessResponse(c.Writer, res, "")
}

func NewStatsHandler(
	statsService *statistic.Service,
	validator *validator.Validate,
) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
		validator:    validator,
	}
}
