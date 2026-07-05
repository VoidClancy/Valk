package main

import (
	"context"
	"integration/valkyrie"
	"strings"
	"testing"
)

func TestCreateBasic(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	u, err := db.User.Create(valkyrie.UserCreateInput{
		Email:    "test@example.com",
		PhoneNum: "+123456789",
	}).Exec(ctx)

	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if u.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", u.Email)
	}
	if u.PhoneNum != "+123456789" {
		t.Errorf("expected phoneNum '+123456789', got '%s'", u.PhoneNum)
	}
	if u.Role != valkyrie.UserRole.Student {
		t.Errorf("expected role '%s' (default), got '%s'", valkyrie.UserRole.Student, u.Role)
	}
	if !strings.HasPrefix(u.Id, "c") || len(u.Id) < 10 {
		t.Errorf("expected generated ID to be a CUID, got '%s'", u.Id)
	}

	var dbEmail, dbPhone string
	var dbRole string
	err = db.Raw().QueryRowContext(ctx, "SELECT email, phoneNum, role FROM User WHERE id = ?", u.Id).Scan(&dbEmail, &dbPhone, &dbRole)
	if err != nil {
		t.Fatalf("failed to query database for created user: %v", err)
	}
	if dbEmail != "test@example.com" || dbPhone != "+123456789" || dbRole != "student" {
		t.Errorf("database record mismatch: email=%s, phone=%s, role=%s", dbEmail, dbPhone, dbRole)
	}
}

func TestCreateWithSelect(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	u, err := db.User.Create(
		valkyrie.UserCreateInput{
			Email:    "select@example.com",
			PhoneNum: "+999999999",
		}).Select(valkyrie.UserSelect{
		Id:    true,
		Email: true,
		Profile: &valkyrie.ProfileSelect{
			Id:  true,
			Bio: true,
		},
	}).Exec(ctx)

	if err != nil {
		t.Fatalf("failed to create user with select: %v", err)
	}

	if u.Email != "select@example.com" {
		t.Errorf("expected email 'select@example.com', got '%s'", u.Email)
	}
	if u.Id == "" {
		t.Errorf("expected non-empty selected ID")
	}
	if u.PhoneNum != "" {
		t.Errorf("expected unselected phoneNum to be empty string, got '%s'", u.PhoneNum)
	}
	if u.Role != "" {
		t.Errorf("expected unselected role to be empty string, got '%s'", u.Role)
	}
}

func TestCreateWithOmit(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	u, err := db.User.Create(valkyrie.UserCreateInput{
		Email:    "omit@example.com",
		PhoneNum: "+888888888",
	}).Omit(valkyrie.UserOmit{
		PhoneNum: true,
	}).Exec(ctx)

	if err != nil {
		t.Fatalf("failed to create user with omit: %v", err)
	}

	if u.Email != "omit@example.com" {
		t.Errorf("expected email 'omit@example.com', got '%s'", u.Email)
	}
	if u.Id == "" {
		t.Errorf("expected non-empty ID")
	}
	if u.Role != valkyrie.UserRole.Student {
		t.Errorf("expected non-omitted role '%s', got '%s'", valkyrie.UserRole.Student, u.Role)
	}
	if u.PhoneNum != "" {
		t.Errorf("expected omitted phoneNum to be empty, got '%s'", u.PhoneNum)
	}
}

func TestCreateWithCustomEnum(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	adminRole := valkyrie.UserRole.Admin

	u, err := db.User.Create(valkyrie.UserCreateInput{
		Email:    "admin@example.com",
		PhoneNum: "+000000000",
		Role:     &adminRole,
	}).Exec(ctx)

	if err != nil {
		t.Fatalf("failed to create admin: %v", err)
	}

	if u.Role != valkyrie.UserRole.Admin {
		t.Errorf("expected role '%s', got '%s'", valkyrie.UserRole.Admin, u.Role)
	}
}

func TestCreateValidation(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.User.Create(valkyrie.UserCreateInput{
		// no email
		PhoneNum: "+123456789",
	}).Exec(ctx)
	if err == nil {
		t.Fatal("expected error creating user with empty required email, got nil")
	}
	if !strings.Contains(err.Error(), "field Email is required") {
		t.Errorf("expected error message to contain 'field Email is required', got: %v", err)
	}

	invalidRole := valkyrie.UserRoleType("INVALID_ROLE")
	_, err = db.User.Create(valkyrie.UserCreateInput{
		Email:    "invalid_role@example.com",
		PhoneNum: "+123456789",
		Role:     &invalidRole,
	}).Exec(ctx)
	if err == nil {
		t.Fatal("expected error creating user with invalid enum role, got nil")
	}
	if !strings.Contains(err.Error(), "invalid enum value \"INVALID_ROLE\" for field Role") {
		t.Errorf("expected error message to contain 'invalid enum value \"INVALID_ROLE\" for field Role', got: %v", err)
	}
}
