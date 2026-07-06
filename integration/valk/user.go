package valk

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

var _ = time.Time{}
var _ = fmt.Sprintf
var _ = strings.Join
var _ = context.Background
var _ = sql.LevelDefault
var _ = slices.Contains[[]string, string]
var _ = utf8.ValidString

// User represents the database model
type User struct {
	Id           string       `db:"id" json:"id"`
	Email        string       `db:"email" json:"email"`
	PhoneNum     string       `db:"phoneNum" json:"phoneNum"`
	Password     *string      `db:"password" json:"password"`
	Role         UserRoleType `db:"role" json:"role"`
	ReferredById *string      `db:"referredById" json:"referredById"`
	Profile      *Profile     `json:"profile,omitempty"`
	Posts        []*Post      `json:"posts,omitempty"`
	Comments     []*Comment   `json:"comments,omitempty"`
	ReferredBy   *User        `json:"referredBy,omitempty"`
	Referrals    []*User      `json:"referrals,omitempty"`
}

// UserCreate represents the input structure for creation
type UserCreate struct {
	Id           *string       `json:"id"`
	Email        string        `json:"email"`
	PhoneNum     string        `json:"phoneNum"`
	Password     *string       `json:"password"`
	Role         *UserRoleType `json:"role"`
	ReferredById *string       `json:"referredById"`
}

// UserSelect specifies which fields to include
type UserSelect struct {
	Id           bool           `json:"id"`
	Email        bool           `json:"email"`
	PhoneNum     bool           `json:"phoneNum"`
	Password     bool           `json:"password"`
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
	Password     bool         `json:"password"`
	Role         bool         `json:"role"`
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
	"referredById",
}

func (q *Queries) selectUserCols(selects *UserSelect, omits *UserOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return userDefaultCols
	}

	anySelected := selects != nil && (selects.Id || selects.Email || selects.PhoneNum || selects.Password || selects.Role || selects.ReferredById || selects.Profile != nil || selects.Posts != nil || selects.Comments != nil || selects.ReferredBy != nil || selects.Referrals != nil)

	specs := []colSpec{
		{"id", selects != nil && selects.Id, omits != nil && omits.Id, false},
		{"email", selects != nil && selects.Email, omits != nil && omits.Email, false},
		{"phoneNum", selects != nil && selects.PhoneNum, omits != nil && omits.PhoneNum, false},
		{"password", selects != nil && selects.Password, omits != nil && omits.Password, false},
		{"role", selects != nil && selects.Role, omits != nil && omits.Role, false},
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

func (input UserCreate) Validate() error {
	errs := &ValidationError{}
	if input.Id != nil {
		val := *input.Id
		if strings.Contains(val, "\x00") {
			errs.Add("id", val, "safety", "string cannot contain null bytes")
		}
		if !utf8.ValidString(val) {
			errs.Add("id", val, "safety", "string must be valid UTF-8")
		}
	}
	if input.Email == "" {
		errs.Add("email", input.Email, "required", "field Email is required")
	}
	if strings.Contains(input.Email, "\x00") {
		errs.Add("email", input.Email, "safety", "string cannot contain null bytes")
	}
	if !utf8.ValidString(input.Email) {
		errs.Add("email", input.Email, "safety", "string must be valid UTF-8")
	}
	if input.PhoneNum == "" {
		errs.Add("phoneNum", input.PhoneNum, "required", "field PhoneNum is required")
	}
	if strings.Contains(input.PhoneNum, "\x00") {
		errs.Add("phoneNum", input.PhoneNum, "safety", "string cannot contain null bytes")
	}
	if !utf8.ValidString(input.PhoneNum) {
		errs.Add("phoneNum", input.PhoneNum, "safety", "string must be valid UTF-8")
	}
	if input.Role != nil {
		if !input.Role.IsValid() {
			errs.Add("role", *input.Role, "enum", fmt.Sprintf("invalid enum value %q for field Role", *input.Role))
		}
	}

	if errs.HasErrors() {
		return *errs
	}
	return nil
}

var UserColOrder = []string{
	"id",
	"email",
	"phoneNum",
	"password",
	"role",
	"referredById",
}

func (s *UserSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Profile != nil || s.Posts != nil || s.Comments != nil || s.ReferredBy != nil || s.Referrals != nil
}

func (d *UserDelegate) Create(input UserCreate) *CreateBuilder[User, UserCreate, UserSelect, UserOmit] {
	return &CreateBuilder[User, UserCreate, UserSelect, UserOmit]{
		client:   d.client,
		input:    input,
		execFunc: d.client.executeUserCreate,
	}
}

func (q *Queries) executeUserCreate(ctx context.Context, input UserCreate, selects *UserSelect, omits *UserOmit) (*User, error) {
	if q.User.beforeCreate != nil {
		if err := q.User.beforeCreate(ctx, &input); err != nil {
			return nil, err
		}
	}

	if err := input.Validate(); err != nil {
		return nil, err
	}
	m := q.UserInputToMap(input)
	cols, vals := mapToColsVals(m, UserColOrder)

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

func (q *Queries) UserInputToMap(input UserCreate) map[string]any {
	m := make(map[string]any)
	if input.Id != nil {
		m["id"] = *input.Id
	} else {
		m["id"] = generateCUID()
	}
	m["email"] = input.Email
	m["phoneNum"] = input.PhoneNum
	if input.Password != nil {
		m["password"] = *input.Password
	}
	if input.Role != nil {
		m["role"] = *input.Role
	}
	if input.ReferredById != nil {
		m["referredById"] = *input.ReferredById
	}
	return m
}

func (d *UserDelegate) CreateMany(inputs []UserCreate) *CreateManyBuilder[User, UserCreate] {
	return &CreateManyBuilder[User, UserCreate]{
		client:   d.client,
		inputs:   inputs,
		execFunc: d.client.executeUserCreateMany,
	}
}

func (d *UserDelegate) CreateManyAndReturn(inputs []UserCreate) *CreateManyAndReturnBuilder[User, UserCreate, UserSelect, UserOmit] {
	return &CreateManyAndReturnBuilder[User, UserCreate, UserSelect, UserOmit]{
		client:   d.client,
		inputs:   inputs,
		execFunc: d.client.executeUserCreateManyAndReturn,
	}
}

func (q *Queries) executeUserCreateMany(ctx context.Context, inputs []UserCreate) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}
	for i, input := range inputs {
		if err := input.Validate(); err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
	}

	if q.dialect.SupportsBulkInsert() {
		rowMaps := make([]map[string]any, len(inputs))
		for i, input := range inputs {
			rowMaps[i] = q.UserInputToMap(input)
		}
		query, vals := buildBulkInsertSQL(q.dialect, "User", rowMaps, UserColOrder, nil)
		res, err := q.exec(ctx, query, vals...)
		if err != nil {
			return 0, err
		}
		return res.RowsAffected()
	}

	var count int64
	err := q.transaction(ctx, func(txQ *Queries) error {
		for _, input := range inputs {
			_, err := txQ.executeUserCreate(ctx, input, nil, nil)
			if err != nil {
				return err
			}
			count++
		}
		return nil
	})
	return count, err
}

func (q *Queries) executeUserCreateManyAndReturn(ctx context.Context, inputs []UserCreate, selects *UserSelect, omits *UserOmit) ([]*User, error) {
	if len(inputs) == 0 {
		return nil, nil
	}
	for i, input := range inputs {
		if err := input.Validate(); err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
	}

	hasRelations := selects.hasAnyRelation()
	returningCols := q.selectUserCols(selects, omits)

	if q.dialect.SupportsBulkInsert() {
		rowMaps := make([]map[string]any, len(inputs))
		for i, input := range inputs {
			rowMaps[i] = q.UserInputToMap(input)
		}
		query, vals := buildBulkInsertSQL(q.dialect, "User", rowMaps, UserColOrder, returningCols)
		var records []*User
		err := q.transaction(ctx, func(txQ *Queries) error {
			rows, err := txQ.query(ctx, query, vals...)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var record User
				if err := rows.Scan(record.ScanFields(returningCols)...); err != nil {
					return err
				}
				records = append(records, &record)
			}
			if err := rows.Err(); err != nil {
				return err
			}
			if hasRelations {
				return txQ.loadUserRelations(ctx, records, selects)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return records, nil
	}

	// Fallback to loop inside transaction
	var records []*User
	err := q.transaction(ctx, func(txQ *Queries) error {
		for _, input := range inputs {
			res, err := txQ.executeUserCreate(ctx, input, nil, nil)
			if err != nil {
				return err
			}
			records = append(records, res)
		}

		if hasRelations {
			return txQ.loadUserRelations(ctx, records, selects)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return records, nil
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
