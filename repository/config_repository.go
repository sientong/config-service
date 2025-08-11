package repository

import (
	"config-service/model/domain"
	"context"
	"database/sql"
)

type ConfigRepository interface {
	GetLatest(ctx context.Context, tx *sql.Tx, config domain.ConfigRecord) (domain.ConfigRecord, error)
	GetByVersion(ctx context.Context, tx *sql.Tx, config domain.ConfigRecord) (domain.ConfigRecord, error)
	ListVersions(ctx context.Context, tx *sql.Tx, config domain.ConfigRecord) []domain.ConfigRecord
	CreateNewVersion(ctx context.Context, tx *sql.Tx, config domain.ConfigRecord) domain.ConfigRecord
}
