package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ConfigController interface {
	CreateConfig(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateConfig(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	RollbackConfig(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FetchConfig(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	ListVersions(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
