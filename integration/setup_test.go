package main

import (
	"context"
	"database/sql"
	"integration/valkyrie"
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

func query(sqlite, postgres string) string {
	if getActiveProvider() == "postgres" {
		return postgres
	}
	return sqlite
}

func getActiveProvider() string {
	content, err := os.ReadFile("schema.prisma")
	if err != nil {
		return "sqlite"
	}
	s := string(content)
	if strings.Contains(s, `provider = "postgres"`) || strings.Contains(s, `provider = "postgresql"`) {
		return "postgres"
	}
	return "sqlite"
}

func getPostgresDSN() string {
	if url := os.Getenv("PG_DATABASE_URL"); url != "" {
		return url
	}

	// Try local docker-compose default DSN first
	localDSN := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", localDSN)
	if err == nil {
		err = db.Ping()
		db.Close()
		if err == nil {
			return localDSN
		}
	}

	// Try CI default DSN
	return "postgres://testuser:testpassword@localhost:5432/valkyrie_test?sslmode=disable"
}

func setupTestDB(t *testing.T) (*valkyrie.DB, func()) {
	ctx := context.Background()

	provider := getActiveProvider()
	var dsn string

	if provider == "postgres" {
		dsn = getPostgresDSN()

		// Reset the postgres schema so we start fresh every time
		resetDB, err := sql.Open("postgres", dsn)
		if err != nil {
			if t != nil {
				t.Fatalf("failed to open database for reset: %v", err)
			} else {
				panic(err)
			}
		}
		_, err = resetDB.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
		resetDB.Close()
		if err != nil {
			if t != nil {
				t.Fatalf("failed to reset postgres database: %v", err)
			} else {
				panic(err)
			}
		}
	} else {
		dsn = "file::memory:?cache=shared&_pragma=foreign_keys(1)"
	}

	db, err := valkyrie.Open(provider, dsn)
	if err != nil {
		if t != nil {
			t.Fatalf("failed to open database (provider: %s, dsn: %s): %v", provider, dsn, err)
		} else {
			panic(err)
		}
	}

	if err := db.RunMigrations(ctx); err != nil {
		db.Close()
		if t != nil {
			t.Fatalf("failed to run migrations: %v", err)
		} else {
			panic(err)
		}
	}

	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}
