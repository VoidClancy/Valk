package main

import (
	"testing"
)

func TestConnectionAndMigration(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	err := db.Raw().Ping()
	if err != nil {
		t.Fatalf("failed to ping raw database: %v", err)
	}
}
