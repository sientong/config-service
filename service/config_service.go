package service

import (
	"config-service/model/web"
	"context"
)

type ConfigService interface {
	CreateConfig(ctx context.Context, schema, name string, request web.ConfigCreateRequest) web.ConfigResponse
	UpdateConfig(ctx context.Context, schema, name string, request web.ConfigUpdateRequest) web.ConfigResponse
	RollbackConfig(ctx context.Context, schema, name string, request web.ConfigRollbackRequest) web.ConfigResponse
	FetchConfig(ctx context.Context, schema, name string, request web.ConfigFetchRequest) web.ConfigResponse
	ListVersions(ctx context.Context, schema, name string) web.ConfigResponses
}
