package router

import (
	"fmt"
	"net/http"

	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/controller"
	"github.com/julienschmidt/httprouter"
)

type HttpRouter struct {
	*httprouter.Router
}

func NewHttpRouter() Router {
	router := &HttpRouter{
		httprouter.New(),
	}

	return router
}

func (hr *HttpRouter) Serve() http.Handler {
	hr.NotFound = http.HandlerFunc(controller.NewNotFound().NotFound)

	return hr
}

func (hr *HttpRouter) Get(path string, handle Handle) {
	hr.GET(path, handleAdapter(http.HandlerFunc(handle)))
}

func (hr *HttpRouter) Post(path string, handle Handle) {
	hr.POST(path, handleAdapter(http.HandlerFunc(handle)))
}

func (hr *HttpRouter) Put(path string, handle Handle) {
	hr.PUT(path, handleAdapter(http.HandlerFunc(handle)))
}

func (hr *HttpRouter) Delete(path string, handle Handle) {
	hr.DELETE(path, handleAdapter(http.HandlerFunc(handle)))
}

func (hr *HttpRouter) PathFormat(format string, a ...any) string {
	argsPath := []any{}

	for _, arg := range a {
		argsPath = append(argsPath, fmt.Sprintf(":%s", arg))
	}

	return fmt.Sprintf(format, argsPath...)
}

func handleAdapter(h http.Handler) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		// // Take the context out from the request
		// ctx := req.Context()

		// // Get new context with key-value "params" -> "httprouter.Params"
		// ctx = context.WithValue(ctx, "params", ps)

		// // Get new http.Request with the new context
		// req = req.WithContext(ctx)

		// Call your original http.Handler
		h.ServeHTTP(rw, req)
	}
}
