package route

import (
	"net/http"

	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/router"
	"github.com/go-openapi/runtime/middleware"
)

func SwaggerRoute(appRouter router.Router) {
	swaggerRedocOpts := middleware.RedocOpts{Path: "api/docs", SpecURL: "/swagger.yaml"}
	swaggerRedoc := middleware.Redoc(swaggerRedocOpts, nil)

	appRouter.Get("/api/docs", swaggerRedoc.ServeHTTP)
	appRouter.Get("/swagger.yaml", http.FileServer(http.Dir("./")).ServeHTTP)
}
