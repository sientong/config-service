package service

import (
	"config-service/helper"
	"config-service/model/domain"
	"config-service/model/web"
	"config-service/repository"
	"context"
	"database/sql"

	"github.com/go-playground/validator"
)

type ConfigServiceImpl struct {
	ConfigRepository repository.ConfigRepository
	DB               *sql.DB
	Validate         *validator.Validate
}

func NewConfigService(configRepository repository.ConfigRepository, DB *sql.DB, validate *validator.Validate) ConfigService {
	return &ConfigServiceImpl{
		ConfigRepository: configRepository,
		DB:               DB,
		Validate:         validate,
	}
}

func (service *ConfigServiceImpl) CreateConfig(ctx context.Context, schema, name string, request web.ConfigCreateRequest) web.ConfigResponse {
	// Validate incoming request payload
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	// Validate against schema
	helper.ValidateAgainstSchema(schema, request.Data)

	// Start transaction
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// Create domain model
	configRecord := domain.ConfigRecord{
		Schema: schema,
		Name:   name,
		Data:   request.Data, // assuming Data is a map or json.RawMessage
	}

	// Check whether config name exist
	latest := service.ConfigRepository.GetLatest(ctx, tx, configRecord)
	if latest.Schema == schema && latest.Name == name && latest.Version > 0 {
		helper.PanicIfError(helper.ValidationError{Msg: "config name already exist"})
	}

	newVersion := 1
	configRecord.Version = newVersion

	// Save new config version
	configRecord = service.ConfigRepository.CreateNewVersion(ctx, tx, configRecord)

	// Map domain model to web response
	return helper.ToConfigResponse(configRecord)
}

func (service *ConfigServiceImpl) UpdateConfig(ctx context.Context, schema, name string, request web.ConfigUpdateRequest) web.ConfigResponse {
	// Validate incoming request payload
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	// Validate against schema
	helper.ValidateAgainstSchema(schema, request.Data)

	// Start transaction
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// Create domain model
	configRecord := domain.ConfigRecord{
		Schema: schema,
		Name:   name,
		Data:   request.Data, // assuming Data is a map or json.RawMessage
	}

	// Check whether config name exist
	latest := service.ConfigRepository.GetLatest(ctx, tx, configRecord)
	if latest.Version == 0 {
		helper.PanicIfError(helper.ValidationError{Msg: "config name doesn't exist"})
	}

	newVersion := latest.Version + 1
	configRecord.Version = newVersion

	// Save new config version
	configRecord = service.ConfigRepository.CreateNewVersion(ctx, tx, configRecord)

	return helper.ToConfigResponse(configRecord)
}

func (service *ConfigServiceImpl) FetchConfig(ctx context.Context, schema, name string, request web.ConfigFetchRequest) web.ConfigResponse {
	// Validate incoming request payload
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	// Validate schema existence
	helper.ValidateSchemaExistence(schema)

	// Start transaction
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// Create domain model
	configRecord := domain.ConfigRecord{
		Schema: schema,
		Name:   name,
	}

	var fetchData domain.ConfigRecord
	if *request.Version == 0 {
		configRecord.Version = *request.Version
		fetchData := service.ConfigRepository.GetLatest(ctx, tx, configRecord)
		if fetchData.Version == 0 {
			helper.PanicIfError(helper.ValidationError{Msg: "no data found"})
		}
	} else {
		fetchData = service.ConfigRepository.GetByVersion(ctx, tx, configRecord)
		if fetchData.Version == 0 && fetchData.Name == "" && fetchData.Schema == "" {
			helper.PanicIfError(helper.ValidationError{Msg: "config name doesn't exist"})
		}
	}

	return helper.ToConfigResponse(fetchData)
}

func (service *ConfigServiceImpl) RollbackConfig(ctx context.Context, schema, name string, request web.ConfigRollbackRequest) web.ConfigResponse {
	// Validate incoming request payload
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	// Validate schema existence
	helper.ValidateSchemaExistence(schema)

	// Start transaction
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// Create domain model
	configRecord := domain.ConfigRecord{
		Schema:  schema,
		Name:    name,
		Version: request.Version, // assuming Data is a map or json.RawMessage
	}

	// Check whether fetched version exist
	fetchData := service.ConfigRepository.GetByVersion(ctx, tx, configRecord)
	if fetchData.Version == 0 && fetchData.Name == "" && fetchData.Schema == "" {
		helper.PanicIfError(helper.ValidationError{Msg: "fetched version doesn't exist"})
	}

	// Get latest version
	latest := service.ConfigRepository.GetLatest(ctx, tx, configRecord)
	if latest.Version == 0 {
		helper.PanicIfError(helper.ValidationError{Msg: "invalid request to get latest data"})
	}

	// Rollback to specified version
	fetchData.Version = latest.Version + 1
	rollbackData := service.ConfigRepository.CreateNewVersion(ctx, tx, fetchData)

	return helper.ToConfigResponse(rollbackData)
}

func (service *ConfigServiceImpl) ListVersions(ctx context.Context, schema, name string) web.ConfigResponses {

	// Validate schema existence
	helper.ValidateSchemaExistence(schema)

	// Start transaction
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// Create domain model
	configRecord := domain.ConfigRecord{
		Schema: schema,
		Name:   name,
	}

	configRecords := service.ConfigRepository.ListVersions(ctx, tx, configRecord)

	return helper.ToConfigResponses(schema, name, configRecords)
}
