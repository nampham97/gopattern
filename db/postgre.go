// File: db/db.go
package db

import (
	"GoPattern/config"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func InitDB(cfg config.Config) error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	var err error
	db, err = sqlx.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open DB: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to connect to DB: %w", err)
	}

	fmt.Println("✅ DB connection established")
	return nil
}

func GetDB() *sqlx.DB {
	if db == nil {
		panic("Database connection is not initialized. Call InitDB first.")
	}
	return db
}
