package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

func RunMigrations(db *sql.DB, schemaPath string) error {
	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("read schema file: %w", err)
	}

	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		return fmt.Errorf("execute schema: %w", err)
	}

	return nil
}

func AutoMigrate(db *sql.DB) error {
	schemaPath := filepath.Join("db", "schema.sql")
	
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		return fmt.Errorf("schema file not found: %s", schemaPath)
	}

	return RunMigrations(db, schemaPath)
}
