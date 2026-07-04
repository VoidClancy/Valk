package valkyrie

import (
	"context"
	"crypto/rand"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
)

var _ = time.Time{}
var _ = json.RawMessage{}
var _ = strings.Join
var _ = uuid.New
var _ = rand.Read
var _ = strconv.AppendUint

//go:embed migrations/*.sql
var migrationsFS embed.FS

func generateCUID() string {
	now := uint64(time.Now().UnixMilli())
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	buf := make([]byte, 1, 33)
	buf[0] = 'c'
	buf = strconv.AppendUint(buf, now, 16)

	const hextable = "0123456789abcdef"
	for _, v := range b {
		buf = append(buf, hextable[v>>4], hextable[v&0x0f])
	}
	return string(buf)
}

func generateUUID() string {
	return uuid.New().String()
}

type UserRoleType string

const (
	UserRoleTypeAdmin   UserRoleType = "ADMIN"
	UserRoleTypeStudent UserRoleType = "student"
	UserRoleTypeTeacher UserRoleType = "TEACHER"
)

type userRoleNamespace struct {
	Admin   UserRoleType
	Student UserRoleType
	Teacher UserRoleType
}

var UserRole = userRoleNamespace{
	Admin:   UserRoleTypeAdmin,
	Student: UserRoleTypeStudent,
	Teacher: UserRoleTypeTeacher,
}

type Dialect interface {
	Quote(ident string) string
	BindVar(idx int) string
	SupportsReturning() bool
}
type sqliteDialect struct{}

func (sqliteDialect) Quote(ident string) string { return `"` + ident + `"` }
func (sqliteDialect) BindVar(idx int) string    { return "?" }
func (sqliteDialect) SupportsReturning() bool   { return true }

type DBTX interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

type Queries struct {
	db             DBTX
	provider       string
	dialect        Dialect
	User           *UserDelegate
	Profile        *ProfileDelegate
	Post           *PostDelegate
	Comment        *CommentDelegate
	Category       *CategoryDelegate
	CategoryToPost *CategoryToPostDelegate
	UserRole       userRoleNamespace
}

type DB struct {
	*Queries
	sqlDB *sql.DB
}

// Open opens a database connection and returns a client.
func Open(provider, dataSourceName string) (*DB, error) {
	sqlDB, err := sql.Open(provider, dataSourceName)
	if err != nil {
		return nil, err
	}

	q := &Queries{
		db:       sqlDB,
		provider: provider,
		dialect:  sqliteDialect{},
		UserRole: UserRole,
	}
	q.User = &UserDelegate{client: q}
	q.Profile = &ProfileDelegate{client: q}
	q.Post = &PostDelegate{client: q}
	q.Comment = &CommentDelegate{client: q}
	q.Category = &CategoryDelegate{client: q}
	q.CategoryToPost = &CategoryToPostDelegate{client: q}
	return &DB{
		Queries: q,
		sqlDB:   sqlDB,
	}, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.sqlDB.Close()
}

// Raw returns the underlying *sql.DB connection pool.
func (db *DB) Raw() *sql.DB {
	return db.sqlDB
}

// RunMigrations runs all pending migrations from the embedded folder.
func (db *DB) RunMigrations(ctx context.Context) error {
	log.Println("Running migrations...")
	if err := goose.SetDialect(db.provider); err != nil {
		return err
	}
	goose.SetLogger(goose.NopLogger())
	goose.SetBaseFS(migrationsFS)
	err := goose.UpContext(ctx, db.sqlDB, "migrations")
	if err != nil {
		log.Printf("Migrations failed: %v", err)
		return err
	}
	log.Println("Migrations completed successfully.")
	return nil
}

func (q *Queries) bindVars(count int) string {
	if count <= 0 {
		return ""
	}
	var sb strings.Builder
	sb.Grow(count * 3)
	for i := 0; i < count; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(q.dialect.BindVar(i + 1))
	}
	return sb.String()
}

func (q *Queries) query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	log.Printf("[%s] SQL Query: %s | Args: %v", strings.ToUpper(q.provider), query, args)
	res, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("[%s] SQL Error: %v | Query: %s | Args: %v", strings.ToUpper(q.provider), err, query, args)
	}
	return res, err
}

func (q *Queries) queryRow(ctx context.Context, query string, args ...any) *sql.Row {
	log.Printf("[%s] SQL QueryRow: %s | Args: %v", strings.ToUpper(q.provider), query, args)
	return q.db.QueryRowContext(ctx, query, args...)
}

func (q *Queries) exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	log.Printf("[%s] SQL Exec: %s | Args: %v", strings.ToUpper(q.provider), query, args)
	res, err := q.db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Printf("[%s] SQL Error: %v | Query: %s | Args: %v", strings.ToUpper(q.provider), err, query, args)
	}
	return res, err
}

type Tx struct {
	*Queries
	tx *sql.Tx
}

// BeginTx starts a database transaction.
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	sqlTx, err := db.sqlDB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	q := &Queries{
		db:       sqlTx,
		provider: db.provider,
		dialect:  db.dialect,
		UserRole: db.UserRole,
	}
	q.User = &UserDelegate{client: q}
	q.Profile = &ProfileDelegate{client: q}
	q.Post = &PostDelegate{client: q}
	q.Comment = &CommentDelegate{client: q}
	q.Category = &CategoryDelegate{client: q}
	q.CategoryToPost = &CategoryToPostDelegate{client: q}
	return &Tx{
		Queries: q,
		tx:      sqlTx,
	}, nil
}

// Commit commits the transaction.
func (tx *Tx) Commit() error {
	return tx.tx.Commit()
}

// Rollback aborts the transaction.
func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}

// Raw returns the underlying *sql.Tx transaction handle.
func (tx *Tx) Raw() *sql.Tx {
	return tx.tx
}

// Transaction starts a transaction and runs the provided function.
// If the function returns an error, the transaction is rolled back.
// Otherwise, the transaction is committed.
func (db *DB) Transaction(ctx context.Context, fn func(tx *Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after Rollback
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %w (rollback failed: %v)", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type CreateBuilder[M any, I any, S any, O any] struct {
	client   *Queries
	input    I
	execFunc func(ctx context.Context, input I, s *S, o *O) (*M, error)
}

func (b *CreateBuilder[M, I, S, O]) Select(s S) *CreateSelectBuilder[M, I, S, O] {
	return &CreateSelectBuilder[M, I, S, O]{builder: b, selects: s}
}

func (b *CreateBuilder[M, I, S, O]) Omit(o O) *CreateOmitBuilder[M, I, S, O] {
	return &CreateOmitBuilder[M, I, S, O]{builder: b, omits: o}
}

func (b *CreateBuilder[M, I, S, O]) Exec(ctx context.Context) (*M, error) {
	return b.execFunc(ctx, b.input, nil, nil)
}

type CreateSelectBuilder[M any, I any, S any, O any] struct {
	builder *CreateBuilder[M, I, S, O]
	selects S
}

func (b *CreateSelectBuilder[M, I, S, O]) Exec(ctx context.Context) (*M, error) {
	return b.builder.execFunc(ctx, b.builder.input, &b.selects, nil)
}

type CreateOmitBuilder[M any, I any, S any, O any] struct {
	builder *CreateBuilder[M, I, S, O]
	omits   O
}

func (b *CreateOmitBuilder[M, I, S, O]) Exec(ctx context.Context) (*M, error) {
	return b.builder.execFunc(ctx, b.builder.input, nil, &b.omits)
}
func executeInsert[M any](
	ctx context.Context,
	q *Queries,
	table string,
	cols []string,
	vals []any,
	returningCols []string,
	idCol string,
	scanFunc func(record *M, cols []string) []any,
) (*M, error) {
	var sb strings.Builder
	sb.Grow(128 + len(table) + len(cols)*15 + len(returningCols)*15)

	sb.WriteString("INSERT INTO ")
	sb.WriteString(q.dialect.Quote(table))
	sb.WriteString(" (")
	for i, col := range cols {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(col)
	}
	sb.WriteString(") VALUES (")
	for i := range cols {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(q.dialect.BindVar(i + 1))
	}
	sb.WriteString(")")

	if q.dialect.SupportsReturning() && len(returningCols) > 0 {
		sb.WriteString(" RETURNING ")
		for i, col := range returningCols {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(q.dialect.Quote(col))
		}
	}
	query := sb.String()

	var res M
	if q.dialect.SupportsReturning() {
		row := q.queryRow(ctx, query, vals...)

		scanTargets := scanFunc(&res, returningCols)
		if err := row.Scan(scanTargets...); err != nil {
			return nil, err
		}
		return &res, nil
	}

	// Fallback for dialects without RETURNING (MySQL)
	result, err := q.exec(ctx, query, vals...)
	if err != nil {
		return nil, err
	}

	var idVal any
	quotedIdCol := q.dialect.Quote(idCol)
	for i, c := range cols {
		if c == quotedIdCol {
			idVal = vals[i]
			break
		}
	}
	if idVal == nil {
		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}
		idVal = lastID
	}

	var selectSb strings.Builder
	selectSb.Grow(64 + len(returningCols)*15 + len(table) + len(idCol))
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		selectSb.WriteString(q.dialect.Quote(col))
	}
	selectSb.WriteString(" FROM ")
	selectSb.WriteString(q.dialect.Quote(table))
	selectSb.WriteString(" WHERE ")
	selectSb.WriteString(q.dialect.Quote(idCol))
	selectSb.WriteString(" = ?")

	row := q.queryRow(ctx, selectSb.String(), idVal)
	scanTargets := scanFunc(&res, returningCols)
	if err := row.Scan(scanTargets...); err != nil {
		return nil, err
	}
	return &res, nil
}

func loadRelation[P any, C any](
	ctx context.Context,
	q *Queries,
	parents []*P,
	parentKey func(*P) (string, bool),
	table string,
	fkCol string,
	returningCols []string,
	scan func(*sql.Rows, *C) error,
	childKey func(*C) (string, bool),
	assign func(*P, []*C),
) ([]*C, error) {
	var parentKeys []any
	for _, p := range parents {
		if p == nil {
			continue
		}
		if key, ok := parentKey(p); ok {
			parentKeys = append(parentKeys, key)
		}
	}
	if len(parentKeys) == 0 {
		return nil, nil
	}

	var sb strings.Builder
	sb.Grow(128 + len(returningCols)*15 + len(table) + len(fkCol) + len(parentKeys)*3)
	sb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(q.dialect.Quote(col))
	}
	sb.WriteString(" FROM ")
	sb.WriteString(q.dialect.Quote(table))
	sb.WriteString(" WHERE ")
	sb.WriteString(q.dialect.Quote(fkCol))
	sb.WriteString(" IN (")
	sb.WriteString(q.bindVars(len(parentKeys)))
	sb.WriteString(")")
	query := sb.String()

	rows, err := q.query(ctx, query, parentKeys...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	childMap := make(map[string][]*C)
	var allChildren []*C

	for rows.Next() {
		var child C
		if err := scan(rows, &child); err != nil {
			return nil, err
		}
		if key, ok := childKey(&child); ok {
			childMap[key] = append(childMap[key], &child)
		}
		allChildren = append(allChildren, &child)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, p := range parents {
		if p == nil {
			continue
		}
		if key, ok := parentKey(p); ok {
			assign(p, childMap[key])
		}
	}

	return allChildren, nil
}
