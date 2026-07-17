package main

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "modernc.org/sqlite"
)

const userSchema = `CREATE TABLE IF NOT EXISTS "User" (
	"id" TEXT PRIMARY KEY,
	"email" TEXT NOT NULL UNIQUE,
	"phoneNum" TEXT NOT NULL UNIQUE,
	"password" TEXT,
	"role" TEXT NOT NULL DEFAULT 'student',
	"roleOptional" TEXT,
	"loginCount" INTEGER NOT NULL DEFAULT 0,
	"referredById" TEXT
)`

const seedCount = 1000

type DBTX interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

func openDB(b *testing.B) *sql.DB {
	b.Helper()
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		b.Fatal(err)
	}
	return db
}

func createSchema(db DBTX) {
	_, err := db.ExecContext(context.Background(), userSchema)
	if err != nil {
		panic(fmt.Sprintf("create schema: %v", err))
	}
}

func seedData(db DBTX, prefix string) {
	for i := range seedCount {
		_, err := db.ExecContext(context.Background(),
			`INSERT INTO "User" ("id", "email", "phoneNum", "password", "role") VALUES (?, ?, ?, ?, ?)`,
			fmt.Sprintf("%s-id-%d", prefix, i),
			fmt.Sprintf("%s-user-%d@example.com", prefix, i),
			fmt.Sprintf("%s-phone-%d", prefix, i),
			nil, "student",
		)
		if err != nil {
			panic(fmt.Sprintf("seed %s: %v", prefix, err))
		}
	}
}
