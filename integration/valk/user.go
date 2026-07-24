package valk

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
)

// User represents the database model
type User struct {
	Id           string        `db:"id" json:"id"`
	Email        string        `db:"email" json:"email"`
	PhoneNum     string        `db:"phoneNum" json:"phoneNum"`
	Password     *string       `db:"password" json:"password,omitempty"`
	Role         UserRoleType  `db:"role" json:"role"`
	RoleOptional *UserRoleType `db:"roleOptional" json:"roleOptional,omitempty"`
	LoginCount   int32         `db:"loginCount" json:"loginCount"`
	ReferredById *string       `db:"referredById" json:"referredById,omitempty"`
	Profile      *Profile      `json:"profile,omitempty"`
	Posts        []*Post       `json:"posts,omitempty"`
	Comments     []*Comment    `json:"comments,omitempty"`
	ReferredBy   *User         `json:"referredBy,omitempty"`
	Referrals    []*User       `json:"referrals,omitempty"`
}

// UserCreate is used for hooks only — the Create API uses FieldAssignment
type UserCreate struct {
	Id           *string       `json:"id"`
	Email        string        `json:"email"`
	PhoneNum     string        `json:"phoneNum"`
	Password     *string       `json:"password"`
	Role         *UserRoleType `json:"role"`
	RoleOptional *UserRoleType `json:"roleOptional"`
	LoginCount   *int32        `json:"loginCount"`
	ReferredById *string       `json:"referredById"`
}

// colMask returns a bit mask of columns that are set
func (s *UserCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	mask |= 1 << 1
	mask |= 1 << 2
	if s.Password != nil {
		mask |= 1 << 3
	}
	if s.Role != nil {
		mask |= 1 << 4
	}
	if s.RoleOptional != nil {
		mask |= 1 << 5
	}
	if s.LoginCount != nil {
		mask |= 1 << 6
	}
	if s.ReferredById != nil {
		mask |= 1 << 7
	}
	return mask
}

type UserSelect struct {
	Id           bool               `json:"id"`
	Email        bool               `json:"email"`
	PhoneNum     bool               `json:"phoneNum"`
	Password     bool               `json:"password"`
	Role         bool               `json:"role"`
	RoleOptional bool               `json:"roleOptional"`
	LoginCount   bool               `json:"loginCount"`
	ReferredById bool               `json:"referredById"`
	Profile      *ProfileSelect     `json:"profile,omitempty"`
	Posts        PostSelectQuery    `json:"posts,omitempty"`
	Comments     CommentSelectQuery `json:"comments,omitempty"`
	ReferredBy   *UserSelect        `json:"referredBy,omitempty"`
	Referrals    UserSelectQuery    `json:"referrals,omitempty"`
}

type UserOmit struct {
	Id           bool `json:"id"`
	Email        bool `json:"email"`
	PhoneNum     bool `json:"phoneNum"`
	Password     bool `json:"password"`
	Role         bool `json:"role"`
	RoleOptional bool `json:"roleOptional"`
	LoginCount   bool `json:"loginCount"`
	ReferredById bool `json:"referredById"`
}

type UserSelectQuery interface {
	GetRelationParams() (*UserSelect, *UserOmit, QueryParams[User])
}

func (s *UserSelect) GetRelationParams() (*UserSelect, *UserOmit, QueryParams[User]) {
	return s, nil, QueryParams[User]{}
}

type UserQueryBuilder struct {
	selects *UserSelect
	omits   *UserOmit
	where   []PredicateOf[User]
	take    *int
	skip    *int
	orderBy []OrderBy[User]
	cursor  UniquePredicate[User]
}

func (b *UserQueryBuilder) Where(preds ...PredicateOf[User]) *UserQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *UserQueryBuilder) Take(limit int) *UserQueryBuilder {
	b.take = &limit
	return b
}

func (b *UserQueryBuilder) Skip(offset int) *UserQueryBuilder {
	b.skip = &offset
	return b
}

func (b *UserQueryBuilder) OrderBy(orders ...OrderBy[User]) *UserQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *UserQueryBuilder) Cursor(where UniquePredicate[User]) *UserQueryBuilder {
	b.cursor = where
	return b
}

func (b *UserQueryBuilder) Select(s UserSelect) *UserQueryBuilder {
	b.selects = &s
	return b
}

func (b *UserQueryBuilder) Omit(o UserOmit) *UserQueryBuilder {
	b.omits = &o
	return b
}

func (b *UserQueryBuilder) GetRelationParams() (*UserSelect, *UserOmit, QueryParams[User]) {
	if b == nil {
		return nil, nil, QueryParams[User]{}
	}
	return b.selects, b.omits, QueryParams[User]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
}

type UserCreateQuery = func(ctx context.Context, args *UserCreate) (*User, error)
type UserCreateManyQuery = func(ctx context.Context, args []*UserCreate) (int64, error)
type UserCreateManyAndReturnQuery = func(ctx context.Context, args []*UserCreate) ([]*User, error)
type UserFindUniqueQuery = func(ctx context.Context, where UniquePredicate[User], additional []PredicateOf[User], selects *UserSelect, omits *UserOmit) (*User, error)
type UserFindFirstQuery = func(ctx context.Context, params QueryParams[User], selects *UserSelect, omits *UserOmit) (*User, error)
type UserFindManyQuery = func(ctx context.Context, params QueryParams[User], selects *UserSelect, omits *UserOmit) ([]*User, error)
type UserDeleteManyQuery = func(ctx context.Context, preds []PredicateOf[User]) (int64, error)
type UserDeleteQuery = func(ctx context.Context, where UniquePredicate[User], selects *UserSelect, omits *UserOmit) (*User, error)
type UserCountQuery = func(ctx context.Context, params QueryParams[User]) (int64, error)
type UserUpdateQuery = func(ctx context.Context, where UniquePredicate[User], additional []PredicateOf[User], assignments []FieldAssignment, selects *UserSelect, omits *UserOmit) (*User, error)
type UserUpdateManyQuery = func(ctx context.Context, preds []PredicateOf[User], assignments []FieldAssignment) (int64, error)
type UserUpdateManyAndReturnQuery = func(ctx context.Context, preds []PredicateOf[User], assignments []FieldAssignment, selects *UserSelect, omits *UserOmit) ([]*User, error)

type UserExtension struct {
	Create              func(ctx context.Context, input *UserCreate, next UserCreateQuery) (*User, error)
	CreateMany          func(ctx context.Context, inputs []*UserCreate, next UserCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, inputs []*UserCreate, next UserCreateManyAndReturnQuery) ([]*User, error)
	FindUnique          func(ctx context.Context, where UniquePredicate[User], additional []PredicateOf[User], selects *UserSelect, omits *UserOmit, next UserFindUniqueQuery) (*User, error)
	FindFirst           func(ctx context.Context, params QueryParams[User], selects *UserSelect, omits *UserOmit, next UserFindFirstQuery) (*User, error)
	FindMany            func(ctx context.Context, params QueryParams[User], selects *UserSelect, omits *UserOmit, next UserFindManyQuery) ([]*User, error)
	DeleteMany          func(ctx context.Context, preds []PredicateOf[User], next UserDeleteManyQuery) (int64, error)
	Delete              func(ctx context.Context, where UniquePredicate[User], selects *UserSelect, omits *UserOmit, next UserDeleteQuery) (*User, error)
	Count               func(ctx context.Context, params QueryParams[User], next UserCountQuery) (int64, error)
	Update              func(ctx context.Context, where UniquePredicate[User], additional []PredicateOf[User], assignments []FieldAssignment, selects *UserSelect, omits *UserOmit, next UserUpdateQuery) (*User, error)
	UpdateMany          func(ctx context.Context, preds []PredicateOf[User], assignments []FieldAssignment, next UserUpdateManyQuery) (int64, error)
	UpdateManyAndReturn func(ctx context.Context, preds []PredicateOf[User], assignments []FieldAssignment, selects *UserSelect, omits *UserOmit, next UserUpdateManyAndReturnQuery) ([]*User, error)
}

type UserDelegate struct {
	client     *Queries
	extensions []UserExtension
}

func (d *UserDelegate) Use(exts ...UserExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *User) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "id":
			targets[i] = &m.Id
		case "email":
			targets[i] = &m.Email
		case "phoneNum":
			targets[i] = &m.PhoneNum
		case "password":
			targets[i] = &m.Password
		case "role":
			targets[i] = &m.Role
		case "roleOptional":
			targets[i] = &m.RoleOptional
		case "loginCount":
			targets[i] = &m.LoginCount
		case "referredById":
			targets[i] = &m.ReferredById
		}
	}
	return targets
}

var userDefaultCols = []string{
	"id",
	"email",
	"phoneNum",
	"password",
	"role",
	"roleOptional",
	"loginCount",
	"referredById",
}

var userPKCols = []string{
	"id",
}

var userUniqueCols = []string{
	"id",
	"email",
	"phoneNum",
}

func selectUserCols(selects *UserSelect, omits *UserOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return userDefaultCols
	}

	anySelected := selects != nil && (selects.Id || selects.Email || selects.PhoneNum || selects.Password || selects.Role || selects.RoleOptional || selects.LoginCount || selects.ReferredById || selects.Profile != nil || selects.Posts != nil || selects.Comments != nil || selects.ReferredBy != nil || selects.Referrals != nil)

	specs := []colSpec{
		{"id", selects != nil && selects.Id, omits != nil && omits.Id, selects != nil && selects.hasAnyRelation()},
		{"email", selects != nil && selects.Email, omits != nil && omits.Email, false},
		{"phoneNum", selects != nil && selects.PhoneNum, omits != nil && omits.PhoneNum, false},
		{"password", selects != nil && selects.Password, omits != nil && omits.Password, false},
		{"role", selects != nil && selects.Role, omits != nil && omits.Role, false},
		{"roleOptional", selects != nil && selects.RoleOptional, omits != nil && omits.RoleOptional, false},
		{"loginCount", selects != nil && selects.LoginCount, omits != nil && omits.LoginCount, false},
		{"referredById", selects != nil && selects.ReferredById, omits != nil && omits.ReferredById, selects != nil && selects.ReferredBy != nil},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *UserSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Profile != nil || s.Posts != nil || s.Comments != nil || s.ReferredBy != nil || s.Referrals != nil
}

type UserCreateBuilder struct {
	*CreateBuilder[User, UserSelect, UserOmit]
}

func (b *UserCreateBuilder) OnConflict(target UniqueConstraintTarget) *UserConflictBuilder[UserCreateBuilder] {
	return &UserConflictBuilder[UserCreateBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *UserCreateBuilder) SetId(v string) *UserCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}
func (b *UserCreateBuilder) SetEmail(v string) *UserCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "email", Val: v})
	return b
}
func (b *UserCreateBuilder) SetPhoneNum(v string) *UserCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "phoneNum", Val: v})
	return b
}
func (b *UserCreateBuilder) SetPassword(v string) *UserCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "password", Val: v})
	return b
}
func (b *UserCreateBuilder) SetRole(v UserRoleType) *UserCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "role", Val: v})
	return b
}
func (b *UserCreateBuilder) SetRoleOptional(v UserRoleType) *UserCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "roleOptional", Val: v})
	return b
}
func (b *UserCreateBuilder) SetLoginCount(v int32) *UserCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "loginCount", Val: v})
	return b
}
func (b *UserCreateBuilder) SetReferredById(v string) *UserCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "referredById", Val: v})
	return b
}

func (b *UserCreateBuilder) Assignments(assignments ...FieldAssignmentOf[User]) *UserCreateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment{Col: a.Col, Val: a.Val})
	}
	return b
}

func (d *UserDelegate) Create() *UserCreateBuilder {
	return &UserCreateBuilder{
		CreateBuilder: &CreateBuilder[User, UserSelect, UserOmit]{
			execFunc: d.executeCreate,
		},
	}
}

const (
	providedUserId           uint64 = 1 << 0
	providedUserEmail        uint64 = 1 << 1
	providedUserPhoneNum     uint64 = 1 << 2
	providedUserPassword     uint64 = 1 << 3
	providedUserRole         uint64 = 1 << 4
	providedUserRoleOptional uint64 = 1 << 5
	providedUserLoginCount   uint64 = 1 << 6
	providedUserReferredById uint64 = 1 << 7
)

func assignmentsToUserCreate(assignments []FieldAssignment) (UserCreate, error) {
	var input UserCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "id":
			provided |= providedUserId
			if v, ok := a.Val.(string); ok {
				input.Id = &v
				ValidateString(&errs, "id", v, false, 0, false, false)
			} else {
				errs.Add("id", a.Val, "type", "field id must be of type string")
			}
		case "email":
			provided |= providedUserEmail
			if v, ok := a.Val.(string); ok {
				input.Email = v
				ValidateString(&errs, "email", v, true, 0, false, false)
			} else {
				errs.Add("email", a.Val, "type", "field email must be of type string")
			}
		case "phoneNum":
			provided |= providedUserPhoneNum
			if v, ok := a.Val.(string); ok {
				input.PhoneNum = v
				ValidateString(&errs, "phoneNum", v, true, 0, false, false)
			} else {
				errs.Add("phoneNum", a.Val, "type", "field phoneNum must be of type string")
			}
		case "password":
			provided |= providedUserPassword
			if v, ok := a.Val.(string); ok {
				input.Password = &v
				ValidateString(&errs, "password", v, false, 0, false, false)
			} else {
				errs.Add("password", a.Val, "type", "field password must be of type string")
			}
		case "role":
			provided |= providedUserRole
			if v, ok := a.Val.(UserRoleType); ok {
				input.Role = &v
				if !v.IsValid() {
					errs.Add("role", v, "enum", fmt.Sprintf("invalid enum value %q for field role", v))
				}
			} else {
				errs.Add("role", a.Val, "type", "field role must be of type UserRoleType")
			}
		case "roleOptional":
			provided |= providedUserRoleOptional
			if v, ok := a.Val.(UserRoleType); ok {
				input.RoleOptional = &v
				if !v.IsValid() {
					errs.Add("roleOptional", v, "enum", fmt.Sprintf("invalid enum value %q for field roleOptional", v))
				}
			} else {
				errs.Add("roleOptional", a.Val, "type", "field roleOptional must be of type UserRoleType")
			}
		case "loginCount":
			provided |= providedUserLoginCount
			if v, ok := a.Val.(int32); ok {
				input.LoginCount = &v
				ValidateInt32(&errs, "loginCount", v, "")
			} else {
				errs.Add("loginCount", a.Val, "type", "field loginCount must be of type int32")
			}
		case "referredById":
			provided |= providedUserReferredById
			if v, ok := a.Val.(string); ok {
				input.ReferredById = &v
				ValidateString(&errs, "referredById", v, false, 0, false, false)
			} else {
				errs.Add("referredById", a.Val, "type", "field referredById must be of type string")
			}
		}
	}
	if provided&providedUserEmail == 0 {
		errs.Add("email", "", "required", "field Email is required")
	}
	if provided&providedUserPhoneNum == 0 {
		errs.Add("phoneNum", "", "required", "field PhoneNum is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *UserCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 8)
	vals = make([]any, 0, 8)
	cols = append(cols, "id")
	if s.Id != nil {
		vals = append(vals, *s.Id)
	} else {
		vals = append(vals, generateCUID())
	}
	cols = append(cols, "email")
	vals = append(vals, s.Email)
	cols = append(cols, "phoneNum")
	vals = append(vals, s.PhoneNum)
	if s.Password != nil {
		cols = append(cols, "password")
		vals = append(vals, *s.Password)
	}
	if s.Role != nil {
		cols = append(cols, "role")
		vals = append(vals, *s.Role)
	}
	if s.RoleOptional != nil {
		cols = append(cols, "roleOptional")
		vals = append(vals, *s.RoleOptional)
	}
	if s.LoginCount != nil {
		cols = append(cols, "loginCount")
		vals = append(vals, *s.LoginCount)
	}
	if s.ReferredById != nil {
		cols = append(cols, "referredById")
		vals = append(vals, *s.ReferredById)
	}
	return
}

func partitionUserInputs(dialect Dialect, inputs []*UserCreate) [][]*UserCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*UserCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*UserCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*UserCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*UserCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*UserCreate{inputs}
}

func (d *UserDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *UserSelect, omits *UserOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*User, error) {
	input, err := assignmentsToUserCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectUserCols(selects, omits)

	if len(d.extensions) == 0 {
		hasRelations := selects.hasAnyRelation()
		if hasRelations {
			var res *User
			err = d.client.transaction(ctx, func(txQ *Queries) error {
				var err error
				res, err = txQ.User.runCreate(ctx, cols, vals, returningCols, userPKCols, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.User.loadRelations(ctx, []*User{res}, selects)
			})
			return res, err
		}
		return d.runCreate(ctx, cols, vals, returningCols, userPKCols, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args *UserCreate) (*User, error) {
		cols, vals := args.ToColsVals()
		returningCols := selectUserCols(selects, omits)

		hasRelations := selects.hasAnyRelation()
		var res *User
		var err error
		if hasRelations {
			err = d.client.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = txQ.User.runCreate(c, cols, vals, returningCols, userPKCols, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.User.loadRelations(c, []*User{res}, selects)
			})
		} else {
			res, err = d.runCreate(c, cols, vals, returningCols, userPKCols, conflictTarget, conflictAction)
		}
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Create != nil {
			next, hook := curr, ext.Create
			curr = func(c context.Context, input *UserCreate) (*User, error) {
				return hook(c, input, next)
			}
		}
	}

	return curr(ctx, &input)
}

type UserCreateManyBuilder struct {
	*CreateManyBuilder[User]
}

func (b *UserCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *UserConflictBuilder[UserCreateManyBuilder] {
	return &UserConflictBuilder[UserCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type UserCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[User, UserSelect, UserOmit]
}

func (b *UserCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *UserConflictBuilder[UserCreateManyAndReturnBuilder] {
	return &UserConflictBuilder[UserCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (d *UserDelegate) CreateMany(builders ...*UserCreateBuilder) *UserCreateManyBuilder {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &UserCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[User]{
			records:  records,
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *UserDelegate) CreateManyAndReturn(builders ...*UserCreateBuilder) *UserCreateManyAndReturnBuilder {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &UserCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[User, UserSelect, UserOmit]{
			records:  records,
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func (d *UserDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs := make([]*UserCreate, len(records))
	for i, rec := range records {
		input, err := assignmentsToUserCreate(rec.Assignments)
		if err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &input
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args []*UserCreate) (int64, error) {
		return d.runCreateMany(c, args, conflictTarget, conflictAction)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.CreateMany != nil {
			next, hook := curr, ext.CreateMany
			curr = func(c context.Context, inputs []*UserCreate) (int64, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
}

func (d *UserDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *UserSelect, omits *UserOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*User, error) {
	inputs := make([]*UserCreate, len(records))
	for i, rec := range records {
		input, err := assignmentsToUserCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &input
	}

	if len(d.extensions) == 0 {
		hasRelations := selects != nil && selects.hasAnyRelation()
		if hasRelations {
			var res []*User
			err := d.client.transaction(ctx, func(txQ *Queries) error {
				var err error
				res, err = txQ.User.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.User.loadRelations(ctx, res, selects)
			})
			return res, err
		}
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args []*UserCreate) ([]*User, error) {
		hasRelations := selects != nil && selects.hasAnyRelation()
		if hasRelations {
			var res []*User
			err := d.client.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = txQ.User.runCreateManyAndReturn(c, args, selects, omits, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.User.loadRelations(c, res, selects)
			})
			return res, err
		}
		return d.runCreateManyAndReturn(c, args, selects, omits, conflictTarget, conflictAction)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.CreateManyAndReturn != nil {
			next, hook := curr, ext.CreateManyAndReturn
			curr = func(c context.Context, inputs []*UserCreate) ([]*User, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
}

func (d *UserDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	pkCols []string,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*User, error) {
	query, clauseArgs := buildSingleInsertSQL(d.client, "User", cols, returningCols, pkCols, conflictTarget, conflictAction, len(vals))
	if len(clauseArgs) > 0 {
		vals = append(vals, clauseArgs...)
	}

	var res User
	if d.client.dialect.SupportsInsertReturning {
		rows, err := d.client.query(ctx, query, vals...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		if rows.Next() {
			if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
				return nil, err
			}
			return &res, nil
		}
		return nil, rows.Err()
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, pkCols)
}

func (d *UserDelegate) runCreateFallback(ctx context.Context, query string, vals []any, cols []string, returningCols []string, pkCols []string) (*User, error) {
	result, err := d.client.exec(ctx, query, vals...)
	if err != nil {
		return nil, err
	}

	var pkVals []any
	for _, pkCol := range pkCols {
		var val any
		for i, c := range cols {
			if c == pkCol {
				val = vals[i]
				break
			}
		}
		if val == nil && len(pkCols) == 1 {
			lastID, err := result.LastInsertId()
			if err != nil {
				return nil, err
			}
			val = lastID
		}
		pkVals = append(pkVals, val)
	}

	var selectSb strings.Builder
	selectSb.Grow(64 + len(returningCols)*15 + len("User") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "User")
	selectSb.WriteString(" WHERE ")
	for i, pkCol := range pkCols {
		if i > 0 {
			selectSb.WriteString(" AND ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, pkCol)
		selectSb.WriteString(" = ")
		d.client.dialect.WritePlaceholder(&selectSb, i+1)
	}

	rows, err := d.client.query(ctx, selectSb.String(), pkVals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res User
	if rows.Next() {
		if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
			return nil, err
		}
		return &res, nil
	}
	return nil, rows.Err()
}

func (d *UserDelegate) buildBulkInsertSQL(q *Queries, batch []*UserCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 8)
	for i, c := range userDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "User")
	sb.WriteString(" (")
	for i, col := range cols {
		if i > 0 {
			sb.WriteString(", ")
		}
		q.dialect.WriteQuotedIdent(&sb, col)
	}
	sb.WriteString(") VALUES ")

	paramIdx := paramStartIdx
	for ri, input := range batch {
		if ri > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("(")
		for j, col := range cols {
			if j > 0 {
				sb.WriteString(", ")
			}
			writeDefault := false
			switch col {
			case "id":
				if input.Id != nil {
					vals = append(vals, *input.Id)
				} else {
					vals = append(vals, generateCUID())
				}
			case "email":
				vals = append(vals, input.Email)
			case "phoneNum":
				vals = append(vals, input.PhoneNum)
			case "password":
				if input.Password != nil {
					vals = append(vals, *input.Password)
				} else {
					writeDefault = true
				}
			case "role":
				if input.Role != nil {
					vals = append(vals, *input.Role)
				} else {
					writeDefault = true
				}
			case "roleOptional":
				if input.RoleOptional != nil {
					vals = append(vals, *input.RoleOptional)
				} else {
					writeDefault = true
				}
			case "loginCount":
				if input.LoginCount != nil {
					vals = append(vals, *input.LoginCount)
				} else {
					writeDefault = true
				}
			case "referredById":
				if input.ReferredById != nil {
					vals = append(vals, *input.ReferredById)
				} else {
					writeDefault = true
				}
			}
			if writeDefault {
				sb.WriteString("DEFAULT")
			} else {
				q.dialect.WritePlaceholder(&sb, paramIdx)
				paramIdx++
			}
		}
		sb.WriteString(")")
	}
	queryStr = sb.String()
	return cols, vals, queryStr
}

func (d *UserDelegate) runCreateMany(ctx context.Context, inputs []*UserCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionUserInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)

		var conflictCols []string
		if conflictTarget != nil {
			conflictCols = conflictTarget.UniqueColumns()
		}
		var nonConflictCols []string
		if conflictAction != nil && conflictAction.Type == ConflictActionUpdateNewValues {
			nonConflictCols = computeNonConflictCols(cols, conflictCols, userPKCols)
		}
		clause, clauseArgs := d.client.dialect.BuildConflictClause(conflictCols, conflictAction, nonConflictCols, len(vals)+1)
		queryStr += clause
		vals = append(vals, clauseArgs...)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return 0, err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		count += affected
	}
	return count, nil
}

func (d *UserDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*UserCreate,
	selects *UserSelect,
	omits *UserOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*User, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionUserInputs(d.client.dialect, inputs)
	returningCols := selectUserCols(selects, omits)
	hasRelations := selects != nil && selects.hasAnyRelation()

	recordsOut := make([]*User, 0, len(inputs))

	runBatch := func(txQ *Queries, batch []*UserCreate) error {
		cols, vals, queryStr := d.buildBulkInsertSQL(txQ, batch, 1)

		var conflictCols []string
		if conflictTarget != nil {
			conflictCols = conflictTarget.UniqueColumns()
		}
		var nonConflictCols []string
		if conflictAction != nil && conflictAction.Type == ConflictActionUpdateNewValues {
			nonConflictCols = computeNonConflictCols(cols, conflictCols, userPKCols)
		}
		clause, clauseArgs := txQ.dialect.BuildConflictClause(conflictCols, conflictAction, nonConflictCols, len(vals)+1)
		queryStr += clause
		vals = append(vals, clauseArgs...)

		if txQ.dialect.SupportsInsertReturning && len(returningCols) > 0 {
			var retSb strings.Builder
			retSb.Grow(12 + len(returningCols)*15)
			retSb.WriteString(" RETURNING ")
			for i, col := range returningCols {
				if i > 0 {
					retSb.WriteString(", ")
				}
				txQ.dialect.WriteQuotedIdent(&retSb, col)
			}
			queryStr += retSb.String()
			rows, err := txQ.query(ctx, queryStr, vals...)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var res User
				if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
					return err
				}
				recordsOut = append(recordsOut, &res)
			}
			return rows.Err()
		}

		// Fallback for dialects without RETURNING (MySQL)
		result, err := txQ.exec(ctx, queryStr, vals...)
		if err != nil {
			return err
		}

		// We need to fetch the inserted records for this batch
		// Note: MySQL bulk inserts only return the ID of the FIRST inserted row
		lastID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// Query back the rows by IDs (assuming autoincrement ID and single PK)
		// If composite PK, it's more complex, but this is a standard fallback
		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("User") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			txQ.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		txQ.dialect.WriteQuotedIdent(&selectSb, "User")
		selectSb.WriteString(" WHERE ")
		txQ.dialect.WriteQuotedIdent(&selectSb, userPKCols[0])
		selectSb.WriteString(" >= ")
		txQ.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		txQ.dialect.WriteQuotedIdent(&selectSb, userPKCols[0])
		selectSb.WriteString(" < ")
		txQ.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := txQ.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var res User
			if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
				return err
			}
			recordsOut = append(recordsOut, &res)
		}
		return rows.Err()
	}

	// Always wrap in transaction if we have multiple batches OR if we need to load relations
	if len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning {
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			for _, batch := range batches {
				if err := runBatch(txQ, batch); err != nil {
					return err
				}
			}
			if hasRelations {
				return txQ.User.loadRelations(ctx, recordsOut, selects)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		if err := runBatch(d.client, batches[0]); err != nil {
			return nil, err
		}
	}

	return recordsOut, nil
}

type UserConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *UserConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *UserConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *UserConflictBuilder[B]) Update(fn func(u *UserUpsert)) *B {
	var up ConflictUpdate
	u := newUserUpsert(&up)
	fn(u)
	cb.setAction(ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}, cb.conflictTarget)
	return cb.builder
}

type UserUpsert struct {
	Id           fieldUpsert[string]
	Email        fieldUpsert[string]
	PhoneNum     fieldUpsert[string]
	Password     fieldUpsert[*string]
	Role         fieldUpsert[string]
	RoleOptional fieldUpsert[*string]
	LoginCount   numericFieldUpsert[int32]
	ReferredById fieldUpsert[*string]
}

func newUserUpsert(up *ConflictUpdate) *UserUpsert {
	return &UserUpsert{
		Id:           fieldUpsert[string]{column: "id", update: up},
		Email:        fieldUpsert[string]{column: "email", update: up},
		PhoneNum:     fieldUpsert[string]{column: "phoneNum", update: up},
		Password:     fieldUpsert[*string]{column: "password", update: up},
		Role:         fieldUpsert[string]{column: "role", update: up},
		RoleOptional: fieldUpsert[*string]{column: "roleOptional", update: up},
		LoginCount: numericFieldUpsert[int32]{
			fieldUpsert: fieldUpsert[int32]{column: "loginCount", update: up},
			tableName:   "User",
		},
		ReferredById: fieldUpsert[*string]{column: "referredById", update: up},
	}
}

type UserUpdateBuilder struct {
	*UpdateBuilder[User, UserSelect, UserOmit]
}

type UserUpdateManyBuilder struct {
	*UpdateManyBuilder[User]
}

type UserUpdateManyAndReturnBuilder struct {
	*UpdateManyAndReturnBuilder[User, UserSelect, UserOmit]
}

func (b *UserUpdateBuilder) SetId(v string) *UserUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}

func (b *UserUpdateManyBuilder) SetId(v string) *UserUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}

func (b *UserUpdateManyAndReturnBuilder) SetId(v string) *UserUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}
func (b *UserUpdateBuilder) SetEmail(v string) *UserUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "email", Val: v})
	return b
}

func (b *UserUpdateManyBuilder) SetEmail(v string) *UserUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "email", Val: v})
	return b
}

func (b *UserUpdateManyAndReturnBuilder) SetEmail(v string) *UserUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "email", Val: v})
	return b
}
func (b *UserUpdateBuilder) SetPhoneNum(v string) *UserUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "phoneNum", Val: v})
	return b
}

func (b *UserUpdateManyBuilder) SetPhoneNum(v string) *UserUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "phoneNum", Val: v})
	return b
}

func (b *UserUpdateManyAndReturnBuilder) SetPhoneNum(v string) *UserUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "phoneNum", Val: v})
	return b
}
func (b *UserUpdateBuilder) SetPassword(v string) *UserUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "password", Val: v})
	return b
}

func (b *UserUpdateManyBuilder) SetPassword(v string) *UserUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "password", Val: v})
	return b
}

func (b *UserUpdateManyAndReturnBuilder) SetPassword(v string) *UserUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "password", Val: v})
	return b
}
func (b *UserUpdateBuilder) SetRole(v UserRoleType) *UserUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "role", Val: v})
	return b
}

func (b *UserUpdateManyBuilder) SetRole(v UserRoleType) *UserUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "role", Val: v})
	return b
}

func (b *UserUpdateManyAndReturnBuilder) SetRole(v UserRoleType) *UserUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "role", Val: v})
	return b
}
func (b *UserUpdateBuilder) SetRoleOptional(v UserRoleType) *UserUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "roleOptional", Val: v})
	return b
}

func (b *UserUpdateManyBuilder) SetRoleOptional(v UserRoleType) *UserUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "roleOptional", Val: v})
	return b
}

func (b *UserUpdateManyAndReturnBuilder) SetRoleOptional(v UserRoleType) *UserUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "roleOptional", Val: v})
	return b
}
func (b *UserUpdateBuilder) SetLoginCount(v int32) *UserUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "loginCount", Val: v})
	return b
}

func (b *UserUpdateManyBuilder) SetLoginCount(v int32) *UserUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "loginCount", Val: v})
	return b
}

func (b *UserUpdateManyAndReturnBuilder) SetLoginCount(v int32) *UserUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "loginCount", Val: v})
	return b
}
func (b *UserUpdateBuilder) SetReferredById(v string) *UserUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "referredById", Val: v})
	return b
}

func (b *UserUpdateManyBuilder) SetReferredById(v string) *UserUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "referredById", Val: v})
	return b
}

func (b *UserUpdateManyAndReturnBuilder) SetReferredById(v string) *UserUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "referredById", Val: v})
	return b
}

func (b *UserUpdateBuilder) Assignments(assignments ...FieldAssignmentOf[User]) *UserUpdateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment{Col: a.Col, Val: a.Val})
	}
	return b
}

func (b *UserUpdateManyBuilder) Assignments(assignments ...FieldAssignmentOf[User]) *UserUpdateManyBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment{Col: a.Col, Val: a.Val})
	}
	return b
}

func (b *UserUpdateManyAndReturnBuilder) Assignments(assignments ...FieldAssignmentOf[User]) *UserUpdateManyAndReturnBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment{Col: a.Col, Val: a.Val})
	}
	return b
}

func (d *UserDelegate) Update(where UniquePredicate[User], additional ...PredicateOf[User]) *UserUpdateBuilder {
	return &UserUpdateBuilder{
		UpdateBuilder: &UpdateBuilder[User, UserSelect, UserOmit]{
			where:      where,
			additional: additional,
			execFunc:   d.executeUpdate,
		},
	}
}

func (d *UserDelegate) UpdateMany(preds ...PredicateOf[User]) *UserUpdateManyBuilder {
	return &UserUpdateManyBuilder{
		UpdateManyBuilder: &UpdateManyBuilder[User]{
			where:    preds,
			execFunc: d.executeUpdateMany,
		},
	}
}

func (d *UserDelegate) UpdateManyAndReturn(preds ...PredicateOf[User]) *UserUpdateManyAndReturnBuilder {
	return &UserUpdateManyAndReturnBuilder{
		UpdateManyAndReturnBuilder: &UpdateManyAndReturnBuilder[User, UserSelect, UserOmit]{
			where:    preds,
			execFunc: d.executeUpdateManyAndReturn,
		},
	}
}

func (d *UserDelegate) buildUpdateSQL(preds []PredicateOf[User], assignments []FieldAssignment, returningCols []string) (string, []any) {
	whereClause, predVals, _ := CompilePredicates(d.client.dialect, preds, len(assignments)+1)

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	d.client.dialect.WriteQuotedIdent(&sb, "User")
	sb.WriteString(" SET ")

	setVals := make([]any, 0, len(assignments)+len(predVals))
	for i, a := range assignments {
		if i > 0 {
			sb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&sb, a.Col)
		sb.WriteString(" = ")
		d.client.dialect.WritePlaceholder(&sb, i+1)
		setVals = append(setVals, a.Val)
	}

	if whereClause != "" {
		sb.WriteString(" WHERE ")
		sb.WriteString(whereClause)
		setVals = append(setVals, predVals...)
	}

	if len(returningCols) > 0 && d.client.dialect.SupportsUpdateReturning {
		sb.WriteString(" RETURNING ")
		for i, col := range returningCols {
			if i > 0 {
				sb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&sb, col)
		}
	}

	return sb.String(), setVals
}

// -----------------------------------------------------------------------------
// Update
// -----------------------------------------------------------------------------

func (d *UserDelegate) executeUpdate(ctx context.Context, where UniquePredicate[User], additional []PredicateOf[User], assignments []FieldAssignment, selects *UserSelect, omits *UserOmit) (*User, error) {
	if len(d.extensions) == 0 {
		return d.runUpdate(ctx, where, additional, assignments, selects, omits)
	}

	curr := func(c context.Context, w UniquePredicate[User], add []PredicateOf[User], a []FieldAssignment, s *UserSelect, o *UserOmit) (*User, error) {
		return d.runUpdate(c, w, add, a, s, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Update != nil {
			next, hook := curr, ext.Update
			curr = func(c context.Context, w UniquePredicate[User], add []PredicateOf[User], a []FieldAssignment, s *UserSelect, o *UserOmit) (*User, error) {
				return hook(c, w, add, a, s, o, next)
			}
		}
	}

	return curr(ctx, where, additional, assignments, selects, omits)
}

func (d *UserDelegate) runUpdate(ctx context.Context, where UniquePredicate[User], additional []PredicateOf[User], assignments []FieldAssignment, selects *UserSelect, omits *UserOmit) (*User, error) {
	if len(assignments) == 0 {
		return d.runFindUnique(ctx, where, additional, selects, omits)
	}

	if err := where.Validate(); err != nil {
		return nil, err
	}
	for _, pr := range additional {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return nil, err
			}
		}
	}

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (!d.client.dialect.SupportsUpdateReturning || hasRelations) && !d.client.inTx()

	if useTx {
		var res *User
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.User.runUpdate(ctx, where, additional, assignments, selects, omits)
			} else {
				res, err = txQ.User.runUpdateFallback(ctx, where, additional, assignments, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectUserCols(selects, omits, userPKCols...)
	allPreds := append([]PredicateOf[User]{where}, additional...)
	query, setVals := d.buildUpdateSQL(allPreds, assignments, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		err := rows.Err()
		rows.Close()
		if err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	}

	var res User
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, scanErr
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*User{&res}, selects); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *UserDelegate) execUpdateStmt(ctx context.Context, preds []PredicateOf[User], assignments []FieldAssignment) (int64, error) {
	if len(assignments) == 0 {
		return 0, nil
	}

	for _, pr := range preds {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return 0, err
			}
		}
	}

	query, setVals := d.buildUpdateSQL(preds, assignments, nil)
	result, err := d.client.exec(ctx, query, setVals...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (d *UserDelegate) runUpdateFallback(ctx context.Context, where UniquePredicate[User], additional []PredicateOf[User], assignments []FieldAssignment, selects *UserSelect, omits *UserOmit) (*User, error) {
	allPreds := append([]PredicateOf[User]{where}, additional...)
	affected, err := d.execUpdateStmt(ctx, allPreds, assignments)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, sql.ErrNoRows
	}
	return d.runFindUnique(ctx, where, additional, selects, omits)
}

// -----------------------------------------------------------------------------
// UpdateMany
// -----------------------------------------------------------------------------

func (d *UserDelegate) executeUpdateMany(ctx context.Context, preds []PredicateOf[User], assignments []FieldAssignment) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runUpdateMany(ctx, preds, assignments)
	}

	curr := func(c context.Context, p []PredicateOf[User], a []FieldAssignment) (int64, error) {
		return d.runUpdateMany(c, p, a)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.UpdateMany != nil {
			next, hook := curr, ext.UpdateMany
			curr = func(c context.Context, p []PredicateOf[User], a []FieldAssignment) (int64, error) {
				return hook(c, p, a, next)
			}
		}
	}

	return curr(ctx, preds, assignments)
}

func (d *UserDelegate) runUpdateMany(ctx context.Context, preds []PredicateOf[User], assignments []FieldAssignment) (int64, error) {
	return d.execUpdateStmt(ctx, preds, assignments)
}

// -----------------------------------------------------------------------------
// UpdateManyAndReturn
// -----------------------------------------------------------------------------

func (d *UserDelegate) executeUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[User], assignments []FieldAssignment, selects *UserSelect, omits *UserOmit) ([]*User, error) {
	if len(d.extensions) == 0 {
		return d.runUpdateManyAndReturn(ctx, preds, assignments, selects, omits)
	}

	curr := func(c context.Context, p []PredicateOf[User], a []FieldAssignment, s *UserSelect, o *UserOmit) ([]*User, error) {
		return d.runUpdateManyAndReturn(c, p, a, s, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.UpdateManyAndReturn != nil {
			next, hook := curr, ext.UpdateManyAndReturn
			curr = func(c context.Context, p []PredicateOf[User], a []FieldAssignment, s *UserSelect, o *UserOmit) ([]*User, error) {
				return hook(c, p, a, s, o, next)
			}
		}
	}

	return curr(ctx, preds, assignments, selects, omits)
}

func (d *UserDelegate) runUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[User], assignments []FieldAssignment, selects *UserSelect, omits *UserOmit) ([]*User, error) {
	if len(assignments) == 0 {
		return d.runFindMany(ctx, QueryParams[User]{Where: preds}, selects, omits)
	}

	for _, pr := range preds {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return nil, err
			}
		}
	}

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (!d.client.dialect.SupportsUpdateReturning || hasRelations) && !d.client.inTx()

	if useTx {
		var res []*User
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.User.runUpdateManyAndReturn(ctx, preds, assignments, selects, omits)
			} else {
				res, err = txQ.User.runUpdateManyAndReturnFallback(ctx, preds, assignments, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectUserCols(selects, omits, userPKCols...)
	query, setVals := d.buildUpdateSQL(preds, assignments, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	results := make([]*User, 0)
	for rows.Next() {
		var res User
		if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
			rows.Close()
			return nil, err
		}
		results = append(results, &res)
	}
	rowsErr := rows.Err()
	rows.Close()
	if rowsErr != nil {
		return nil, rowsErr
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, results, selects); err != nil {
			return nil, err
		}
	}

	return results, nil
}

func (d *UserDelegate) runUpdateManyAndReturnFallback(ctx context.Context, preds []PredicateOf[User], assignments []FieldAssignment, selects *UserSelect, omits *UserOmit) ([]*User, error) {
	affected, err := d.execUpdateStmt(ctx, preds, assignments)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return []*User{}, nil
	}
	return d.runFindMany(ctx, QueryParams[User]{Where: preds}, selects, omits)
}
func (d *UserDelegate) FindUnique(where UniquePredicate[User], additional ...PredicateOf[User]) *FindUniqueBuilder[User, UserSelect, UserOmit] {
	return &FindUniqueBuilder[User, UserSelect, UserOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *UserDelegate) FindFirst(preds ...PredicateOf[User]) *FindFirstBuilder[User, UserSelect, UserOmit] {
	return &FindFirstBuilder[User, UserSelect, UserOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *UserDelegate) FindMany(preds ...PredicateOf[User]) *FindManyBuilder[User, UserSelect, UserOmit] {
	return &FindManyBuilder[User, UserSelect, UserOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *UserDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[User], additional []PredicateOf[User], selects *UserSelect, omits *UserOmit) (*User, error) {
	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, where, additional, selects, omits)
	}

	curr := func(c context.Context, w UniquePredicate[User], add []PredicateOf[User], sel *UserSelect, o *UserOmit) (*User, error) {
		return d.runFindUnique(c, w, add, sel, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindUnique != nil {
			next, hook := curr, ext.FindUnique
			curr = func(c context.Context, w UniquePredicate[User], add []PredicateOf[User], sel *UserSelect, o *UserOmit) (*User, error) {
				return hook(c, w, add, sel, o, next)
			}
		}
	}

	return curr(ctx, where, additional, selects, omits)
}

func (d *UserDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[User],
	selects *UserSelect,
	omits *UserOmit,
) (*User, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	curr := func(c context.Context, p QueryParams[User], sel *UserSelect, o *UserOmit) (*User, error) {
		return d.runFindFirst(c, p, sel, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindFirst != nil {
			next, hook := curr, ext.FindFirst
			curr = func(c context.Context, p QueryParams[User], sel *UserSelect, o *UserOmit) (*User, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}

func (d *UserDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[User],
	selects *UserSelect,
	omits *UserOmit,
) ([]*User, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	curr := func(c context.Context, p QueryParams[User], sel *UserSelect, o *UserOmit) ([]*User, error) {
		return d.runFindMany(c, p, sel, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindMany != nil {
			next, hook := curr, ext.FindMany
			curr = func(c context.Context, p QueryParams[User], sel *UserSelect, o *UserOmit) ([]*User, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}

func (d *UserDelegate) runFindUnique(ctx context.Context, where UniquePredicate[User], additional []PredicateOf[User], selects *UserSelect, omits *UserOmit) (*User, error) {
	if err := where.Validate(); err != nil {
		return nil, err
	}
	for _, p := range additional {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}
	allPreds := append([]PredicateOf[User]{where}, additional...)
	whereClause, vals, _ := CompilePredicates(d.client.dialect, allPreds)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	returningCols := selectUserCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*User{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *UserDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[User],
	selects *UserSelect,
	omits *UserOmit,
) (*User, error) {
	for _, pr := range params.Where {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return nil, err
			}
		}
	}
	whereClause, vals, nextIdx := CompilePredicates(d.client.dialect, params.Where)
	isCursorQuery := (params.Cursor.Data.Column != "" || len(params.Cursor.Data.Children) > 0)
	if isCursorQuery {
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, userPKCols, userUniqueCols, "User", nextIdx, params.Take)
		if err != nil {
			return nil, err
		}
		if cClause != "" {
			if whereClause == "" {
				whereClause = cClause
			} else {
				whereClause = "(" + whereClause + ") AND " + cClause
			}
			vals = append(vals, cVals...)
		}
	}
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, userPKCols, userUniqueCols, isCursorQuery, params.Take)
	returningCols := selectUserCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, orderByClause, vals, returningCols, params.Skip)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*User{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *UserDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[User],
	selects *UserSelect,
	omits *UserOmit,
) ([]*User, error) {
	for _, pr := range params.Where {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return nil, err
			}
		}
	}
	whereClause, vals, nextIdx := CompilePredicates(d.client.dialect, params.Where)
	isCursorQuery := (params.Cursor.Data.Column != "" || len(params.Cursor.Data.Children) > 0)
	if isCursorQuery {
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, userPKCols, userUniqueCols, "User", nextIdx, params.Take)
		if err != nil {
			return nil, err
		}
		if cClause != "" {
			if whereClause == "" {
				whereClause = cClause
			} else {
				whereClause = "(" + whereClause + ") AND " + cClause
			}
			vals = append(vals, cVals...)
		}
	}
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, userPKCols, userUniqueCols, isCursorQuery, params.Take)
	returningCols := selectUserCols(selects, omits)

	results, err := d.queryMany(ctx, whereClause, orderByClause, vals, returningCols, params.Take, params.Skip)
	if err != nil || len(results) == 0 {
		return results, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, results, selects); err != nil {
			return nil, err
		}
	}
	return results, nil
}

func (d *UserDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*User, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "User", returningCols, whereClause, orderByClause, &limitOne, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		return nil, nil
	}

	var res User
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (d *UserDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*User, error) {
	query := buildSelectSQL(d.client, "User", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*User, 0)
	for rows.Next() {
		var res User
		if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
			return nil, err
		}
		results = append(results, &res)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if take != nil && *take < 0 {
		reverseSlice(results)
	}
	return results, nil
}
func (d *UserDelegate) DeleteMany(preds ...PredicateOf[User]) *DeleteManyBuilder[User] {
	return &DeleteManyBuilder[User]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *UserDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[User]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	curr := func(c context.Context, p []PredicateOf[User]) (int64, error) {
		return d.runDeleteMany(c, p)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.DeleteMany != nil {
			next, hook := curr, ext.DeleteMany
			curr = func(c context.Context, p []PredicateOf[User]) (int64, error) {
				return hook(c, p, next)
			}
		}
	}

	return curr(ctx, preds)
}

func (d *UserDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[User]) (int64, error) {
	for _, pr := range preds {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return 0, err
			}
		}
	}

	whereClause, vals, _ := CompilePredicates(d.client.dialect, preds)

	var sb strings.Builder
	sb.WriteString("DELETE FROM ")
	d.client.dialect.WriteQuotedIdent(&sb, "User")
	if whereClause != "" {
		sb.WriteString(" WHERE ")
		sb.WriteString(whereClause)
	}

	result, err := d.client.exec(ctx, sb.String(), vals...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (d *UserDelegate) Delete(where UniquePredicate[User]) *DeleteBuilder[User, UserSelect, UserOmit] {
	return &DeleteBuilder[User, UserSelect, UserOmit]{
		where:    where,
		execFunc: d.executeDelete,
	}
}

func (d *UserDelegate) executeDelete(ctx context.Context, where UniquePredicate[User], selects *UserSelect, omits *UserOmit) (*User, error) {
	if len(d.extensions) == 0 {
		return d.runDelete(ctx, where, selects, omits)
	}

	curr := func(c context.Context, w UniquePredicate[User], s *UserSelect, o *UserOmit) (*User, error) {
		return d.runDelete(c, w, s, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Delete != nil {
			next, hook := curr, ext.Delete
			curr = func(c context.Context, w UniquePredicate[User], s *UserSelect, o *UserOmit) (*User, error) {
				return hook(c, w, s, o, next)
			}
		}
	}

	return curr(ctx, where, selects, omits)
}

func (d *UserDelegate) runDelete(ctx context.Context, where UniquePredicate[User], selects *UserSelect, omits *UserOmit) (*User, error) {
	if err := where.Validate(); err != nil {
		return nil, err
	}

	returningCols := selectUserCols(selects, omits, userPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *User
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.User.executeFindUnique(ctx, where, nil, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return sql.ErrNoRows
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "User")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[User]
			pkPreds = append(pkPreds, Predicate[User]{
				Data: PredicateData{
					Column:   "id",
					Operator: "=",
					Value:    res.Id,
				},
			})

			whereClause, vals, _ := CompilePredicates(txQ.dialect, pkPreds)
			deleteSb.WriteString(whereClause)

			_, err = txQ.exec(ctx, deleteSb.String(), vals...)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	// Dialect supports RETURNING, and no relations need loading: run direct DELETE ... RETURNING
	var sb strings.Builder
	sb.WriteString("DELETE FROM ")
	d.client.dialect.WriteQuotedIdent(&sb, "User")

	whereClause, vals, _ := CompilePredicates(d.client.dialect, []PredicateOf[User]{where})
	if whereClause != "" {
		sb.WriteString(" WHERE ")
		sb.WriteString(whereClause)
	}

	sb.WriteString(" RETURNING ")
	for i, col := range returningCols {
		if i > 0 {
			sb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&sb, col)
	}

	rows, err := d.client.query(ctx, sb.String(), vals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	}

	var row User
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, err
	}
	return &row, nil
}
func (d *UserDelegate) Count(preds ...PredicateOf[User]) *CountBuilder[User] {
	return &CountBuilder[User]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *UserDelegate) executeCount(ctx context.Context, params QueryParams[User]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	curr := func(c context.Context, p QueryParams[User]) (int64, error) {
		return d.runCount(c, p)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Count != nil {
			next, hook := curr, ext.Count
			curr = func(c context.Context, p QueryParams[User]) (int64, error) {
				return hook(c, p, next)
			}
		}
	}

	return curr(ctx, params)
}

func (d *UserDelegate) runCount(ctx context.Context, params QueryParams[User]) (int64, error) {
	for _, pr := range params.Where {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return 0, err
			}
		}
	}

	whereClause, vals, _ := CompilePredicates(d.client.dialect, params.Where)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}

	var query string
	if params.Take != nil || params.Skip != nil {
		var subQuery strings.Builder
		subQuery.WriteString("SELECT 1 FROM ")
		d.client.dialect.WriteQuotedIdent(&subQuery, "User")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "User")
		if whereClause != "" {
			sb.WriteString(whereClause)
		}
		query = sb.String()
	}

	rows, err := d.client.query(ctx, query, vals...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var count int64
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	return count, nil
}
func (d *UserDelegate) loadRelations(ctx context.Context, records []*User, selects *UserSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Profile != nil {
		relationSelects, relationOmits, relationParams := selects.Profile.GetRelationParams()
		returningCols := selectProfileCols(relationSelects, relationOmits, "userId")
		// Inverse holds the FK: Profile.userId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *User) string { return p.Id }),
			"Profile",
			"userId",
			returningCols,
			scanInto(returningCols, (*Profile).ScanFields),
			directKey(func(c *Profile) string { return c.UserId }),
			setOne(func(p *User, c *Profile) { p.Profile = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading profile: %w", err)
		}
		if err := d.client.Profile.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Posts != nil {
		relationSelects, relationOmits, relationParams := selects.Posts.GetRelationParams()
		returningCols := selectPostCols(relationSelects, relationOmits, "authorId")
		// Inverse holds the FK: Post.authorId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *User) string { return p.Id }),
			"Post",
			"authorId",
			returningCols,
			scanInto(returningCols, (*Post).ScanFields),
			directKey(func(c *Post) string { return c.AuthorId }),
			appendMany(func(p *User) *[]*Post { return &p.Posts }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading posts: %w", err)
		}
		if err := d.client.Post.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Comments != nil {
		relationSelects, relationOmits, relationParams := selects.Comments.GetRelationParams()
		returningCols := selectCommentCols(relationSelects, relationOmits, "authorId")
		// Inverse holds the FK: Comment.authorId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *User) string { return p.Id }),
			"Comment",
			"authorId",
			returningCols,
			scanInto(returningCols, (*Comment).ScanFields),
			directKey(func(c *Comment) string { return c.AuthorId }),
			appendMany(func(p *User) *[]*Comment { return &p.Comments }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading comments: %w", err)
		}
		if err := d.client.Comment.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.ReferredBy != nil {
		relationSelects, relationOmits, relationParams := selects.ReferredBy.GetRelationParams()
		returningCols := selectUserCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: User.referredById
		allChildren, err := loadRelation(
			ctx, d.client, records,
			optionalKey(func(p *User) *string { return p.ReferredById }),
			"User",
			"id",
			returningCols,
			scanInto(returningCols, (*User).ScanFields),
			directKey(func(c *User) string { return c.Id }),
			setOne(func(p *User, c *User) { p.ReferredBy = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading referredBy: %w", err)
		}
		if err := d.client.User.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Referrals != nil {
		relationSelects, relationOmits, relationParams := selects.Referrals.GetRelationParams()
		returningCols := selectUserCols(relationSelects, relationOmits, "referredById")
		// Inverse holds the FK: User.referredById
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *User) string { return p.Id }),
			"User",
			"referredById",
			returningCols,
			scanInto(returningCols, (*User).ScanFields),
			optionalKey(func(c *User) *string { return c.ReferredById }),
			appendMany(func(p *User) *[]*User { return &p.Referrals }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading referrals: %w", err)
		}
		if err := d.client.User.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
