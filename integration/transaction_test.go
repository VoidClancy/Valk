package main

import (
	"context"
	"errors"
	"fmt"
	"integration/valkyrie"
	"testing"
)

func TestTransactionCommit(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	err := db.Transaction(ctx, func(tx *valkyrie.Tx) error {
		_, err := tx.Raw().ExecContext(ctx,
			query(
				`INSERT INTO "User" ("id", "email", "phoneNum", "role") VALUES (?, ?, ?, ?)`,
				`INSERT INTO "User" ("id", "email", "phoneNum", "role") VALUES ($1, $2, $3, $4)`,
			),
			"u1", "user1@example.com", "123456789", "student",
		)
		return err
	})
	if err != nil {
		t.Fatalf("expected transaction to commit, got error: %v", err)
	}

	var email string
	err = db.Raw().QueryRowContext(ctx, query(
		`SELECT "email" FROM "User" WHERE "id" = ?`,
		`SELECT "email" FROM "User" WHERE "id" = $1`,
	), "u1").Scan(&email)
	if err != nil {
		t.Fatalf("failed to query committed user: %v", err)
	}
	if email != "user1@example.com" {
		t.Errorf("expected email 'user1@example.com', got '%s'", email)
	}
}

func TestTransactionRollbackOnError(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	expectedErr := errors.New("something went wrong")
	err := db.Transaction(ctx, func(tx *valkyrie.Tx) error {
		_, err := tx.Raw().ExecContext(ctx,
			query(
				`INSERT INTO "User" ("id", "email", "phoneNum", "role") VALUES (?, ?, ?, ?)`,
				`INSERT INTO "User" ("id", "email", "phoneNum", "role") VALUES ($1, $2, $3, $4)`,
			),
			"u2", "user2@example.com", "987654321", "ADMIN",
		)
		if err != nil {
			return err
		}
		return expectedErr //to force a rollback
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error '%v', got '%v'", expectedErr, err)
	}

	var count int
	err = db.Raw().QueryRowContext(ctx, query(
		`SELECT COUNT(*) FROM "User" WHERE "id" = ?`,
		`SELECT COUNT(*) FROM "User" WHERE "id" = $1`,
	), "u2").Scan(&count)
	if err != nil {
		t.Fatalf("failed to query database: %v", err)
	}
	if count != 0 {
		t.Errorf("expected user count to be 0 after rollback, got %d", count)
	}
}

func TestTransactionRollbackOnPanic(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	panicVal := "unrecoverable panic"
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected transaction block to panic, but it completed normally")
		}
		if fmt.Sprintf("%v", r) != panicVal {
			t.Errorf("expected panic value '%v', got '%v'", panicVal, r)
		}

		var count int
		err := db.Raw().QueryRowContext(ctx, query(
			`SELECT COUNT(*) FROM "User" WHERE "id" = ?`,
			`SELECT COUNT(*) FROM "User" WHERE "id" = $1`,
		), "u3").Scan(&count)
		if err != nil {
			t.Fatalf("failed to query database after panic rollback: %v", err)
		}
		if count != 0 {
			t.Errorf("expected user count to be 0 after panic rollback, got %d", count)
		}
	}()

	_ = db.Transaction(ctx, func(tx *valkyrie.Tx) error {
		_, err := tx.Raw().ExecContext(ctx,
			query(
				`INSERT INTO "User" ("id", "email", "phoneNum", "role") VALUES (?, ?, ?, ?)`,
				`INSERT INTO "User" ("id", "email", "phoneNum", "role") VALUES ($1, $2, $3, $4)`,
			),
			"u3", "user3@example.com", "555555555", "TEACHER",
		)
		if err != nil {
			return err
		}
		panic(panicVal)
	})
}

func TestManualTransactionCommitAndRollback(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	tx1, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("failed to begin tx1: %v", err)
	}
	_, err = tx1.Raw().ExecContext(ctx,
		query(
			`INSERT INTO "User" ("id", "email", "phoneNum", "role") VALUES (?, ?, ?, ?)`,
			`INSERT INTO "User" ("id", "email", "phoneNum", "role") VALUES ($1, $2, $3, $4)`,
		),
		"u4", "user4@example.com", "444444444", "student",
	)
	if err != nil {
		tx1.Rollback()
		t.Fatalf("insert failed: %v", err)
	}
	if err := tx1.Commit(); err != nil {
		t.Fatalf("commit failed: %v", err)
	}

	var count1 int
	err = db.Raw().QueryRowContext(ctx, query(
		`SELECT COUNT(*) FROM "User" WHERE "id" = ?`,
		`SELECT COUNT(*) FROM "User" WHERE "id" = $1`,
	), "u4").Scan(&count1)
	if err != nil || count1 != 1 {
		t.Errorf("expected user to be committed, count=%d, err=%v", count1, err)
	}

	tx2, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("failed to begin tx2: %v", err)
	}
	_, err = tx2.Raw().ExecContext(ctx,
		query(
			`INSERT INTO "User" ("id", "email", "phoneNum", "role") VALUES (?, ?, ?, ?)`,
			`INSERT INTO "User" ("id", "email", "phoneNum", "role") VALUES ($1, $2, $3, $4)`,
		),
		"u5", "user5@example.com", "555555555", "student",
	)
	if err != nil {
		tx2.Rollback()
		t.Fatalf("insert failed: %v", err)
	}
	if err := tx2.Rollback(); err != nil {
		t.Fatalf("rollback failed: %v", err)
	}

	var count2 int
	err = db.Raw().QueryRowContext(ctx, query(
		`SELECT COUNT(*) FROM "User" WHERE "id" = ?`,
		`SELECT COUNT(*) FROM "User" WHERE "id" = $1`,
	), "u5").Scan(&count2)
	if err != nil || count2 != 0 {
		t.Errorf("expected user to be rolled back, count=%d, err=%v", count2, err)
	}
}

func TestTransactionNestedErrorWrapping(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	err := db.Transaction(ctx, func(tx *valkyrie.Tx) error {
		// invalid SQL to trigger a real SQL error
		_, err := tx.Raw().ExecContext(ctx, "INSERT INTO NonExistentTable (id) VALUES (1)")
		return err
	})

	if err == nil {
		t.Fatalf("expected error from invalid query, got nil")
	}
}
