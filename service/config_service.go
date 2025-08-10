package service

import (
	"config-service/model/web"
	"context"
)

type ConfigService interface {
	CreateConfig(ctx context.Context, request web.ConfigCreateRequest) web.ConfigResponse
	UpdateConfig(ctx context.Context, request web.ConfigUpdateRequest) web.ConfigResponse
	RollbackConfig(ctx context.Context, request web.ConfigRollbackRequest) web.ConfigResponse
	FetchConfig(ctx context.Context, request web.ConfigFetchRequest) web.ConfigResponse
	ListVersions(ctx context.Context, request web.ConfigListVersionsRequest) web.ConfigResponses
}
