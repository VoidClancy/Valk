package main

import (
	"benchmark/valk"
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type DialectConfig struct {
	Name                string
	Driver              string
	DSN                 string
	QuoteChar           byte
	PlaceholderFmt      string // "" means "?"
	SupportsReturning   bool
	ConflictKeyword     string // "ON CONFLICT" or "ON DUPLICATE KEY"
	ConflictIgnore      string // "DO NOTHING" or ""
	ConflictUpdate      string // "DO UPDATE SET" or "UPDATE"
	ConflictExcluded    string // "EXCLUDED." or "VALUES("
	ConflictExcludedEnd string // "" or ")"
	SchemaQuery         string
}

func (d DialectConfig) Quote(ident string) string {
	return string(d.QuoteChar) + ident + string(d.QuoteChar)
}

func (d DialectConfig) BindVar(idx int) string {
	if d.PlaceholderFmt == "" {
		return "?"
	}
	return d.PlaceholderFmt + strconv.Itoa(idx)
}

func (d DialectConfig) Placeholders(startIdx, count int) string {
	var sb strings.Builder
	sb.WriteByte('(')
	for i := 0; i < count; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(d.BindVar(startIdx + i))
	}
	sb.WriteByte(')')
	return sb.String()
}

const sqliteSchema = ""

const postgresSchema = ""

const seedCount = 1000

type DBTX interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

var (
	activeDialect DialectConfig

	rawQueryCreate              string
	rawQueryCreateMany          string
	rawQueryCreateManyAndReturn string
	rawQueryFindUnique          string
	rawQueryFindFirst           string
	rawQueryFindMany            string
	rawQueryUpsert              string
	rawQueryDeepRelation        string
)

func TestMain(m *testing.M) {
	initDialect()
	os.Exit(m.Run())
}

func initDialect() {
	dbType := os.Getenv("BENCH_DB")
	if dbType == "" {
		dbType = "sqlite"
	}

	switch dbType {
	case "postgres", "postgresql":
		activeDialect = DialectConfig{
			Name:                "postgres",
			Driver:              "postgres",
			DSN:                 getPostgresDSN(),
			QuoteChar:           '"',
			PlaceholderFmt:      "$",
			SupportsReturning:   true,
			ConflictKeyword:     "ON CONFLICT",
			ConflictIgnore:      "DO NOTHING",
			ConflictUpdate:      "DO UPDATE SET",
			ConflictExcluded:    "EXCLUDED.",
			ConflictExcludedEnd: "",
			SchemaQuery:         postgresSchema,
		}
	default: // sqlite
		activeDialect = DialectConfig{
			Name:                "sqlite",
			Driver:              "sqlite3", // Using github.com/mattn/go-sqlite3 Cgo driver
			DSN:                 "file:benchmark?mode=memory&cache=shared",
			QuoteChar:           '"',
			PlaceholderFmt:      "",
			SupportsReturning:   true,
			ConflictKeyword:     "ON CONFLICT",
			ConflictIgnore:      "DO NOTHING",
			ConflictUpdate:      "DO UPDATE SET",
			ConflictExcluded:    "EXCLUDED.",
			ConflictExcludedEnd: "",
			SchemaQuery:         sqliteSchema,
		}
	}

	initQueries(activeDialect)
}

func getPostgresDSN() string {
	if url := os.Getenv("PG_DATABASE_URL"); url != "" {
		return url
	}
	localDSN := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", localDSN)
	if err == nil {
		err = db.Ping()
		db.Close()
		if err == nil {
			return localDSN
		}
	}
	return "postgres://testuser:testpassword@localhost:5432/valk_test?sslmode=disable"
}

func initQueries(d DialectConfig) {
	rawQueryCreate = fmt.Sprintf(
		`INSERT INTO %s (%s, %s, %s, %s) VALUES (%s, %s, %s, %s)`,
		d.Quote("User"),
		d.Quote("id"), d.Quote("email"), d.Quote("phoneNum"), d.Quote("role"),
		d.BindVar(1), d.BindVar(2), d.BindVar(3), d.BindVar(4),
	)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`INSERT INTO %s (%s, %s, %s, %s, %s) VALUES `,
		d.Quote("User"),
		d.Quote("id"), d.Quote("email"), d.Quote("phoneNum"), d.Quote("role"), d.Quote("loginCount"),
	))
	for i := 0; i < 10; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(d.Placeholders(i*5+1, 5))
	}
	rawQueryCreateMany = sb.String()

	if d.SupportsReturning {
		rawQueryCreateManyAndReturn = rawQueryCreateMany + fmt.Sprintf(
			` RETURNING %s, %s, %s, %s, %s`,
			d.Quote("id"), d.Quote("email"), d.Quote("phoneNum"), d.Quote("role"), d.Quote("loginCount"),
		)
	}

	rawQueryFindUnique = fmt.Sprintf(
		`SELECT %s, %s, %s, %s, %s, %s, %s, %s FROM %s WHERE %s = %s`,
		d.Quote("id"), d.Quote("email"), d.Quote("phoneNum"), d.Quote("password"), d.Quote("role"), d.Quote("roleOptional"), d.Quote("loginCount"), d.Quote("referredById"),
		d.Quote("User"),
		d.Quote("email"),
		d.BindVar(1),
	)

	rawQueryFindFirst = fmt.Sprintf(
		`SELECT %s, %s, %s, %s, %s, %s, %s, %s FROM %s WHERE %s = %s LIMIT 1`,
		d.Quote("id"), d.Quote("email"), d.Quote("phoneNum"), d.Quote("password"), d.Quote("role"), d.Quote("roleOptional"), d.Quote("loginCount"), d.Quote("referredById"),
		d.Quote("User"),
		d.Quote("email"),
		d.BindVar(1),
	)

	rawQueryFindMany = fmt.Sprintf(
		`SELECT %s, %s, %s, %s, %s, %s, %s, %s FROM %s ORDER BY %s LIMIT 10 OFFSET %s`,
		d.Quote("id"), d.Quote("email"), d.Quote("phoneNum"), d.Quote("password"), d.Quote("role"), d.Quote("roleOptional"), d.Quote("loginCount"), d.Quote("referredById"),
		d.Quote("User"),
		d.Quote("id"),
		d.BindVar(1),
	)

	switch d.ConflictKeyword {
	case "ON CONFLICT":
		rawQueryUpsert = fmt.Sprintf(
			`INSERT INTO %s (%s, %s, %s, %s, %s) VALUES (%s, %s, %s, %s, %s) `+
				`ON CONFLICT(%s) DO UPDATE SET %s = EXCLUDED.%s, %s = EXCLUDED.%s, %s = EXCLUDED.%s`,
			d.Quote("User"),
			d.Quote("id"), d.Quote("email"), d.Quote("phoneNum"), d.Quote("role"), d.Quote("loginCount"),
			d.BindVar(1), d.BindVar(2), d.BindVar(3), d.BindVar(4), d.BindVar(5),
			d.Quote("email"),
			d.Quote("phoneNum"), d.Quote("phoneNum"),
			d.Quote("role"), d.Quote("role"),
			d.Quote("loginCount"), d.Quote("loginCount"),
		)
	case "ON DUPLICATE KEY":
		rawQueryUpsert = fmt.Sprintf(
			`INSERT INTO %s (%s, %s, %s, %s, %s) VALUES (%s, %s, %s, %s, %s) `+
				`ON DUPLICATE KEY UPDATE %s = VALUES(%s), %s = VALUES(%s), %s = VALUES(%s)`,
			d.Quote("User"),
			d.Quote("id"), d.Quote("email"), d.Quote("phoneNum"), d.Quote("role"), d.Quote("loginCount"),
			d.BindVar(1), d.BindVar(2), d.BindVar(3), d.BindVar(4), d.BindVar(5),
			d.Quote("phoneNum"), d.Quote("phoneNum"),
			d.Quote("role"), d.Quote("role"),
			d.Quote("loginCount"), d.Quote("loginCount"),
		)
	}

	userCols := []string{"id", "email", "phoneNum", "password", "role", "roleOptional", "loginCount", "referredById"}
	var selectCols []string
	for _, col := range userCols {
		selectCols = append(selectCols, fmt.Sprintf("u.%s", d.Quote(col)))
	}
	for _, col := range userCols {
		selectCols = append(selectCols, fmt.Sprintf("r.%s", d.Quote(col)))
	}
	for _, col := range userCols {
		selectCols = append(selectCols, fmt.Sprintf("rr.%s", d.Quote(col)))
	}

	rawQueryDeepRelation = fmt.Sprintf(
		`SELECT %s FROM %s u LEFT JOIN %s r ON u.%s = r.%s LEFT JOIN %s rr ON r.%s = rr.%s WHERE u.%s = %s`,
		strings.Join(selectCols, ", "),
		d.Quote("User"), d.Quote("User"), d.Quote("referredById"), d.Quote("id"),
		d.Quote("User"), d.Quote("referredById"), d.Quote("id"),
		d.Quote("email"),
		d.BindVar(1),
	)
}

func openDB(b *testing.B) *sql.DB {
	b.Helper()
	db, err := sql.Open(activeDialect.Driver, activeDialect.DSN)
	if err != nil {
		b.Fatal(err)
	}
	db.SetMaxOpenConns(80)
	db.SetMaxIdleConns(80)
	return db
}

func createSQLiteSchema(sqlDB *sql.DB) {
	ctx := context.Background()
	stmts := []string{
		`DROP TABLE IF EXISTS "Profile"`,
		`DROP TABLE IF EXISTS "Comment"`,
		`DROP TABLE IF EXISTS "CategoryToPost"`,
		`DROP TABLE IF EXISTS "Category"`,
		`DROP TABLE IF EXISTS "Post"`,
		`DROP TABLE IF EXISTS "User"`,
		`CREATE TABLE "User" (
			"id" text NOT NULL PRIMARY KEY,
			"email" text NOT NULL UNIQUE,
			"phoneNum" text NOT NULL UNIQUE,
			"password" text,
			"role" text NOT NULL DEFAULT 'STUDENT',
			"roleOptional" text,
			"loginCount" integer NOT NULL DEFAULT 0,
			"referredById" text REFERENCES "User"("id")
		)`,
		`CREATE TABLE "Post" (
			"id" text NOT NULL PRIMARY KEY,
			"title" text NOT NULL,
			"content" text,
			"published" integer NOT NULL DEFAULT 0,
			"authorId" text NOT NULL REFERENCES "User"("id")
		)`,
		`CREATE TABLE "Category" (
			"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
			"name" text NOT NULL UNIQUE
		)`,
		`CREATE TABLE "CategoryToPost" (
			"postId" text NOT NULL REFERENCES "Post"("id"),
			"categoryId" integer NOT NULL REFERENCES "Category"("id"),
			PRIMARY KEY ("postId", "categoryId")
		)`,
		`CREATE TABLE "Comment" (
			"id" text NOT NULL PRIMARY KEY,
			"textify" integer NOT NULL,
			"dummy3" text NOT NULL,
			"dummy1" integer NOT NULL,
			"dummy2" text NOT NULL,
			"postId" text NOT NULL REFERENCES "Post"("id"),
			"authorId" text NOT NULL REFERENCES "User"("id"),
			"meta" text
		)`,
		`CREATE TABLE "Profile" (
			"id" text NOT NULL PRIMARY KEY,
			"bio" text,
			"userId" text NOT NULL UNIQUE REFERENCES "User"("id") ON DELETE CASCADE,
			"createdAt" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE UNIQUE INDEX "emailPhone" ON "User"("email", "phoneNum")`,
	}
	for _, stmt := range stmts {
		if _, err := sqlDB.ExecContext(ctx, stmt); err != nil {
			panic(fmt.Sprintf("SQLite schema: %s: %v", stmt, err))
		}
	}
}

func resetPostgres(db *sql.DB) {
	_, _ = db.ExecContext(context.Background(), `
		SELECT pg_terminate_backend(pid)
		FROM pg_stat_activity
		WHERE datname = current_database()
		  AND pid <> pg_backend_pid()
	`)
	_, _ = db.ExecContext(context.Background(), "DROP EXTENSION IF EXISTS citext CASCADE; DROP EXTENSION IF EXISTS hstore CASCADE; DROP EXTENSION IF EXISTS ltree CASCADE;")
	_, err := db.ExecContext(context.Background(), "DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	if err != nil {
		panic(fmt.Sprintf("failed to reset postgres database: %v", err))
	}
}

func createSchema(db DBTX) {
	var sqlDB *sql.DB
	if concrete, ok := db.(*sql.DB); ok {
		sqlDB = concrete
	} else {
		panic("expected *sql.DB")
	}

	if err := sqlDB.PingContext(context.Background()); err != nil {
		panic(err)
	}

	if activeDialect.Name == "postgres" {
		resetPostgres(sqlDB)

		valkDB, err := valk.Open(activeDialect.Driver, activeDialect.DSN)
		if err != nil {
			panic(fmt.Sprintf("failed to open valk for migrations: %v", err))
		}
		defer valkDB.Close()

		if err := valkDB.RunMigrations(context.Background()); err != nil {
			panic(fmt.Sprintf("run migrations: %v", err))
		}
	} else {
		createSQLiteSchema(sqlDB)
	}
}

func seedData(db DBTX, prefix string) {
	for i := range seedCount {
		_, err := db.ExecContext(context.Background(),
			fmt.Sprintf(`INSERT INTO %s (%s, %s, %s, %s, %s) VALUES (%s, %s, %s, %s, %s)`,
				activeDialect.Quote("User"),
				activeDialect.Quote("id"), activeDialect.Quote("email"), activeDialect.Quote("phoneNum"), activeDialect.Quote("password"), activeDialect.Quote("role"),
				activeDialect.BindVar(1), activeDialect.BindVar(2), activeDialect.BindVar(3), activeDialect.BindVar(4), activeDialect.BindVar(5),
			),
			fmt.Sprintf("%s-id-%d", prefix, i),
			fmt.Sprintf("%s-user-%d@example.com", prefix, i),
			fmt.Sprintf("%s-phone-%d", prefix, i),
			nil, "STUDENT",
		)
		if err != nil {
			panic(fmt.Sprintf("seed %s: %v", prefix, err))
		}
	}
}

func seedRelations(db DBTX, prefix string) {
	for i := 0; i < 500; i++ {
		// 1. Parent referrer
		parentID := fmt.Sprintf("%s-parent-id-%d", prefix, i)
		parentEmail := fmt.Sprintf("%s-parent-%d@example.com", prefix, i)
		parentPhone := fmt.Sprintf("%s-parent-phone-%d", prefix, i)
		_, err := db.ExecContext(context.Background(),
			fmt.Sprintf(`INSERT INTO %s (%s, %s, %s, %s, %s) VALUES (%s, %s, %s, %s, %s)`,
				activeDialect.Quote("User"),
				activeDialect.Quote("id"), activeDialect.Quote("email"), activeDialect.Quote("phoneNum"), activeDialect.Quote("password"), activeDialect.Quote("role"),
				activeDialect.BindVar(1), activeDialect.BindVar(2), activeDialect.BindVar(3), activeDialect.BindVar(4), activeDialect.BindVar(5),
			),
			parentID, parentEmail, parentPhone, nil, "STUDENT",
		)
		if err != nil {
			panic(err)
		}

		// 2. Child referred
		childID := fmt.Sprintf("%s-child-id-%d", prefix, i)
		childEmail := fmt.Sprintf("%s-child-%d@example.com", prefix, i)
		childPhone := fmt.Sprintf("%s-child-phone-%d", prefix, i)
		_, err = db.ExecContext(context.Background(),
			fmt.Sprintf(`INSERT INTO %s (%s, %s, %s, %s, %s, %s) VALUES (%s, %s, %s, %s, %s, %s)`,
				activeDialect.Quote("User"),
				activeDialect.Quote("id"), activeDialect.Quote("email"), activeDialect.Quote("phoneNum"), activeDialect.Quote("password"), activeDialect.Quote("role"), activeDialect.Quote("referredById"),
				activeDialect.BindVar(1), activeDialect.BindVar(2), activeDialect.BindVar(3), activeDialect.BindVar(4), activeDialect.BindVar(5), activeDialect.BindVar(6),
			),
			childID, childEmail, childPhone, nil, "STUDENT", parentID,
		)
		if err != nil {
			panic(err)
		}

		// 3. Grandchild referred
		grandchildID := fmt.Sprintf("%s-grand-id-%d", prefix, i)
		grandchildEmail := fmt.Sprintf("%s-grand-%d@example.com", prefix, i)
		grandchildPhone := fmt.Sprintf("%s-grand-phone-%d", prefix, i)
		_, err = db.ExecContext(context.Background(),
			fmt.Sprintf(`INSERT INTO %s (%s, %s, %s, %s, %s, %s) VALUES (%s, %s, %s, %s, %s, %s)`,
				activeDialect.Quote("User"),
				activeDialect.Quote("id"), activeDialect.Quote("email"), activeDialect.Quote("phoneNum"), activeDialect.Quote("password"), activeDialect.Quote("role"), activeDialect.Quote("referredById"),
				activeDialect.BindVar(1), activeDialect.BindVar(2), activeDialect.BindVar(3), activeDialect.BindVar(4), activeDialect.BindVar(5), activeDialect.BindVar(6),
			),
			grandchildID, grandchildEmail, grandchildPhone, nil, "STUDENT", childID,
		)
		if err != nil {
			panic(err)
		}
	}
}
