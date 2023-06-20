package router

import (
	"fmt"
	"net/http"

	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/controller"
	"github.com/gorilla/mux"
)

var muxDispatcher = mux.NewRouter()

type MuxRouter struct{}

func NewMuxRouter() Router {
	return &MuxRouter{}
}

func (*MuxRouter) Serve() http.Handler {
	muxDispatcher.NotFoundHandler = http.HandlerFunc(controller.NewNotFound().NotFound)

	return muxDispatcher
}

func (*MuxRouter) Get(path string, handle Handle) {
	muxDispatcher.HandleFunc(path, handle).Methods(http.MethodGet)
}

func (*MuxRouter) Post(path string, handle Handle) {
	muxDispatcher.HandleFunc(path, handle).Methods(http.MethodPost)
}

func (*MuxRouter) Put(path string, handle Handle) {
	muxDispatcher.HandleFunc(path, handle).Methods(http.MethodPut)
}

func (*MuxRouter) Delete(path string, handle Handle) {
	muxDispatcher.HandleFunc(path, handle).Methods(http.MethodDelete)
}

func (*MuxRouter) PathFormat(format string, a ...any) string {
	argsPath := []any{}

	for _, arg := range a {
		argsPath = append(argsPath, fmt.Sprintf("{%s}", arg))
	}

	return fmt.Sprintf(format, argsPath...)

}
