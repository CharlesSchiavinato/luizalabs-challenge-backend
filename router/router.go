package router

import (
	"net/http"
)

type Handle func(http.ResponseWriter, *http.Request)

type Router interface {
	Serve() http.Handler
	PathFormat(format string, a ...any) string
	Get(path string, handle Handle)
	Post(path string, handle Handle)
	Put(path string, handle Handle)
	Delete(path string, handle Handle)
}
