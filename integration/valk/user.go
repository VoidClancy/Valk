package valk

import (
	"context"
	"fmt"
	"slices"
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

// UserSelect specifies which fields to include
type UserSelect struct {
	Id           bool               `json:"id"`
	Email        bool               `json:"email"`
	PhoneNum     bool               `json:"phoneNum"`
	Password     bool               `json:"password"`
	Role         bool               `json:"role"`
	RoleOptional bool               `json:"roleOptional"`
	LoginCount   bool               `json:"loginCount"`
	ReferredById bool               `json:"referredById"`
	Profile      ProfileSelectQuery `json:"profile,omitempty"`
	Posts        PostSelectQuery    `json:"posts,omitempty"`
	Comments     CommentSelectQuery `json:"comments,omitempty"`
	ReferredBy   UserSelectQuery    `json:"referredBy,omitempty"`
	Referrals    UserSelectQuery    `json:"referrals,omitempty"`
}

// UserOmit specifies which fields to exclude
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

// UserQueryBuilder builds a query for the relation User
type UserQueryBuilder struct {
	selects *UserSelect
	omits   *UserOmit
	where   []PredicateOf[User]
	take    *int
	skip    *int
	orderBy []OrderBy[User]
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
	}
}

type UserCreateQuery = func(ctx context.Context, args *UserCreate) (*User, error)
type UserCreateManyQuery = func(ctx context.Context, args []*UserCreate) (int64, error)
type UserCreateManyAndReturnQuery = func(ctx context.Context, args []*UserCreate) ([]*User, error)
type UserFindUniqueQuery = func(ctx context.Context, where UniquePredicate[User], additional []PredicateOf[User], selects *UserSelect, omits *UserOmit) (*User, error)
type UserFindFirstQuery = func(ctx context.Context, params QueryParams[User], selects *UserSelect, omits *UserOmit) (*User, error)
type UserFindManyQuery = func(ctx context.Context, params QueryParams[User], selects *UserSelect, omits *UserOmit) ([]*User, error)

type UserExtension struct {
	Create              func(ctx context.Context, input *UserCreate, next UserCreateQuery) (*User, error)
	CreateMany          func(ctx context.Context, inputs []*UserCreate, next UserCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, inputs []*UserCreate, next UserCreateManyAndReturnQuery) ([]*User, error)
	FindUnique          func(ctx context.Context, where UniquePredicate[User], additional []PredicateOf[User], selects *UserSelect, omits *UserOmit, next UserFindUniqueQuery) (*User, error)
	FindFirst           func(ctx context.Context, params QueryParams[User], selects *UserSelect, omits *UserOmit, next UserFindFirstQuery) (*User, error)
	FindMany            func(ctx context.Context, params QueryParams[User], selects *UserSelect, omits *UserOmit, next UserFindManyQuery) ([]*User, error)
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

func (d *UserDelegate) Create(assignments ...FieldAssignment) *UserCreateBuilder {
	return &UserCreateBuilder{
		CreateBuilder: &CreateBuilder[User, UserSelect, UserOmit]{
			assignments: assignments,
			execFunc:    d.executeCreate,
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

func (s *UserCreate) ToRowMap() map[string]any {
	cols, vals := s.ToColsVals()
	m := make(map[string]any, len(cols))
	for i, c := range cols {
		m[c] = vals[i]
	}
	return m
}

func (d *UserDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *UserSelect, omits *UserOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*User, error) {
	input, err := assignmentsToUserCreate(assignments)
	if err != nil {
		return nil, err
	}

	curr := func(c context.Context, args *UserCreate) (*User, error) {
		cols, vals := args.ToColsVals()

		returningCols := selectUserCols(selects, omits)

		scanFunc := func(res *User, cols []string) []any {
			return res.ScanFields(cols)
		}

		pkCols := []string{
			"id",
		}

		hasRelations := selects.hasAnyRelation()

		var res *User
		var err error
		if hasRelations {
			err = d.client.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = executeInsert(c, txQ, "User", cols, vals, returningCols, pkCols, scanFunc, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.User.loadRelations(c, []*User{res}, selects)
			})
		} else {
			res, err = executeInsert(c, d.client, "User", cols, vals, returningCols, pkCols, scanFunc, conflictTarget, conflictAction)
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

	curr := func(c context.Context, args []*UserCreate) (int64, error) {
		rowMaps := make([]map[string]any, len(args))
		for i, input := range args {
			rowMaps[i] = input.ToRowMap()
		}

		pkCols := []string{
			"id",
		}

		return executeCreateMany(c, d.client, rowMaps, "User", userDefaultCols, pkCols, conflictTarget, conflictAction)
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

	curr := func(c context.Context, args []*UserCreate) ([]*User, error) {
		rowMaps := make([]map[string]any, len(args))
		for i, input := range args {
			rowMaps[i] = input.ToRowMap()
		}

		pkCols := []string{
			"id",
		}

		return executeCreateManyAndReturn(c, d.client, rowMaps, "User", userDefaultCols, selects, omits,
			selectUserCols,
			func(ctx context.Context, txQ *Queries, results []*User, sel *UserSelect) error {
				return txQ.User.loadRelations(ctx, results, sel)
			},
			(*User).ScanFields,
			(*UserSelect).hasAnyRelation,
			pkCols,
			conflictTarget,
			conflictAction,
		)
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
	curr := func(c context.Context, w UniquePredicate[User], add []PredicateOf[User], sel *UserSelect, o *UserOmit) (*User, error) {
		if err := w.Validate(); err != nil {
			return nil, err
		}
		for _, p := range add {
			if p != nil {
				if err := p.Validate(); err != nil {
					return nil, err
				}
			}
		}
		allPreds := append([]PredicateOf[User]{w}, add...)
		whereClause, vals := CompilePredicates(d.client.dialect, allPreds)
		if whereClause != "" {
			whereClause = " WHERE " + whereClause
		}
		returningCols := selectUserCols(sel, o)
		return executeSingleWithRelations(c, d.client, "User", whereClause, vals, returningCols,
			func(res *User, cols []string) []any { return res.ScanFields(cols) },
			sel.hasAnyRelation(),
			func(ctx context.Context, txQ *Queries, results []*User) error {
				return txQ.User.loadRelations(ctx, results, sel)
			},
			nil,
		)
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
	curr := func(c context.Context, p QueryParams[User], sel *UserSelect, o *UserOmit) (*User, error) {
		for _, pr := range p.Where {
			if pr != nil {
				if err := pr.Validate(); err != nil {
					return nil, err
				}
			}
		}
		whereClause, vals := CompilePredicates(d.client.dialect, p.Where)
		if whereClause != "" {
			whereClause = " WHERE " + whereClause
		}
		returningCols := selectUserCols(sel, o)
		return executeSingleWithRelations(c, d.client, "User", whereClause, vals, returningCols,
			func(res *User, cols []string) []any { return res.ScanFields(cols) },
			sel.hasAnyRelation(),
			func(ctx context.Context, txQ *Queries, results []*User) error {
				return txQ.User.loadRelations(ctx, results, sel)
			},
			p.Skip,
		)
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
	curr := func(c context.Context, p QueryParams[User], sel *UserSelect, o *UserOmit) ([]*User, error) {
		for _, pr := range p.Where {
			if pr != nil {
				if err := pr.Validate(); err != nil {
					return nil, err
				}
			}
		}
		whereClause, vals := CompilePredicates(d.client.dialect, p.Where)
		if whereClause != "" {
			whereClause = " WHERE " + whereClause
		}
		returningCols := selectUserCols(sel, o)
		return executeManyWithRelations(c, d.client, "User", whereClause, vals, returningCols,
			func(res *User, cols []string) []any { return res.ScanFields(cols) },
			sel.hasAnyRelation(),
			func(ctx context.Context, txQ *Queries, results []*User) error {
				return txQ.User.loadRelations(ctx, results, sel)
			},
			p.Take,
			p.Skip,
		)
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
