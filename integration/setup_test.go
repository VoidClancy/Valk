package main

import (
	"context"
	"integration/valkyrie"
	"testing"
)

func setupTestDB(t *testing.T) (*valkyrie.DB, func()) {
	ctx := context.Background()
	db, err := valkyrie.Open("sqlite", "file::memory:?cache=shared&_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	if err := db.RunMigrations(ctx); err != nil {
		db.Close()
		t.Fatalf("failed to run migrations: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}
