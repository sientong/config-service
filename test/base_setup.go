package test

import (
	"config-service/app"
	"config-service/controller"
	"config-service/helper"
	"config-service/repository"
	"config-service/service"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-playground/validator"
)

func setupTestDB() *sql.DB {
	db, err := sql.Open("sqlite3", "config_database_testing.db")
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

	return db
}

func setupRouter(db *sql.DB) http.Handler {
	validate := validator.New()
	configRepository := repository.NewConfigRepository()
	configService := service.NewConfigService(configRepository, db, validate)
	configController := controller.NewConfigController(configService)
	router := app.NewRouter(configController)

	return router
}

func truncateConfigs(db *sql.DB) {
	db.Exec("DELETE from configs")
}

func performRequest(method, path string, body io.Reader, truncateData bool) (map[string]interface{}, *http.Response) {

	db := setupTestDB()
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
