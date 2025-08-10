package repository

import (
	"config-service/helper"
	"config-service/model/domain"
	"context"
	"database/sql"
	"encoding/json"
)

type ConfigRepositoryImpl struct{}

func NewConfigRepository() ConfigRepository {
	return &ConfigRepositoryImpl{}
}

func (repository *ConfigRepositoryImpl) GetLatest(ctx context.Context, tx *sql.Tx, config domain.ConfigRecord) domain.ConfigRecord {
	SQL := "SELECT schema, name, version, data, created_at FROM configs WHERE schema = ? AND name = ? ORDER BY version DESC LIMIT 1"
	rows, err := tx.QueryContext(ctx, SQL, config.Schema, config.Name)
	helper.PanicIfError(err)
	defer rows.Close()

	var dataStr string
	configRecord := domain.ConfigRecord{}
	if rows.Next() {
		err = rows.Scan(&configRecord.Schema, &configRecord.Name, &configRecord.Version, &dataStr, &configRecord.CreatedAt)
		helper.PanicIfError(err)

		err = json.Unmarshal([]byte(dataStr), &configRecord.Data)
		helper.PanicIfError(err)
	}

	return configRecord
}

func (repository *ConfigRepositoryImpl) GetByVersion(ctx context.Context, tx *sql.Tx, config domain.ConfigRecord) domain.ConfigRecord {

	if config.Version == 0 {
		configRecord := repository.GetLatest(ctx, tx, config)
		return configRecord
	}

	SQL := "SELECT schema, name, version, data, created_at FROM configs WHERE schema = ? AND name = ? AND version = ? ORDER BY version DESC LIMIT 1"
	rows, err := tx.QueryContext(ctx, SQL, config.Schema, config.Name, config.Version)
	helper.PanicIfError(err)
	defer rows.Close()

	var dataStr string
	configRecord := domain.ConfigRecord{}
	if rows.Next() {
		err = rows.Scan(&configRecord.Schema, &configRecord.Name, &configRecord.Version, &dataStr, &configRecord.CreatedAt)
		helper.PanicIfError(err)

		err = json.Unmarshal([]byte(dataStr), &configRecord.Data)
		helper.PanicIfError(err)
	}

	return configRecord
}

func (repository *ConfigRepositoryImpl) ListVersions(ctx context.Context, tx *sql.Tx, config domain.ConfigRecord) []domain.ConfigRecord {

	SQL := "SELECT schema, name, version, data, created_at FROM configs WHERE schema = ? AND name = ? ORDER BY version ASC"
	rows, err := tx.QueryContext(ctx, SQL, config.Schema, config.Name, config.Version)
	helper.PanicIfError(err)
	defer rows.Close()

	var dataStr string
	configRecords := []domain.ConfigRecord{}
	for rows.Next() {
		configRecord := domain.ConfigRecord{}
		err = rows.Scan(&configRecord.Schema, &configRecord.Name, &configRecord.Version, &dataStr, &configRecord.CreatedAt)
		helper.PanicIfError(err)

		err = json.Unmarshal([]byte(dataStr), &configRecord.Data)
		helper.PanicIfError(err)
		configRecords = append(configRecords, configRecord)
	}

	return configRecords
}

func (repository *ConfigRepositoryImpl) CreateNewVersion(ctx context.Context, tx *sql.Tx, config domain.ConfigRecord) domain.ConfigRecord {

	dataJSON, err := json.Marshal(config.Data)
	helper.PanicIfError(err)

	SQL := "INSERT INTO configs (schema, name, version, data) VALUES (?, ?, ?, ?)"
	_, err = tx.ExecContext(ctx, SQL, config.Schema, config.Name, config.Version, string(dataJSON))
	helper.PanicIfError(err)

	helper.PanicIfError(err)

	return config
}
