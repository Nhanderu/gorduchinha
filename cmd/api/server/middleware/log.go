package middleware

import (
	"strconv"
	"time"

	"github.com/paemuri/gorduchinha/app/logger"
	"github.com/paemuri/gorduchinha/cmd/api/server/handler"
	"github.com/valyala/fasthttp"
)

func Logger(log logger.Logger) RequestMiddleware {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {

			start := time.Now()
			next(ctx)
			end := time.Now()
			latency := end.Sub(start)

			log.WithFields(map[string]interface{}{

				"conn-id":    ctx.ConnID(),
				"request-id": ctx.ID(),
				"context":    ctx.String(),

				"host":       ctx.Host(),
				"method":     string(ctx.Method()),
				"path":       string(ctx.Path()),
				"query":      string(ctx.URI().QueryString()),
				"user-agent": string(ctx.UserAgent()),
				"referer":    string(ctx.Referer()),
				"remote-ip":  ctx.RemoteIP().String(),

				"status":    ctx.Response.StatusCode(),
				"bytes-in":  ctx.Request.Header.ContentLength(),
				"bytes-out": len(ctx.Response.Body()),

				"type":          "request-handle",
				"time":          end.Format(time.RFC3339),
				"latency":       strconv.FormatInt(latency.Nanoseconds()/1000, 10),
				"latency-human": latency.String(),
				"error-code":    ctx.UserValue(handler.ErrorCodeContextKey),
				"error-message": ctx.UserValue(handler.ErrorMessageContextKey),
			}).Infof("Handled request.")
		}
	}
}
