package valkyrie

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

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

type DBTX interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

type Queries struct {
	db             DBTX
	provider       string
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
	if err := goose.SetDialect(db.provider); err != nil {
		return err
	}
	goose.SetLogger(goose.NopLogger())
	goose.SetBaseFS(migrationsFS)
	return goose.UpContext(ctx, db.sqlDB, "migrations")
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

type UserDelegate struct {
	client *Queries
}
type ProfileDelegate struct {
	client *Queries
}
type PostDelegate struct {
	client *Queries
}
type CommentDelegate struct {
	client *Queries
}
type CategoryDelegate struct {
	client *Queries
}
type CategoryToPostDelegate struct {
	client *Queries
}
