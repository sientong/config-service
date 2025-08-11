package test

import (
	"config-service/model/domain"
	"config-service/model/web"
	"config-service/service"
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockConfigRepository struct {
	mock.Mock
}

func init() {
	err := domain.LoadSchemas("../schemas")
	if err != nil {
		log.Fatalf("Error loading schemas: %v", err)
	}
}

func (m *mockConfigRepository) GetLatest(ctx context.Context, tx *sql.Tx, record domain.ConfigRecord) domain.ConfigRecord {
	args := m.Called(ctx, tx, record)
	return args.Get(0).(domain.ConfigRecord)
}

func (m *mockConfigRepository) CreateNewVersion(ctx context.Context, tx *sql.Tx, record domain.ConfigRecord) domain.ConfigRecord {
	args := m.Called(ctx, tx, record)
	return args.Get(0).(domain.ConfigRecord)
}

func (m *mockConfigRepository) GetByVersion(ctx context.Context, tx *sql.Tx, record domain.ConfigRecord) domain.ConfigRecord {
	args := m.Called(ctx, tx, record)
	return args.Get(0).(domain.ConfigRecord)
}

func (m *mockConfigRepository) ListVersions(ctx context.Context, tx *sql.Tx, record domain.ConfigRecord) []domain.ConfigRecord {
	args := m.Called(ctx, tx, record)
	return args.Get(0).([]domain.ConfigRecord)
}

func TestCreateConfig(t *testing.T) {

	db, sqlmock := fakeDB(t)
	repo := new(mockConfigRepository)
	validate := validator.New()

	svc := service.NewConfigService(repo, db, validate)

	req := web.ConfigCreateRequest{
		Data: map[string]interface{}{"max_limit": 500, "enabled": true},
	}

	sqlmock.ExpectBegin()
	sqlmock.ExpectCommit()

	repo.On("GetLatest", mock.Anything, mock.Anything, mock.Anything).Return(domain.ConfigRecord{})
	repo.On("CreateNewVersion", mock.Anything, mock.Anything, mock.Anything).Return(domain.ConfigRecord{
		Schema:  "payment_config",
		Name:    "payment",
		Version: 1,
		Data:    map[string]interface{}{"max_limit": 500, "enabled": true},
	})

	resp := svc.CreateConfig(context.Background(), "payment_config", "payment", req)
	assert.Equal(t, 1, resp.Version)
	repo.AssertExpectations(t)
}

func TestUpdateConfig(t *testing.T) {
	db, sqlmock := fakeDB(t)
	repo := new(mockConfigRepository)
	validate := validator.New()

	svc := service.NewConfigService(repo, db, validate)

	req := web.ConfigUpdateRequest{
		Data: map[string]interface{}{"max_limit": 500, "enabled": true},
	}

	sqlmock.ExpectBegin()
	sqlmock.ExpectCommit()

	repo.On("GetLatest", mock.Anything, mock.Anything, mock.Anything).Return(domain.ConfigRecord{
		Schema:  "payment_config",
		Name:    "payment",
		Version: 1,
	})
	repo.On("CreateNewVersion", mock.Anything, mock.Anything, mock.Anything).Return(domain.ConfigRecord{
		Schema:  "payment_config",
		Name:    "payment",
		Version: 2,
		Data:    map[string]interface{}{"max_limit": 500, "enabled": true},
	})

	resp := svc.UpdateConfig(context.Background(), "payment_config", "payment", req)
	assert.Equal(t, 2, resp.Version)
	repo.AssertExpectations(t)
}

func TestFetchConfigLatest(t *testing.T) {
	db, sqlmock := fakeDB(t)
	validate := validator.New()
	repo := new(mockConfigRepository)

	svc := service.NewConfigService(repo, db, validate)

	version := 0
	req := web.ConfigFetchRequest{Version: &version}

	sqlmock.ExpectBegin()
	sqlmock.ExpectCommit()

	repo.On("GetLatest", mock.Anything, mock.Anything, mock.Anything).Return(domain.ConfigRecord{
		Schema:  "payment_config",
		Name:    "payment",
		Version: 5,
	})

	resp := svc.FetchConfig(context.Background(), "payment_config", "payment", req)
	assert.Equal(t, 5, resp.Version)
	repo.AssertExpectations(t)
}

func TestFetchConfigByVersion(t *testing.T) {

	db, sqlmock := fakeDB(t)
	validate := validator.New()
	repo := new(mockConfigRepository)

	svc := service.NewConfigService(repo, db, validate)

	version := 3
	req := web.ConfigFetchRequest{Version: &version}

	sqlmock.ExpectBegin()
	sqlmock.ExpectCommit()

	repo.On("GetByVersion", mock.Anything, mock.Anything, mock.Anything).Return(domain.ConfigRecord{
		Schema:  "payment_config",
		Name:    "payment",
		Version: 3,
	})

	resp := svc.FetchConfig(context.Background(), "payment_config", "payment", req)
	assert.Equal(t, 3, resp.Version)
	repo.AssertExpectations(t)
}

func TestRollbackConfig(t *testing.T) {

	db, sqlmock := fakeDB(t)
	validate := validator.New()
	repo := new(mockConfigRepository)

	svc := service.NewConfigService(repo, db, validate)

	req := web.ConfigRollbackRequest{Version: 2}

	sqlmock.ExpectBegin()
	sqlmock.ExpectCommit()

	repo.On("GetByVersion", mock.Anything, mock.Anything, mock.Anything).Return(domain.ConfigRecord{
		Schema:  "payment_config",
		Name:    "payment",
		Version: 2,
	})
	repo.On("GetLatest", mock.Anything, mock.Anything, mock.Anything).Return(domain.ConfigRecord{
		Schema:  "payment_config",
		Name:    "payment",
		Version: 5,
	})
	repo.On("CreateNewVersion", mock.Anything, mock.Anything, mock.Anything).Return(domain.ConfigRecord{
		Schema:  "payment_config",
		Name:    "payment",
		Version: 6,
	})

	resp := svc.RollbackConfig(context.Background(), "payment_config", "payment", req)
	assert.Equal(t, 6, resp.Version)
	repo.AssertExpectations(t)
}

func TestListVersionService(t *testing.T) {
	db, sqlmock := fakeDB(t)
	repo := new(mockConfigRepository)
	validate := validator.New()

	svc := service.NewConfigService(repo, db, validate)

	sqlmock.ExpectBegin()
	sqlmock.ExpectCommit()

	repo.On("ListVersions", mock.Anything, mock.Anything, mock.Anything).Return([]domain.ConfigRecord{
		{Schema: "payment_config", Name: "payment", Version: 1},
		{Schema: "payment_config", Name: "payment", Version: 2},
	})

	resp := svc.ListVersions(context.Background(), "payment_config", "payment")
	assert.Len(t, resp.ConfigVersions, 2)
	repo.AssertExpectations(t)
}
