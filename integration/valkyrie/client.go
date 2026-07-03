package valkyrie

import (
	"context"
	"crypto/rand"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
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

//go:embed migrations/*.sql
var migrationsFS embed.FS

func generateCUID() string {
	now := time.Now().UnixMilli()
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("c%x%x", now, b)
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

type postgresDialect struct{}

func (postgresDialect) Quote(ident string) string { return `"` + ident + `"` }
func (postgresDialect) BindVar(idx int) string    { return fmt.Sprintf("$%d", idx) }
func (postgresDialect) SupportsReturning() bool   { return true }

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

	var d Dialect
	switch provider {
	case "postgres", "postgresql":
		d = postgresDialect{}
	case "sqlite", "sqlite3":
		d = sqliteDialect{}

	default:
		sqlDB.Close()
		return nil, fmt.Errorf("unsupported database provider: %s", provider)
	}

	q := &Queries{
		db:       sqlDB,
		provider: provider,
		dialect:  d,
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
	placeholders := make([]string, len(cols))
	for i := range cols {
		placeholders[i] = q.dialect.BindVar(i + 1)
	}

	var res M
	quotedReturningCols := make([]string, len(returningCols))
	for i, col := range returningCols {
		quotedReturningCols[i] = q.dialect.Quote(col)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		q.dialect.Quote(table),
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
	)

	if q.dialect.SupportsReturning() {
		query += " RETURNING " + strings.Join(quotedReturningCols, ", ")
		row := q.db.QueryRowContext(ctx, query, vals...)

		scanTargets := scanFunc(&res, returningCols)
		if err := row.Scan(scanTargets...); err != nil {
			return nil, err
		}
		return &res, nil
	}

	// Fallback for dialects without RETURNING (MySQL)
	result, err := q.db.ExecContext(ctx, query, vals...)
	if err != nil {
		return nil, err
	}

	var idVal any
	for i, c := range cols {
		if c == q.dialect.Quote(idCol) {
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

	selectQuery := fmt.Sprintf("SELECT %s FROM %s WHERE %s = ?",
		strings.Join(quotedReturningCols, ", "),
		q.dialect.Quote(table),
		q.dialect.Quote(idCol),
	)
	row := q.db.QueryRowContext(ctx, selectQuery, idVal)
	scanTargets := scanFunc(&res, returningCols)
	if err := row.Scan(scanTargets...); err != nil {
		return nil, err
	}
	return &res, nil
}

// User represents the database model
type User struct {
	Id           string       `db:"id" json:"id"`
	Email        string       `db:"email" json:"email"`
	PhoneNum     string       `db:"phoneNum" json:"phoneNum"`
	Role         UserRoleType `db:"role" json:"role"`
	ReferredById *string      `db:"referredById" json:"referredById"`
	Profile      *Profile     `json:"profile,omitempty"`
	Posts        []Post       `json:"posts,omitempty"`
	Comments     []Comment    `json:"comments,omitempty"`
	ReferredBy   *User        `json:"referredBy,omitempty"`
	Referrals    []User       `json:"referrals,omitempty"`
}

// UserCreateInput represents the input structure for creation
type UserCreateInput struct {
	Id           *string       `json:"id"`
	Email        string        `json:"email"`
	PhoneNum     string        `json:"phoneNum"`
	Role         *UserRoleType `json:"role"`
	ReferredById *string       `json:"referredById"`
}

// UserSelect specifies which fields to include
type UserSelect struct {
	Id           bool           `json:"id"`
	Email        bool           `json:"email"`
	PhoneNum     bool           `json:"phoneNum"`
	Role         bool           `json:"role"`
	ReferredById bool           `json:"referredById"`
	Profile      *ProfileSelect `json:"profile,omitempty"`
	Posts        *PostSelect    `json:"posts,omitempty"`
	Comments     *CommentSelect `json:"comments,omitempty"`
	ReferredBy   *UserSelect    `json:"referredBy,omitempty"`
	Referrals    *UserSelect    `json:"referrals,omitempty"`
}

// UserOmit specifies which fields to exclude
type UserOmit struct {
	Id           bool         `json:"id"`
	Email        bool         `json:"email"`
	PhoneNum     bool         `json:"phoneNum"`
	Role         bool         `json:"role"`
	ReferredById bool         `json:"referredById"`
	Profile      *ProfileOmit `json:"profile,omitempty"`
	Posts        *PostOmit    `json:"posts,omitempty"`
	Comments     *CommentOmit `json:"comments,omitempty"`
	ReferredBy   *UserOmit    `json:"referredBy,omitempty"`
	Referrals    *UserOmit    `json:"referrals,omitempty"`
}

type UserDelegate struct {
	client *Queries
}

func (d *UserDelegate) Create(input UserCreateInput) *CreateBuilder[User, UserCreateInput, UserSelect, UserOmit] {
	return &CreateBuilder[User, UserCreateInput, UserSelect, UserOmit]{
		client:   d.client,
		input:    input,
		execFunc: d.client.executeUserCreate,
	}
}

func (q *Queries) executeUserCreate(ctx context.Context, input UserCreateInput, selects *UserSelect, omits *UserOmit) (*User, error) {
	var cols []string
	var vals []any
	if input.Id != nil {
		cols = append(cols, q.dialect.Quote("id"))
		vals = append(vals, *input.Id)
	} else {
		cols = append(cols, q.dialect.Quote("id"))
		vals = append(vals, generateCUID())
	}
	cols = append(cols, q.dialect.Quote("email"))
	vals = append(vals, input.Email)
	cols = append(cols, q.dialect.Quote("phoneNum"))
	vals = append(vals, input.PhoneNum)
	if input.Role != nil {
		cols = append(cols, q.dialect.Quote("role"))
		vals = append(vals, *input.Role)
	}
	if input.ReferredById != nil {
		cols = append(cols, q.dialect.Quote("referredById"))
		vals = append(vals, *input.ReferredById)
	}

	var returningCols []string
	{
		include := true
		if selects != nil {
			include = false
			if selects.Id {
				include = true
			}
		} else if omits != nil {
			if omits.Id {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "id")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if selects.Email {
				include = true
			}
		} else if omits != nil {
			if omits.Email {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "email")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if selects.PhoneNum {
				include = true
			}
		} else if omits != nil {
			if omits.PhoneNum {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "phoneNum")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if selects.Role {
				include = true
			}
		} else if omits != nil {
			if omits.Role {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "role")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if selects.ReferredById {
				include = true
			}
		} else if omits != nil {
			if omits.ReferredById {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "referredById")
		}
	}

	if len(returningCols) == 0 {
		returningCols = append(returningCols, "id")
		returningCols = append(returningCols, "email")
		returningCols = append(returningCols, "phoneNum")
		returningCols = append(returningCols, "role")
		returningCols = append(returningCols, "referredById")
	}

	scanFunc := func(res *User, cols []string) []any {
		targets := make([]any, len(cols))
		for i, col := range cols {
			switch col {
			case "id":
				targets[i] = &res.Id
			case "email":
				targets[i] = &res.Email
			case "phoneNum":
				targets[i] = &res.PhoneNum
			case "role":
				targets[i] = &res.Role
			case "referredById":
				targets[i] = &res.ReferredById
			}
		}
		return targets
	}

	idCol := "id"

	return executeInsert(ctx, q, "User", cols, vals, returningCols, idCol, scanFunc)
}

// Profile represents the database model
type Profile struct {
	Id     string  `db:"id" json:"id"`
	Bio    *string `db:"bio" json:"bio"`
	UserId string  `db:"userId" json:"userId"`
	User   *User   `json:"user,omitempty"`
}

// ProfileCreateInput represents the input structure for creation
type ProfileCreateInput struct {
	Id     *string `json:"id"`
	Bio    *string `json:"bio"`
	UserId string  `json:"userId"`
}

// ProfileSelect specifies which fields to include
type ProfileSelect struct {
	Id     bool        `json:"id"`
	Bio    bool        `json:"bio"`
	UserId bool        `json:"userId"`
	User   *UserSelect `json:"user,omitempty"`
}

// ProfileOmit specifies which fields to exclude
type ProfileOmit struct {
	Id     bool      `json:"id"`
	Bio    bool      `json:"bio"`
	UserId bool      `json:"userId"`
	User   *UserOmit `json:"user,omitempty"`
}

type ProfileDelegate struct {
	client *Queries
}

func (d *ProfileDelegate) Create(input ProfileCreateInput) *CreateBuilder[Profile, ProfileCreateInput, ProfileSelect, ProfileOmit] {
	return &CreateBuilder[Profile, ProfileCreateInput, ProfileSelect, ProfileOmit]{
		client:   d.client,
		input:    input,
		execFunc: d.client.executeProfileCreate,
	}
}

func (q *Queries) executeProfileCreate(ctx context.Context, input ProfileCreateInput, selects *ProfileSelect, omits *ProfileOmit) (*Profile, error) {
	var cols []string
	var vals []any
	if input.Id != nil {
		cols = append(cols, q.dialect.Quote("id"))
		vals = append(vals, *input.Id)
	} else {
		cols = append(cols, q.dialect.Quote("id"))
		vals = append(vals, generateCUID())
	}
	if input.Bio != nil {
		cols = append(cols, q.dialect.Quote("bio"))
		vals = append(vals, *input.Bio)
	}
	cols = append(cols, q.dialect.Quote("userId"))
	vals = append(vals, input.UserId)

	var returningCols []string
	{
		include := true
		if selects != nil {
			include = false
			if selects.Id {
				include = true
			}
		} else if omits != nil {
			if omits.Id {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "id")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if selects.Bio {
				include = true
			}
		} else if omits != nil {
			if omits.Bio {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "bio")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if selects.UserId {
				include = true
			}
		} else if omits != nil {
			if omits.UserId {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "userId")
		}
	}

	if len(returningCols) == 0 {
		returningCols = append(returningCols, "id")
		returningCols = append(returningCols, "bio")
		returningCols = append(returningCols, "userId")
	}

	scanFunc := func(res *Profile, cols []string) []any {
		targets := make([]any, len(cols))
		for i, col := range cols {
			switch col {
			case "id":
				targets[i] = &res.Id
			case "bio":
				targets[i] = &res.Bio
			case "userId":
				targets[i] = &res.UserId
			}
		}
		return targets
	}

	idCol := "id"

	return executeInsert(ctx, q, "Profile", cols, vals, returningCols, idCol, scanFunc)
}

// Post represents the database model
type Post struct {
	Id         string           `db:"id" json:"id"`
	Title      string           `db:"title" json:"title"`
	Content    *string          `db:"content" json:"content"`
	Published  bool             `db:"published" json:"published"`
	AuthorId   string           `db:"authorId" json:"authorId"`
	Author     *User            `json:"author,omitempty"`
	Comments   []Comment        `json:"comments,omitempty"`
	Categories []CategoryToPost `json:"categories,omitempty"`
}

// PostCreateInput represents the input structure for creation
type PostCreateInput struct {
	Id        *string `json:"id"`
	Title     string  `json:"title"`
	Content   *string `json:"content"`
	Published *bool   `json:"published"`
	AuthorId  string  `json:"authorId"`
}

// PostSelect specifies which fields to include
type PostSelect struct {
	Id         bool                  `json:"id"`
	Title      bool                  `json:"title"`
	Content    bool                  `json:"content"`
	Published  bool                  `json:"published"`
	AuthorId   bool                  `json:"authorId"`
	Author     *UserSelect           `json:"author,omitempty"`
	Comments   *CommentSelect        `json:"comments,omitempty"`
	Categories *CategoryToPostSelect `json:"categories,omitempty"`
}

// PostOmit specifies which fields to exclude
type PostOmit struct {
	Id         bool                `json:"id"`
	Title      bool                `json:"title"`
	Content    bool                `json:"content"`
	Published  bool                `json:"published"`
	AuthorId   bool                `json:"authorId"`
	Author     *UserOmit           `json:"author,omitempty"`
	Comments   *CommentOmit        `json:"comments,omitempty"`
	Categories *CategoryToPostOmit `json:"categories,omitempty"`
}

type PostDelegate struct {
	client *Queries
}

func (d *PostDelegate) Create(input PostCreateInput) *CreateBuilder[Post, PostCreateInput, PostSelect, PostOmit] {
	return &CreateBuilder[Post, PostCreateInput, PostSelect, PostOmit]{
		client:   d.client,
		input:    input,
		execFunc: d.client.executePostCreate,
	}
}

func (q *Queries) executePostCreate(ctx context.Context, input PostCreateInput, selects *PostSelect, omits *PostOmit) (*Post, error) {
	var cols []string
	var vals []any
	if input.Id != nil {
		cols = append(cols, q.dialect.Quote("id"))
		vals = append(vals, *input.Id)
	} else {
		cols = append(cols, q.dialect.Quote("id"))
		vals = append(vals, generateCUID())
	}
	cols = append(cols, q.dialect.Quote("title"))
	vals = append(vals, input.Title)
	if input.Content != nil {
		cols = append(cols, q.dialect.Quote("content"))
		vals = append(vals, *input.Content)
	}
	if input.Published != nil {
		cols = append(cols, q.dialect.Quote("published"))
		vals = append(vals, *input.Published)
	}
	cols = append(cols, q.dialect.Quote("authorId"))
	vals = append(vals, input.AuthorId)

	var returningCols []string
	{
		include := true
		if selects != nil {
			include = false
			if selects.Id {
				include = true
			}
		} else if omits != nil {
			if omits.Id {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "id")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if selects.Title {
				include = true
			}
		} else if omits != nil {
			if omits.Title {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "title")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if selects.Content {
				include = true
			}
		} else if omits != nil {
			if omits.Content {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "content")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if selects.Published {
				include = true
			}
		} else if omits != nil {
			if omits.Published {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "published")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if selects.AuthorId {
				include = true
			}
		} else if omits != nil {
			if omits.AuthorId {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "authorId")
		}
	}

	if len(returningCols) == 0 {
		returningCols = append(returningCols, "id")
		returningCols = append(returningCols, "title")
		returningCols = append(returningCols, "content")
		returningCols = append(returningCols, "published")
		returningCols = append(returningCols, "authorId")
	}

	scanFunc := func(res *Post, cols []string) []any {
		targets := make([]any, len(cols))
		for i, col := range cols {
			switch col {
			case "id":
				targets[i] = &res.Id
			case "title":
				targets[i] = &res.Title
			case "content":
				targets[i] = &res.Content
			case "published":
				targets[i] = &res.Published
			case "authorId":
				targets[i] = &res.AuthorId
			}
		}
		return targets
	}

	idCol := "id"

	return executeInsert(ctx, q, "Post", cols, vals, returningCols, idCol, scanFunc)
}

// Comment represents the database model
type Comment struct {
	Id       string `db:"id" json:"id"`
	Textify  int32  `db:"textify" json:"textify"`
	Dummy3   string `db:"dummy3" json:"dummy3"`
	Dummy1   int32  `db:"dummy1" json:"dummy1"`
	Dummy2   string `db:"dummy2" json:"dummy2"`
	PostId   string `db:"postId" json:"postId"`
	AuthorId string `db:"authorId" json:"authorId"`
	Post     *Post  `json:"post,omitempty"`
	Author   *User  `json:"author,omitempty"`
}

// CommentCreateInput represents the input structure for creation
type CommentCreateInput struct {
	Id       *string `json:"id"`
	Textify  int32   `json:"textify"`
	Dummy3   string  `json:"dummy3"`
	Dummy1   int32   `json:"dummy1"`
	Dummy2   string  `json:"dummy2"`
	PostId   string  `json:"postId"`
	AuthorId string  `json:"authorId"`
}

// CommentSelect specifies which fields to include
type CommentSelect struct {
	Id       bool        `json:"id"`
	Textify  bool        `json:"textify"`
	Dummy3   bool        `json:"dummy3"`
	Dummy1   bool        `json:"dummy1"`
	Dummy2   bool        `json:"dummy2"`
	PostId   bool        `json:"postId"`
	AuthorId bool        `json:"authorId"`
	Post     *PostSelect `json:"post,omitempty"`
	Author   *UserSelect `json:"author,omitempty"`
}

// CommentOmit specifies which fields to exclude
type CommentOmit struct {
	Id       bool      `json:"id"`
	Textify  bool      `json:"textify"`
	Dummy3   bool      `json:"dummy3"`
	Dummy1   bool      `json:"dummy1"`
	Dummy2   bool      `json:"dummy2"`
	PostId   bool      `json:"postId"`
	AuthorId bool      `json:"authorId"`
	Post     *PostOmit `json:"post,omitempty"`
	Author   *UserOmit `json:"author,omitempty"`
}

type CommentDelegate struct {
	client *Queries
}

func (d *CommentDelegate) Create(input CommentCreateInput) *CreateBuilder[Comment, CommentCreateInput, CommentSelect, CommentOmit] {
	return &CreateBuilder[Comment, CommentCreateInput, CommentSelect, CommentOmit]{
		client:   d.client,
		input:    input,
		execFunc: d.client.executeCommentCreate,
	}
}

func (q *Queries) executeCommentCreate(ctx context.Context, input CommentCreateInput, selects *CommentSelect, omits *CommentOmit) (*Comment, error) {
	var cols []string
	var vals []any
	if input.Id != nil {
		cols = append(cols, q.dialect.Quote("id"))
		vals = append(vals, *input.Id)
	} else {
		cols = append(cols, q.dialect.Quote("id"))
		vals = append(vals, generateCUID())
	}
	cols = append(cols, q.dialect.Quote("textify"))
	vals = append(vals, input.Textify)
	cols = append(cols, q.dialect.Quote("dummy3"))
	vals = append(vals, input.Dummy3)
	cols = append(cols, q.dialect.Quote("dummy1"))
	vals = append(vals, input.Dummy1)
	cols = append(cols, q.dialect.Quote("dummy2"))
	vals = append(vals, input.Dummy2)
	cols = append(cols, q.dialect.Quote("postId"))
	vals = append(vals, input.PostId)
	cols = append(cols, q.dialect.Quote("authorId"))
	vals = append(vals, input.AuthorId)

	var returningCols []string
	{
		include := true
		if selects != nil {
			include = false
			if selects.Id {
				include = true
			}
		} else if omits != nil {
			if omits.Id {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "id")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if selects.Textify {
				include = true
			}
		} else if omits != nil {
			if omits.Textify {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "textify")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if selects.Dummy3 {
				include = true
			}
		} else if omits != nil {
			if omits.Dummy3 {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "dummy3")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if selects.Dummy1 {
				include = true
			}
		} else if omits != nil {
			if omits.Dummy1 {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "dummy1")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if selects.Dummy2 {
				include = true
			}
		} else if omits != nil {
			if omits.Dummy2 {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "dummy2")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if selects.PostId {
				include = true
			}
		} else if omits != nil {
			if omits.PostId {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "postId")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if selects.AuthorId {
				include = true
			}
		} else if omits != nil {
			if omits.AuthorId {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "authorId")
		}
	}

	if len(returningCols) == 0 {
		returningCols = append(returningCols, "id")
		returningCols = append(returningCols, "textify")
		returningCols = append(returningCols, "dummy3")
		returningCols = append(returningCols, "dummy1")
		returningCols = append(returningCols, "dummy2")
		returningCols = append(returningCols, "postId")
		returningCols = append(returningCols, "authorId")
	}

	scanFunc := func(res *Comment, cols []string) []any {
		targets := make([]any, len(cols))
		for i, col := range cols {
			switch col {
			case "id":
				targets[i] = &res.Id
			case "textify":
				targets[i] = &res.Textify
			case "dummy3":
				targets[i] = &res.Dummy3
			case "dummy1":
				targets[i] = &res.Dummy1
			case "dummy2":
				targets[i] = &res.Dummy2
			case "postId":
				targets[i] = &res.PostId
			case "authorId":
				targets[i] = &res.AuthorId
			}
		}
		return targets
	}

	idCol := "id"

	return executeInsert(ctx, q, "Comment", cols, vals, returningCols, idCol, scanFunc)
}

// Category represents the database model
type Category struct {
	Id    int32            `db:"id" json:"id"`
	Name  string           `db:"name" json:"name"`
	Posts []CategoryToPost `json:"posts,omitempty"`
}

// CategoryCreateInput represents the input structure for creation
type CategoryCreateInput struct {
	Id   *int32 `json:"id"`
	Name string `json:"name"`
}

// CategorySelect specifies which fields to include
type CategorySelect struct {
	Id    bool                  `json:"id"`
	Name  bool                  `json:"name"`
	Posts *CategoryToPostSelect `json:"posts,omitempty"`
}

// CategoryOmit specifies which fields to exclude
type CategoryOmit struct {
	Id    bool                `json:"id"`
	Name  bool                `json:"name"`
	Posts *CategoryToPostOmit `json:"posts,omitempty"`
}

type CategoryDelegate struct {
	client *Queries
}

func (d *CategoryDelegate) Create(input CategoryCreateInput) *CreateBuilder[Category, CategoryCreateInput, CategorySelect, CategoryOmit] {
	return &CreateBuilder[Category, CategoryCreateInput, CategorySelect, CategoryOmit]{
		client:   d.client,
		input:    input,
		execFunc: d.client.executeCategoryCreate,
	}
}

func (q *Queries) executeCategoryCreate(ctx context.Context, input CategoryCreateInput, selects *CategorySelect, omits *CategoryOmit) (*Category, error) {
	var cols []string
	var vals []any
	if input.Id != nil {
		cols = append(cols, q.dialect.Quote("id"))
		vals = append(vals, *input.Id)
	} else {
		cols = append(cols, q.dialect.Quote("id"))
	}
	cols = append(cols, q.dialect.Quote("name"))
	vals = append(vals, input.Name)

	var returningCols []string
	{
		include := true
		if selects != nil {
			include = false
			if selects.Id {
				include = true
			}
		} else if omits != nil {
			if omits.Id {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "id")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if selects.Name {
				include = true
			}
		} else if omits != nil {
			if omits.Name {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "name")
		}
	}

	if len(returningCols) == 0 {
		returningCols = append(returningCols, "id")
		returningCols = append(returningCols, "name")
	}

	scanFunc := func(res *Category, cols []string) []any {
		targets := make([]any, len(cols))
		for i, col := range cols {
			switch col {
			case "id":
				targets[i] = &res.Id
			case "name":
				targets[i] = &res.Name
			}
		}
		return targets
	}

	idCol := "id"

	return executeInsert(ctx, q, "Category", cols, vals, returningCols, idCol, scanFunc)
}

// CategoryToPost represents the database model
type CategoryToPost struct {
	PostId     string    `db:"postId" json:"postId"`
	CategoryId int32     `db:"categoryId" json:"categoryId"`
	Post       *Post     `json:"post,omitempty"`
	Category   *Category `json:"category,omitempty"`
}

// CategoryToPostCreateInput represents the input structure for creation
type CategoryToPostCreateInput struct {
	PostId     string `json:"postId"`
	CategoryId int32  `json:"categoryId"`
}

// CategoryToPostSelect specifies which fields to include
type CategoryToPostSelect struct {
	PostId     bool            `json:"postId"`
	CategoryId bool            `json:"categoryId"`
	Post       *PostSelect     `json:"post,omitempty"`
	Category   *CategorySelect `json:"category,omitempty"`
}

// CategoryToPostOmit specifies which fields to exclude
type CategoryToPostOmit struct {
	PostId     bool          `json:"postId"`
	CategoryId bool          `json:"categoryId"`
	Post       *PostOmit     `json:"post,omitempty"`
	Category   *CategoryOmit `json:"category,omitempty"`
}

type CategoryToPostDelegate struct {
	client *Queries
}

func (d *CategoryToPostDelegate) Create(input CategoryToPostCreateInput) *CreateBuilder[CategoryToPost, CategoryToPostCreateInput, CategoryToPostSelect, CategoryToPostOmit] {
	return &CreateBuilder[CategoryToPost, CategoryToPostCreateInput, CategoryToPostSelect, CategoryToPostOmit]{
		client:   d.client,
		input:    input,
		execFunc: d.client.executeCategoryToPostCreate,
	}
}

func (q *Queries) executeCategoryToPostCreate(ctx context.Context, input CategoryToPostCreateInput, selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error) {
	var cols []string
	var vals []any
	cols = append(cols, q.dialect.Quote("postId"))
	vals = append(vals, input.PostId)
	cols = append(cols, q.dialect.Quote("categoryId"))
	vals = append(vals, input.CategoryId)

	var returningCols []string
	{
		include := true
		if selects != nil {
			include = false
			if selects.PostId {
				include = true
			}
		} else if omits != nil {
			if omits.PostId {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "postId")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if selects.CategoryId {
				include = true
			}
		} else if omits != nil {
			if omits.CategoryId {
				include = false
			}
		}
		if include {
			returningCols = append(returningCols, "categoryId")
		}
	}

	if len(returningCols) == 0 {
		returningCols = append(returningCols, "postId")
		returningCols = append(returningCols, "categoryId")
	}

	scanFunc := func(res *CategoryToPost, cols []string) []any {
		targets := make([]any, len(cols))
		for i, col := range cols {
			switch col {
			case "postId":
				targets[i] = &res.PostId
			case "categoryId":
				targets[i] = &res.CategoryId
			}
		}
		return targets
	}

	idCol := ""

	return executeInsert(ctx, q, "CategoryToPost", cols, vals, returningCols, idCol, scanFunc)
}
