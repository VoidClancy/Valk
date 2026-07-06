package valkyrie

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
	client       *Queries
	beforeCreate func(context.Context, *ProfileCreateInput) error
	afterCreate  func(context.Context, *Profile) error
}

func (d *ProfileDelegate) BeforeCreate(hook func(context.Context, *ProfileCreateInput) error) {
	d.beforeCreate = hook
}

func (d *ProfileDelegate) AfterCreate(hook func(context.Context, *Profile) error) {
	d.afterCreate = hook
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
		}
	}
	return targets
}

var profileDefaultCols = []string{
	"id",
	"bio",
	"userId",
}

func (q *Queries) selectProfileCols(selects *ProfileSelect, omits *ProfileOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return profileDefaultCols
	}

	anySelected := selects != nil && (selects.Id || selects.Bio || selects.UserId || selects.User != nil)

	specs := []colSpec{
		{"id", selects != nil && selects.Id, omits != nil && omits.Id, false},
		{"bio", selects != nil && selects.Bio, omits != nil && omits.Bio, false},
		{"userId", selects != nil && selects.UserId, omits != nil && omits.UserId, selects != nil && selects.User != nil},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (input ProfileCreateInput) Validate() error {
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
	if input.UserId == "" {
		errs.Add("userId", input.UserId, "required", "field UserId is required")
	}
	if strings.Contains(input.UserId, "\x00") {
		errs.Add("userId", input.UserId, "safety", "string cannot contain null bytes")
	}
	if !utf8.ValidString(input.UserId) {
		errs.Add("userId", input.UserId, "safety", "string must be valid UTF-8")
	}

	if errs.HasErrors() {
		return *errs
	}
	return nil
}

var ProfileColOrder = []string{
	"id",
	"bio",
	"userId",
}

func (s *ProfileSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.User != nil
}

func (d *ProfileDelegate) Create(input ProfileCreateInput) *CreateBuilder[Profile, ProfileCreateInput, ProfileSelect, ProfileOmit] {
	return &CreateBuilder[Profile, ProfileCreateInput, ProfileSelect, ProfileOmit]{
		client:   d.client,
		input:    input,
		execFunc: d.client.executeProfileCreate,
	}
}

func (q *Queries) executeProfileCreate(ctx context.Context, input ProfileCreateInput, selects *ProfileSelect, omits *ProfileOmit) (*Profile, error) {
	if q.Profile.beforeCreate != nil {
		if err := q.Profile.beforeCreate(ctx, &input); err != nil {
			return nil, err
		}
	}

	if err := input.Validate(); err != nil {
		return nil, err
	}
	m := q.ProfileInputToMap(input)
	cols, vals := mapToColsVals(m, ProfileColOrder)

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
		if err := q.Profile.afterCreate(ctx, res); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (q *Queries) ProfileInputToMap(input ProfileCreateInput) map[string]any {
	m := make(map[string]any)
	if input.Id != nil {
		m["id"] = *input.Id
	} else {
		m["id"] = generateCUID()
	}
	if input.Bio != nil {
		m["bio"] = *input.Bio
	}
	m["userId"] = input.UserId
	return m
}

func (d *ProfileDelegate) CreateMany(inputs []ProfileCreateInput) *CreateManyBuilder[Profile, ProfileCreateInput] {
	return &CreateManyBuilder[Profile, ProfileCreateInput]{
		client:   d.client,
		inputs:   inputs,
		execFunc: d.client.executeProfileCreateMany,
	}
}

func (d *ProfileDelegate) CreateManyAndReturn(inputs []ProfileCreateInput) *CreateManyAndReturnBuilder[Profile, ProfileCreateInput, ProfileSelect, ProfileOmit] {
	return &CreateManyAndReturnBuilder[Profile, ProfileCreateInput, ProfileSelect, ProfileOmit]{
		client:   d.client,
		inputs:   inputs,
		execFunc: d.client.executeProfileCreateManyAndReturn,
	}
}

func (q *Queries) executeProfileCreateMany(ctx context.Context, inputs []ProfileCreateInput) (int64, error) {
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
			rowMaps[i] = q.ProfileInputToMap(input)
		}
		query, vals := buildBulkInsertSQL(q.dialect, "Profile", rowMaps, ProfileColOrder, nil)
		res, err := q.exec(ctx, query, vals...)
		if err != nil {
			return 0, err
		}
		return res.RowsAffected()
	}

	var count int64
	err := q.transaction(ctx, func(txQ *Queries) error {
		for _, input := range inputs {
			_, err := txQ.executeProfileCreate(ctx, input, nil, nil)
			if err != nil {
				return err
			}
			count++
		}
		return nil
	})
	return count, err
}

func (q *Queries) executeProfileCreateManyAndReturn(ctx context.Context, inputs []ProfileCreateInput, selects *ProfileSelect, omits *ProfileOmit) ([]*Profile, error) {
	if len(inputs) == 0 {
		return nil, nil
	}
	for i, input := range inputs {
		if err := input.Validate(); err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
	}

	hasRelations := selects.hasAnyRelation()
	returningCols := q.selectProfileCols(selects, omits)

	if q.dialect.SupportsBulkInsert() {
		rowMaps := make([]map[string]any, len(inputs))
		for i, input := range inputs {
			rowMaps[i] = q.ProfileInputToMap(input)
		}
		query, vals := buildBulkInsertSQL(q.dialect, "Profile", rowMaps, ProfileColOrder, returningCols)
		var records []*Profile
		err := q.transaction(ctx, func(txQ *Queries) error {
			rows, err := txQ.query(ctx, query, vals...)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var record Profile
				if err := rows.Scan(record.ScanFields(returningCols)...); err != nil {
					return err
				}
				records = append(records, &record)
			}
			if err := rows.Err(); err != nil {
				return err
			}
			if hasRelations {
				return txQ.loadProfileRelations(ctx, records, selects)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return records, nil
	}

	// Fallback to loop inside transaction
	var records []*Profile
	err := q.transaction(ctx, func(txQ *Queries) error {
		for _, input := range inputs {
			res, err := txQ.executeProfileCreate(ctx, input, nil, nil)
			if err != nil {
				return err
			}
			records = append(records, res)
		}

		if hasRelations {
			return txQ.loadProfileRelations(ctx, records, selects)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return records, nil
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
