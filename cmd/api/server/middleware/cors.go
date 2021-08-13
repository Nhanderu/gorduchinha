package middleware

import (
	"github.com/paemuri/fastcors"
	"github.com/valyala/fasthttp"
)

func CORS(serverAuthClientsURLs []string) RequestMiddleware {
	return RequestMiddleware(fastcors.New(
		fastcors.SetAllowedOrigins(serverAuthClientsURLs),
		fastcors.SetAllowedHeaders([]string{fasthttp.HeaderOrigin, fasthttp.HeaderAccept, fasthttp.HeaderContentType}),
		fastcors.SetAllowedMethods([]string{fasthttp.MethodGet, fasthttp.MethodPost, fasthttp.MethodPut}),
		fastcors.SetAllowCredentials(false),
		fastcors.SetMaxAge(60*60*24*30),
	))
}
