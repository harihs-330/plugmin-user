package apikit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"user/errortools"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status  string         `json:"status,omitempty"`
	Message string         `json:"message,omitempty"`
	Data    map[string]any `json:"data,omitempty"`
	Error   error          `json:"error,omitempty"`
}

func (resp *APIResponse) AddData(key string, data any) {
	resp.Data[key] = data
}

func SuccessGenerator(data any, key, msg string) APIResponse {
	res := APIResponse{
		Status:  "success",
		Message: msg,
	}
	if data != nil {
		res.Data = map[string]any{
			key: data,
		}
	}

	return res
}

func ErrorGenerator(ctx *gin.Context, err error) {
	res := APIResponse{
		Status: "failure",
		Error:  err,
	}
	code := http.StatusBadRequest
	SrvErr := &errortools.Error{}
	if errors.As(err, &SrvErr) {
		if code = SrvErr.Code.HTTPCode(); code == 0 {
			code = http.StatusBadRequest
		}
	}
	ctx.AbortWithStatusJSON(code, res)
}

// function to handle binding error
func BindingError(ctx *gin.Context, err error) {
	var (
		unmarshalTypeErr *json.UnmarshalTypeError
		syntaxErr        *json.SyntaxError
	)
	resErr := errortools.Init()
	switch {
	case errors.As(err, &unmarshalTypeErr):
		field := unmarshalTypeErr.Field
		expectedType := unmarshalTypeErr.Type
		actualValue := unmarshalTypeErr.Value
		resErr.AddValidationError(field, errortools.InvalidFieldType, expectedType, actualValue)
	case errors.As(err, &syntaxErr):
		resErr = errortools.New(errortools.BindingError,
			errortools.WithDetail(fmt.Sprintf("%s : %d", syntaxErr, syntaxErr.Offset)))
	default:
		resErr = errortools.New(errortools.BindingError, errortools.WithDetail(err.Error()))
	}
	ErrorGenerator(ctx, resErr)
}
