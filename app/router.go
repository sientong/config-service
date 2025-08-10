package app

import (
	"config-service/controller"
	"config-service/exception"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(configController controller.ConfigController) *httprouter.Router {
	router := httprouter.New()

	router.POST("/configs/:schema/:name", configController.CreateConfig)
	router.PUT("/configs/:schema/:name", configController.UpdateConfig)
	router.POST("/configs/:schema/:name/rollback", configController.RollbackConfig)
	router.GET("/configs/:schema/:name", configController.FetchConfig)
	router.GET("/configs/:schema/:name/versions", configController.ListVersions)

	router.PanicHandler = exception.ErrorHandler

	return router
}
