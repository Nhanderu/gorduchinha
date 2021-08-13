package server

import (
	"github.com/paemuri/gorduchinha/cmd/api/server/handler"
	"github.com/paemuri/gorduchinha/cmd/api/server/middleware"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type r struct {
	router *router.Router
	prefix string
	mw     []middleware.RequestMiddleware
}

func newRouter(corsMiddleware middleware.RequestMiddleware) *r {

	router := router.New()
	router.HandleOPTIONS = true
	// TODO: organize CORS
	router.GlobalOPTIONS = corsMiddleware(func(ctx *fasthttp.RequestCtx) {})
	router.HandleMethodNotAllowed = true
	router.MethodNotAllowed = handler.MethodNotAllowed()
	router.NotFound = handler.PageNotFound()
	router.PanicHandler = handler.Panic()

	return &r{
		router: router,
		mw:     make([]middleware.RequestMiddleware, 0),
	}
}

func (root *r) requestHandler() fasthttp.RequestHandler {
	return root.router.Handler
}

func (root *r) group(prefix string, mws ...middleware.RequestMiddleware) *r {
	return &r{
		router: root.router,
		prefix: root.prefix + prefix,
		mw:     append(root.mw, mws...),
	}
}

func (root *r) handle(method, path string, handler fasthttp.RequestHandler, mws ...middleware.RequestMiddleware) {
	p := root.prefix + path
	handler = middleware.Use(handler, root.mw...)
	handler = middleware.Use(handler, mws...)
	root.router.Handle(method, p, handler)
}
