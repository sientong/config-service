package controller

import (
	"config-service/model/domain"
	"config-service/model/web"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type SchemaControllerImpl struct{}

func NewSchemaController() SchemaController {
	return &SchemaControllerImpl{}
}

// GetSchema godoc
// @Summary Get JSON Schema by name
// @Description Returns a schema from the loaded set
// @Tags schemas
// @Produce json
// @Param name path string true "Schema Name"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {string} string "Schema not found"
// @Router /schemas/{name} [get]
func (controller *SchemaControllerImpl) GetSchema(c *gin.Context) {
	name := c.Param("name")
	filePath := fmt.Sprintf("%s/%s", "schemas", name)

	data, err := os.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("schema file %s not found", name),
		})
		return
	}

	// Serve the file contents as JSON/YAML depending on your file type
	c.Data(http.StatusOK, "application/json", data)
}

// ListSchemas godoc
// @Summary List of loaded schemas
// @Description Returns schema list from loaded set
// @Tags schemas
// @Produce json
// @Success 200 {object} web.SchemaResponse
// @Router /schemas [get]
func (controller *SchemaControllerImpl) ListSchemas(c *gin.Context) {

	var schemas []web.SchemaResponse
	for schema := range domain.Schemas {
		SchemaResponse := web.SchemaResponse{
			Name:      schema,
			Path:      schema + ".json",
			Directory: domain.SchemaDir + "/" + schema + ".json",
		}
		schemas = append(schemas, SchemaResponse)
	}

	c.JSON(http.StatusOK, schemas)
}
