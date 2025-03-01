package api

import (
	"encoding/json"
	"net/http"
)

var (
	ContentTypeMessage     = "Content-Type"
	ApplicationJsonMessage = "application/json"
)

type ResponseSuccess struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ResponseFailure struct {
	Success bool      `json:"success" example:"false"`
	Error   ErrorCode `json:"error,omitempty"`
}

type ErrorCode struct {
	Code    int    `json:"code" example:"404"`
	Message string `json:"message" example:"item not found"`
}

func MakeSuccessResponse(w http.ResponseWriter, data interface{}, message string) {
	responseJson := ResponseSuccess{
		Success: true,
		Message: message,
		Data:    data,
	}
	jData, err := json.Marshal(&responseJson)
	if err != nil {
		MakeErrorResponseWithCode(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set(ContentTypeMessage, ApplicationJsonMessage)
	w.Write(jData)
}

func MakeSuccessResponseWithoutBody(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

func MakeErrorResponseWithCode(w http.ResponseWriter, code int, err error) {
	responseJson := ResponseFailure{
		Success: false,
		Error: ErrorCode{
			Code:    code,
			Message: err.Error(),
		},
	}
	jData, _ := json.Marshal(&responseJson)
	w.WriteHeader(code)
	w.Header().Set(ContentTypeMessage, ApplicationJsonMessage)
	w.Write(jData)
}
