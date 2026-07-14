package valk

import (
	"context"
	"crypto/rand"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq/hstore"
	"github.com/pressly/goose/v3"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var _ = time.Time{}
var _ = hstore.Hstore{}
var _ = net.ParseIP
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
func generateUUID7() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New().String()
	}
	return id.String()
}
func generateCUID2() string {
	now := uint64(time.Now().UnixMilli())
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	buf := make([]byte, 0, 34)
	buf = strconv.AppendUint(buf, now, 36)
	const chars = "0123456789abcdefghijklmnopqrstuvwxyz"
	for _, v := range b {
		buf = append(buf, chars[v&0x1f])
		buf = append(buf, chars[(v>>5)&0x1f])
	}
	return string(buf)
}
func generateULID() string {
	now := uint64(time.Now().UnixMilli())
	b := make([]byte, 10)
	_, _ = rand.Read(b)
	bits := make([]byte, 16)
	bits[0] = byte(now >> 40)
	bits[1] = byte(now >> 32)
	bits[2] = byte(now >> 24)
	bits[3] = byte(now >> 16)
	bits[4] = byte(now >> 8)
	bits[5] = byte(now)
	copy(bits[6:], b)

	const crockford = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	var buf [26]byte
	for i := 25; i >= 0; i-- {
		bitOff := i * 5
		byteIdx := bitOff / 8
		byteOff := uint(bitOff % 8)
		val := uint16(bits[byteIdx]) >> byteOff
		if byteIdx+1 < 16 {
			val |= uint16(bits[byteIdx+1]) << (8 - byteOff)
		}
		buf[i] = crockford[val&0x1f]
	}
	return string(buf[:])
}
func generateNanoID() string {
	const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz-"
	b := make([]byte, 21)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = alphabet[b[i]&0x3f]
	}
	return string(b)
}

type UserRoleType string

const (
	// Admin maps to "ADMIN"
	UserRoleTypeAdmin UserRoleType = "ADMIN"
	// Student maps to "student"
	UserRoleTypeStudent UserRoleType = "student"
	// Teacher maps to "TEACHER"
	UserRoleTypeTeacher UserRoleType = "TEACHER"
)

type userRoleNamespace struct {
	// Admin maps to "ADMIN"
	Admin UserRoleType
	// Student maps to "student"
	Student UserRoleType
	// Teacher maps to "TEACHER"
	Teacher UserRoleType
}

// UserRole enum values:
//
//	ADMIN   ADMIN
//	STUDENT student
//	TEACHER TEACHER
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
	FormatLimitOffset(take *int, skip *int) string
}

type DBTX interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
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

type FieldAssignment struct {
	Col string
	Val any
}

type RecordInput struct {
	Assignments []FieldAssignment
}

func ToHstore(m map[string]*string) hstore.Hstore {
	result := hstore.Hstore{Map: make(map[string]sql.NullString, len(m))}
	for k, v := range m {
		if v == nil {
			result.Map[k] = sql.NullString{Valid: false}
		} else {
			result.Map[k] = sql.NullString{String: *v, Valid: true}
		}
	}
	return result
}

type HstoreScan struct {
	P **map[string]*string
}

func (s HstoreScan) Scan(src any) error {
	var h hstore.Hstore
	if err := h.Scan(src); err != nil {
		return err
	}
	if h.Map == nil {
		*s.P = nil
		return nil
	}
	m := make(map[string]*string, len(h.Map))
	for k, v := range h.Map {
		if v.Valid {
			val := v.String
			m[k] = &val
		} else {
			m[k] = nil
		}
	}
	*s.P = &m
	return nil
}

func ValidateString(errs *ValidationError, fieldName string, val string, isRequired bool, maxLen int, isBit bool, isInet bool) {
	if isRequired && val == "" {
		errs.Add(fieldName, val, "required", fmt.Sprintf("field %s is required", fieldName))
	}
	if strings.Contains(val, "\x00") {
		errs.Add(fieldName, val, "safety", "string cannot contain null bytes")
	}
	if !utf8.ValidString(val) {
		errs.Add(fieldName, val, "safety", "string must be valid UTF-8")
	}
	if maxLen > 0 && utf8.RuneCountInString(val) > maxLen {
		errs.Add(fieldName, val, "length", fmt.Sprintf("string exceeds maximum length of %d characters", maxLen))
	}
	if isBit {
		if strings.IndexFunc(val, func(r rune) bool { return r != '0' && r != '1' }) >= 0 {
			errs.Add(fieldName, val, "format", "bit string must contain only '0' and '1'")
		}
	}
	if isInet {
		if net.ParseIP(val) == nil {
			if _, _, err := net.ParseCIDR(val); err != nil {
				errs.Add(fieldName, val, "format", fmt.Sprintf("field %s must be a valid IP address", fieldName))
			}
		}
	}
}

func ValidateInt32(errs *ValidationError, fieldName string, val int32, rule string) {
	switch rule {
	case "SmallInt":
		if val < -32768 || val > 32767 {
			errs.Add(fieldName, val, "range", "value is out of range for SmallInt (-32768 to 32767)")
		}
	case "TinyInt":
		if val < -128 || val > 127 {
			errs.Add(fieldName, val, "range", "value is out of range for TinyInt (-128 to 127)")
		}
	case "Oid":
		if val < 0 {
			errs.Add(fieldName, val, "range", "value is out of range for Oid (must be non-negative)")
		}
	}
}

func ValidateInt64(errs *ValidationError, fieldName string, val int64, rule string) {
	switch rule {
	case "SmallInt":
		if val < -32768 || val > 32767 {
			errs.Add(fieldName, val, "range", "value is out of range for SmallInt (-32768 to 32767)")
		}
	case "TinyInt":
		if val < -128 || val > 127 {
			errs.Add(fieldName, val, "range", "value is out of range for TinyInt (-128 to 127)")
		}
	case "Oid":
		if val < 0 {
			errs.Add(fieldName, val, "range", "value is out of range for Oid (must be non-negative)")
		}
	}
}

func ValidateInt(errs *ValidationError, fieldName string, val int, rule string) {
	switch rule {
	case "SmallInt":
		if val < -32768 || val > 32767 {
			errs.Add(fieldName, val, "range", "value is out of range for SmallInt (-32768 to 32767)")
		}
	case "TinyInt":
		if val < -128 || val > 127 {
			errs.Add(fieldName, val, "range", "value is out of range for TinyInt (-128 to 127)")
		}
	case "Oid":
		if val < 0 {
			errs.Add(fieldName, val, "range", "value is out of range for Oid (must be non-negative)")
		}
	}
}

type OrderDirection string

const (
	Asc  OrderDirection = "ASC"
	Desc OrderDirection = "DESC"
)

type OrderBy struct {
	Field     string
	Direction OrderDirection
}

type QueryParams[M any] struct {
	Where   []PredicateOf[M]
	Take    *int
	Skip    *int
	OrderBy []OrderBy
}

type postgresDialect struct{}

func (postgresDialect) Quote(ident string) string { return `"` + ident + `"` }
func (postgresDialect) BindVar(idx int) string    { return fmt.Sprintf("$%d", idx) }
func (postgresDialect) SupportsReturning() bool   { return true }
func (postgresDialect) SupportsBulkInsert() bool  { return true }
func (postgresDialect) FormatLimitOffset(take *int, skip *int) string {
	if take != nil {
		if skip != nil {
			return fmt.Sprintf(" LIMIT %d OFFSET %d", *take, *skip)
		}
		return fmt.Sprintf(" LIMIT %d", *take)
	}
	if skip != nil {
		return fmt.Sprintf(" OFFSET %d", *skip)
	}
	return ""
}

type Queries struct {
	db       DBTX
	provider string
	dialect  Dialect
	// User provides CRUD operations for User.
	//
	//   id           string   default: cuid()
	//   email        string   required
	//   phoneNum     string   required
	//   password     string   optional
	//   role         UserRole default: STUDENT
	//   roleOptional UserRole optional
	//   referredById string   optional
	User *UserDelegate
	// Profile provides CRUD operations for Profile.
	//
	//   id        string    default: cuid()
	//   bio       string    optional
	//   userId    string    required
	//   createdAt time.Time default: now()
	Profile *ProfileDelegate
	// Post provides CRUD operations for Post.
	//
	//   id        string default: cuid()
	//   title     string required
	//   content   string optional
	//   published bool   default: false
	//   authorId  string required
	Post *PostDelegate
	// Comment provides CRUD operations for Comment.
	//
	//   id       string          default: cuid()
	//   textify  int32           required
	//   dummy3   string          required
	//   dummy1   int32           required
	//   dummy2   string          required
	//   postId   string          required
	//   authorId string          required
	//   meta     json.RawMessage optional
	Comment *CommentDelegate
	// Category provides CRUD operations for Category.
	//
	//   id   int32  default: autoincrement()
	//   name string required
	Category *CategoryDelegate
	// CategoryToPost provides CRUD operations for CategoryToPost.
	//
	//   postId     string required
	//   categoryId int32  required
	CategoryToPost *CategoryToPostDelegate
	// DefaultsTest provides CRUD operations for DefaultsTest.
	//
	//   uuid4      string    default: uuid()
	//   uuid7      string    default: uuid()
	//   uuidNoArgs string    default: uuid()
	//   cuid1      string    default: cuid()
	//   cuid2      string    default: cuid()
	//   cuidNoArgs string    default: cuid()
	//   ulid       string    default: ulid()
	//   nanoid     string    default: nanoid()
	//   now        time.Time default: now()
	DefaultsTest *DefaultsTestDelegate
	// AllFieldsSoFar provides CRUD operations for AllFieldsSoFar.
	//
	//   id              int32              default: autoincrement()
	//   stringReq       string             required
	//   stringOpt       string             optional
	//   stringDefault   string             default: default
	//   stringVarchar   string             required
	//   stringChar      string             required
	//   bitVal          string             required
	//   varBitVal       string             required
	//   inetVal         string             required
	//   xmlVal          string             required
	//   cuidDefault     string             default: cuid()
	//   cuid1Default    string             default: cuid()
	//   cuid2Default    string             default: cuid()
	//   uuidDefault     string             default: uuid()
	//   uuid4Default    string             default: uuid()
	//   uuid7Default    string             default: uuid()
	//   ulidDefault     string             default: ulid()
	//   nanoidDefault   string             default: nanoid()
	//   uuidDb          string             required
	//   intReq          int32              required
	//   intOpt          int32              optional
	//   intDefault      int32              default: 42
	//   integerVal      int32              required
	//   smallInt        int32              required
	//   tinyInt         int32              required
	//   oidVal          int32              required
	//   bigIntReq       int64              required
	//   bigIntOpt       int64              optional
	//   floatReq        float64            required
	//   floatOpt        float64            optional
	//   realVal         float64            required
	//   decimalReq      string             required
	//   decimalOpt      string             optional
	//   decimalPrecise  string             required
	//   moneyVal        string             required
	//   boolReq         bool               required
	//   boolOpt         bool               optional
	//   boolDefault     bool               default: false
	//   dateTimeReq     time.Time          required
	//   dateTimeOpt     time.Time          optional
	//   dateTimeDefault time.Time          default: now()
	//   updatedAt       time.Time          required
	//   dateTimeTz      time.Time          required
	//   timestampVal    time.Time          required
	//   timeVal         time.Time          required
	//   timetzVal       time.Time          required
	//   jsonReq         json.RawMessage    required
	//   jsonOpt         json.RawMessage    optional
	//   jsonVal         json.RawMessage    required
	//   bytesReq        []byte             required
	//   bytesOpt        []byte             optional
	//   hstoreField     map[string]*string optional
	//   ltreeField      string             required
	//   citextField     string             optional
	AllFieldsSoFar *AllFieldsSoFarDelegate
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
		dialect:  postgresDialect{},
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
	q.DefaultsTest = &DefaultsTestDelegate{client: q}
	q.AllFieldsSoFar = &AllFieldsSoFarDelegate{client: q}
}

func (q *Queries) copyHooksFrom(other *Queries) {
	q.User.beforeCreate = other.User.beforeCreate
	q.User.afterCreate = other.User.afterCreate
	q.User.afterCreateMany = other.User.afterCreateMany
	q.Profile.beforeCreate = other.Profile.beforeCreate
	q.Profile.afterCreate = other.Profile.afterCreate
	q.Profile.afterCreateMany = other.Profile.afterCreateMany
	q.Post.beforeCreate = other.Post.beforeCreate
	q.Post.afterCreate = other.Post.afterCreate
	q.Post.afterCreateMany = other.Post.afterCreateMany
	q.Comment.beforeCreate = other.Comment.beforeCreate
	q.Comment.afterCreate = other.Comment.afterCreate
	q.Comment.afterCreateMany = other.Comment.afterCreateMany
	q.Category.beforeCreate = other.Category.beforeCreate
	q.Category.afterCreate = other.Category.afterCreate
	q.Category.afterCreateMany = other.Category.afterCreateMany
	q.CategoryToPost.beforeCreate = other.CategoryToPost.beforeCreate
	q.CategoryToPost.afterCreate = other.CategoryToPost.afterCreate
	q.CategoryToPost.afterCreateMany = other.CategoryToPost.afterCreateMany
	q.DefaultsTest.beforeCreate = other.DefaultsTest.beforeCreate
	q.DefaultsTest.afterCreate = other.DefaultsTest.afterCreate
	q.DefaultsTest.afterCreateMany = other.DefaultsTest.afterCreateMany
	q.AllFieldsSoFar.beforeCreate = other.AllFieldsSoFar.beforeCreate
	q.AllFieldsSoFar.afterCreate = other.AllFieldsSoFar.afterCreate
	q.AllFieldsSoFar.afterCreateMany = other.AllFieldsSoFar.afterCreateMany
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
	log.Printf("[%s] SQL Begin Transaction", strings.ToUpper(q.provider))
	tx, err := starter.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			log.Printf("[%s] SQL Rollback Transaction", strings.ToUpper(q.provider))
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
		log.Printf("[%s] SQL Rollback Transaction", strings.ToUpper(q.provider))
		_ = tx.Rollback()
		return err
	}
	log.Printf("[%s] SQL Commit Transaction", strings.ToUpper(q.provider))
	return tx.Commit()
}

type PredicateData struct {
	Column    string
	Operator  string
	Value     any
	IsLogical bool
	Children  []PredicateData
}

type PredicateOf[M any] interface {
	ToPredicateData() PredicateData
	Validate() error
	phantom(M)
}

type Predicate[M any] struct {
	Data PredicateData
}

func (p Predicate[M]) ToPredicateData() PredicateData {
	return p.Data
}

func (p Predicate[M]) Validate() error {
	return p.Data.Validate()
}

func (p Predicate[M]) phantom(M) {}

type UniquePredicate[M any] struct {
	Data PredicateData
}

func (p UniquePredicate[M]) ToPredicateData() PredicateData {
	return p.Data
}

func (p UniquePredicate[M]) Validate() error {
	if p.Data.Column == "" && len(p.Data.Children) == 0 {
		return fmt.Errorf("at least one unique field must be set for FindUnique")
	}
	return p.Data.Validate()
}

func (p UniquePredicate[M]) phantom(M) {}

func validateValue(col string, val any) error {
	switch v := val.(type) {
	case string:
		if strings.Contains(v, "\x00") {
			return &ValidationError{
				Errors: []FieldError{
					{Field: col, Value: v, Rule: "safety", Msg: "string cannot contain null bytes"},
				},
			}
		}
		if !utf8.ValidString(v) {
			return &ValidationError{
				Errors: []FieldError{
					{Field: col, Value: v, Rule: "safety", Msg: "string must be valid UTF-8"},
				},
			}
		}
	case []string:
		for _, s := range v {
			if err := validateValue(col, s); err != nil {
				return err
			}
		}
	case []any:
		for _, item := range v {
			if err := validateValue(col, item); err != nil {
				return err
			}
		}
	}
	return nil
}

func (pd PredicateData) Validate() error {
	if pd.IsLogical {
		for _, child := range pd.Children {
			if err := child.Validate(); err != nil {
				return err
			}
		}
		return nil
	}
	return validateValue(pd.Column, pd.Value)
}

func And[M any](preds ...PredicateOf[M]) PredicateOf[M] {
	var children []PredicateData
	for _, p := range preds {
		if p != nil {
			children = append(children, p.ToPredicateData())
		}
	}
	return Predicate[M]{
		Data: PredicateData{
			IsLogical: true,
			Operator:  "AND",
			Children:  children,
		},
	}
}

func Or[M any](preds ...PredicateOf[M]) PredicateOf[M] {
	var children []PredicateData
	for _, p := range preds {
		if p != nil {
			children = append(children, p.ToPredicateData())
		}
	}
	return Predicate[M]{
		Data: PredicateData{
			IsLogical: true,
			Operator:  "OR",
			Children:  children,
		},
	}
}

func Not[M any](pred PredicateOf[M]) PredicateOf[M] {
	var children []PredicateData
	if pred != nil {
		children = append(children, pred.ToPredicateData())
	}
	return Predicate[M]{
		Data: PredicateData{
			IsLogical: true,
			Operator:  "NOT",
			Children:  children,
		},
	}
}

type Field[M any, T any] struct {
	Column string
}

func (f Field[M, T]) Set(val T) FieldAssignment {
	return FieldAssignment{Col: f.Column, Val: val}
}

func (f Field[M, T]) EQ(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "=",
			Value:    val,
		},
	}
}

func (f Field[M, T]) NEQ(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "!=",
			Value:    val,
		},
	}
}

func (f Field[M, T]) GT(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">",
			Value:    val,
		},
	}
}

func (f Field[M, T]) GTE(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">=",
			Value:    val,
		},
	}
}

func (f Field[M, T]) LT(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<",
			Value:    val,
		},
	}
}

func (f Field[M, T]) LTE(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<=",
			Value:    val,
		},
	}
}

func (f Field[M, T]) In(vals []T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IN",
			Value:    vals,
		},
	}
}

func (f Field[M, T]) IsNull() Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NULL",
		},
	}
}

func (f Field[M, T]) IsNotNull() Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NOT NULL",
		},
	}
}

func (f Field[M, T]) Asc() OrderBy {
	return OrderBy{Field: f.Column, Direction: Asc}
}

func (f Field[M, T]) Desc() OrderBy {
	return OrderBy{Field: f.Column, Direction: Desc}
}

type UniqueField[M any, T any] struct {
	Column string
}

func (f UniqueField[M, T]) Set(val T) FieldAssignment {
	return FieldAssignment{Col: f.Column, Val: val}
}

func (f UniqueField[M, T]) EQ(val T) UniquePredicate[M] {
	return UniquePredicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "=",
			Value:    val,
		},
	}
}

func (f UniqueField[M, T]) NEQ(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "!=",
			Value:    val,
		},
	}
}

func (f UniqueField[M, T]) GT(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">",
			Value:    val,
		},
	}
}

func (f UniqueField[M, T]) GTE(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">=",
			Value:    val,
		},
	}
}

func (f UniqueField[M, T]) LT(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<",
			Value:    val,
		},
	}
}

func (f UniqueField[M, T]) LTE(val T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<=",
			Value:    val,
		},
	}
}

func (f UniqueField[M, T]) In(vals []T) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IN",
			Value:    vals,
		},
	}
}

func (f UniqueField[M, T]) IsNull() Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NULL",
		},
	}
}

func (f UniqueField[M, T]) IsNotNull() Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NOT NULL",
		},
	}
}

func (f UniqueField[M, T]) Asc() OrderBy {
	return OrderBy{Field: f.Column, Direction: Asc}
}

func (f UniqueField[M, T]) Desc() OrderBy {
	return OrderBy{Field: f.Column, Direction: Desc}
}

type StringField[M any] struct {
	Column string
}

func (f StringField[M]) Set(val string) FieldAssignment {
	return FieldAssignment{Col: f.Column, Val: val}
}

func (f StringField[M]) EQ(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "=",
			Value:    val,
		},
	}
}

func (f StringField[M]) NEQ(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "!=",
			Value:    val,
		},
	}
}

func (f StringField[M]) GT(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">",
			Value:    val,
		},
	}
}

func (f StringField[M]) GTE(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">=",
			Value:    val,
		},
	}
}

func (f StringField[M]) LT(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<",
			Value:    val,
		},
	}
}

func (f StringField[M]) LTE(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<=",
			Value:    val,
		},
	}
}

func (f StringField[M]) In(vals []string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IN",
			Value:    vals,
		},
	}
}

func (f StringField[M]) Like(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "LIKE",
			Value:    val,
		},
	}
}

func (f StringField[M]) Contains(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "LIKE",
			Value:    "%" + val + "%",
		},
	}
}

func (f StringField[M]) IsNull() Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NULL",
		},
	}
}

func (f StringField[M]) IsNotNull() Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NOT NULL",
		},
	}
}

func (f StringField[M]) Asc() OrderBy {
	return OrderBy{Field: f.Column, Direction: Asc}
}

func (f StringField[M]) Desc() OrderBy {
	return OrderBy{Field: f.Column, Direction: Desc}
}

type StringUniqueField[M any] struct {
	Column string
}

func (f StringUniqueField[M]) Set(val string) FieldAssignment {
	return FieldAssignment{Col: f.Column, Val: val}
}

func (f StringUniqueField[M]) EQ(val string) UniquePredicate[M] {
	return UniquePredicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "=",
			Value:    val,
		},
	}
}

func (f StringUniqueField[M]) NEQ(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "!=",
			Value:    val,
		},
	}
}

func (f StringUniqueField[M]) GT(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">",
			Value:    val,
		},
	}
}

func (f StringUniqueField[M]) GTE(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: ">=",
			Value:    val,
		},
	}
}

func (f StringUniqueField[M]) LT(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<",
			Value:    val,
		},
	}
}

func (f StringUniqueField[M]) LTE(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "<=",
			Value:    val,
		},
	}
}

func (f StringUniqueField[M]) In(vals []string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IN",
			Value:    vals,
		},
	}
}

func (f StringUniqueField[M]) Like(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "LIKE",
			Value:    val,
		},
	}
}

func (f StringUniqueField[M]) Contains(val string) Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "LIKE",
			Value:    "%" + val + "%",
		},
	}
}

func (f StringUniqueField[M]) IsNull() Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NULL",
		},
	}
}

func (f StringUniqueField[M]) IsNotNull() Predicate[M] {
	return Predicate[M]{
		Data: PredicateData{
			Column:   f.Column,
			Operator: "IS NOT NULL",
		},
	}
}

func (f StringUniqueField[M]) Asc() OrderBy {
	return OrderBy{Field: f.Column, Direction: Asc}
}

func (f StringUniqueField[M]) Desc() OrderBy {
	return OrderBy{Field: f.Column, Direction: Desc}
}

func CompilePredicates[M any](dialect Dialect, preds []PredicateOf[M]) (string, []any) {
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
	log.Printf("[%s] SQL Begin Transaction", strings.ToUpper(db.provider))
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
	log.Printf("[%s] SQL Commit Transaction", strings.ToUpper(tx.provider))
	return tx.tx.Commit()
}

// Rollback aborts the transaction.
func (tx *Tx) Rollback() error {
	log.Printf("[%s] SQL Rollback Transaction", strings.ToUpper(tx.provider))
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

type CreateBuilder[M any, S any, O any] struct {
	client      *Queries
	assignments []FieldAssignment
	execFunc    func(ctx context.Context, assignments []FieldAssignment, s *S, o *O) (*M, error)
}

func (b *CreateBuilder[M, S, O]) Select(s S) *CreateSelectBuilder[M, S, O] {
	return &CreateSelectBuilder[M, S, O]{builder: b, selects: s}
}

func (b *CreateBuilder[M, S, O]) Omit(o O) *CreateOmitBuilder[M, S, O] {
	return &CreateOmitBuilder[M, S, O]{builder: b, omits: o}
}

func (b *CreateBuilder[M, S, O]) Exec(ctx context.Context) (*M, error) {
	return b.execFunc(ctx, b.assignments, nil, nil)
}

type CreateSelectBuilder[M any, S any, O any] struct {
	builder *CreateBuilder[M, S, O]
	selects S
}

func (b *CreateSelectBuilder[M, S, O]) Exec(ctx context.Context) (*M, error) {
	return b.builder.execFunc(ctx, b.builder.assignments, &b.selects, nil)
}

type CreateOmitBuilder[M any, S any, O any] struct {
	builder *CreateBuilder[M, S, O]
	omits   O
}

func (b *CreateOmitBuilder[M, S, O]) Exec(ctx context.Context) (*M, error) {
	return b.builder.execFunc(ctx, b.builder.assignments, nil, &b.omits)
}

type CreateManyBuilder[M any] struct {
	client   *Queries
	records  []RecordInput
	execFunc func(ctx context.Context, records []RecordInput) (int64, error)
}

func (b *CreateManyBuilder[M]) Exec(ctx context.Context) (int64, error) {
	return b.execFunc(ctx, b.records)
}

type CreateManyAndReturnBuilder[M any, S any, O any] struct {
	client   *Queries
	records  []RecordInput
	execFunc func(ctx context.Context, records []RecordInput, s *S, o *O) ([]*M, error)
}

func (b *CreateManyAndReturnBuilder[M, S, O]) Select(s S) *CreateManyAndReturnSelectBuilder[M, S, O] {
	return &CreateManyAndReturnSelectBuilder[M, S, O]{builder: b, selects: s}
}

func (b *CreateManyAndReturnBuilder[M, S, O]) Omit(o O) *CreateManyAndReturnOmitBuilder[M, S, O] {
	return &CreateManyAndReturnOmitBuilder[M, S, O]{builder: b, omits: o}
}

func (b *CreateManyAndReturnBuilder[M, S, O]) Exec(ctx context.Context) ([]*M, error) {
	return b.execFunc(ctx, b.records, nil, nil)
}

type CreateManyAndReturnSelectBuilder[M any, S any, O any] struct {
	builder *CreateManyAndReturnBuilder[M, S, O]
	selects S
}

func (b *CreateManyAndReturnSelectBuilder[M, S, O]) Exec(ctx context.Context) ([]*M, error) {
	return b.builder.execFunc(ctx, b.builder.records, &b.selects, nil)
}

type CreateManyAndReturnOmitBuilder[M any, S any, O any] struct {
	builder *CreateManyAndReturnBuilder[M, S, O]
	omits   O
}

func (b *CreateManyAndReturnOmitBuilder[M, S, O]) Exec(ctx context.Context) ([]*M, error) {
	return b.builder.execFunc(ctx, b.builder.records, nil, &b.omits)
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

func executeCreateMany(
	ctx context.Context,
	q *Queries,
	rowMaps []map[string]any,
	tableName string,
	colOrder []string,
) (int64, error) {
	if len(rowMaps) == 0 {
		return 0, nil
	}

	if q.dialect.SupportsBulkInsert() {
		query, vals := buildBulkInsertSQL(q.dialect, tableName, rowMaps, colOrder, nil)
		res, err := q.exec(ctx, query, vals...)
		if err != nil {
			return 0, err
		}
		return res.RowsAffected()
	}

	var count int64
	err := q.transaction(ctx, func(txQ *Queries) error {
		for _, rowMap := range rowMaps {
			query, vals := buildBulkInsertSQL(txQ.dialect, tableName, []map[string]any{rowMap}, colOrder, nil)
			_, err := txQ.exec(ctx, query, vals...)
			if err != nil {
				return err
			}
			count++
		}
		return nil
	})
	return count, err
}

func executeCreateManyAndReturn[M any, S any, O any](
	ctx context.Context,
	q *Queries,
	rowMaps []map[string]any,
	tableName string,
	colOrder []string,
	selects *S,
	omits *O,
	selectColsFn func(*S, *O, ...string) []string,
	loadRelationsFn func(context.Context, []*M, *S) error,
	scanFunc func(*M, []string) []any,
	hasRelationsFn func(*S) bool,
	idCol string,
) ([]*M, error) {
	if len(rowMaps) == 0 {
		return nil, nil
	}

	hasRelations := selects != nil && hasRelationsFn(selects)
	returningCols := selectColsFn(selects, omits)

	if q.dialect.SupportsBulkInsert() {
		query, vals := buildBulkInsertSQL(q.dialect, tableName, rowMaps, colOrder, returningCols)
		recordsOut := make([]*M, 0)
		err := q.transaction(ctx, func(txQ *Queries) error {
			rows, err := txQ.query(ctx, query, vals...)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var record M
				if err := rows.Scan(scanFunc(&record, returningCols)...); err != nil {
					return err
				}
				recordsOut = append(recordsOut, &record)
			}
			if err := rows.Err(); err != nil {
				return err
			}
			if hasRelations {
				return loadRelationsFn(ctx, recordsOut, selects)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return recordsOut, nil
	}

	recordsOut := make([]*M, 0)
	err := q.transaction(ctx, func(txQ *Queries) error {
		for _, rowMap := range rowMaps {
			cols, vals := mapToColsVals(rowMap, colOrder)
			res, err := executeInsert(ctx, txQ, tableName, cols, vals, returningCols, idCol, scanFunc)
			if err != nil {
				return err
			}
			recordsOut = append(recordsOut, res)
		}

		if hasRelations {
			return loadRelationsFn(ctx, recordsOut, selects)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return recordsOut, nil
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
	params QueryParams[C],
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

	// Prepend parent ID checks to filters using Predicate[C]
	allPreds := append([]PredicateOf[C]{
		Predicate[C]{
			Data: PredicateData{
				Column:    fkCol,
				Operator:  "IN",
				Value:     parentKeys,
				IsLogical: false,
			},
		},
	}, params.Where...)

	whereClause, vals := CompilePredicates(q.dialect, allPreds)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}

	query := compileRelationSQL(q.dialect, table, fkCol, returningCols, whereClause, params)

	rows, err := q.query(ctx, query, vals...)
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

func compileRelationSQL[M any](dialect Dialect, table, fkCol string, cols []string, where string, params QueryParams[M]) string {
	if params.Take != nil || params.Skip != nil {
		return compilePartitionedRelationSQL(dialect, table, fkCol, cols, where, params)
	}
	return compileSimpleRelationSQL(dialect, table, cols, where, params)
}

func compilePartitionedRelationSQL[M any](dialect Dialect, table, fkCol string, cols []string, where string, params QueryParams[M]) string {
	var innerSb strings.Builder
	innerSb.WriteString("SELECT ")
	for i, col := range cols {
		if i > 0 {
			innerSb.WriteString(", ")
		}
		innerSb.WriteString(dialect.Quote(col))
	}
	innerSb.WriteString(", ROW_NUMBER() OVER (PARTITION BY ")
	innerSb.WriteString(dialect.Quote(fkCol))
	innerSb.WriteString(" ORDER BY ")
	if len(params.OrderBy) > 0 {
		for i, ord := range params.OrderBy {
			if i > 0 {
				innerSb.WriteString(", ")
			}
			innerSb.WriteString(dialect.Quote(ord.Field))
			innerSb.WriteString(" ")
			innerSb.WriteString(string(ord.Direction))
		}
	} else {
		innerSb.WriteString(dialect.Quote("id"))
		innerSb.WriteString(" ASC")
	}
	innerSb.WriteString(") as row_num FROM ")
	innerSb.WriteString(dialect.Quote(table))
	innerSb.WriteString(where)

	var outerSb strings.Builder
	outerSb.WriteString("SELECT ")
	for i, col := range cols {
		if i > 0 {
			outerSb.WriteString(", ")
		}
		outerSb.WriteString(dialect.Quote(col))
	}
	outerSb.WriteString(" FROM (")
	outerSb.WriteString(innerSb.String())
	outerSb.WriteString(") t WHERE ")

	if params.Take != nil && params.Skip != nil {
		outerSb.WriteString(fmt.Sprintf("row_num > %d AND row_num <= %d", *params.Skip, *params.Skip+*params.Take))
	} else if params.Take != nil {
		outerSb.WriteString(fmt.Sprintf("row_num <= %d", *params.Take))
	} else if params.Skip != nil {
		outerSb.WriteString(fmt.Sprintf("row_num > %d", *params.Skip))
	}
	return outerSb.String()
}

func compileSimpleRelationSQL[M any](dialect Dialect, table string, cols []string, where string, params QueryParams[M]) string {
	var sb strings.Builder
	sb.WriteString("SELECT ")
	for i, col := range cols {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(dialect.Quote(col))
	}
	sb.WriteString(" FROM ")
	sb.WriteString(dialect.Quote(table))
	sb.WriteString(where)
	if len(params.OrderBy) > 0 {
		sb.WriteString(" ORDER BY ")
		for i, ord := range params.OrderBy {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(dialect.Quote(ord.Field))
			sb.WriteString(" ")
			sb.WriteString(string(ord.Direction))
		}
	}
	return sb.String()
}

type FindUniqueBuilder[M any, S any, O any] struct {
	client     *Queries
	where      UniquePredicate[M]
	additional []PredicateOf[M]
	execFunc   func(ctx context.Context, where UniquePredicate[M], additional []PredicateOf[M], s *S, o *O) (*M, error)
}

func (b *FindUniqueBuilder[M, S, O]) Select(s S) *FindUniqueSelectBuilder[M, S, O] {
	return &FindUniqueSelectBuilder[M, S, O]{builder: b, selects: s}
}

func (b *FindUniqueBuilder[M, S, O]) Omit(o O) *FindUniqueOmitBuilder[M, S, O] {
	return &FindUniqueOmitBuilder[M, S, O]{builder: b, omits: o}
}

func (b *FindUniqueBuilder[M, S, O]) Exec(ctx context.Context) (*M, error) {
	return b.execFunc(ctx, b.where, b.additional, nil, nil)
}

type FindUniqueSelectBuilder[M any, S any, O any] struct {
	builder *FindUniqueBuilder[M, S, O]
	selects S
}

func (b *FindUniqueSelectBuilder[M, S, O]) Exec(ctx context.Context) (*M, error) {
	return b.builder.execFunc(ctx, b.builder.where, b.builder.additional, &b.selects, nil)
}

type FindUniqueOmitBuilder[M any, S any, O any] struct {
	builder *FindUniqueBuilder[M, S, O]
	omits   O
}

func (b *FindUniqueOmitBuilder[M, S, O]) Exec(ctx context.Context) (*M, error) {
	return b.builder.execFunc(ctx, b.builder.where, b.builder.additional, nil, &b.omits)
}

type FindFirstBuilder[M any, S any, O any] struct {
	client   *Queries
	where    []PredicateOf[M]
	skip     *int
	execFunc func(ctx context.Context, params QueryParams[M], s *S, o *O) (*M, error)
}

func (b *FindFirstBuilder[M, S, O]) Skip(offset int) *FindFirstBuilder[M, S, O] {
	b.skip = &offset
	return b
}

func (b *FindFirstBuilder[M, S, O]) Select(s S) *FindFirstSelectBuilder[M, S, O] {
	return &FindFirstSelectBuilder[M, S, O]{builder: b, selects: s}
}

func (b *FindFirstBuilder[M, S, O]) Omit(o O) *FindFirstOmitBuilder[M, S, O] {
	return &FindFirstOmitBuilder[M, S, O]{builder: b, omits: o}
}

func (b *FindFirstBuilder[M, S, O]) Exec(ctx context.Context) (*M, error) {
	params := QueryParams[M]{
		Where: b.where,
		Skip:  b.skip,
	}
	return b.execFunc(ctx, params, nil, nil)
}

type FindFirstSelectBuilder[M any, S any, O any] struct {
	builder *FindFirstBuilder[M, S, O]
	selects S
}

func (b *FindFirstSelectBuilder[M, S, O]) Exec(ctx context.Context) (*M, error) {
	params := QueryParams[M]{
		Where: b.builder.where,
		Skip:  b.builder.skip,
	}
	return b.builder.execFunc(ctx, params, &b.selects, nil)
}

type FindFirstOmitBuilder[M any, S any, O any] struct {
	builder *FindFirstBuilder[M, S, O]
	omits   O
}

func (b *FindFirstOmitBuilder[M, S, O]) Exec(ctx context.Context) (*M, error) {
	params := QueryParams[M]{
		Where: b.builder.where,
		Skip:  b.builder.skip,
	}
	return b.builder.execFunc(ctx, params, nil, &b.omits)
}

type FindManyBuilder[M any, S any, O any] struct {
	client   *Queries
	where    []PredicateOf[M]
	take     *int
	skip     *int
	execFunc func(ctx context.Context, params QueryParams[M], s *S, o *O) ([]*M, error)
}

func (b *FindManyBuilder[M, S, O]) Take(limit int) *FindManyBuilder[M, S, O] {
	b.take = &limit
	return b
}

func (b *FindManyBuilder[M, S, O]) Skip(offset int) *FindManyBuilder[M, S, O] {
	b.skip = &offset
	return b
}

func (b *FindManyBuilder[M, S, O]) Select(s S) *FindManySelectBuilder[M, S, O] {
	return &FindManySelectBuilder[M, S, O]{builder: b, selects: s}
}

func (b *FindManyBuilder[M, S, O]) Omit(o O) *FindManyOmitBuilder[M, S, O] {
	return &FindManyOmitBuilder[M, S, O]{builder: b, omits: o}
}

func (b *FindManyBuilder[M, S, O]) Exec(ctx context.Context) ([]*M, error) {
	params := QueryParams[M]{
		Where: b.where,
		Take:  b.take,
		Skip:  b.skip,
	}
	return b.execFunc(ctx, params, nil, nil)
}

type FindManySelectBuilder[M any, S any, O any] struct {
	builder *FindManyBuilder[M, S, O]
	selects S
}

func (b *FindManySelectBuilder[M, S, O]) Exec(ctx context.Context) ([]*M, error) {
	params := QueryParams[M]{
		Where: b.builder.where,
		Take:  b.builder.take,
		Skip:  b.builder.skip,
	}
	return b.builder.execFunc(ctx, params, &b.selects, nil)
}

type FindManyOmitBuilder[M any, S any, O any] struct {
	builder *FindManyBuilder[M, S, O]
	omits   O
}

func (b *FindManyOmitBuilder[M, S, O]) Exec(ctx context.Context) ([]*M, error) {
	params := QueryParams[M]{
		Where: b.builder.where,
		Take:  b.builder.take,
		Skip:  b.builder.skip,
	}
	return b.builder.execFunc(ctx, params, nil, &b.omits)
}

func executeFindOne[M any](
	ctx context.Context,
	q *Queries,
	table string,
	whereClause string,
	whereVals []any,
	returningCols []string,
	scanFunc func(record *M, cols []string) []any,
	skip *int,
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
	limitOne := 1
	sb.WriteString(q.dialect.FormatLimitOffset(&limitOne, skip))

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
	take *int,
	skip *int,
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
	sb.WriteString(q.dialect.FormatLimitOffset(take, skip))

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
	skip *int,
) (*M, error) {
	if !hasRelations {
		return executeFindOne(ctx, q, table, whereClause, whereVals, returningCols, scanFunc, skip)
	}

	var res *M
	err := q.transaction(ctx, func(txQ *Queries) error {
		var err error
		res, err = executeFindOne(ctx, txQ, table, whereClause, whereVals, returningCols, scanFunc, skip)
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
	take *int,
	skip *int,
) ([]*M, error) {
	if !hasRelations {
		return executeFindMany(ctx, q, table, whereClause, whereVals, returningCols, scanFunc, take, skip)
	}

	results := make([]*M, 0)
	err := q.transaction(ctx, func(txQ *Queries) error {
		var err error
		results, err = executeFindMany(ctx, txQ, table, whereClause, whereVals, returningCols, scanFunc, take, skip)
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
