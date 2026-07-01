package generator

import (
	"strings"
	"testing"
	"valkyrie/schema"
)

func TestGenerateClient(t *testing.T) {
	input := `
	datasource db {
		provider = "postgresql"
	}

	enum Role {
		USER
		ADMIN
	}

	model User {
		id    Int    @id @default(autoincrement())
		email String @unique
		role  Role   @default(USER)
		
		@@map("users")
	}

	model Post {
		id        String   @id @default(uuid())
		title     String
		published Boolean  @default(false)
		authorId  Int
		author    User     @relation(fields: [authorId], references: [id])
	}
	`

	sch, errs := schema.ParseSchema(input)
	if len(errs) > 0 {
		t.Fatalf("parser errors: %v", errs)
	}

	code, err := GenerateClient(*sch, "client", "migrations/*.sql", "migrations")
	if err != nil {
		t.Fatalf("failed to generate client: %v\nCode output:\n%s", err, code)
	}

	// Verify package name
	if !strings.Contains(code, "package client") {
		t.Errorf("expected package client, got:\n%s", code)
	}

	// Verify client and delegates
	if !strings.Contains(code, "type DB struct {") {
		t.Errorf("expected DB struct, got:\n%s", code)
	}
	if !strings.Contains(code, "User") || !strings.Contains(code, "*UserDelegate") {
		t.Errorf("expected User delegate field on DB, got:\n%s", code)
	}
	if !strings.Contains(code, "Post") || !strings.Contains(code, "*PostDelegate") {
		t.Errorf("expected Post delegate field on DB, got:\n%s", code)
	}
	if !strings.Contains(code, "type UserDelegate struct {") {
		t.Errorf("expected UserDelegate struct, got:\n%s", code)
	}
	if !strings.Contains(code, "type PostDelegate struct {") {
		t.Errorf("expected PostDelegate struct, got:\n%s", code)
	}

	// Verify enums and namespaces (Step 2)
	if !strings.Contains(code, "type RoleType string") {
		t.Errorf("expected RoleType type definition, got:\n%s", code)
	}
	if !strings.Contains(code, "RoleTypeUser") || !strings.Contains(code, "\"USER\"") {
		t.Errorf("expected RoleTypeUser constant, got:\n%s", code)
	}
	if !strings.Contains(code, "RoleTypeAdmin") || !strings.Contains(code, "\"ADMIN\"") {
		t.Errorf("expected RoleTypeAdmin constant, got:\n%s", code)
	}
	if !strings.Contains(code, "type roleNamespace struct {") {
		t.Errorf("expected roleNamespace struct definition, got:\n%s", code)
	}
	if !strings.Contains(code, "var Role = roleNamespace{") {
		t.Errorf("expected Role namespace variable declaration, got:\n%s", code)
	}
	if !strings.Contains(code, "User:") || !strings.Contains(code, "RoleTypeUser") {
		t.Errorf("expected User: RoleTypeUser mapping in Role namespace, got:\n%s", code)
	}
	if !strings.Contains(code, "Admin:") || !strings.Contains(code, "RoleTypeAdmin") {
		t.Errorf("expected Admin: RoleTypeAdmin mapping in Role namespace, got:\n%s", code)
	}
}
