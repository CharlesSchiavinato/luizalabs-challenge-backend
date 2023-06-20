package route

import (
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/controller"
	"github.com/CharlesSchiavinato/luizalabs-challenge-backend/usecase"
)

func OrderRoute(params *RouteParameters) {
	usecaseOrder := usecase.NewOrder(params.Repository, params.Cache)
	controllerOrder := controller.NewOrder(params.Log, usecaseOrder)

	pathApiOrder := "/api/order"
	paramID := params.AppRouter.PathFormat("/%s", "order_id")

	params.AppRouter.Get(pathApiOrder+paramID, controllerOrder.GetDetailsByOrderID)
	params.AppRouter.Get(pathApiOrder, controllerOrder.ListDetails)

	params.AppRouter.Post(pathApiOrder+"/legacy/import", controllerOrder.LegacyImport)
}
