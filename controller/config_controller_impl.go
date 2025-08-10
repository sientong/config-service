package controller

import (
	"config-service/helper"
	"config-service/model/web"
	"config-service/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ConfigControllerImpl struct {
	ConfigService service.ConfigService
}

func NewConfigController(configService service.ConfigService) ConfigController {
	return &ConfigControllerImpl{
		ConfigService: configService,
	}
}

func (controller *ConfigControllerImpl) CreateConfig(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	schema := params.ByName("schema")
	name := params.ByName("name")

	var rawData map[string]interface{}
	helper.ReadFromRequestBody(request, &rawData)
	configCreateRequest := web.ConfigCreateRequest{
		Schema: schema,
		Name:   name,
		Data:   rawData,
	}

	configResponse := controller.ConfigService.CreateConfig(request.Context(), configCreateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   configResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ConfigControllerImpl) UpdateConfig(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	schema := params.ByName("schema")
	name := params.ByName("name")

	var rawData map[string]interface{}
	helper.ReadFromRequestBody(request, &rawData)
	configUpdateRequest := web.ConfigUpdateRequest{
		Schema: schema,
		Name:   name,
		Data:   rawData,
	}

	configResponse := controller.ConfigService.UpdateConfig(request.Context(), configUpdateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   configResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ConfigControllerImpl) RollbackConfig(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	schema := params.ByName("schema")
	name := params.ByName("name")

	var rawData map[string]interface{}
	var version int
	helper.ReadFromRequestBody(request, &rawData)
	if v, ok := rawData["version"].(float64); ok {
		version = int(v)
	} else {
		helper.PanicIfError(helper.ValidationError{Msg: "missing version"})
	}

	configUpdateRequest := web.ConfigRollbackRequest{
		Schema:  schema,
		Name:    name,
		Version: version,
	}

	configResponse := controller.ConfigService.RollbackConfig(request.Context(), configUpdateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   configResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ConfigControllerImpl) FetchConfig(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	schema := params.ByName("schema")
	name := params.ByName("name")

	var version int
	if request.GetBody != nil {
		var rawData map[string]interface{}
		helper.ReadFromRequestBody(request, &rawData)
		if v, ok := rawData["version"].(float64); ok {
			version = int(v)
		} else {
			version = 0
		}
	} else {
		version = 0
	}

	ConfigFetchRequest := web.ConfigFetchRequest{
		Schema:  schema,
		Name:    name,
		Version: version,
	}

	configResponse := controller.ConfigService.FetchConfig(request.Context(), ConfigFetchRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   configResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *ConfigControllerImpl) ListVersions(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	schema := params.ByName("schema")
	name := params.ByName("name")

	configListVersionsRequest := web.ConfigListVersionsRequest{
		Schema: schema,
		Name:   name,
	}

	configResponse := controller.ConfigService.ListVersions(request.Context(), configListVersionsRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   configResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}
