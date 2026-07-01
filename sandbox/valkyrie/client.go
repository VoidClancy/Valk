package valkyrie

import (
	"context"
	"database/sql"
	"embed"
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

type DB struct {
	DB             *sql.DB
	provider       string
	User           *UserDelegate
	Profile        *ProfileDelegate
	Post           *PostDelegate
	Comment        *CommentDelegate
	Category       *CategoryDelegate
	CategoryToPost *CategoryToPostDelegate
	UserRole       userRoleNamespace
}

// Open opens a database connection and returns a client.
func Open(provider, dataSourceName string) (*DB, error) {
	sqlDB, err := sql.Open(provider, dataSourceName)
	if err != nil {
		return nil, err
	}
	db := &DB{
		DB:       sqlDB,
		provider: provider,
		UserRole: UserRole,
	}
	db.User = &UserDelegate{client: db}
	db.Profile = &ProfileDelegate{client: db}
	db.Post = &PostDelegate{client: db}
	db.Comment = &CommentDelegate{client: db}
	db.Category = &CategoryDelegate{client: db}
	db.CategoryToPost = &CategoryToPostDelegate{client: db}
	return db, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.DB.Close()
}

// RunMigrations runs all pending migrations from the embedded folder.
func (db *DB) RunMigrations(ctx context.Context) error {
	if err := goose.SetDialect(db.provider); err != nil {
		return err
	}
	goose.SetLogger(goose.NopLogger())
	goose.SetBaseFS(migrationsFS)
	return goose.UpContext(ctx, db.DB, "migrations")
}

type UserDelegate struct {
	client *DB
}
type ProfileDelegate struct {
	client *DB
}
type PostDelegate struct {
	client *DB
}
type CommentDelegate struct {
	client *DB
}
type CategoryDelegate struct {
	client *DB
}
type CategoryToPostDelegate struct {
	client *DB
}
