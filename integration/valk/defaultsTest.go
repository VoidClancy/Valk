package valk

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

// DefaultsTest represents the database model
type DefaultsTest struct {
	Uuid4      string    `db:"uuid4" json:"uuid4"`
	Uuid7      string    `db:"uuid7" json:"uuid7"`
	UuidNoArgs string    `db:"uuidNoArgs" json:"uuidNoArgs"`
	Cuid1      string    `db:"cuid1" json:"cuid1"`
	Cuid2      string    `db:"cuid2" json:"cuid2"`
	CuidNoArgs string    `db:"cuidNoArgs" json:"cuidNoArgs"`
	Ulid       string    `db:"ulid" json:"ulid"`
	Nanoid     string    `db:"nanoid" json:"nanoid"`
	Now        time.Time `db:"now" json:"now"`
}

// DefaultsTestCreate is used for hooks only — the Create API uses FieldAssignment
type DefaultsTestCreate struct {
	Uuid4      *string    `json:"uuid4"`
	Uuid7      *string    `json:"uuid7"`
	UuidNoArgs *string    `json:"uuidNoArgs"`
	Cuid1      *string    `json:"cuid1"`
	Cuid2      *string    `json:"cuid2"`
	CuidNoArgs *string    `json:"cuidNoArgs"`
	Ulid       *string    `json:"ulid"`
	Nanoid     *string    `json:"nanoid"`
	Now        *time.Time `json:"now"`
}

// DefaultsTestSelect specifies which fields to include
type DefaultsTestSelect struct {
	Uuid4      bool `json:"uuid4"`
	Uuid7      bool `json:"uuid7"`
	UuidNoArgs bool `json:"uuidNoArgs"`
	Cuid1      bool `json:"cuid1"`
	Cuid2      bool `json:"cuid2"`
	CuidNoArgs bool `json:"cuidNoArgs"`
	Ulid       bool `json:"ulid"`
	Nanoid     bool `json:"nanoid"`
	Now        bool `json:"now"`
}

// DefaultsTestOmit specifies which fields to exclude
type DefaultsTestOmit struct {
	Uuid4      bool `json:"uuid4"`
	Uuid7      bool `json:"uuid7"`
	UuidNoArgs bool `json:"uuidNoArgs"`
	Cuid1      bool `json:"cuid1"`
	Cuid2      bool `json:"cuid2"`
	CuidNoArgs bool `json:"cuidNoArgs"`
	Ulid       bool `json:"ulid"`
	Nanoid     bool `json:"nanoid"`
	Now        bool `json:"now"`
}

type DefaultsTestDelegate struct {
	client          *Queries
	beforeCreate    func(context.Context, *DefaultsTestCreate) error
	afterCreate     func(context.Context, []*DefaultsTest) error
	afterCreateMany func(context.Context, []DefaultsTestCreate, int64) error
}

func (d *DefaultsTestDelegate) BeforeCreate(hook func(context.Context, *DefaultsTestCreate) error) {
	d.beforeCreate = hook
}

func (d *DefaultsTestDelegate) AfterCreate(hook func(context.Context, []*DefaultsTest) error) {
	d.afterCreate = hook
}

func (d *DefaultsTestDelegate) AfterCreateMany(hook func(context.Context, []DefaultsTestCreate, int64) error) {
	d.afterCreateMany = hook
}

func (m *DefaultsTest) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "uuid4":
			targets[i] = &m.Uuid4
		case "uuid7":
			targets[i] = &m.Uuid7
		case "uuidNoArgs":
			targets[i] = &m.UuidNoArgs
		case "cuid1":
			targets[i] = &m.Cuid1
		case "cuid2":
			targets[i] = &m.Cuid2
		case "cuidNoArgs":
			targets[i] = &m.CuidNoArgs
		case "ulid":
			targets[i] = &m.Ulid
		case "nanoid":
			targets[i] = &m.Nanoid
		case "now":
			targets[i] = &m.Now
		}
	}
	return targets
}

var defaultsTestDefaultCols = []string{
	"uuid4",
	"uuid7",
	"uuidNoArgs",
	"cuid1",
	"cuid2",
	"cuidNoArgs",
	"ulid",
	"nanoid",
	"now",
}

func (q *Queries) selectDefaultsTestCols(selects *DefaultsTestSelect, omits *DefaultsTestOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return defaultsTestDefaultCols
	}

	anySelected := selects != nil && (selects.Uuid4 || selects.Uuid7 || selects.UuidNoArgs || selects.Cuid1 || selects.Cuid2 || selects.CuidNoArgs || selects.Ulid || selects.Nanoid || selects.Now)

	specs := []colSpec{
		{"uuid4", selects != nil && selects.Uuid4, omits != nil && omits.Uuid4, selects != nil && selects.hasAnyRelation()},
		{"uuid7", selects != nil && selects.Uuid7, omits != nil && omits.Uuid7, false},
		{"uuidNoArgs", selects != nil && selects.UuidNoArgs, omits != nil && omits.UuidNoArgs, false},
		{"cuid1", selects != nil && selects.Cuid1, omits != nil && omits.Cuid1, false},
		{"cuid2", selects != nil && selects.Cuid2, omits != nil && omits.Cuid2, false},
		{"cuidNoArgs", selects != nil && selects.CuidNoArgs, omits != nil && omits.CuidNoArgs, false},
		{"ulid", selects != nil && selects.Ulid, omits != nil && omits.Ulid, false},
		{"nanoid", selects != nil && selects.Nanoid, omits != nil && omits.Nanoid, false},
		{"now", selects != nil && selects.Now, omits != nil && omits.Now, false},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

var DefaultsTestColOrder = []string{
	"uuid4",
	"uuid7",
	"uuidNoArgs",
	"cuid1",
	"cuid2",
	"cuidNoArgs",
	"ulid",
	"nanoid",
	"now",
}

func (s *DefaultsTestSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return false
}

func (d *DefaultsTestDelegate) Create(assignments ...FieldAssignment) *CreateBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit] {
	return &CreateBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit]{
		client:      d.client,
		assignments: assignments,
		execFunc:    d.client.executeDefaultsTestCreate,
	}
}

func validateDefaultsTestCreate(assignments []FieldAssignment) error {
	errs := &ValidationError{}

	provided := make(map[string]bool)
	for _, a := range assignments {
		provided[a.Col] = true
		switch a.Col {
		case "uuid4":
			if v, ok := a.Val.(string); ok {
				if strings.Contains(v, "\x00") {
					errs.Add("uuid4", v, "safety", "string cannot contain null bytes")
				}
				if !utf8.ValidString(v) {
					errs.Add("uuid4", v, "safety", "string must be valid UTF-8")
				}
			}
		case "uuid7":
			if v, ok := a.Val.(string); ok {
				if strings.Contains(v, "\x00") {
					errs.Add("uuid7", v, "safety", "string cannot contain null bytes")
				}
				if !utf8.ValidString(v) {
					errs.Add("uuid7", v, "safety", "string must be valid UTF-8")
				}
			}
		case "uuidNoArgs":
			if v, ok := a.Val.(string); ok {
				if strings.Contains(v, "\x00") {
					errs.Add("uuidNoArgs", v, "safety", "string cannot contain null bytes")
				}
				if !utf8.ValidString(v) {
					errs.Add("uuidNoArgs", v, "safety", "string must be valid UTF-8")
				}
			}
		case "cuid1":
			if v, ok := a.Val.(string); ok {
				if strings.Contains(v, "\x00") {
					errs.Add("cuid1", v, "safety", "string cannot contain null bytes")
				}
				if !utf8.ValidString(v) {
					errs.Add("cuid1", v, "safety", "string must be valid UTF-8")
				}
			}
		case "cuid2":
			if v, ok := a.Val.(string); ok {
				if strings.Contains(v, "\x00") {
					errs.Add("cuid2", v, "safety", "string cannot contain null bytes")
				}
				if !utf8.ValidString(v) {
					errs.Add("cuid2", v, "safety", "string must be valid UTF-8")
				}
			}
		case "cuidNoArgs":
			if v, ok := a.Val.(string); ok {
				if strings.Contains(v, "\x00") {
					errs.Add("cuidNoArgs", v, "safety", "string cannot contain null bytes")
				}
				if !utf8.ValidString(v) {
					errs.Add("cuidNoArgs", v, "safety", "string must be valid UTF-8")
				}
			}
		case "ulid":
			if v, ok := a.Val.(string); ok {
				if strings.Contains(v, "\x00") {
					errs.Add("ulid", v, "safety", "string cannot contain null bytes")
				}
				if !utf8.ValidString(v) {
					errs.Add("ulid", v, "safety", "string must be valid UTF-8")
				}
			}
		case "nanoid":
			if v, ok := a.Val.(string); ok {
				if strings.Contains(v, "\x00") {
					errs.Add("nanoid", v, "safety", "string cannot contain null bytes")
				}
				if !utf8.ValidString(v) {
					errs.Add("nanoid", v, "safety", "string must be valid UTF-8")
				}
			}
		case "now":
			if _, ok := a.Val.(time.Time); !ok {
				errs.Add("now", a.Val, "type", "field now must be of type time.Time")
			}
		}
	}

	if errs.HasErrors() {
		return *errs
	}
	return nil
}

func assignmentsToDefaultsTestCreate(assignments []FieldAssignment) DefaultsTestCreate {
	var input DefaultsTestCreate
	for _, a := range assignments {
		switch a.Col {
		case "uuid4":
			if v, ok := a.Val.(string); ok {
				input.Uuid4 = &v
			}
		case "uuid7":
			if v, ok := a.Val.(string); ok {
				input.Uuid7 = &v
			}
		case "uuidNoArgs":
			if v, ok := a.Val.(string); ok {
				input.UuidNoArgs = &v
			}
		case "cuid1":
			if v, ok := a.Val.(string); ok {
				input.Cuid1 = &v
			}
		case "cuid2":
			if v, ok := a.Val.(string); ok {
				input.Cuid2 = &v
			}
		case "cuidNoArgs":
			if v, ok := a.Val.(string); ok {
				input.CuidNoArgs = &v
			}
		case "ulid":
			if v, ok := a.Val.(string); ok {
				input.Ulid = &v
			}
		case "nanoid":
			if v, ok := a.Val.(string); ok {
				input.Nanoid = &v
			}
		case "now":
			if v, ok := a.Val.(time.Time); ok {
				input.Now = &v
			}
		}
	}
	return input
}

func (s *DefaultsTestCreate) ToRowMap() map[string]any {
	m := make(map[string]any, 9)
	if s.Uuid4 != nil {
		m["uuid4"] = *s.Uuid4
	} else {
		m["uuid4"] = generateUUID()
	}
	if s.Uuid7 != nil {
		m["uuid7"] = *s.Uuid7
	} else {
		m["uuid7"] = generateUUID7()
	}
	if s.UuidNoArgs != nil {
		m["uuidNoArgs"] = *s.UuidNoArgs
	} else {
		m["uuidNoArgs"] = generateUUID()
	}
	if s.Cuid1 != nil {
		m["cuid1"] = *s.Cuid1
	} else {
		m["cuid1"] = generateCUID()
	}
	if s.Cuid2 != nil {
		m["cuid2"] = *s.Cuid2
	} else {
		m["cuid2"] = generateCUID2()
	}
	if s.CuidNoArgs != nil {
		m["cuidNoArgs"] = *s.CuidNoArgs
	} else {
		m["cuidNoArgs"] = generateCUID()
	}
	if s.Ulid != nil {
		m["ulid"] = *s.Ulid
	} else {
		m["ulid"] = generateULID()
	}
	if s.Nanoid != nil {
		m["nanoid"] = *s.Nanoid
	} else {
		m["nanoid"] = generateNanoID()
	}
	if s.Now != nil {
		m["now"] = *s.Now
	} else {
		m["now"] = time.Now()
	}
	return m
}

func (q *Queries) executeDefaultsTestCreate(ctx context.Context, assignments []FieldAssignment, selects *DefaultsTestSelect, omits *DefaultsTestOmit) (*DefaultsTest, error) {
	input := assignmentsToDefaultsTestCreate(assignments)

	if q.DefaultsTest.beforeCreate != nil {
		if err := q.DefaultsTest.beforeCreate(ctx, &input); err != nil {
			return nil, err
		}
	}

	if err := validateDefaultsTestCreate(assignments); err != nil {
		return nil, err
	}

	var cols []string
	var vals []any
	if input.Uuid4 != nil {
		cols = append(cols, "uuid4")
		vals = append(vals, *input.Uuid4)
	} else {
		cols = append(cols, "uuid4")
		vals = append(vals, generateUUID())
	}
	if input.Uuid7 != nil {
		cols = append(cols, "uuid7")
		vals = append(vals, *input.Uuid7)
	} else {
		cols = append(cols, "uuid7")
		vals = append(vals, generateUUID7())
	}
	if input.UuidNoArgs != nil {
		cols = append(cols, "uuidNoArgs")
		vals = append(vals, *input.UuidNoArgs)
	} else {
		cols = append(cols, "uuidNoArgs")
		vals = append(vals, generateUUID())
	}
	if input.Cuid1 != nil {
		cols = append(cols, "cuid1")
		vals = append(vals, *input.Cuid1)
	} else {
		cols = append(cols, "cuid1")
		vals = append(vals, generateCUID())
	}
	if input.Cuid2 != nil {
		cols = append(cols, "cuid2")
		vals = append(vals, *input.Cuid2)
	} else {
		cols = append(cols, "cuid2")
		vals = append(vals, generateCUID2())
	}
	if input.CuidNoArgs != nil {
		cols = append(cols, "cuidNoArgs")
		vals = append(vals, *input.CuidNoArgs)
	} else {
		cols = append(cols, "cuidNoArgs")
		vals = append(vals, generateCUID())
	}
	if input.Ulid != nil {
		cols = append(cols, "ulid")
		vals = append(vals, *input.Ulid)
	} else {
		cols = append(cols, "ulid")
		vals = append(vals, generateULID())
	}
	if input.Nanoid != nil {
		cols = append(cols, "nanoid")
		vals = append(vals, *input.Nanoid)
	} else {
		cols = append(cols, "nanoid")
		vals = append(vals, generateNanoID())
	}
	if input.Now != nil {
		cols = append(cols, "now")
		vals = append(vals, *input.Now)
	} else {
		cols = append(cols, "now")
		vals = append(vals, time.Now())
	}

	returningCols := q.selectDefaultsTestCols(selects, omits)

	scanFunc := func(res *DefaultsTest, cols []string) []any {
		return res.ScanFields(cols)
	}

	idCol := "uuid4"

	hasRelations := selects.hasAnyRelation()

	var res *DefaultsTest
	var err error
	if hasRelations {
		err = q.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = executeInsert(ctx, txQ, "DefaultsTest", cols, vals, returningCols, idCol, scanFunc)
			if err != nil {
				return err
			}
			return txQ.loadDefaultsTestRelations(ctx, []*DefaultsTest{res}, selects)
		})
	} else {
		res, err = executeInsert(ctx, q, "DefaultsTest", cols, vals, returningCols, idCol, scanFunc)
	}
	if err != nil {
		return nil, err
	}

	if q.DefaultsTest.afterCreate != nil {
		if err := q.DefaultsTest.afterCreate(ctx, []*DefaultsTest{res}); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (d *DefaultsTestDelegate) CreateMany(records ...RecordInput) *CreateManyBuilder[DefaultsTest] {
	return &CreateManyBuilder[DefaultsTest]{
		client:   d.client,
		records:  records,
		execFunc: d.client.executeDefaultsTestCreateMany,
	}
}

func (d *DefaultsTestDelegate) CreateManyAndReturn(records ...RecordInput) *CreateManyAndReturnBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit] {
	return &CreateManyAndReturnBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit]{
		client:   d.client,
		records:  records,
		execFunc: d.client.executeDefaultsTestCreateManyAndReturn,
	}
}

func (q *Queries) executeDefaultsTestCreateMany(ctx context.Context, records []RecordInput) (int64, error) {
	rowMaps := make([]map[string]any, len(records))
	inputs := make([]DefaultsTestCreate, len(records))
	for i, rec := range records {
		if err := validateDefaultsTestCreate(rec.Assignments); err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		input := assignmentsToDefaultsTestCreate(rec.Assignments)
		if q.DefaultsTest.beforeCreate != nil {
			if err := q.DefaultsTest.beforeCreate(ctx, &input); err != nil {
				return 0, err
			}
		}
		rowMaps[i] = input.ToRowMap()
		inputs[i] = input
	}
	count, err := executeCreateMany(ctx, q, rowMaps, "DefaultsTest", DefaultsTestColOrder)
	if err != nil {
		return 0, err
	}
	if q.DefaultsTest.afterCreateMany != nil {
		if err := q.DefaultsTest.afterCreateMany(ctx, inputs, count); err != nil {
			return 0, err
		}
	}
	return count, nil
}

func (q *Queries) executeDefaultsTestCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *DefaultsTestSelect, omits *DefaultsTestOmit) ([]*DefaultsTest, error) {
	rowMaps := make([]map[string]any, len(records))
	idCol := "uuid4"
	for i, rec := range records {
		if err := validateDefaultsTestCreate(rec.Assignments); err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		input := assignmentsToDefaultsTestCreate(rec.Assignments)
		if q.DefaultsTest.beforeCreate != nil {
			if err := q.DefaultsTest.beforeCreate(ctx, &input); err != nil {
				return nil, err
			}
		}
		rowMaps[i] = input.ToRowMap()
	}
	results, err := executeCreateManyAndReturn(ctx, q, rowMaps, "DefaultsTest", DefaultsTestColOrder, selects, omits,
		q.selectDefaultsTestCols,
		q.loadDefaultsTestRelations,
		(*DefaultsTest).ScanFields,
		(*DefaultsTestSelect).hasAnyRelation,
		idCol,
	)
	if err != nil {
		return nil, err
	}
	if q.DefaultsTest.afterCreate != nil {
		if err := q.DefaultsTest.afterCreate(ctx, results); err != nil {
			return nil, err
		}
	}
	return results, nil
}
func (d *DefaultsTestDelegate) FindUnique(where UniquePredicate) *FindUniqueBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit] {
	return &FindUniqueBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit]{
		client:   d.client,
		where:    where,
		execFunc: d.client.executeDefaultsTestFindUnique,
	}
}

func (d *DefaultsTestDelegate) FindFirst(preds ...Predicate) *FindFirstBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit] {
	return &FindFirstBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeDefaultsTestFindFirst,
	}
}

func (d *DefaultsTestDelegate) FindMany(preds ...Predicate) *FindManyBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit] {
	return &FindManyBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeDefaultsTestFindMany,
	}
}

func (q *Queries) executeDefaultsTestFindUnique(ctx context.Context, where UniquePredicate, selects *DefaultsTestSelect, omits *DefaultsTestOmit) (*DefaultsTest, error) {
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
	returningCols := q.selectDefaultsTestCols(selects, omits)
	return executeSingleWithRelations(ctx, q, "DefaultsTest", whereClause, vals, returningCols,
		func(res *DefaultsTest, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*DefaultsTest) error {
			return txQ.loadDefaultsTestRelations(ctx, results, selects)
		},
	)
}

func (q *Queries) executeDefaultsTestFindFirst(ctx context.Context, where []Predicate, selects *DefaultsTestSelect, omits *DefaultsTestOmit) (*DefaultsTest, error) {
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
	returningCols := q.selectDefaultsTestCols(selects, omits)
	return executeSingleWithRelations(ctx, q, "DefaultsTest", whereClause, vals, returningCols,
		func(res *DefaultsTest, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*DefaultsTest) error {
			return txQ.loadDefaultsTestRelations(ctx, results, selects)
		},
	)
}

func (q *Queries) executeDefaultsTestFindMany(ctx context.Context, where []Predicate, selects *DefaultsTestSelect, omits *DefaultsTestOmit) ([]*DefaultsTest, error) {
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
	returningCols := q.selectDefaultsTestCols(selects, omits)
	return executeManyWithRelations(ctx, q, "DefaultsTest", whereClause, vals, returningCols,
		func(res *DefaultsTest, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*DefaultsTest) error {
			return txQ.loadDefaultsTestRelations(ctx, results, selects)
		},
	)
}
func (q *Queries) loadDefaultsTestRelations(ctx context.Context, records []*DefaultsTest, selects *DefaultsTestSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}

	return nil
}
