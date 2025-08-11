package controller

import (
	"config-service/helper"
	"config-service/model/web"
	"config-service/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConfigControllerImpl struct {
	configService service.ConfigService
}

func NewConfigController(configService service.ConfigService) ConfigController {
	return &ConfigControllerImpl{
		configService: configService,
	}
}

// CreateConfig godoc
// @Summary Create a new configuration
// @Description Create configuration with given schema and name
// @Tags configs
// @Accept json
// @Produce json
// @Param schema path string true "Schema name"
// @Param name path string true "Configuration name"
// @Param request body web.ConfigCreateRequest true "Config data"
// @Success 201 {object} web.ConfigResponse
// @Failure 400 {object} web.WebResponse
// @Router /configs/{schema}/{name} [post]
func (c *ConfigControllerImpl) CreateConfig(ctx *gin.Context) {
	schema := ctx.Param("schema")
	name := ctx.Param("name")

	var rawData map[string]interface{}
	err := ctx.ShouldBindJSON(&rawData)
	helper.PanicIfError(err)

	req := web.ConfigCreateRequest{
		Data: rawData,
	}

	result := c.configService.CreateConfig(ctx.Request.Context(), schema, name, req)

	ctx.JSON(http.StatusCreated, result)
}

// UpdateConfig godoc
// @Summary Update configuration
// @Description Update configuration with given schema and name
// @Tags configs
// @Accept json
// @Produce json
// @Param schema path string true "Schema name"
// @Param name path string true "Configuration name"
// @Param request body web.ConfigUpdateRequest true "Config data"
// @Success 200 {object} web.ConfigResponse
// @Failure 400 {object} web.WebResponse
// @Failure 404 {object} web.WebResponse
// @Router /configs/{schema}/{name} [put]
func (c *ConfigControllerImpl) UpdateConfig(ctx *gin.Context) {
	schema := ctx.Param("schema")
	name := ctx.Param("name")

	var rawData map[string]interface{}
	err := ctx.ShouldBindJSON(&rawData)
	helper.PanicIfError(err)

	req := web.ConfigUpdateRequest{
		Data: rawData,
	}

	result := c.configService.UpdateConfig(ctx.Request.Context(), schema, name, req)

	ctx.JSON(http.StatusOK, result)
}

// RollbackConfig godoc
// @Summary Rollback configuration to previous version
// @Tags configs
// @Produce json
// @Param schema path string true "Schema name"
// @Param name path string true "Configuration name"
// @Param request body web.ConfigRollbackRequest true "Config data"
// @Success 200 {object} web.ConfigResponse
// @Failure 500 {object} web.WebResponse
// @Failure 404 {object} web.WebResponse
// @Router /configs/{schema}/{name}/rollback [post]
func (c *ConfigControllerImpl) RollbackConfig(ctx *gin.Context) {
	schema := ctx.Param("schema")
	name := ctx.Param("name")

	var rawData map[string]interface{}
	var version int
	err := ctx.ShouldBindJSON(&rawData)
	helper.PanicIfError(err)
	if v, ok := rawData["version"].(float64); ok {
		version = int(v)
	} else {
		helper.PanicIfError(helper.ValidationError{Msg: "missing version"})
	}

	req := web.ConfigRollbackRequest{
		Version: version,
	}

	result := c.configService.RollbackConfig(ctx.Request.Context(), schema, name, req)

	ctx.JSON(http.StatusOK, result)
}

// FetchConfig godoc
// @Summary Fetch configuration
// @Tags configs
// @Produce json
// @Param schema path string true "Schema name"
// @Param name path string true "Configuration name"
// @Param request body web.ConfigFetchRequest false "Config data"
// @Success 200 {object} web.ConfigResponse
// @Failure 404 {object} web.WebResponse
// @Router /configs/{schema}/{name} [get]
func (c *ConfigControllerImpl) FetchConfig(ctx *gin.Context) {
	schema := ctx.Param("schema")
	name := ctx.Param("name")

	var version int

	// Check if request has a body
	if ctx.Request.ContentLength > 0 {
		var req web.ConfigFetchRequest
		err := ctx.ShouldBindJSON(&req)
		helper.PanicIfError(err)

		if req.Version != nil {
			version = *req.Version // user sent a value
		} else {
			version = 0 // user didn’t send version
		}
	} else {
		version = 0
	}

	req := web.ConfigFetchRequest{
		Version: &version,
	}

	result := c.configService.FetchConfig(ctx.Request.Context(), schema, name, req)

	ctx.JSON(http.StatusOK, result)
}

// ListVersions godoc
// @Summary List configuration versions
// @Tags configs
// @Produce json
// @Param schema path string true "Schema name"
// @Param name path string true "Configuration name"
// @Success 200 {array} web.ConfigResponses
// @Failure 404 {object} web.WebResponse
// @Router /configs/{schema}/{name}/versions [get]
func (c *ConfigControllerImpl) ListVersions(ctx *gin.Context) {
	schema := ctx.Param("schema")
	name := ctx.Param("name")

	result := c.configService.ListVersions(ctx.Request.Context(), schema, name)

	ctx.JSON(http.StatusOK, result)
}
