package controller

import (
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
