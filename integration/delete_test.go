package main

import (
	"context"
	"integration/valk"
	"integration/valk/user"
	"testing"
)

func TestDeleteMany(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.CreateMany(
		db.User.Create().SetEmail("delete1@example.com").SetPhoneNum("+111").SetRole(valk.UserRole.Student),
		db.User.Create().SetEmail("delete2@example.com").SetPhoneNum("+222").SetRole(valk.UserRole.Student),
		db.User.Create().SetEmail("keep1@example.com").SetPhoneNum("+333").SetRole(valk.UserRole.Admin),
	).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to seed users: %v", err)
	}

	count, err := db.User.DeleteMany(user.Role.EQ(valk.UserRole.Student)).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to delete students: %v", err)
	}

	if count != 2 {
		t.Errorf("expected 2 deleted users, got %d", count)
	}

	var remainingCount int
	err = db.Raw().QueryRowContext(ctx, query(
		`SELECT COUNT(*) FROM "User"`,
		`SELECT COUNT(*) FROM "User"`,
	)).Scan(&remainingCount)
	if err != nil {
		t.Fatalf("failed to count remaining users: %v", err)
	}

	if remainingCount != 1 {
		t.Errorf("expected 1 remaining user in DB, got %d", remainingCount)
	}

	var remainingEmail string
	err = db.Raw().QueryRowContext(ctx, query(
		`SELECT "email" FROM "User"`,
		`SELECT "email" FROM "User"`,
	)).Scan(&remainingEmail)
	if err != nil {
		t.Fatalf("failed to get remaining user: %v", err)
	}

	if remainingEmail != "keep1@example.com" {
		t.Errorf("expected remaining user to be keep1@example.com, got %s", remainingEmail)
	}
}

func TestDeleteMany_NoMatches(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.CreateMany(
		db.User.Create().SetEmail("keep1@example.com").SetPhoneNum("+333").SetRole(valk.UserRole.Admin),
	).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to seed users: %v", err)
	}

	count, err := db.User.DeleteMany(user.Role.EQ(valk.UserRole.Student)).Exec(ctx)
	if err != nil {
		t.Fatalf("failed to delete students: %v", err)
	}

	if count != 0 {
		t.Errorf("expected 0 deleted users, got %d", count)
	}
}
