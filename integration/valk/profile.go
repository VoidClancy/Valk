package valk

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

// Profile represents the database model
type Profile struct {
	Id        string    `db:"id" json:"id"`
	Bio       *string   `db:"bio" json:"bio,omitempty"`
	UserId    string    `db:"userId" json:"userId"`
	CreatedAt time.Time `db:"createdAt" json:"createdAt"`
	User      *User     `json:"user,omitempty"`
}

// ProfileCreate is used for hooks only — the Create API uses FieldAssignment
type ProfileCreate struct {
	Id        *string    `json:"id"`
	Bio       *string    `json:"bio"`
	UserId    string     `json:"userId"`
	CreatedAt *time.Time `json:"createdAt"`
}

// ProfileSelect specifies which fields to include
type ProfileSelect struct {
	Id        bool        `json:"id"`
	Bio       bool        `json:"bio"`
	UserId    bool        `json:"userId"`
	CreatedAt bool        `json:"createdAt"`
	User      *UserSelect `json:"user,omitempty"`
}

// ProfileOmit specifies which fields to exclude
type ProfileOmit struct {
	Id        bool      `json:"id"`
	Bio       bool      `json:"bio"`
	UserId    bool      `json:"userId"`
	CreatedAt bool      `json:"createdAt"`
	User      *UserOmit `json:"user,omitempty"`
}

type ProfileDelegate struct {
	client          *Queries
	beforeCreate    func(context.Context, *ProfileCreate) error
	afterCreate     func(context.Context, []*Profile) error
	afterCreateMany func(context.Context, []ProfileCreate, int64) error
}

func (d *ProfileDelegate) BeforeCreate(hook func(context.Context, *ProfileCreate) error) {
	d.beforeCreate = hook
}

func (d *ProfileDelegate) AfterCreate(hook func(context.Context, []*Profile) error) {
	d.afterCreate = hook
}

func (d *ProfileDelegate) AfterCreateMany(hook func(context.Context, []ProfileCreate, int64) error) {
	d.afterCreateMany = hook
}

func (m *Profile) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "id":
			targets[i] = &m.Id
		case "bio":
			targets[i] = &m.Bio
		case "userId":
			targets[i] = &m.UserId
		case "createdAt":
			targets[i] = &m.CreatedAt
		}
	}
	return targets
}

var profileDefaultCols = []string{
	"id",
	"bio",
	"userId",
	"createdAt",
}

func (q *Queries) selectProfileCols(selects *ProfileSelect, omits *ProfileOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return profileDefaultCols
	}

	anySelected := selects != nil && (selects.Id || selects.Bio || selects.UserId || selects.CreatedAt || selects.User != nil)

	specs := []colSpec{
		{"id", selects != nil && selects.Id, omits != nil && omits.Id, selects != nil && selects.hasAnyRelation()},
		{"bio", selects != nil && selects.Bio, omits != nil && omits.Bio, false},
		{"userId", selects != nil && selects.UserId, omits != nil && omits.UserId, selects != nil && selects.User != nil},
		{"createdAt", selects != nil && selects.CreatedAt, omits != nil && omits.CreatedAt, false},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

var ProfileColOrder = []string{
	"id",
	"bio",
	"userId",
	"createdAt",
}

func (s *ProfileSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.User != nil
}

func (d *ProfileDelegate) Create(assignments ...FieldAssignment) *CreateBuilder[Profile, ProfileSelect, ProfileOmit] {
	return &CreateBuilder[Profile, ProfileSelect, ProfileOmit]{
		client:      d.client,
		assignments: assignments,
		execFunc:    d.client.executeProfileCreate,
	}
}

func validateProfileCreate(assignments []FieldAssignment) error {
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
		case "bio":
			if _, ok := a.Val.(string); !ok {
				errs.Add("bio", a.Val, "type", "field bio must be of type string")
			}
		case "userId":
			if v, ok := a.Val.(string); ok {
				if v == "" {
					errs.Add("userId", v, "required", "field userId is required")
				}
				if strings.Contains(v, "\x00") {
					errs.Add("userId", v, "safety", "string cannot contain null bytes")
				}
				if !utf8.ValidString(v) {
					errs.Add("userId", v, "safety", "string must be valid UTF-8")
				}
			}
		case "createdAt":
			if _, ok := a.Val.(time.Time); !ok {
				errs.Add("createdAt", a.Val, "type", "field createdAt must be of type time.Time")
			}
		}
	}
	if !provided["userId"] {
		errs.Add("userId", "", "required", "field UserId is required")
	}

	if errs.HasErrors() {
		return *errs
	}
	return nil
}

func assignmentsToProfileCreate(assignments []FieldAssignment) ProfileCreate {
	var input ProfileCreate
	for _, a := range assignments {
		switch a.Col {
		case "id":
			if v, ok := a.Val.(string); ok {
				input.Id = &v
			}
		case "bio":
			if v, ok := a.Val.(string); ok {
				input.Bio = &v
			}
		case "userId":
			if v, ok := a.Val.(string); ok {
				input.UserId = v
			}
		case "createdAt":
			if v, ok := a.Val.(time.Time); ok {
				input.CreatedAt = &v
			}
		}
	}
	return input
}

func (s *ProfileCreate) ToRowMap() map[string]any {
	m := make(map[string]any, 4)
	if s.Id != nil {
		m["id"] = *s.Id
	} else {
		m["id"] = generateCUID()
	}
	if s.Bio != nil {
		m["bio"] = *s.Bio
	}
	m["userId"] = s.UserId
	if s.CreatedAt != nil {
		m["createdAt"] = *s.CreatedAt
	} else {
		m["createdAt"] = time.Now()
	}
	return m
}

func (q *Queries) executeProfileCreate(ctx context.Context, assignments []FieldAssignment, selects *ProfileSelect, omits *ProfileOmit) (*Profile, error) {
	input := assignmentsToProfileCreate(assignments)

	if q.Profile.beforeCreate != nil {
		if err := q.Profile.beforeCreate(ctx, &input); err != nil {
			return nil, err
		}
	}

	if err := validateProfileCreate(assignments); err != nil {
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
	if input.Bio != nil {
		cols = append(cols, "bio")
		vals = append(vals, *input.Bio)
	}
	cols = append(cols, "userId")
	vals = append(vals, input.UserId)
	if input.CreatedAt != nil {
		cols = append(cols, "createdAt")
		vals = append(vals, *input.CreatedAt)
	} else {
		cols = append(cols, "createdAt")
		vals = append(vals, time.Now())
	}

	returningCols := q.selectProfileCols(selects, omits)

	scanFunc := func(res *Profile, cols []string) []any {
		return res.ScanFields(cols)
	}

	idCol := "id"

	hasRelations := selects.hasAnyRelation()

	var res *Profile
	var err error
	if hasRelations {
		err = q.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = executeInsert(ctx, txQ, "Profile", cols, vals, returningCols, idCol, scanFunc)
			if err != nil {
				return err
			}
			return txQ.loadProfileRelations(ctx, []*Profile{res}, selects)
		})
	} else {
		res, err = executeInsert(ctx, q, "Profile", cols, vals, returningCols, idCol, scanFunc)
	}
	if err != nil {
		return nil, err
	}

	if q.Profile.afterCreate != nil {
		if err := q.Profile.afterCreate(ctx, []*Profile{res}); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (d *ProfileDelegate) CreateMany(records ...RecordInput) *CreateManyBuilder[Profile] {
	return &CreateManyBuilder[Profile]{
		client:   d.client,
		records:  records,
		execFunc: d.client.executeProfileCreateMany,
	}
}

func (d *ProfileDelegate) CreateManyAndReturn(records ...RecordInput) *CreateManyAndReturnBuilder[Profile, ProfileSelect, ProfileOmit] {
	return &CreateManyAndReturnBuilder[Profile, ProfileSelect, ProfileOmit]{
		client:   d.client,
		records:  records,
		execFunc: d.client.executeProfileCreateManyAndReturn,
	}
}

func (q *Queries) executeProfileCreateMany(ctx context.Context, records []RecordInput) (int64, error) {
	rowMaps := make([]map[string]any, len(records))
	inputs := make([]ProfileCreate, len(records))
	for i, rec := range records {
		if err := validateProfileCreate(rec.Assignments); err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		input := assignmentsToProfileCreate(rec.Assignments)
		if q.Profile.beforeCreate != nil {
			if err := q.Profile.beforeCreate(ctx, &input); err != nil {
				return 0, err
			}
		}
		rowMaps[i] = input.ToRowMap()
		inputs[i] = input
	}
	count, err := executeCreateMany(ctx, q, rowMaps, "Profile", ProfileColOrder)
	if err != nil {
		return 0, err
	}
	if q.Profile.afterCreateMany != nil {
		if err := q.Profile.afterCreateMany(ctx, inputs, count); err != nil {
			return 0, err
		}
	}
	return count, nil
}

func (q *Queries) executeProfileCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *ProfileSelect, omits *ProfileOmit) ([]*Profile, error) {
	rowMaps := make([]map[string]any, len(records))
	idCol := "id"
	for i, rec := range records {
		if err := validateProfileCreate(rec.Assignments); err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		input := assignmentsToProfileCreate(rec.Assignments)
		if q.Profile.beforeCreate != nil {
			if err := q.Profile.beforeCreate(ctx, &input); err != nil {
				return nil, err
			}
		}
		rowMaps[i] = input.ToRowMap()
	}
	results, err := executeCreateManyAndReturn(ctx, q, rowMaps, "Profile", ProfileColOrder, selects, omits,
		q.selectProfileCols,
		q.loadProfileRelations,
		(*Profile).ScanFields,
		(*ProfileSelect).hasAnyRelation,
		idCol,
	)
	if err != nil {
		return nil, err
	}
	if q.Profile.afterCreate != nil {
		if err := q.Profile.afterCreate(ctx, results); err != nil {
			return nil, err
		}
	}
	return results, nil
}
func (d *ProfileDelegate) FindUnique(where UniquePredicate) *FindUniqueBuilder[Profile, ProfileSelect, ProfileOmit] {
	return &FindUniqueBuilder[Profile, ProfileSelect, ProfileOmit]{
		client:   d.client,
		where:    where,
		execFunc: d.client.executeProfileFindUnique,
	}
}

func (d *ProfileDelegate) FindFirst(preds ...Predicate) *FindFirstBuilder[Profile, ProfileSelect, ProfileOmit] {
	return &FindFirstBuilder[Profile, ProfileSelect, ProfileOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeProfileFindFirst,
	}
}

func (d *ProfileDelegate) FindMany(preds ...Predicate) *FindManyBuilder[Profile, ProfileSelect, ProfileOmit] {
	return &FindManyBuilder[Profile, ProfileSelect, ProfileOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeProfileFindMany,
	}
}

func (q *Queries) executeProfileFindUnique(ctx context.Context, where UniquePredicate, selects *ProfileSelect, omits *ProfileOmit) (*Profile, error) {
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
	returningCols := q.selectProfileCols(selects, omits)
	return executeSingleWithRelations(ctx, q, "Profile", whereClause, vals, returningCols,
		func(res *Profile, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*Profile) error {
			return txQ.loadProfileRelations(ctx, results, selects)
		},
	)
}

func (q *Queries) executeProfileFindFirst(ctx context.Context, where []Predicate, selects *ProfileSelect, omits *ProfileOmit) (*Profile, error) {
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
	returningCols := q.selectProfileCols(selects, omits)
	return executeSingleWithRelations(ctx, q, "Profile", whereClause, vals, returningCols,
		func(res *Profile, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*Profile) error {
			return txQ.loadProfileRelations(ctx, results, selects)
		},
	)
}

func (q *Queries) executeProfileFindMany(ctx context.Context, where []Predicate, selects *ProfileSelect, omits *ProfileOmit) ([]*Profile, error) {
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
	returningCols := q.selectProfileCols(selects, omits)
	return executeManyWithRelations(ctx, q, "Profile", whereClause, vals, returningCols,
		func(res *Profile, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*Profile) error {
			return txQ.loadProfileRelations(ctx, results, selects)
		},
	)
}
func (q *Queries) loadProfileRelations(ctx context.Context, records []*Profile, selects *ProfileSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.User != nil {
		returningCols := q.selectUserCols(selects.User, nil, "id")
		// Current model holds the FK: Profile.userId
		allChildren, err := loadRelation(
			ctx, q, records,
			directKey(func(p *Profile) string { return p.UserId }),
			"User",
			"id",
			returningCols,
			scanInto(returningCols, (*User).ScanFields),
			directKey(func(c *User) string { return c.Id }),
			setOne(func(p *Profile, c *User) { p.User = c }),
		)
		if err != nil {
			return fmt.Errorf("loading user: %w", err)
		}
		if err := q.loadUserRelations(ctx, allChildren, selects.User); err != nil {
			return err
		}
	}

	return nil
}
