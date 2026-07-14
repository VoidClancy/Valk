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

type UserDelegate struct {
	client          *Queries
	beforeCreate    func(context.Context, *UserCreate) error
	afterCreate     func(context.Context, []*User) error
	afterCreateMany func(context.Context, []UserCreate, int64) error
}

func (d *UserDelegate) BeforeCreate(hook func(context.Context, *UserCreate) error) {
	d.beforeCreate = hook
}

func (d *UserDelegate) AfterCreate(hook func(context.Context, []*User) error) {
	d.afterCreate = hook
}

func (d *UserDelegate) AfterCreateMany(hook func(context.Context, []UserCreate, int64) error) {
	d.afterCreateMany = hook
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
	"referredById",
}

func (q *Queries) selectUserCols(selects *UserSelect, omits *UserOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return userDefaultCols
	}

	anySelected := selects != nil && (selects.Id || selects.Email || selects.PhoneNum || selects.Password || selects.Role || selects.RoleOptional || selects.ReferredById || selects.Profile != nil || selects.Posts != nil || selects.Comments != nil || selects.ReferredBy != nil || selects.Referrals != nil)

	specs := []colSpec{
		{"id", selects != nil && selects.Id, omits != nil && omits.Id, selects != nil && selects.hasAnyRelation()},
		{"email", selects != nil && selects.Email, omits != nil && omits.Email, false},
		{"phoneNum", selects != nil && selects.PhoneNum, omits != nil && omits.PhoneNum, false},
		{"password", selects != nil && selects.Password, omits != nil && omits.Password, false},
		{"role", selects != nil && selects.Role, omits != nil && omits.Role, false},
		{"roleOptional", selects != nil && selects.RoleOptional, omits != nil && omits.RoleOptional, false},
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

var UserColOrder = []string{
	"id",
	"email",
	"phoneNum",
	"password",
	"role",
	"roleOptional",
	"referredById",
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
func (b *UserCreateBuilder) SetReferredById(v string) *UserCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "referredById", Val: v})
	return b
}

func (d *UserDelegate) Create(assignments ...FieldAssignment) *UserCreateBuilder {
	return &UserCreateBuilder{
		CreateBuilder: &CreateBuilder[User, UserSelect, UserOmit]{
			client:      d.client,
			assignments: assignments,
			execFunc:    d.client.executeUserCreate,
		},
	}
}

func validateUserCreate(assignments []FieldAssignment) error {
	errs := &ValidationError{}

	provided := make(map[string]bool)
	for _, a := range assignments {
		provided[a.Col] = true
		switch a.Col {
		case "id":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "id", v, false, 0, false, false)
			} else {
				errs.Add("id", a.Val, "type", "field id must be of type string")
			}
		case "email":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "email", v, true, 0, false, false)
			} else {
				errs.Add("email", a.Val, "type", "field email must be of type string")
			}
		case "phoneNum":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "phoneNum", v, true, 0, false, false)
			} else {
				errs.Add("phoneNum", a.Val, "type", "field phoneNum must be of type string")
			}
		case "password":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "password", v, false, 0, false, false)
			} else {
				errs.Add("password", a.Val, "type", "field password must be of type string")
			}
		case "role":
			if v, ok := a.Val.(UserRoleType); ok {
				if !v.IsValid() {
					errs.Add("role", v, "enum", fmt.Sprintf("invalid enum value %q for field role", v))
				}
			} else {
				errs.Add("role", a.Val, "type", "field role must be of type UserRoleType")
			}
		case "roleOptional":
			if v, ok := a.Val.(UserRoleType); ok {
				if !v.IsValid() {
					errs.Add("roleOptional", v, "enum", fmt.Sprintf("invalid enum value %q for field roleOptional", v))
				}
			} else {
				errs.Add("roleOptional", a.Val, "type", "field roleOptional must be of type UserRoleType")
			}
		case "referredById":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "referredById", v, false, 0, false, false)
			} else {
				errs.Add("referredById", a.Val, "type", "field referredById must be of type string")
			}
		}
	}
	if !provided["email"] {
		errs.Add("email", "", "required", "field Email is required")
	}
	if !provided["phoneNum"] {
		errs.Add("phoneNum", "", "required", "field PhoneNum is required")
	}

	if errs.HasErrors() {
		return *errs
	}
	return nil
}

func assignmentsToUserCreate(assignments []FieldAssignment) UserCreate {
	var input UserCreate
	for _, a := range assignments {
		switch a.Col {
		case "id":
			if v, ok := a.Val.(string); ok {
				input.Id = &v
			}
		case "email":
			if v, ok := a.Val.(string); ok {
				input.Email = v
			}
		case "phoneNum":
			if v, ok := a.Val.(string); ok {
				input.PhoneNum = v
			}
		case "password":
			if v, ok := a.Val.(string); ok {
				input.Password = &v
			}
		case "role":
			if v, ok := a.Val.(UserRoleType); ok {
				input.Role = &v
			}
		case "roleOptional":
			if v, ok := a.Val.(UserRoleType); ok {
				input.RoleOptional = &v
			}
		case "referredById":
			if v, ok := a.Val.(string); ok {
				input.ReferredById = &v
			}
		}
	}
	return input
}

func (s *UserCreate) ToRowMap() map[string]any {
	m := make(map[string]any, 7)
	if s.Id != nil {
		m["id"] = *s.Id
	} else {
		m["id"] = generateCUID()
	}
	m["email"] = s.Email
	m["phoneNum"] = s.PhoneNum
	if s.Password != nil {
		m["password"] = *s.Password
	}
	if s.Role != nil {
		m["role"] = *s.Role
	}
	if s.RoleOptional != nil {
		m["roleOptional"] = *s.RoleOptional
	}
	if s.ReferredById != nil {
		m["referredById"] = *s.ReferredById
	}
	return m
}

func (q *Queries) executeUserCreate(ctx context.Context, assignments []FieldAssignment, selects *UserSelect, omits *UserOmit) (*User, error) {
	input := assignmentsToUserCreate(assignments)

	if q.User.beforeCreate != nil {
		if err := q.User.beforeCreate(ctx, &input); err != nil {
			return nil, err
		}
	}

	if err := validateUserCreate(assignments); err != nil {
		return nil, err
	}

	rowMap := input.ToRowMap()
	var cols []string
	var vals []any
	for _, col := range UserColOrder {
		if val, ok := rowMap[col]; ok {
			cols = append(cols, col)
			vals = append(vals, val)
		}
	}

	returningCols := q.selectUserCols(selects, omits)

	scanFunc := func(res *User, cols []string) []any {
		return res.ScanFields(cols)
	}

	idCol := "id"

	hasRelations := selects.hasAnyRelation()

	var res *User
	var err error
	if hasRelations {
		err = q.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = executeInsert(ctx, txQ, "User", cols, vals, returningCols, idCol, scanFunc)
			if err != nil {
				return err
			}
			return txQ.loadUserRelations(ctx, []*User{res}, selects)
		})
	} else {
		res, err = executeInsert(ctx, q, "User", cols, vals, returningCols, idCol, scanFunc)
	}
	if err != nil {
		return nil, err
	}

	if q.User.afterCreate != nil {
		if err := q.User.afterCreate(ctx, []*User{res}); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (d *UserDelegate) CreateMany(builders ...*UserCreateBuilder) *CreateManyBuilder[User] {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &CreateManyBuilder[User]{
		client:   d.client,
		records:  records,
		execFunc: d.client.executeUserCreateMany,
	}
}

func (d *UserDelegate) CreateManyAndReturn(builders ...*UserCreateBuilder) *CreateManyAndReturnBuilder[User, UserSelect, UserOmit] {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &CreateManyAndReturnBuilder[User, UserSelect, UserOmit]{
		client:   d.client,
		records:  records,
		execFunc: d.client.executeUserCreateManyAndReturn,
	}
}

func (q *Queries) executeUserCreateMany(ctx context.Context, records []RecordInput) (int64, error) {
	rowMaps := make([]map[string]any, len(records))
	inputs := make([]UserCreate, len(records))
	for i, rec := range records {
		if err := validateUserCreate(rec.Assignments); err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		input := assignmentsToUserCreate(rec.Assignments)
		if q.User.beforeCreate != nil {
			if err := q.User.beforeCreate(ctx, &input); err != nil {
				return 0, err
			}
		}
		rowMaps[i] = input.ToRowMap()
		inputs[i] = input
	}
	count, err := executeCreateMany(ctx, q, rowMaps, "User", UserColOrder)
	if err != nil {
		return 0, err
	}
	if q.User.afterCreateMany != nil {
		if err := q.User.afterCreateMany(ctx, inputs, count); err != nil {
			return 0, err
		}
	}
	return count, nil
}

func (q *Queries) executeUserCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *UserSelect, omits *UserOmit) ([]*User, error) {
	rowMaps := make([]map[string]any, len(records))
	idCol := "id"
	for i, rec := range records {
		if err := validateUserCreate(rec.Assignments); err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		input := assignmentsToUserCreate(rec.Assignments)
		if q.User.beforeCreate != nil {
			if err := q.User.beforeCreate(ctx, &input); err != nil {
				return nil, err
			}
		}
		rowMaps[i] = input.ToRowMap()
	}
	results, err := executeCreateManyAndReturn(ctx, q, rowMaps, "User", UserColOrder, selects, omits,
		q.selectUserCols,
		q.loadUserRelations,
		(*User).ScanFields,
		(*UserSelect).hasAnyRelation,
		idCol,
	)
	if err != nil {
		return nil, err
	}
	if q.User.afterCreate != nil {
		if err := q.User.afterCreate(ctx, results); err != nil {
			return nil, err
		}
	}
	return results, nil
}
func (d *UserDelegate) FindUnique(where UniquePredicate[User], additional ...PredicateOf[User]) *FindUniqueBuilder[User, UserSelect, UserOmit] {
	return &FindUniqueBuilder[User, UserSelect, UserOmit]{
		client:     d.client,
		where:      where,
		additional: additional,
		execFunc:   d.client.executeUserFindUnique,
	}
}

func (d *UserDelegate) FindFirst(preds ...PredicateOf[User]) *FindFirstBuilder[User, UserSelect, UserOmit] {
	return &FindFirstBuilder[User, UserSelect, UserOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeUserFindFirst,
	}
}

func (d *UserDelegate) FindMany(preds ...PredicateOf[User]) *FindManyBuilder[User, UserSelect, UserOmit] {
	return &FindManyBuilder[User, UserSelect, UserOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeUserFindMany,
	}
}

func (q *Queries) executeUserFindUnique(ctx context.Context, where UniquePredicate[User], additional []PredicateOf[User], selects *UserSelect, omits *UserOmit) (*User, error) {
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
	whereClause, vals := CompilePredicates(q.dialect, allPreds)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	returningCols := q.selectUserCols(selects, omits)
	return executeSingleWithRelations(ctx, q, "User", whereClause, vals, returningCols,
		func(res *User, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*User) error {
			return txQ.loadUserRelations(ctx, results, selects)
		},
		nil,
	)
}

func (q *Queries) executeUserFindFirst(
	ctx context.Context,
	params QueryParams[User],
	selects *UserSelect,
	omits *UserOmit,
) (*User, error) {
	for _, p := range params.Where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}
	whereClause, vals := CompilePredicates(q.dialect, params.Where)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	returningCols := q.selectUserCols(selects, omits)
	return executeSingleWithRelations(ctx, q, "User", whereClause, vals, returningCols,
		func(res *User, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*User) error {
			return txQ.loadUserRelations(ctx, results, selects)
		},
		params.Skip,
	)
}

func (q *Queries) executeUserFindMany(
	ctx context.Context,
	params QueryParams[User],
	selects *UserSelect,
	omits *UserOmit,
) ([]*User, error) {
	for _, p := range params.Where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}
	whereClause, vals := CompilePredicates(q.dialect, params.Where)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	returningCols := q.selectUserCols(selects, omits)
	return executeManyWithRelations(ctx, q, "User", whereClause, vals, returningCols,
		func(res *User, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*User) error {
			return txQ.loadUserRelations(ctx, results, selects)
		},
		params.Take,
		params.Skip,
	)
}
func (q *Queries) loadUserRelations(ctx context.Context, records []*User, selects *UserSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Profile != nil {
		relationSelects, relationOmits, relationParams := selects.Profile.GetRelationParams()
		returningCols := q.selectProfileCols(relationSelects, relationOmits, "userId")
		// Inverse holds the FK: Profile.userId
		allChildren, err := loadRelation(
			ctx, q, records,
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
		if err := q.loadProfileRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Posts != nil {
		relationSelects, relationOmits, relationParams := selects.Posts.GetRelationParams()
		returningCols := q.selectPostCols(relationSelects, relationOmits, "authorId")
		// Inverse holds the FK: Post.authorId
		allChildren, err := loadRelation(
			ctx, q, records,
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
		if err := q.loadPostRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Comments != nil {
		relationSelects, relationOmits, relationParams := selects.Comments.GetRelationParams()
		returningCols := q.selectCommentCols(relationSelects, relationOmits, "authorId")
		// Inverse holds the FK: Comment.authorId
		allChildren, err := loadRelation(
			ctx, q, records,
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
		if err := q.loadCommentRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.ReferredBy != nil {
		relationSelects, relationOmits, relationParams := selects.ReferredBy.GetRelationParams()
		returningCols := q.selectUserCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: User.referredById
		allChildren, err := loadRelation(
			ctx, q, records,
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
		if err := q.loadUserRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Referrals != nil {
		relationSelects, relationOmits, relationParams := selects.Referrals.GetRelationParams()
		returningCols := q.selectUserCols(relationSelects, relationOmits, "referredById")
		// Inverse holds the FK: User.referredById
		allChildren, err := loadRelation(
			ctx, q, records,
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
		if err := q.loadUserRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
