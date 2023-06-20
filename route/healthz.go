package route

import (
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/controller"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/usecase"
)

func HealthzRoute(params *RouteParameters) {
	healthzUseCase := usecase.NewHealthz(params.Repository, params.Cache)
	controllerHealthz := controller.NewHealthz(params.Log, healthzUseCase)

	params.AppRouter.Get("/api/healthz", controllerHealthz.Check)
}
