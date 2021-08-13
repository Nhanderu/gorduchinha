package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/paemuri/gorduchinha/app/constant"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

const (
	ErrorCodeContextKey    = "error-code"
	ErrorMessageContextKey = "error-message"
)

var (
	errorStatusMap = map[string]int{
		constant.ErrorCodePageNotFound:        http.StatusNotFound,
		constant.ErrorCodeMethodNotAllowed:    http.StatusMethodNotAllowed,
		constant.ErrorCodeCacheMiss:           http.StatusInternalServerError,
		constant.ErrorCodeTooManyRequests:     http.StatusTooManyRequests,
		constant.ErrorCodeInvalidRequestBody:  http.StatusBadRequest,
		constant.ErrorCodeRequestBodyTooLarge: http.StatusRequestEntityTooLarge,
		constant.ErrorCodeEntityNotFound:      http.StatusNotFound,
		constant.ErrorCodeInvalidQueryKey:     http.StatusForbidden,
		constant.ErrorCodeInternal:            http.StatusInternalServerError,
	}
)

type resultWrapper struct {
	Success bool                 `json:"success"`
	Data    interface{}          `json:"data,omitempty"`
	Errors  []resultWrapperError `json:"errors,omitempty"`
}

type resultWrapperError struct {
	Code  string `json:"code"`
	Field string `json:"field,omitempty"`
}

func HandleError(ctx *fasthttp.RequestCtx, err error) {

	if err == nil {
		return
	}

	ctx.SetUserValue(ErrorMessageContextKey, err.Error())
	err = errors.Cause(err)

	switch e := err.(type) {

	case constant.AppError:
		statusCode, ok := errorStatusMap[e.Code]
		if !ok {
			statusCode = http.StatusInternalServerError
		}

		respondError(ctx, statusCode, e.Code, e.Field, e.Error())
		return

	default:
		respondError(
			ctx,
			http.StatusInternalServerError,
			constant.ErrorCodeInternal,
			"",
			fmt.Sprintf("unmapped error: %s", e.Error()),
		)
		return
	}

}

func respondError(ctx *fasthttp.RequestCtx, status int, code string, field string, message string) {
	ctx.SetUserValue(ErrorCodeContextKey, code)
	ctx.SetUserValue(ErrorMessageContextKey, message)

	errors := make([]resultWrapperError, 0)
	errors = append(errors, resultWrapperError{
		Code:  code,
		Field: field,
	})

	respondJSON(ctx, status, resultWrapper{
		Success: false,
		Errors:  errors,
	})
}

func respondOK(ctx *fasthttp.RequestCtx, data interface{}) {
	respondJSON(ctx, http.StatusOK, resultWrapper{
		Success: true,
		Data:    data,
	})
}

func respondJSON(ctx *fasthttp.RequestCtx, code int, result interface{}) {
	ctx.Response.Header.Add("Content-Encoding", "gzip")
	ctx.Response.Header.Add("X-XSS-Protection", "1; mode=block")
	ctx.Response.Header.Add("X-Content-Type-Options", "nosniff")
	ctx.Response.Header.Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	ctx.SetContentType("app/json; charset=UTF-8")
	ctx.SetStatusCode(code)
	b, _ := json.Marshal(result)
	fasthttp.WriteGzipLevel(ctx, b, fasthttp.CompressBestSpeed)
}
