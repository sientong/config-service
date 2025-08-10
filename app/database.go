package app

import (
	"config-service/helper"
	"database/sql"
	"log"
)

func NewDB() *sql.DB {
	db, err := sql.Open("sqlite3", "config_database.db")
	helper.PanicIfError(err)

	log.Println("Connected to SQLite successfully!")

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS configs (
		schema TEXT NOT NULL,
        name TEXT NOT NULL,
        version INTEGER NOT NULL,
        data TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`)
	helper.PanicIfError(err)

	log.Println("Created configs table successfully!")

	return db
}
