package handler

import (
	"fmt"

	"github.com/paemuri/gorduchinha/app/constant"
	"github.com/valyala/fasthttp"
)

func MethodNotAllowed() func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		HandleError(ctx, constant.NewErrorMethodNotAllowed())
	}
}

func PageNotFound() func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		HandleError(ctx, constant.NewErrorPageNotFound())
	}
}

func Panic() func(*fasthttp.RequestCtx, interface{}) {
	return func(ctx *fasthttp.RequestCtx, data interface{}) {
		HandleError(ctx, fmt.Errorf("[PANIC RECOVERED]\n%v", data))
	}
}
