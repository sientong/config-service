package test

import (
	"config-service/app"
	"config-service/controller"
	"config-service/helper"
	"config-service/model/domain"
	"config-service/repository"
	"config-service/service"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator"
)

var db *sql.DB

func setupTestDB() {
	var err error
	db, err = sql.Open("sqlite3", "config_database_testing.db")
	helper.PanicIfError(err)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS configs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		schema TEXT NOT NULL,
        name TEXT NOT NULL,
        version INTEGER NOT NULL,
        data TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`)
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

}

func setupRouter(db *sql.DB) http.Handler {
	validate := validator.New()
	configRepository := repository.NewConfigRepository()
	configService := service.NewConfigService(configRepository, db, validate)
	configController := controller.NewConfigController(configService)
	schemaController := controller.NewSchemaController()

	router := app.NewRouter(configController, schemaController)

	return router
}

func truncateConfigs(db *sql.DB) {
	db.Exec("DELETE from configs")
	db.Exec("VACUUM")
}

func init() {
	setupTestDB()

	err := domain.LoadSchemas("../schemas")
	if err != nil {
		log.Fatalf("Error loading schemas: %v", err)
	}
}

func performRequest(method, path string, body io.Reader, truncateData bool) (map[string]interface{}, *http.Response) {

	if truncateData {
		truncateConfigs(db)
	}
	router := setupRouter(db)

	req := httptest.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	res := rec.Result()
	var jsonBody map[string]interface{}
	_ = json.NewDecoder(res.Body).Decode(&jsonBody)
	return jsonBody, res
}

func fakeDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	return db, mock
}
