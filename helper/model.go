package helper

import (
	"config-service/model/domain"
	"config-service/model/web"

	"github.com/xeipuuv/gojsonschema"
)

func ToConfigResponse(config domain.ConfigRecord) web.ConfigResponse {
	return web.ConfigResponse{
		Schema:    config.Schema,
		Name:      config.Name,
		Version:   config.Version,
		Data:      config.Data,
		CreatedAt: config.CreatedAt,
	}
}

func ToConfigResponses(schema string, name string, configRecords []domain.ConfigRecord) web.ConfigResponses {
	configResponses := make([]web.ConfigResponse, 0)
	for _, config := range configRecords {
		schema = config.Schema
		name = config.Name
		configResponses = append(configResponses, ToConfigResponse(config))
	}

	return web.ConfigResponses{
		Schema:         schema,
		Name:           name,
		ConfigVersions: configResponses,
	}
}

func ValidateSchemaExistence(schemaName string) string {
	schema, ok := domain.Schemas[schemaName]
	if !ok {
		PanicIfError(ValidationError{Msg: "unknown schema"})
	}
	return schema
}

func ValidateAgainstSchema(schemaName string, data map[string]interface{}) {

	schema := ValidateSchemaExistence(schemaName)
	schemaLoader := gojsonschema.NewStringLoader(schema)

	// Load data as JSON document
	docLoader := gojsonschema.NewGoLoader(data)

	res, err := gojsonschema.Validate(schemaLoader, docLoader)
	PanicIfError(err)

	if !res.Valid() {
		msgs := ""
		for _, e := range res.Errors() {
			msgs += e.String() + "; "
		}
		PanicIfError(ValidationError{Msg: "schema validation failed: " + msgs})
	}

}
