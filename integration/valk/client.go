package valk

import (
	"context"
	"crypto/rand"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
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

// FieldError represents a single validation failure on a specific field.
type FieldError struct {
	Field string
	Value any
	Rule  string
	Msg   string
}

func (e FieldError) Error() string {
	return fmt.Sprintf("field %s: %s (value: %v, rule: %s)", e.Field, e.Msg, e.Value, e.Rule)
}

// ValidationError collects multiple validation errors during an operation.
type ValidationError struct {
	Errors []FieldError
}

func (e ValidationError) Error() string {
	var msgs []string
	for _, err := range e.Errors {
		msgs = append(msgs, err.Error())
	}
	return fmt.Sprintf("validation failed: %s", strings.Join(msgs, "; "))
}

func (e *ValidationError) Add(field string, value any, rule string, msg string) {
	e.Errors = append(e.Errors, FieldError{
		Field: field,
		Value: value,
		Rule:  rule,
		Msg:   msg,
	})
}

func (e *ValidationError) HasErrors() bool {
	return len(e.Errors) > 0
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

func (e UserRoleType) IsValid() bool {
	switch e {
	case UserRoleTypeAdmin, UserRoleTypeStudent, UserRoleTypeTeacher:
		return true
	}
	return false
}

type Dialect interface {
	Quote(ident string) string
	BindVar(idx int) string
	SupportsReturning() bool
	SupportsBulkInsert() bool
}
type sqliteDialect struct{}

func (sqliteDialect) Quote(ident string) string { return `"` + ident + `"` }
func (sqliteDialect) BindVar(idx int) string    { return "?" }
func (sqliteDialect) SupportsReturning() bool   { return true }
func (sqliteDialect) SupportsBulkInsert() bool  { return false }

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
	q.initDelegates()
	return &DB{
		Queries: q,
		sqlDB:   sqlDB,
	}, nil
}

func (q *Queries) initDelegates() {
	q.User = &UserDelegate{client: q}
	q.Profile = &ProfileDelegate{client: q}
	q.Post = &PostDelegate{client: q}
	q.Comment = &CommentDelegate{client: q}
	q.Category = &CategoryDelegate{client: q}
	q.CategoryToPost = &CategoryToPostDelegate{client: q}
}

func (q *Queries) copyHooksFrom(other *Queries) {
	q.User.beforeCreate = other.User.beforeCreate
	q.User.afterCreate = other.User.afterCreate
	q.Profile.beforeCreate = other.Profile.beforeCreate
	q.Profile.afterCreate = other.Profile.afterCreate
	q.Post.beforeCreate = other.Post.beforeCreate
	q.Post.afterCreate = other.Post.afterCreate
	q.Comment.beforeCreate = other.Comment.beforeCreate
	q.Comment.afterCreate = other.Comment.afterCreate
	q.Category.beforeCreate = other.Category.beforeCreate
	q.Category.afterCreate = other.Category.afterCreate
	q.CategoryToPost.beforeCreate = other.CategoryToPost.beforeCreate
	q.CategoryToPost.afterCreate = other.CategoryToPost.afterCreate
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
	err := goose.UpContext(ctx, db.sqlDB, "migrations")
	if err != nil {
		return err
	}
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
	res, err := q.db.QueryContext(ctx, query, args...)
	return res, err
}

func (q *Queries) queryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return q.db.QueryRowContext(ctx, query, args...)
}

func (q *Queries) exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	res, err := q.db.ExecContext(ctx, query, args...)
	return res, err
}

func (q *Queries) transaction(ctx context.Context, fn func(txQ *Queries) error) error {
	if _, ok := q.db.(*sql.Tx); ok {
		return fn(q)
	}

	starter, ok := q.db.(interface {
		BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	})
	if !ok {
		return fn(q)
	}
	tx, err := starter.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	txQueries := &Queries{
		db:       tx,
		provider: q.provider,
		dialect:  q.dialect,
		UserRole: q.UserRole,
	}
	txQueries.initDelegates()
	txQueries.copyHooksFrom(q)

	if err := fn(txQueries); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

type PredicateData struct {
	Column    string
	Operator  string
	Value     any
	IsLogical bool
	Children  []PredicateData
}

type Predicate interface {
	ToPredicateData() PredicateData
}

type UniquePredicate interface {
	Predicate
	IsUnique()
	Validate() error
}

type StandardPredicate struct {
	Data PredicateData
}

func (sp StandardPredicate) ToPredicateData() PredicateData {
	return sp.Data
}

func And(preds ...Predicate) Predicate {
	var children []PredicateData
	for _, p := range preds {
		if p != nil {
			children = append(children, p.ToPredicateData())
		}
	}
	return StandardPredicate{
		Data: PredicateData{
			IsLogical: true,
			Operator:  "AND",
			Children:  children,
		},
	}
}

func Or(preds ...Predicate) Predicate {
	var children []PredicateData
	for _, p := range preds {
		if p != nil {
			children = append(children, p.ToPredicateData())
		}
	}
	return StandardPredicate{
		Data: PredicateData{
			IsLogical: true,
			Operator:  "OR",
			Children:  children,
		},
	}
}

func Not(pred Predicate) Predicate {
	var children []PredicateData
	if pred != nil {
		children = append(children, pred.ToPredicateData())
	}
	return StandardPredicate{
		Data: PredicateData{
			IsLogical: true,
			Operator:  "NOT",
			Children:  children,
		},
	}
}

type Field[T any] struct {
	Column string
}

func (f Field[T]) EQ(val T) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "=",
			Value:    val,
		},
	}
}

func (f Field[T]) NEQ(val T) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "!=",
			Value:    val,
		},
	}
}

func (f Field[T]) GT(val T) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">",
			Value:    val,
		},
	}
}

func (f Field[T]) GTE(val T) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">=",
			Value:    val,
		},
	}
}

func (f Field[T]) LT(val T) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<",
			Value:    val,
		},
	}
}

func (f Field[T]) LTE(val T) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<=",
			Value:    val,
		},
	}
}

func (f Field[T]) In(vals []T) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IN",
			Value:    vals,
		},
	}
}

func (f Field[T]) IsNull() Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NULL",
		},
	}
}

func (f Field[T]) IsNotNull() Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NOT NULL",
		},
	}
}

type UniqueField[T any] struct {
	Column string
}

type UniqueFieldPredicate struct {
	StandardPredicate
}

func (UniqueFieldPredicate) IsUnique() {}

func (p UniqueFieldPredicate) Validate() error {
	if p.Data.Column == "" {
		return fmt.Errorf("at least one unique field must be set for FindUnique")
	}
	return nil
}

func (f UniqueField[T]) EQ(val T) UniquePredicate {
	return UniqueFieldPredicate{
		StandardPredicate: StandardPredicate{
			Data: PredicateData{
				Column:   f.Column,
				Operator: "=",
				Value:    val,
			},
		},
	}
}

func (f UniqueField[T]) NEQ(val T) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "!=",
			Value:    val,
		},
	}
}

func (f UniqueField[T]) GT(val T) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">",
			Value:    val,
		},
	}
}

func (f UniqueField[T]) GTE(val T) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">=",
			Value:    val,
		},
	}
}

func (f UniqueField[T]) LT(val T) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<",
			Value:    val,
		},
	}
}

func (f UniqueField[T]) LTE(val T) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<=",
			Value:    val,
		},
	}
}

func (f UniqueField[T]) In(vals []T) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IN",
			Value:    vals,
		},
	}
}

func (f UniqueField[T]) IsNull() Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NULL",
		},
	}
}

func (f UniqueField[T]) IsNotNull() Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NOT NULL",
		},
	}
}

type StringField struct {
	Column string
}

func (f StringField) EQ(val string) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "=",
			Value:    val,
		},
	}
}

func (f StringField) NEQ(val string) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "!=",
			Value:    val,
		},
	}
}

func (f StringField) GT(val string) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">",
			Value:    val,
		},
	}
}

func (f StringField) GTE(val string) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">=",
			Value:    val,
		},
	}
}

func (f StringField) LT(val string) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<",
			Value:    val,
		},
	}
}

func (f StringField) LTE(val string) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<=",
			Value:    val,
		},
	}
}

func (f StringField) In(vals []string) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IN",
			Value:    vals,
		},
	}
}

func (f StringField) Like(val string) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "LIKE",
			Value:    val,
		},
	}
}

func (f StringField) Contains(val string) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "LIKE",
			Value:    "%" + val + "%",
		},
	}
}

func (f StringField) IsNull() Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NULL",
		},
	}
}

func (f StringField) IsNotNull() Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NOT NULL",
		},
	}
}

type StringUniqueField struct {
	Column string
}

func (f StringUniqueField) EQ(val string) UniquePredicate {
	return UniqueFieldPredicate{
		StandardPredicate: StandardPredicate{
			Data: PredicateData{
				Column:   f.Column,
				Operator: "=",
				Value:    val,
			},
		},
	}
}

func (f StringUniqueField) NEQ(val string) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "!=",
			Value:    val,
		},
	}
}

func (f StringUniqueField) GT(val string) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">",
			Value:    val,
		},
	}
}

func (f StringUniqueField) GTE(val string) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">=",
			Value:    val,
		},
	}
}

func (f StringUniqueField) LT(val string) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<",
			Value:    val,
		},
	}
}

func (f StringUniqueField) LTE(val string) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<=",
			Value:    val,
		},
	}
}

func (f StringUniqueField) In(vals []string) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IN",
			Value:    vals,
		},
	}
}

func (f StringUniqueField) Like(val string) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "LIKE",
			Value:    val,
		},
	}
}

func (f StringUniqueField) Contains(val string) Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "LIKE",
			Value:    "%" + val + "%",
		},
	}
}

func (f StringUniqueField) IsNull() Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NULL",
		},
	}
}

func (f StringUniqueField) IsNotNull() Predicate {
	return StandardPredicate{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NOT NULL",
		},
	}
}

func CompilePredicates(dialect Dialect, preds []Predicate) (string, []any) {
	if len(preds) == 0 {
		return "", nil
	}
	var data []PredicateData
	for _, p := range preds {
		if p != nil {
			data = append(data, p.ToPredicateData())
		}
	}
	return CompilePredicateData(dialect, data)
}

func CompilePredicateData(dialect Dialect, data []PredicateData) (string, []any) {
	if len(data) == 0 {
		return "", nil
	}
	var parts []string
	var args []any
	var bindIdx = 1

	var compile func(p PredicateData) string
	compile = func(p PredicateData) string {
		if p.IsLogical {
			if len(p.Children) == 0 {
				return ""
			}
			if p.Operator == "NOT" {
				sub := compile(p.Children[0])
				if sub == "" {
					return ""
				}
				return fmt.Sprintf("NOT (%s)", sub)
			}
			var subParts []string
			for _, child := range p.Children {
				sub := compile(child)
				if sub != "" {
					subParts = append(subParts, sub)
				}
			}
			if len(subParts) == 0 {
				return ""
			}
			if len(subParts) == 1 {
				return subParts[0]
			}
			return fmt.Sprintf("(%s)", strings.Join(subParts, " "+p.Operator+" "))
		}

		switch p.Operator {
		case "IS NULL", "IS NOT NULL":
			return fmt.Sprintf("%s %s", dialect.Quote(p.Column), p.Operator)
		case "IN":
			valSlice := unpackSlice(p.Value)
			if len(valSlice) == 0 {
				return "1=0"
			}
			var placeHolders []string
			for range valSlice {
				placeHolders = append(placeHolders, dialect.BindVar(bindIdx))
				bindIdx++
			}
			for _, val := range valSlice {
				args = append(args, val)
			}
			return fmt.Sprintf("%s IN (%s)", dialect.Quote(p.Column), strings.Join(placeHolders, ", "))
		default:
			placeholder := dialect.BindVar(bindIdx)
			bindIdx++
			args = append(args, p.Value)
			return fmt.Sprintf("%s %s %s", dialect.Quote(p.Column), p.Operator, placeholder)
		}
	}

	for _, p := range data {
		part := compile(p)
		if part != "" {
			parts = append(parts, part)
		}
	}

	if len(parts) == 0 {
		return "", nil
	}
	return strings.Join(parts, " AND "), args
}

func unpackSlice(val any) []any {
	if val == nil {
		return nil
	}
	switch v := val.(type) {
	case []string:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []int:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []int32:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []int64:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []float32:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []float64:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []bool:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []time.Time:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case [][]byte:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []json.RawMessage:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	case []any:
		return v
	case []UserRoleType:
		res := make([]any, len(v))
		for i, x := range v {
			res[i] = x
		}
		return res
	default:
		return []any{val}
	}
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
	q.initDelegates()
	q.copyHooksFrom(db.Queries)
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

type CreateManyBuilder[M any, I any] struct {
	client   *Queries
	inputs   []I
	execFunc func(ctx context.Context, inputs []I) (int64, error)
}

func (b *CreateManyBuilder[M, I]) Exec(ctx context.Context) (int64, error) {
	return b.execFunc(ctx, b.inputs)
}

type CreateManyAndReturnBuilder[M any, I any, S any, O any] struct {
	client   *Queries
	inputs   []I
	execFunc func(ctx context.Context, inputs []I, s *S, o *O) ([]*M, error)
}

func (b *CreateManyAndReturnBuilder[M, I, S, O]) Select(s S) *CreateManyAndReturnSelectBuilder[M, I, S, O] {
	return &CreateManyAndReturnSelectBuilder[M, I, S, O]{builder: b, selects: s}
}

func (b *CreateManyAndReturnBuilder[M, I, S, O]) Omit(o O) *CreateManyAndReturnOmitBuilder[M, I, S, O] {
	return &CreateManyAndReturnOmitBuilder[M, I, S, O]{builder: b, omits: o}
}

func (b *CreateManyAndReturnBuilder[M, I, S, O]) Exec(ctx context.Context) ([]*M, error) {
	return b.execFunc(ctx, b.inputs, nil, nil)
}

type CreateManyAndReturnSelectBuilder[M any, I any, S any, O any] struct {
	builder *CreateManyAndReturnBuilder[M, I, S, O]
	selects S
}

func (b *CreateManyAndReturnSelectBuilder[M, I, S, O]) Exec(ctx context.Context) ([]*M, error) {
	return b.builder.execFunc(ctx, b.builder.inputs, &b.selects, nil)
}

type CreateManyAndReturnOmitBuilder[M any, I any, S any, O any] struct {
	builder *CreateManyAndReturnBuilder[M, I, S, O]
	omits   O
}

func (b *CreateManyAndReturnOmitBuilder[M, I, S, O]) Exec(ctx context.Context) ([]*M, error) {
	return b.builder.execFunc(ctx, b.builder.inputs, nil, &b.omits)
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
		sb.WriteString(q.dialect.Quote(col))
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
	for i, c := range cols {
		if c == idCol {
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

type FindUniqueBuilder[M any, S any, O any] struct {
	client   *Queries
	where    UniquePredicate
	execFunc func(ctx context.Context, where UniquePredicate, s *S, o *O) (*M, error)
}

func (b *FindUniqueBuilder[M, S, O]) Select(s S) *FindUniqueSelectBuilder[M, S, O] {
	return &FindUniqueSelectBuilder[M, S, O]{builder: b, selects: s}
}

func (b *FindUniqueBuilder[M, S, O]) Omit(o O) *FindUniqueOmitBuilder[M, S, O] {
	return &FindUniqueOmitBuilder[M, S, O]{builder: b, omits: o}
}

func (b *FindUniqueBuilder[M, S, O]) Exec(ctx context.Context) (*M, error) {
	return b.execFunc(ctx, b.where, nil, nil)
}

type FindUniqueSelectBuilder[M any, S any, O any] struct {
	builder *FindUniqueBuilder[M, S, O]
	selects S
}

func (b *FindUniqueSelectBuilder[M, S, O]) Exec(ctx context.Context) (*M, error) {
	return b.builder.execFunc(ctx, b.builder.where, &b.selects, nil)
}

type FindUniqueOmitBuilder[M any, S any, O any] struct {
	builder *FindUniqueBuilder[M, S, O]
	omits   O
}

func (b *FindUniqueOmitBuilder[M, S, O]) Exec(ctx context.Context) (*M, error) {
	return b.builder.execFunc(ctx, b.builder.where, nil, &b.omits)
}

type FindFirstBuilder[M any, S any, O any] struct {
	client   *Queries
	where    []Predicate
	execFunc func(ctx context.Context, where []Predicate, s *S, o *O) (*M, error)
}

func (b *FindFirstBuilder[M, S, O]) Select(s S) *FindFirstSelectBuilder[M, S, O] {
	return &FindFirstSelectBuilder[M, S, O]{builder: b, selects: s}
}

func (b *FindFirstBuilder[M, S, O]) Omit(o O) *FindFirstOmitBuilder[M, S, O] {
	return &FindFirstOmitBuilder[M, S, O]{builder: b, omits: o}
}

func (b *FindFirstBuilder[M, S, O]) Exec(ctx context.Context) (*M, error) {
	return b.execFunc(ctx, b.where, nil, nil)
}

type FindFirstSelectBuilder[M any, S any, O any] struct {
	builder *FindFirstBuilder[M, S, O]
	selects S
}

func (b *FindFirstSelectBuilder[M, S, O]) Exec(ctx context.Context) (*M, error) {
	return b.builder.execFunc(ctx, b.builder.where, &b.selects, nil)
}

type FindFirstOmitBuilder[M any, S any, O any] struct {
	builder *FindFirstBuilder[M, S, O]
	omits   O
}

func (b *FindFirstOmitBuilder[M, S, O]) Exec(ctx context.Context) (*M, error) {
	return b.builder.execFunc(ctx, b.builder.where, nil, &b.omits)
}

type FindManyBuilder[M any, S any, O any] struct {
	client   *Queries
	where    []Predicate
	execFunc func(ctx context.Context, where []Predicate, s *S, o *O) ([]*M, error)
}

func (b *FindManyBuilder[M, S, O]) Select(s S) *FindManySelectBuilder[M, S, O] {
	return &FindManySelectBuilder[M, S, O]{builder: b, selects: s}
}

func (b *FindManyBuilder[M, S, O]) Omit(o O) *FindManyOmitBuilder[M, S, O] {
	return &FindManyOmitBuilder[M, S, O]{builder: b, omits: o}
}

func (b *FindManyBuilder[M, S, O]) Exec(ctx context.Context) ([]*M, error) {
	return b.execFunc(ctx, b.where, nil, nil)
}

type FindManySelectBuilder[M any, S any, O any] struct {
	builder *FindManyBuilder[M, S, O]
	selects S
}

func (b *FindManySelectBuilder[M, S, O]) Exec(ctx context.Context) ([]*M, error) {
	return b.builder.execFunc(ctx, b.builder.where, &b.selects, nil)
}

type FindManyOmitBuilder[M any, S any, O any] struct {
	builder *FindManyBuilder[M, S, O]
	omits   O
}

func (b *FindManyOmitBuilder[M, S, O]) Exec(ctx context.Context) ([]*M, error) {
	return b.builder.execFunc(ctx, b.builder.where, nil, &b.omits)
}

func executeFindOne[M any](
	ctx context.Context,
	q *Queries,
	table string,
	whereClause string,
	whereVals []any,
	returningCols []string,
	scanFunc func(record *M, cols []string) []any,
) (*M, error) {
	var sb strings.Builder
	sb.Grow(64 + len(returningCols)*15 + len(table) + len(whereClause))
	sb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(q.dialect.Quote(col))
	}
	sb.WriteString(" FROM ")
	sb.WriteString(q.dialect.Quote(table))
	sb.WriteString(whereClause)
	sb.WriteString(" LIMIT 1")

	var res M
	row := q.queryRow(ctx, sb.String(), whereVals...)
	scanTargets := scanFunc(&res, returningCols)
	if err := row.Scan(scanTargets...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func executeFindMany[M any](
	ctx context.Context,
	q *Queries,
	table string,
	whereClause string,
	whereVals []any,
	returningCols []string,
	scanFunc func(record *M, cols []string) []any,
) ([]*M, error) {
	var sb strings.Builder
	sb.Grow(64 + len(returningCols)*15 + len(table) + len(whereClause))
	sb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(q.dialect.Quote(col))
	}
	sb.WriteString(" FROM ")
	sb.WriteString(q.dialect.Quote(table))
	sb.WriteString(whereClause)

	rows, err := q.query(ctx, sb.String(), whereVals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*M, 0)
	for rows.Next() {
		var res M
		scanTargets := scanFunc(&res, returningCols)
		if err := rows.Scan(scanTargets...); err != nil {
			return nil, err
		}
		results = append(results, &res)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func executeSingleWithRelations[M any](
	ctx context.Context,
	q *Queries,
	table string,
	whereClause string,
	whereVals []any,
	returningCols []string,
	scanFunc func(*M, []string) []any,
	hasRelations bool,
	loadRelations func(ctx context.Context, txQ *Queries, results []*M) error,
) (*M, error) {
	if !hasRelations {
		return executeFindOne(ctx, q, table, whereClause, whereVals, returningCols, scanFunc)
	}

	var res *M
	err := q.transaction(ctx, func(txQ *Queries) error {
		var err error
		res, err = executeFindOne(ctx, txQ, table, whereClause, whereVals, returningCols, scanFunc)
		if err != nil || res == nil {
			return err
		}
		return loadRelations(ctx, txQ, []*M{res})
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func executeManyWithRelations[M any](
	ctx context.Context,
	q *Queries,
	table string,
	whereClause string,
	whereVals []any,
	returningCols []string,
	scanFunc func(*M, []string) []any,
	hasRelations bool,
	loadRelations func(ctx context.Context, txQ *Queries, results []*M) error,
) ([]*M, error) {
	if !hasRelations {
		return executeFindMany(ctx, q, table, whereClause, whereVals, returningCols, scanFunc)
	}

	results := make([]*M, 0)
	err := q.transaction(ctx, func(txQ *Queries) error {
		var err error
		results, err = executeFindMany(ctx, txQ, table, whereClause, whereVals, returningCols, scanFunc)
		if err != nil || len(results) == 0 {
			return err
		}
		return loadRelations(ctx, txQ, results)
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

func directKey[T any, K any](get func(*T) K) func(*T) (string, bool) {
	return func(t *T) (string, bool) {
		return fmt.Sprint(get(t)), true
	}
}

func optionalKey[T any, K any](get func(*T) *K) func(*T) (string, bool) {
	return func(t *T) (string, bool) {
		if p := get(t); p != nil {
			return fmt.Sprint(*p), true
		}
		return "", false
	}
}

func setOne[P any, C any](set func(*P, *C)) func(*P, []*C) {
	return func(p *P, children []*C) {
		if len(children) > 0 {
			set(p, children[0])
		}
	}
}

func appendMany[P any, C any](get func(*P) *[]*C) func(*P, []*C) {
	return func(p *P, children []*C) {
		if s := get(p); s != nil {
			*s = append(*s, children...)
		}
	}
}

func scanInto[C any](cols []string, scan func(*C, []string) []any) func(*sql.Rows, *C) error {
	return func(rows *sql.Rows, c *C) error {
		return rows.Scan(scan(c, cols)...)
	}
}

type colSpec struct {
	col      string
	selected bool
	omitted  bool
	forceIn  bool
}

func computeCols(specs []colSpec, hasSelects, anySelected bool) []string {
	var cols []string
	for _, s := range specs {
		include := true
		if hasSelects {
			include = !anySelected || s.selected || s.forceIn
		} else if s.omitted {
			include = false
		}
		if include {
			cols = append(cols, s.col)
		}
	}
	if len(cols) == 0 {
		for _, s := range specs {
			cols = append(cols, s.col)
		}
	}
	return cols
}

func mapToColsVals(m map[string]any, colOrder []string) (cols []string, vals []any) {
	for _, c := range colOrder {
		if v, ok := m[c]; ok {
			cols = append(cols, c)
			vals = append(vals, v)
		}
	}
	return
}

func buildBulkInsertSQL(dialect Dialect, table string, rowMaps []map[string]any, colOrder []string, returningCols []string) (string, []any) {
	colsSet := make(map[string]bool)
	for _, rMap := range rowMaps {
		for col := range rMap {
			colsSet[col] = true
		}
	}
	var cols []string
	for _, c := range colOrder {
		if colsSet[c] {
			cols = append(cols, c)
		}
	}

	var vals []any
	for _, rMap := range rowMaps {
		for _, col := range cols {
			vals = append(vals, rMap[col])
		}
	}

	var sb strings.Builder
	sb.WriteString("INSERT INTO ")
	sb.WriteString(dialect.Quote(table))
	sb.WriteString(" (")
	for i, col := range cols {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(dialect.Quote(col))
	}
	sb.WriteString(") VALUES ")

	paramIndex := 1
	for i := range rowMaps {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("(")
		for j := range cols {
			if j > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(dialect.BindVar(paramIndex))
			paramIndex++
		}
		sb.WriteString(")")
	}

	if dialect.SupportsReturning() && len(returningCols) > 0 {
		sb.WriteString(" RETURNING ")
		for i, col := range returningCols {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(dialect.Quote(col))
		}
	}

	return sb.String(), vals
}
