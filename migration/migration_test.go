package migration

import (
	"strings"
	"testing"
	"valkyrie/schema"
)

func TestGenerateMigrationsPostgres(t *testing.T) {
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
		author    User     @relation(fields: [authorId], references: [id], onDelete: Cascade)
	}
	`

	schemaDef, errs := schema.ParseSchema(input)
	if len(errs) > 0 {
		t.Fatalf("parser errors: %v", errs)
	}

	sql, err := GenerateUpMigrations(schemaDef)
	if err != nil {
		t.Fatalf("failed to generate DDL: %v", err)
	}

	// Verify enums
	if !strings.Contains(sql, "CREATE TYPE \"Role\" AS ENUM") {
		t.Errorf("expected CREATE TYPE for Role enum, got:\n%s", sql)
	}

	// Verify users table
	if !strings.Contains(sql, "CREATE TABLE \"users\"") {
		t.Errorf("expected CREATE TABLE users, got:\n%s", sql)
	}
	if !strings.Contains(sql, "\"id\" SERIAL NOT NULL") && !strings.Contains(sql, "\"id\" BIGSERIAL NOT NULL") {
		t.Errorf("expected id serial/bigserial column, got:\n%s", sql)
	}
	if !strings.Contains(sql, "CONSTRAINT \"users_pkey\" PRIMARY KEY (\"id\")") {
		t.Errorf("expected primary key constraint users_pkey, got:\n%s", sql)
	}
	if !strings.Contains(sql, "\"role\" \"Role\" NOT NULL DEFAULT 'USER'") {
		t.Errorf("expected role enum column, got:\n%s", sql)
	}

	// Verify foreign keys and relation onDelete
	if !strings.Contains(sql, "FOREIGN KEY (\"authorId\") REFERENCES \"users\" (\"id\") ON DELETE CASCADE") {
		t.Errorf("expected foreign key on authorId referencing users(id) ON DELETE CASCADE, got:\n%s", sql)
	}
}

func TestGenerateMigrationsSQLite(t *testing.T) {
	input := `
	datasource db {
		provider = "sqlite"
	
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
	`

	schemaDef, errs := schema.ParseSchema(input)
	if len(errs) > 0 {
		t.Fatalf("parser errors: %v", errs)
	}

	sql, err := GenerateUpMigrations(schemaDef)
	if err != nil {
		t.Fatalf("failed to generate DDL: %v", err)
	}

	// In SQLite, there shouldn't be CREATE TYPE
	if strings.Contains(sql, "CREATE TYPE") {
		t.Errorf("expected no CREATE TYPE for SQLite, got:\n%s", sql)
	}

	// SQLite check constraints
	if !strings.Contains(sql, "CHECK (\"role\" IN ('USER', 'ADMIN'))") {
		t.Errorf("expected SQLite enum CHECK constraint, got:\n%s", sql)
	}
	// SQLite autoincrement
	if !strings.Contains(sql, "\"id\" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT") {
		t.Errorf("expected INTEGER PRIMARY KEY AUTOINCREMENT for SQLite, got:\n%s", sql)
	}
}

func TestGenerateGooseMigrations(t *testing.T) {
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
	`
	schemaDef, errs := schema.ParseSchema(input)
	if len(errs) > 0 {
		t.Fatalf("parser errors: %v", errs)
	}

	sql, err := GenerateMigration(schemaDef)
	if err != nil {
		t.Fatalf("failed to generate Goose DDL: %v", err)
	}

	if !strings.Contains(sql, "-- +goose Up") {
		t.Errorf("expected -- +goose Up directive, got:\n%s", sql)
	}
	if !strings.Contains(sql, "-- +goose Down") {
		t.Errorf("expected -- +goose Down directive, got:\n%s", sql)
	}
	if !strings.Contains(sql, "DROP TABLE IF EXISTS \"users\";") {
		t.Errorf("expected table drop, got:\n%s", sql)
	}
	if !strings.Contains(sql, "DROP TYPE IF EXISTS \"Role\";") {
		t.Errorf("expected enum drop, got:\n%s", sql)
	}
}

func TestViciousMigrationEdgeCases(t *testing.T) {
	schemaTemplate := `
	datasource db {
		provider = "PROVIDER"
	}

	enum CustomRole {
		SUPER_ADMIN  @map("super_admin")
		REGULAR_USER @map("regular_user")
		
		@@map("custom_roles_enum")
	}

	model CustomUser {
		id        Int        @id @default(autoincrement()) @map("user_id")
		email     String     @unique @map("user_email")
		phone     String?    @map("user_phone")
		role      CustomRole @default(REGULAR_USER) @map("user_role")
		
		// Self-referential relation
		managerId Int?       @map("manager_id")
		manager   CustomUser? @relation("UserToManager", fields: [managerId], references: [id], onDelete: SetNull, onUpdate: Cascade)
		
		@@unique([email, phone])
		@@map("users_table")
	}

	model Category {
		id   Int    @id @default(autoincrement()) @map("cat_id")
		name String @unique @map("cat_name")
		
		@@map("categories_table")
	}

	model CategoryToUser {
		userId     Int @map("fk_user_id")
		categoryId Int @map("fk_cat_id")
		
		user     CustomUser @relation(fields: [userId], references: [id], onDelete: Cascade)
		category Category   @relation(fields: [categoryId], references: [id], onDelete: Restrict)
		
		@@id([userId, categoryId])
		@@map("user_categories")
	}
	`

	t.Run("postgresql", func(t *testing.T) {
		input := strings.Replace(schemaTemplate, "PROVIDER", "postgresql", 1)
		schemaDef, errs := schema.ParseSchema(input)
		if len(errs) > 0 {
			t.Fatalf("parser errors: %v", errs)
		}

		sql, err := GenerateUpMigrations(schemaDef)
		if err != nil {
			t.Fatalf("failed to generate Postgres DDL: %v", err)
		}

		// 1. Verify mapped custom enum type and mapped value strings
		if !strings.Contains(sql, "CREATE TYPE \"custom_roles_enum\" AS ENUM") {
			t.Errorf("expected CREATE TYPE for custom_roles_enum, got:\n%s", sql)
		}
		if !strings.Contains(sql, "'super_admin'") || !strings.Contains(sql, "'regular_user'") {
			t.Errorf("expected custom enum values in creation block, got:\n%s", sql)
		}

		// 2. Verify users_table mapped columns & types
		if !strings.Contains(sql, "CREATE TABLE \"users_table\"") {
			t.Errorf("expected CREATE TABLE users_table, got:\n%s", sql)
		}
		if !strings.Contains(sql, "\"user_id\" SERIAL NOT NULL") && !strings.Contains(sql, "\"user_id\" BIGSERIAL NOT NULL") {
			t.Errorf("expected mapped user_id serial/bigserial type, got:\n%s", sql)
		}
		if !strings.Contains(sql, "\"user_role\" \"custom_roles_enum\" NOT NULL DEFAULT 'regular_user'") {
			t.Errorf("expected mapped user_role using custom_roles_enum type and default value, got:\n%s", sql)
		}

		// 3. Verify composite unique constraints on users_table
		if !strings.Contains(sql, "CONSTRAINT \"users_table_user_email_user_phone_key\" UNIQUE (\"user_email\", \"user_phone\")") {
			t.Errorf("expected composite unique key for users_table using mapped column names, got:\n%s", sql)
		}

		// 4. Verify self-referential foreign key constraint on users_table
		if !strings.Contains(sql, "CONSTRAINT \"users_table_manager_id_fkey\" FOREIGN KEY (\"manager_id\") REFERENCES \"users_table\" (\"user_id\") ON DELETE SET NULL ON UPDATE CASCADE") {
			t.Errorf("expected self-referential foreign key on users_table, got:\n%s", sql)
		}

		// 5. Verify user_categories table PK and FKs using mapped column names
		if !strings.Contains(sql, "CONSTRAINT \"user_categories_pkey\" PRIMARY KEY (\"fk_user_id\", \"fk_cat_id\")") {
			t.Errorf("expected composite primary key on user_categories using mapped names, got:\n%s", sql)
		}
		if !strings.Contains(sql, "CONSTRAINT \"user_categories_fk_user_id_fkey\" FOREIGN KEY (\"fk_user_id\") REFERENCES \"users_table\" (\"user_id\") ON DELETE CASCADE") {
			t.Errorf("expected foreign key on user_categories pointing to users_table(user_id) with CASCADE, got:\n%s", sql)
		}
		if !strings.Contains(sql, "CONSTRAINT \"user_categories_fk_cat_id_fkey\" FOREIGN KEY (\"fk_cat_id\") REFERENCES \"categories_table\" (\"cat_id\") ON DELETE RESTRICT") {
			t.Errorf("expected foreign key on user_categories pointing to categories_table(cat_id) with RESTRICT, got:\n%s", sql)
		}

		// 6. Verify Down Goose migrations dropping enums and tables in correct reverse order
		gooseSQL, err := GenerateMigration(schemaDef)
		if err != nil {
			t.Fatalf("failed to generate Goose migrations: %v", err)
		}

		if !strings.Contains(gooseSQL, "DROP TABLE IF EXISTS \"user_categories\";") {
			t.Errorf("expected dropping dependent table user_categories, got:\n%s", gooseSQL)
		}
		if !strings.Contains(gooseSQL, "DROP TYPE IF EXISTS \"custom_roles_enum\";") {
			t.Errorf("expected dropping custom roles enum type in down block, got:\n%s", gooseSQL)
		}
	})

	t.Run("sqlite", func(t *testing.T) {
		input := strings.Replace(schemaTemplate, "PROVIDER", "sqlite", 1)
		schemaDef, errs := schema.ParseSchema(input)
		if len(errs) > 0 {
			t.Fatalf("parser errors: %v", errs)
		}

		sql, err := GenerateUpMigrations(schemaDef)
		if err != nil {
			t.Fatalf("failed to generate SQLite DDL: %v", err)
		}

		// 1. Verify SQLite enums are inline CHECK constraints rather than custom types
		if strings.Contains(sql, "CREATE TYPE") {
			t.Errorf("expected no CREATE TYPE for SQLite, got:\n%s", sql)
		}
		if !strings.Contains(sql, "CHECK (\"user_role\" IN ('super_admin', 'regular_user'))") {
			t.Errorf("expected inline CHECK constraint on SQLite enum column, got:\n%s", sql)
		}

		// 2. Verify single PK autoincrement handling on SQLite
		if !strings.Contains(sql, "\"user_id\" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT") {
			t.Errorf("expected INTEGER PRIMARY KEY AUTOINCREMENT for user_id in users_table, got:\n%s", sql)
		}

		// 3. Verify self-referential foreign key constraint on users_table
		if !strings.Contains(sql, "CONSTRAINT \"users_table_manager_id_fkey\" FOREIGN KEY (\"manager_id\") REFERENCES \"users_table\" (\"user_id\") ON DELETE SET NULL ON UPDATE CASCADE") {
			t.Errorf("expected SQLite self-referential foreign key, got:\n%s", sql)
		}

		// 4. Verify composite primary key and foreign keys on user_categories using mapped column names
		if !strings.Contains(sql, "CONSTRAINT \"user_categories_pkey\" PRIMARY KEY (\"fk_user_id\", \"fk_cat_id\")") {
			t.Errorf("expected composite primary key on user_categories, got:\n%s", sql)
		}
		if !strings.Contains(sql, "CONSTRAINT \"user_categories_fk_user_id_fkey\" FOREIGN KEY (\"fk_user_id\") REFERENCES \"users_table\" (\"user_id\") ON DELETE CASCADE") {
			t.Errorf("expected foreign key referencing users_table(user_id), got:\n%s", sql)
		}

		// 5. Verify goose down migrations have PRAGMA foreign_keys = ON;
		gooseSQL, err := GenerateMigration(schemaDef)
		if err != nil {
			t.Fatalf("failed to generate Goose Down SQL: %v", err)
		}
		if !strings.Contains(gooseSQL, "PRAGMA foreign_keys = ON;") {
			t.Errorf("expected PRAGMA foreign_keys = ON; in Down migration, got:\n%s", gooseSQL)
		}
	})
}
