package repositories

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	// Optional: Setup global test DB if needed
	// For now, individual tests will check for env var
	os.Exit(m.Run())
}

func setupTestDB(t *testing.T) *sql.DB {
	dsn := os.Getenv("TEST_DB_URL")
	if dsn == "" {
		t.Skip("TEST_DB_URL not set, skipping integration test")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}

	return db
}
