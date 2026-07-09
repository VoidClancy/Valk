package valk

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"
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
	Id           bool           `json:"id"`
	Email        bool           `json:"email"`
	PhoneNum     bool           `json:"phoneNum"`
	Password     bool           `json:"password"`
	Role         bool           `json:"role"`
	RoleOptional bool           `json:"roleOptional"`
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
	Password     bool         `json:"password"`
	Role         bool         `json:"role"`
	RoleOptional bool         `json:"roleOptional"`
	ReferredById bool         `json:"referredById"`
	Profile      *ProfileOmit `json:"profile,omitempty"`
	Posts        *PostOmit    `json:"posts,omitempty"`
	Comments     *CommentOmit `json:"comments,omitempty"`
	ReferredBy   *UserOmit    `json:"referredBy,omitempty"`
	Referrals    *UserOmit    `json:"referrals,omitempty"`
}

type UserDelegate struct {
	client       *Queries
	beforeCreate func(context.Context, *UserCreate) error
	afterCreate  func(context.Context, *User) error
}

func (d *UserDelegate) BeforeCreate(hook func(context.Context, *UserCreate) error) {
	d.beforeCreate = hook
}

func (d *UserDelegate) AfterCreate(hook func(context.Context, *User) error) {
	d.afterCreate = hook
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

func (d *UserDelegate) Create(assignments ...FieldAssignment) *CreateBuilder[User, UserSelect, UserOmit] {
	return &CreateBuilder[User, UserSelect, UserOmit]{
		client:      d.client,
		assignments: assignments,
		execFunc:    d.client.executeUserCreate,
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
				if strings.Contains(v, "\x00") {
					errs.Add("id", v, "safety", "string cannot contain null bytes")
				}
				if !utf8.ValidString(v) {
					errs.Add("id", v, "safety", "string must be valid UTF-8")
				}
			}
		case "email":
			if v, ok := a.Val.(string); ok {
				if v == "" {
					errs.Add("email", v, "required", "field email is required")
				}
				if strings.Contains(v, "\x00") {
					errs.Add("email", v, "safety", "string cannot contain null bytes")
				}
				if !utf8.ValidString(v) {
					errs.Add("email", v, "safety", "string must be valid UTF-8")
				}
			}
		case "phoneNum":
			if v, ok := a.Val.(string); ok {
				if v == "" {
					errs.Add("phoneNum", v, "required", "field phoneNum is required")
				}
				if strings.Contains(v, "\x00") {
					errs.Add("phoneNum", v, "safety", "string cannot contain null bytes")
				}
				if !utf8.ValidString(v) {
					errs.Add("phoneNum", v, "safety", "string must be valid UTF-8")
				}
			}
		case "password":
		case "role":
			if v, ok := a.Val.(UserRoleType); ok && !v.IsValid() {
				errs.Add("role", v, "enum", fmt.Sprintf("invalid enum value %q for field role", v))
			}
		case "roleOptional":
			if v, ok := a.Val.(UserRoleType); ok && !v.IsValid() {
				errs.Add("roleOptional", v, "enum", fmt.Sprintf("invalid enum value %q for field roleOptional", v))
			}
		case "referredById":
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

	var cols []string
	var vals []any
	if input.Id != nil {
		cols = append(cols, "id")
		vals = append(vals, *input.Id)
	} else {
		cols = append(cols, "id")
		vals = append(vals, generateCUID())
	}
	cols = append(cols, "email")
	vals = append(vals, input.Email)
	cols = append(cols, "phoneNum")
	vals = append(vals, input.PhoneNum)
	if input.Password != nil {
		cols = append(cols, "password")
		vals = append(vals, *input.Password)
	}
	if input.Role != nil {
		cols = append(cols, "role")
		vals = append(vals, *input.Role)
	}
	if input.RoleOptional != nil {
		cols = append(cols, "roleOptional")
		vals = append(vals, *input.RoleOptional)
	}
	if input.ReferredById != nil {
		cols = append(cols, "referredById")
		vals = append(vals, *input.ReferredById)
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
		if err := q.User.afterCreate(ctx, res); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (d *UserDelegate) CreateMany(records ...RecordInput) *CreateManyBuilder[User] {
	return &CreateManyBuilder[User]{
		client:   d.client,
		records:  records,
		execFunc: d.client.executeUserCreateMany,
	}
}

func (d *UserDelegate) CreateManyAndReturn(records ...RecordInput) *CreateManyAndReturnBuilder[User, UserSelect, UserOmit] {
	return &CreateManyAndReturnBuilder[User, UserSelect, UserOmit]{
		client:   d.client,
		records:  records,
		execFunc: d.client.executeUserCreateManyAndReturn,
	}
}

func (q *Queries) executeUserCreateMany(ctx context.Context, records []RecordInput) (int64, error) {
	rowMaps := make([]map[string]any, len(records))
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
	}
	return executeCreateMany(ctx, q, rowMaps, "User", UserColOrder)
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
	return executeCreateManyAndReturn(ctx, q, rowMaps, "User", UserColOrder, selects, omits,
		q.selectUserCols,
		q.loadUserRelations,
		(*User).ScanFields,
		(*UserSelect).hasAnyRelation,
		idCol,
	)
}
func (d *UserDelegate) FindUnique(where UniquePredicate) *FindUniqueBuilder[User, UserSelect, UserOmit] {
	return &FindUniqueBuilder[User, UserSelect, UserOmit]{
		client:   d.client,
		where:    where,
		execFunc: d.client.executeUserFindUnique,
	}
}

func (d *UserDelegate) FindFirst(preds ...Predicate) *FindFirstBuilder[User, UserSelect, UserOmit] {
	return &FindFirstBuilder[User, UserSelect, UserOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeUserFindFirst,
	}
}

func (d *UserDelegate) FindMany(preds ...Predicate) *FindManyBuilder[User, UserSelect, UserOmit] {
	return &FindManyBuilder[User, UserSelect, UserOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeUserFindMany,
	}
}

func (q *Queries) executeUserFindUnique(ctx context.Context, where UniquePredicate, selects *UserSelect, omits *UserOmit) (*User, error) {
	if where == nil {
		return nil, fmt.Errorf("at least one unique field must be set for FindUnique")
	}
	if err := where.Validate(); err != nil {
		return nil, err
	}
	whereClause, vals := CompilePredicates(q.dialect, []Predicate{where})
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
	)
}

func (q *Queries) executeUserFindFirst(ctx context.Context, where []Predicate, selects *UserSelect, omits *UserOmit) (*User, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}
	whereClause, vals := CompilePredicates(q.dialect, where)
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
	)
}

func (q *Queries) executeUserFindMany(ctx context.Context, where []Predicate, selects *UserSelect, omits *UserOmit) ([]*User, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}
	whereClause, vals := CompilePredicates(q.dialect, where)
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
	)
}
func (q *Queries) loadUserRelations(ctx context.Context, records []*User, selects *UserSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Profile != nil {
		returningCols := q.selectProfileCols(selects.Profile, nil, "userId")
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
		)
		if err != nil {
			return fmt.Errorf("loading profile: %w", err)
		}
		if err := q.loadProfileRelations(ctx, allChildren, selects.Profile); err != nil {
			return err
		}
	}
	if selects.Posts != nil {
		returningCols := q.selectPostCols(selects.Posts, nil, "authorId")
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
		)
		if err != nil {
			return fmt.Errorf("loading posts: %w", err)
		}
		if err := q.loadPostRelations(ctx, allChildren, selects.Posts); err != nil {
			return err
		}
	}
	if selects.Comments != nil {
		returningCols := q.selectCommentCols(selects.Comments, nil, "authorId")
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
		)
		if err != nil {
			return fmt.Errorf("loading comments: %w", err)
		}
		if err := q.loadCommentRelations(ctx, allChildren, selects.Comments); err != nil {
			return err
		}
	}
	if selects.ReferredBy != nil {
		returningCols := q.selectUserCols(selects.ReferredBy, nil, "id")
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
		)
		if err != nil {
			return fmt.Errorf("loading referredBy: %w", err)
		}
		if err := q.loadUserRelations(ctx, allChildren, selects.ReferredBy); err != nil {
			return err
		}
	}
	if selects.Referrals != nil {
		returningCols := q.selectUserCols(selects.Referrals, nil, "referredById")
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
		)
		if err != nil {
			return fmt.Errorf("loading referrals: %w", err)
		}
		if err := q.loadUserRelations(ctx, allChildren, selects.Referrals); err != nil {
			return err
		}
	}

	return nil
}
