package migration

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"

	providers "valk/dbProviders"
	vs "valk/schema"
)

// Helper: Open unique SQLite in-memory database
func openTestSQLite(t *testing.T, dbName string) *sql.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", dbName)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	db.SetMaxOpenConns(1) // Force single connection for consistency and PRAGMAs
	t.Cleanup(func() {
		db.Close()
	})
	return db
}

// Helper: Open PostgreSQL database with isolated schema
func openTestPostgres(t *testing.T, schemaName string) *sql.DB {
	defaultDSN := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", defaultDSN)
	if err != nil {
		t.Skipf("skipping Postgres test; failed to connect to default database: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		t.Skipf("skipping Postgres test; postgres is not reachable: %v", err)
	}

	// Drop schema if it exists from previous crashed runs, then create it
	_, err = db.Exec(fmt.Sprintf(`DROP SCHEMA IF EXISTS %q CASCADE`, schemaName))
	if err != nil {
		t.Fatalf("failed to drop schema %q: %v", schemaName, err)
	}
	_, err = db.Exec(fmt.Sprintf(`CREATE SCHEMA %q`, schemaName))
	if err != nil {
		t.Fatalf("failed to create schema %q: %v", schemaName, err)
	}

	// Now connect with search_path set to our schemaName
	schemaDSN := fmt.Sprintf("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable&search_path=%s", schemaName)
	schemaDB, err := sql.Open("postgres", schemaDSN)
	if err != nil {
		t.Fatalf("failed to connect to schema %q: %v", schemaName, err)
	}

	// Register cleanup to drop the schema
	t.Cleanup(func() {
		schemaDB.Close()

		// Reconnect to default to drop the schema
		db, err := sql.Open("postgres", defaultDSN)
		if err == nil {
			db.Exec(fmt.Sprintf(`DROP SCHEMA IF EXISTS %q CASCADE`, schemaName))
			db.Close()
		}
	})

	return schemaDB
}

// Helper: Parse schema or fail test
func parseOrFail(t *testing.T, input string) *vs.Schema {
	schemaDef, errs := vs.ParseSchema(input)
	if len(errs) > 0 {
		t.Fatalf("parse errors: %v", errs)
	}
	return schemaDef
}

// Helper: Run DiffAndPlan or fail test
func diffOrFail(t *testing.T, db *sql.DB, provider providers.DbProvider, schemaDef *vs.Schema, isInteractive bool) (up, down string) {
	up, down, err := DiffAndPlan(db, provider, schemaDef, isInteractive)
	if err != nil {
		t.Fatalf("DiffAndPlan failed: %v", err)
	}
	return up, down
}

// Helper: Execute SQL statements or fail test
func execOrFail(t *testing.T, db *sql.DB, query string) {
	if strings.TrimSpace(query) == "" {
		return
	}
	stmts := strings.Split(query, ";")
	for _, stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		_, err := db.Exec(stmt)
		if err != nil {
			t.Fatalf("failed to execute SQL statement %q: %v", stmt, err)
		}
	}
}

// Helper: Assert string contains substring
func assertContains(t *testing.T, str, substr string) {
	if !strings.Contains(str, substr) {
		t.Errorf("expected string to contain %q, but got:\n%s", substr, str)
	}
}

// Helper: Assert string does not contain substring
func assertNotContains(t *testing.T, str, substr string) {
	if strings.Contains(str, substr) {
		t.Errorf("expected string NOT to contain %q, but got:\n%s", substr, str)
	}
}

// Helper: Mock stdin input for rename prompts
func mockStdin(t *testing.T, input string) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	oldStdin := os.Stdin
	t.Cleanup(func() {
		os.Stdin = oldStdin
	})
	os.Stdin = r
	w.Write([]byte(input))
	w.Close()
}

// 1. Initial Migration from Empty DB (TestInitialMigration)
func TestInitialMigration(t *testing.T) {
	t.Run("sqlite", func(t *testing.T) {
		db := openTestSQLite(t, "initial_migration_sqlite")
		schemaText := `
		datasource db {
			provider = "sqlite"
		}
		model User {
			id        Int      @id @default(autoincrement())
			email     String   @unique
			active    Boolean  @default(true)
			createdAt DateTime @default(now())
			name      String?
			
			@@map("users")
		}
		model Profile {
			id      Int    @id @default(autoincrement())
			bio     String @map("biography")
			userId  Int    @unique @map("user_id")
		}
		`
		sc := parseOrFail(t, schemaText)
		up, _ := diffOrFail(t, db, providers.Sqlite, sc, false)

		assertContains(t, up, "CREATE TABLE")
		assertContains(t, up, "users")
		assertContains(t, up, "biography")
		assertContains(t, up, "user_id")

		execOrFail(t, db, up)

		// Test insert
		_, err := db.Exec(`INSERT INTO "users" ("email", "name") VALUES ('alice@example.com', 'Alice')`)
		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}

		// Verify default values
		var active bool
		var name string
		err = db.QueryRow(`SELECT "active", "name" FROM "users" WHERE "email" = 'alice@example.com'`).Scan(&active, &name)
		if err != nil {
			t.Fatalf("failed to query inserted user: %v", err)
		}
		if !active || name != "Alice" {
			t.Errorf("expected active=true, name='Alice', got active=%t, name=%q", active, name)
		}

		// Verify unique constraint
		_, err = db.Exec(`INSERT INTO "users" ("email") VALUES ('alice@example.com')`)
		if err == nil {
			t.Errorf("expected unique constraint violation, but insert succeeded")
		}
	})

	t.Run("postgres", func(t *testing.T) {
		db := openTestPostgres(t, "test_initial_migration_pg")
		schemaText := `
		datasource db {
			provider = "postgresql"
		}
		model User {
			id        Int      @id @default(autoincrement())
			email     String   @unique
			active    Boolean  @default(true)
			createdAt DateTime @default(now())
			name      String?
			
			@@map("users")
		}
		model Profile {
			id      Int    @id @default(autoincrement())
			bio     String @map("biography")
			userId  Int    @unique @map("user_id")
		}
		model Post {
			id    String @id @default(uuid())
			title String @default("Untitled")
		}
		`
		sc := parseOrFail(t, schemaText)
		up, _ := diffOrFail(t, db, providers.Postgres, sc, false)

		assertContains(t, up, "CREATE TABLE")
		assertContains(t, up, `"users"`)
		assertContains(t, up, `"biography"`)
		assertContains(t, up, `"user_id"`)
		assertContains(t, up, "gen_random_uuid()")

		execOrFail(t, db, up)

		// Test insert
		_, err := db.Exec(`INSERT INTO "users" ("email", "name") VALUES ('alice@example.com', 'Alice')`)
		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}

		// Verify default values
		var active bool
		var name string
		err = db.QueryRow(`SELECT "active", "name" FROM "users" WHERE "email" = 'alice@example.com'`).Scan(&active, &name)
		if err != nil {
			t.Fatalf("failed to query inserted user: %v", err)
		}
		if !active || name != "Alice" {
			t.Errorf("expected active=true, name='Alice', got active=%t, name=%q", active, name)
		}

		// Verify UUID default
		_, err = db.Exec(`INSERT INTO "Post" ("title") VALUES ('My first post')`)
		if err != nil {
			t.Fatalf("failed to insert post: %v", err)
		}
		var pid, title string
		err = db.QueryRow(`SELECT "id", "title" FROM "Post" WHERE "title" = 'My first post'`).Scan(&pid, &title)
		if err != nil {
			t.Fatalf("failed to query inserted post: %v", err)
		}
		if len(pid) == 0 || title != "My first post" {
			t.Errorf("expected non-empty uuid ID, got %q", pid)
		}
	})
}

// 2. Enum Handling (TestEnumMigration)
func TestEnumMigration(t *testing.T) {
	schemaText := `
	datasource db {
		provider = "PROVIDER"
	}
	enum Role {
		USER
		ADMIN
	}
	enum UserStatus {
		ACTIVE    @map("active_status")
		SUSPENDED @map("suspended_status")
		
		@@map("user_status_enum")
	}
	model Member {
		id     Int        @id @default(autoincrement())
		role   Role       @default(USER)
		status UserStatus @default(ACTIVE)
	}
	`

	t.Run("sqlite", func(t *testing.T) {
		db := openTestSQLite(t, "enum_migration_sqlite")
		sc := parseOrFail(t, strings.Replace(schemaText, "PROVIDER", "sqlite", 1))
		up, _ := diffOrFail(t, db, providers.Sqlite, sc, false)

		assertNotContains(t, up, "CREATE TYPE")
		assertContains(t, up, `CHECK ("role" IN ('USER', 'ADMIN'))`)
		assertContains(t, up, `CHECK ("status" IN ('active_status', 'suspended_status'))`)

		execOrFail(t, db, up)

		// Test insert default values
		_, err := db.Exec(`INSERT INTO "Member" DEFAULT VALUES`)
		if err != nil {
			t.Fatalf("failed to insert default member: %v", err)
		}
		var role, status string
		err = db.QueryRow(`SELECT "role", "status" FROM "Member"`).Scan(&role, &status)
		if err != nil {
			t.Fatalf("failed to query member: %v", err)
		}
		if role != "USER" || status != "active_status" {
			t.Errorf("expected role='USER', status='active_status', got role=%q, status=%q", role, status)
		}

		// Test invalid enum value
		_, err = db.Exec(`INSERT INTO "Member" ("role", "status") VALUES ('INVALID', 'active_status')`)
		if err == nil {
			t.Errorf("expected CHECK constraint violation for invalid role")
		}
	})

	t.Run("postgres", func(t *testing.T) {
		db := openTestPostgres(t, "test_enum_migration_pg")
		sc := parseOrFail(t, strings.Replace(schemaText, "PROVIDER", "postgresql", 1))
		up, _ := diffOrFail(t, db, providers.Postgres, sc, false)

		assertContains(t, up, `CREATE TYPE`)
		assertContains(t, up, `Role`)
		assertContains(t, up, `user_status_enum`)
		assertContains(t, up, `AS ENUM`)
		assertContains(t, up, `'active_status'`)
		assertContains(t, up, `'suspended_status'`)

		execOrFail(t, db, up)

		// Test insert default values
		_, err := db.Exec(`INSERT INTO "Member" DEFAULT VALUES`)
		if err != nil {
			t.Fatalf("failed to insert default member: %v", err)
		}
		var role, status string
		err = db.QueryRow(`SELECT "role", "status" FROM "Member"`).Scan(&role, &status)
		if err != nil {
			t.Fatalf("failed to query member: %v", err)
		}
		if role != "USER" || status != "active_status" {
			t.Errorf("expected role='USER', status='active_status', got role=%q, status=%q", role, status)
		}
	})
}

// 3. Relations & Foreign Keys (TestForeignKeyMigration)
func TestForeignKeyMigration(t *testing.T) {
	schemaText := `
	datasource db {
		provider = "PROVIDER"
	}
	model Node {
		id       Int   @id @default(autoincrement())
		parentId Int?
		parent   Node? @relation("NodeToNode", fields: [parentId], references: [id], onDelete: SetNull)
	}
	model User {
		id    Int    @id @default(autoincrement())
		name  String
	}
	model Post {
		id       Int  @id @default(autoincrement())
		authorId Int
		author   User @relation(fields: [authorId], references: [id], onDelete: Cascade)
	}
	model Order {
		id     Int  @id @default(autoincrement())
		userId Int
		user   User @relation(fields: [userId], references: [id], onDelete: Restrict)
	}
	`

	t.Run("sqlite", func(t *testing.T) {
		db := openTestSQLite(t, "fk_migration_sqlite")
		sc := parseOrFail(t, strings.Replace(schemaText, "PROVIDER", "sqlite", 1))
		up, _ := diffOrFail(t, db, providers.Sqlite, sc, false)

		assertContains(t, up, `parentId`)
		assertContains(t, up, `REFERENCES`)
		assertContains(t, up, `Node`)
		assertContains(t, up, `ON DELETE SET NULL`)

		assertContains(t, up, `authorId`)
		assertContains(t, up, `User`)
		assertContains(t, up, `ON DELETE CASCADE`)

		assertContains(t, up, `userId`)
		assertContains(t, up, `ON DELETE RESTRICT`)

		execOrFail(t, db, up)
	})

	t.Run("postgres", func(t *testing.T) {
		db := openTestPostgres(t, "test_fk_migration_pg")
		sc := parseOrFail(t, strings.Replace(schemaText, "PROVIDER", "postgresql", 1))
		up, _ := diffOrFail(t, db, providers.Postgres, sc, false)

		assertContains(t, up, `parentId`)
		assertContains(t, up, `REFERENCES`)
		assertContains(t, up, `Node`)
		assertContains(t, up, `ON DELETE SET NULL`)

		assertContains(t, up, `authorId`)
		assertContains(t, up, `User`)
		assertContains(t, up, `ON DELETE CASCADE`)

		assertContains(t, up, `userId`)
		assertContains(t, up, `ON DELETE RESTRICT`)

		execOrFail(t, db, up)
	})
}

// 4. Composite Constraints (TestCompositeConstraints)
func TestCompositeConstraints(t *testing.T) {
	schemaText := `
	datasource db {
		provider = "PROVIDER"
	}
	model Group {
		id   Int    @id @default(autoincrement())
		name String
	}
	model User {
		id   Int    @id @default(autoincrement())
		name String
	}
	model GroupMember {
		groupId Int    @map("group_id")
		userId  Int    @map("user_id")
		role    String @map("member_role")

		@@id([groupId, userId])
		@@unique([userId, role], name: "custom_user_role_unique")
		@@unique([groupId, role])
		@@index([groupId, role])
		@@index([role], name: "custom_role_idx")
	}
	`

	t.Run("sqlite", func(t *testing.T) {
		db := openTestSQLite(t, "composite_constraints_sqlite")
		sc := parseOrFail(t, strings.Replace(schemaText, "PROVIDER", "sqlite", 1))
		up, _ := diffOrFail(t, db, providers.Sqlite, sc, false)

		assertContains(t, up, `PRIMARY KEY`)
		assertContains(t, up, `group_id`)
		assertContains(t, up, `user_id`)
		assertContains(t, up, `custom_user_role_unique`)
		assertContains(t, up, `GroupMember_group_id_member_role_key`)
		assertContains(t, up, `GroupMember_group_id_member_role_idx`)
		assertContains(t, up, `custom_role_idx`)

		execOrFail(t, db, up)
	})

	t.Run("postgres", func(t *testing.T) {
		db := openTestPostgres(t, "test_composite_constraints_pg")
		sc := parseOrFail(t, strings.Replace(schemaText, "PROVIDER", "postgresql", 1))
		up, _ := diffOrFail(t, db, providers.Postgres, sc, false)

		assertContains(t, up, `PRIMARY KEY`)
		assertContains(t, up, `group_id`)
		assertContains(t, up, `user_id`)
		assertContains(t, up, `custom_user_role_unique`)
		assertContains(t, up, `GroupMember_group_id_member_role_key`)
		assertContains(t, up, `GroupMember_group_id_member_role_idx`)
		assertContains(t, up, `custom_role_idx`)

		execOrFail(t, db, up)
	})
}

// 5. Incremental Migrations (TestIncrementalMigration)
func TestIncrementalMigration(t *testing.T) {
	db := openTestSQLite(t, "incremental_migration_sqlite")

	// V1
	v1Schema := `
	datasource db {
		provider = "sqlite"
	}
	model User {
		id    Int    @id @default(autoincrement())
		email String
		name  String
	}
	`
	sc1 := parseOrFail(t, v1Schema)
	up1, _ := diffOrFail(t, db, providers.Sqlite, sc1, false)
	assertContains(t, up1, "CREATE TABLE")
	assertContains(t, up1, "User")
	execOrFail(t, db, up1)

	// V2: Add Column
	v2Schema := `
	datasource db {
		provider = "sqlite"
	}
	model User {
		id    Int    @id @default(autoincrement())
		email String
		name  String
		age   Int?
	}
	`
	sc2 := parseOrFail(t, v2Schema)
	up2, down2 := diffOrFail(t, db, providers.Sqlite, sc2, false)
	assertContains(t, up2, "ADD COLUMN")
	assertContains(t, up2, "age")

	// Down from V2 to V1 drops the column. In SQLite, this is a table recreation.
	assertContains(t, down2, "new_User")
	assertNotContains(t, down2, "age")
	execOrFail(t, db, up2)

	// V3: Create dependent table
	v3Schema := `
	datasource db {
		provider = "sqlite"
	}
	model User {
		id    Int    @id @default(autoincrement())
		email String
		name  String
		age   Int?
	}
	model Post {
		id       Int    @id @default(autoincrement())
		title    String
		authorId Int
		author   User   @relation(fields: [authorId], references: [id])
	}
	`
	sc3 := parseOrFail(t, v3Schema)
	up3, down3 := diffOrFail(t, db, providers.Sqlite, sc3, false)
	assertContains(t, up3, "CREATE TABLE")
	assertContains(t, up3, "Post")

	assertContains(t, down3, "DROP TABLE")
	assertContains(t, down3, "Post")
	execOrFail(t, db, up3)

	// V4: Add Unique Index
	v4Schema := `
	datasource db {
		provider = "sqlite"
	}
	model User {
		id    Int    @id @default(autoincrement())
		email String @unique
		name  String
		age   Int?
	}
	model Post {
		id       Int    @id @default(autoincrement())
		title    String
		authorId Int
		author   User   @relation(fields: [authorId], references: [id])
	}
	`
	sc4 := parseOrFail(t, v4Schema)
	up4, down4 := diffOrFail(t, db, providers.Sqlite, sc4, false)
	assertContains(t, up4, "CREATE UNIQUE INDEX")
	assertContains(t, down4, "DROP INDEX")
	execOrFail(t, db, up4)

	// V5: Add Index
	v5Schema := `
	datasource db {
		provider = "sqlite"
	}
	model User {
		id    Int    @id @default(autoincrement())
		email String @unique
		name  String
		age   Int?
		
		@@index([name])
	}
	model Post {
		id       Int    @id @default(autoincrement())
		title    String
		authorId Int
		author   User   @relation(fields: [authorId], references: [id])
	}
	`
	sc5 := parseOrFail(t, v5Schema)
	up5, down5 := diffOrFail(t, db, providers.Sqlite, sc5, false)
	assertContains(t, up5, "CREATE INDEX")
	assertContains(t, down5, "DROP INDEX")
	execOrFail(t, db, up5)
}

// 6. Column Rename Detection (TestRenameDetection)
func TestRenameDetection(t *testing.T) {
	v1Schema := `
	datasource db {
		provider = "sqlite"
	}
	model User {
		id       Int    @id @default(autoincrement())
		phoneNum String
		addr     String
	}
	`
	v2Schema := `
	datasource db {
		provider = "sqlite"
	}
	model User {
		id          Int    @id @default(autoincrement())
		phoneNumber String
		address     String
	}
	`

	t.Run("accept_rename", func(t *testing.T) {
		db := openTestSQLite(t, "rename_accept")
		sc1 := parseOrFail(t, v1Schema)
		up1, _ := diffOrFail(t, db, providers.Sqlite, sc1, false)
		execOrFail(t, db, up1)

		// Mock stdin to accept both renames (phoneNum -> phoneNumber, addr -> address)
		mockStdin(t, "y\ny\n")

		sc2 := parseOrFail(t, v2Schema)
		up2, down2 := diffOrFail(t, db, providers.Sqlite, sc2, true)

		// SQLite uses table recreation. Verify that columns are copied from old to new names.
		assertContains(t, up2, "new_User")
		assertContains(t, up2, "phoneNumber")
		assertContains(t, up2, "phoneNum") // Source in SELECT

		// Verify down SQL contains reverse RENAME (copies phoneNumber back to phoneNum)
		assertContains(t, down2, "new_User")
		assertContains(t, down2, "phoneNum")
		assertContains(t, down2, "phoneNumber") // Source in SELECT

		execOrFail(t, db, up2)
	})

	t.Run("refuse_rename", func(t *testing.T) {
		db := openTestSQLite(t, "rename_refuse")
		sc1 := parseOrFail(t, v1Schema)
		up1, _ := diffOrFail(t, db, providers.Sqlite, sc1, false)
		execOrFail(t, db, up1)

		mockStdin(t, "n\nn\n")

		sc2 := parseOrFail(t, v2Schema)
		up2, _ := diffOrFail(t, db, providers.Sqlite, sc2, true)

		// Since rename was refused, it should recreate User table without copying data from the old phoneNum/addr
		assertContains(t, up2, "new_User")
		assertNotContains(t, up2, "`phoneNum`")
		assertNotContains(t, up2, "`addr`")
	})

	t.Run("non_interactive", func(t *testing.T) {
		db := openTestSQLite(t, "rename_non_interactive")
		sc1 := parseOrFail(t, v1Schema)
		up1, _ := diffOrFail(t, db, providers.Sqlite, sc1, false)
		execOrFail(t, db, up1)

		sc2 := parseOrFail(t, v2Schema)
		up2, _ := diffOrFail(t, db, providers.Sqlite, sc2, false)

		// Non-interactive also defaults to drop & add (recreate without copying)
		assertContains(t, up2, "new_User")
		assertNotContains(t, up2, "`phoneNum`")
		assertNotContains(t, up2, "`addr`")
	})
}

// 7. Destructive Changes (TestDestructiveChanges)
func TestDestructiveChanges(t *testing.T) {
	db := openTestSQLite(t, "destructive_changes")

	v1Schema := `
	datasource db {
		provider = "sqlite"
	}
	model User {
		id     Int    @id @default(autoincrement())
		email  String @unique
		age    Int
		name   String
		posts  Post[]
		
		@@index([name])
	}
	model Post {
		id       Int  @id @default(autoincrement())
		authorId Int
		author   User @relation(fields: [authorId], references: [id])
	}
	`
	sc1 := parseOrFail(t, v1Schema)
	up1, _ := diffOrFail(t, db, providers.Sqlite, sc1, false)
	execOrFail(t, db, up1)

	// V2: Drop age column, Post table, unique constraint on email, index on name
	v2Schema := `
	datasource db {
		provider = "sqlite"
	}
	model User {
		id     Int    @id @default(autoincrement())
		email  String
		name   String
	}
	`
	sc2 := parseOrFail(t, v2Schema)
	up2, _ := diffOrFail(t, db, providers.Sqlite, sc2, false)

	// User table is recreated without "age", and Post table is dropped.
	assertContains(t, up2, "new_User")
	assertContains(t, up2, "DROP TABLE")
	assertContains(t, up2, "User")
	assertContains(t, up2, "Post")

	execOrFail(t, db, up2)
}

// 8. Data Preservation Round-Trip (TestRoundTripDataPreservation)
func TestRoundTripDataPreservation(t *testing.T) {
	db := openTestSQLite(t, "roundtrip_preservation")

	v1Schema := `
	datasource db {
		provider = "sqlite"
	}
	model User {
		id    Int    @id @default(autoincrement())
		email String @unique
		name  String
	}
	`
	sc1 := parseOrFail(t, v1Schema)
	up1, _ := diffOrFail(t, db, providers.Sqlite, sc1, false)
	execOrFail(t, db, up1)

	// Insert test data
	_, err := db.Exec(`INSERT INTO "User" ("email", "name") VALUES ('alice@example.com', 'Alice')`)
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	// V2: Add role column
	v2Schema := `
	datasource db {
		provider = "sqlite"
	}
	model User {
		id    Int    @id @default(autoincrement())
		email String @unique
		name  String
		role  String @default("USER")
	}
	`
	sc2 := parseOrFail(t, v2Schema)
	up2, down2 := diffOrFail(t, db, providers.Sqlite, sc2, false)
	execOrFail(t, db, up2)

	// Verify data survived and role has default value
	var email, name, role string
	err = db.QueryRow(`SELECT "email", "name", "role" FROM "User" WHERE "email" = 'alice@example.com'`).Scan(&email, &name, &role)
	if err != nil {
		t.Fatalf("failed to query user after up migration: %v", err)
	}
	if email != "alice@example.com" || name != "Alice" || role != "USER" {
		t.Errorf("data mismatch after up: got email=%q, name=%q, role=%q", email, name, role)
	}

	// Apply down migration
	execOrFail(t, db, down2)

	// Verify original schema restored (querying role should fail)
	rows, err := db.Query(`SELECT role FROM "User"`)
	if err == nil {
		rows.Close()
		t.Errorf("expected query on dropped column 'role' to fail, but it succeeded")
	}

	// Verify data still intact
	err = db.QueryRow(`SELECT "email", "name" FROM "User" WHERE "email" = 'alice@example.com'`).Scan(&email, &name)
	if err != nil {
		t.Fatalf("failed to query user after down migration: %v", err)
	}
	if email != "alice@example.com" || name != "Alice" {
		t.Errorf("data mismatch after down: got email=%q, name=%q", email, name)
	}
}

// 9. Complex Real-World Schema ("E-Commerce") (TestComplexECommerceSchema)
func TestComplexECommerceSchema(t *testing.T) {
	schemaTemplate := `
	datasource db {
		provider = "PROVIDER"
	}

	enum OrderStatus {
		PENDING @map("pending_payment")
		SHIPPED @map("shipped_to_customer")
		DELIVERED
		
		@@map("order_status_enum")
	}

	model User {
		id        Int      @id @default(autoincrement()) @map("user_id")
		email     String   @unique @map("user_email")
		createdAt DateTime @default(now()) @map("created_at")
		
		// Self-referential manager relation
		managerId Int?       @map("manager_id")
		manager   User?      @relation("UserManager", fields: [managerId], references: [id], onDelete: SetNull)

		@@map("users")
	}

	model Product {
		id        Int     @id @default(autoincrement()) @map("product_id")
		name      String  @map("product_name")
		sku       String  @unique @map("product_sku")
		price     Float   @map("product_price")
		active    Boolean @default(true)

		@@map("products")
	}

	model Order {
		id        ORDER_ID_TYPE   @id ORDER_ID_DEFAULT @map("order_id")
		userId    Int             @map("fk_user_id")
		status    OrderStatus     @default(PENDING) @map("order_status")
		total     Float           @default(0.0) @map("order_total")
		
		user      User            @relation(fields: [userId], references: [id], onDelete: Restrict)

		@@map("orders")
	}

	model OrderItem {
		orderId   ORDER_ID_TYPE   @map("fk_order_id")
		productId Int             @map("fk_product_id")
		quantity  Int             @default(1)
		
		order     Order           @relation(fields: [orderId], references: [id], onDelete: Cascade)
		product   Product         @relation(fields: [productId], references: [id], onDelete: Cascade)

		@@id([orderId, productId])
		@@map("order_items")
	}

	model Category {
		id        Int             @id @default(autoincrement())
		name      String
		
		@@unique([name])
		@@index([name])
	}

	model CategoryToProduct {
		categoryId Int
		productId  Int
		
		@@id([categoryId, productId])
	}
	`

	t.Run("sqlite", func(t *testing.T) {
		db := openTestSQLite(t, "complex_ecommerce_sqlite")
		schemaText := strings.ReplaceAll(schemaTemplate, "PROVIDER", "sqlite")
		schemaText = strings.ReplaceAll(schemaText, "ORDER_ID_TYPE", "Int")
		schemaText = strings.ReplaceAll(schemaText, "ORDER_ID_DEFAULT", "@default(autoincrement())")

		sc := parseOrFail(t, schemaText)
		up, down := diffOrFail(t, db, providers.Sqlite, sc, false)

		execOrFail(t, db, up)
		execOrFail(t, db, down)
	})

	t.Run("postgres", func(t *testing.T) {
		db := openTestPostgres(t, "test_complex_ecommerce_pg")
		schemaText := strings.ReplaceAll(schemaTemplate, "PROVIDER", "postgresql")
		schemaText = strings.ReplaceAll(schemaText, "ORDER_ID_TYPE", "String")
		schemaText = strings.ReplaceAll(schemaText, "ORDER_ID_DEFAULT", "@default(uuid())")

		sc := parseOrFail(t, schemaText)
		up, down := diffOrFail(t, db, providers.Postgres, sc, false)

		execOrFail(t, db, up)
		execOrFail(t, db, down)
	})
}

// 10. Statement Ordering (TestSortMigrationStatements)
func TestSortMigrationStatements(t *testing.T) {
	stmts := []string{
		`CREATE UNIQUE INDEX "User_email_phoneNum_key" ON "User" ("email", "phoneNum");`,
		`DROP INDEX "User_email_phoneNumber_key";`,
		`ALTER TABLE "User" RENAME COLUMN "phoneNumber" TO "phoneNum";`,
		`CREATE TABLE "Post" ("id" INTEGER PRIMARY KEY);`,
		`ALTER TABLE "User" DROP COLUMN "age";`,
		`CREATE INDEX "Post_title_idx" ON "Post" ("title");`,
	}

	sorted := SortMigrationStatements(stmts)
	if len(sorted) != 6 {
		t.Fatalf("expected 6 statements, got %d", len(sorted))
	}

	// Verify that all drops come first
	for i := 0; i < 2; i++ {
		trimmed := strings.ToUpper(sorted[i])
		isDrop := strings.Contains(trimmed, "DROP ")
		if !isDrop {
			t.Errorf("expected statement %d to be a DROP statement, got: %s", i+1, sorted[i])
		}
	}

	// Verify that alters & table creations come in the middle
	for i := 2; i < 4; i++ {
		trimmed := strings.ToUpper(sorted[i])
		isCreateIdx := strings.Contains(trimmed, "CREATE INDEX") || strings.Contains(trimmed, "CREATE UNIQUE INDEX")
		isDrop := strings.Contains(trimmed, "DROP ")
		if isCreateIdx || isDrop {
			t.Errorf("expected statement %d to be an ALTER or CREATE TABLE statement, got: %s", i+1, sorted[i])
		}
	}

	// Verify that index creations come last
	for i := 4; i < 6; i++ {
		trimmed := strings.ToUpper(sorted[i])
		isCreateIdx := strings.Contains(trimmed, "CREATE INDEX") || strings.Contains(trimmed, "CREATE UNIQUE INDEX")
		if !isCreateIdx {
			t.Errorf("expected statement %d to be a CREATE INDEX statement, got: %s", i+1, sorted[i])
		}
	}
}

// 11. No-Op Migration (TestNoOpMigration)
func TestNoOpMigration(t *testing.T) {
	schemaText := `
	datasource db {
		provider = "PROVIDER"
	}
	model User {
		id    Int    @id @default(autoincrement())
		email String @unique
	}
	`

	t.Run("sqlite", func(t *testing.T) {
		db := openTestSQLite(t, "noop_sqlite")
		sc := parseOrFail(t, strings.Replace(schemaText, "PROVIDER", "sqlite", 1))
		up1, _ := diffOrFail(t, db, providers.Sqlite, sc, false)
		execOrFail(t, db, up1)

		up2, down2 := diffOrFail(t, db, providers.Sqlite, sc, false)
		if up2 != "" || down2 != "" {
			t.Errorf("expected no-op migration, got up=%q, down=%q", up2, down2)
		}
	})

	t.Run("postgres", func(t *testing.T) {
		db := openTestPostgres(t, "test_noop_pg")
		sc := parseOrFail(t, strings.Replace(schemaText, "PROVIDER", "postgresql", 1))
		up1, _ := diffOrFail(t, db, providers.Postgres, sc, false)
		execOrFail(t, db, up1)

		up2, down2 := diffOrFail(t, db, providers.Postgres, sc, false)
		if up2 != "" || down2 != "" {
			t.Errorf("expected no-op migration, got up=%q, down=%q", up2, down2)
		}
	})
}

// 12. Idempotency (TestMigrationIdempotency)
func TestMigrationIdempotency(t *testing.T) {
	db := openTestSQLite(t, "idempotency_sqlite")
	schemaText := `
	datasource db {
		provider = "sqlite"
	}
	model User {
		id    Int    @id @default(autoincrement())
		email String @unique
	}
	`
	sc := parseOrFail(t, schemaText)
	up1, _ := diffOrFail(t, db, providers.Sqlite, sc, false)
	execOrFail(t, db, up1)

	// Second DiffAndPlan should return empty strings
	up2, down2 := diffOrFail(t, db, providers.Sqlite, sc, false)
	if up2 != "" || down2 != "" {
		t.Errorf("expected second plan to be empty, got up=%q, down=%q", up2, down2)
	}
}
