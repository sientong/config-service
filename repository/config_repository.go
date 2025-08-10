package repository

import (
	"config-service/model/domain"
	"context"
	"database/sql"
)

type ConfigRepository interface {
	GetLatest(ctx context.Context, tx *sql.Tx, config domain.ConfigRecord) domain.ConfigRecord
	GetByVersion(ctx context.Context, tx *sql.Tx, config domain.ConfigRecord) domain.ConfigRecord
	ListVersions(ctx context.Context, tx *sql.Tx, config domain.ConfigRecord) []domain.ConfigRecord
	CreateNewVersion(ctx context.Context, tx *sql.Tx, config domain.ConfigRecord) domain.ConfigRecord
}
