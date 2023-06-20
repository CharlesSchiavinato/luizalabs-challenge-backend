package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/model"
	logger "github.com/CharlesSchiavinato/luizalabs-challenge-backend/service/logger"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/usecase"
	"github.com/hashicorp/go-hclog"
)

type Healthz struct {
	Log            hclog.Logger
	HealthzUseCase usecase.Healthz
}

func NewHealthz(log hclog.Logger, healthzUseCase usecase.Healthz) *Healthz {
	return &Healthz{
		Log:            log,
		HealthzUseCase: healthzUseCase,
	}
}

func (controllerHealthz *Healthz) Check(rw http.ResponseWriter, req *http.Request) {
	const handlerLogTitle = "Health Check"
	err := controllerHealthz.HealthzUseCase.CheckRepository()

	healthzModel := &model.Healthz{}

	if err != nil {
		logger.LogErrorRequest(controllerHealthz.Log, req, handlerLogTitle, err)
		rw.WriteHeader(http.StatusFailedDependency)
		healthzModel.Database = err.Error()
	} else {
		healthzModel.Database = "OK"
	}

	err = controllerHealthz.HealthzUseCase.CheckCache()

	if err != nil {
		logger.LogErrorRequest(controllerHealthz.Log, req, handlerLogTitle, err)
		rw.WriteHeader(http.StatusFailedDependency)
		healthzModel.Cache = err.Error()
	} else {
		healthzModel.Cache = "OK"
	}

	json.NewEncoder(rw).Encode(healthzModel)
}
