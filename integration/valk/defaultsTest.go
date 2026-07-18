package valk

import (
	"context"
	"fmt"
	"slices"
	"time"
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

type DefaultsTestSelectQuery interface {
	GetRelationParams() (*DefaultsTestSelect, *DefaultsTestOmit, QueryParams[DefaultsTest])
}

func (s *DefaultsTestSelect) GetRelationParams() (*DefaultsTestSelect, *DefaultsTestOmit, QueryParams[DefaultsTest]) {
	return s, nil, QueryParams[DefaultsTest]{}
}

// DefaultsTestQueryBuilder builds a query for the relation DefaultsTest
type DefaultsTestQueryBuilder struct {
	selects *DefaultsTestSelect
	omits   *DefaultsTestOmit
	where   []PredicateOf[DefaultsTest]
	take    *int
	skip    *int
	orderBy []OrderBy[DefaultsTest]
}

func (b *DefaultsTestQueryBuilder) Where(preds ...PredicateOf[DefaultsTest]) *DefaultsTestQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *DefaultsTestQueryBuilder) Take(limit int) *DefaultsTestQueryBuilder {
	b.take = &limit
	return b
}

func (b *DefaultsTestQueryBuilder) Skip(offset int) *DefaultsTestQueryBuilder {
	b.skip = &offset
	return b
}

func (b *DefaultsTestQueryBuilder) OrderBy(orders ...OrderBy[DefaultsTest]) *DefaultsTestQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *DefaultsTestQueryBuilder) Select(s DefaultsTestSelect) *DefaultsTestQueryBuilder {
	b.selects = &s
	return b
}

func (b *DefaultsTestQueryBuilder) Omit(o DefaultsTestOmit) *DefaultsTestQueryBuilder {
	b.omits = &o
	return b
}

func (b *DefaultsTestQueryBuilder) GetRelationParams() (*DefaultsTestSelect, *DefaultsTestOmit, QueryParams[DefaultsTest]) {
	if b == nil {
		return nil, nil, QueryParams[DefaultsTest]{}
	}
	return b.selects, b.omits, QueryParams[DefaultsTest]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
	}
}

type DefaultsTestCreateQuery = func(ctx context.Context, args *DefaultsTestCreate) (*DefaultsTest, error)
type DefaultsTestCreateManyQuery = func(ctx context.Context, args []*DefaultsTestCreate) (int64, error)
type DefaultsTestCreateManyAndReturnQuery = func(ctx context.Context, args []*DefaultsTestCreate) ([]*DefaultsTest, error)
type DefaultsTestFindUniqueQuery = func(ctx context.Context, where UniquePredicate[DefaultsTest], additional []PredicateOf[DefaultsTest], selects *DefaultsTestSelect, omits *DefaultsTestOmit) (*DefaultsTest, error)
type DefaultsTestFindFirstQuery = func(ctx context.Context, params QueryParams[DefaultsTest], selects *DefaultsTestSelect, omits *DefaultsTestOmit) (*DefaultsTest, error)
type DefaultsTestFindManyQuery = func(ctx context.Context, params QueryParams[DefaultsTest], selects *DefaultsTestSelect, omits *DefaultsTestOmit) ([]*DefaultsTest, error)

type DefaultsTestExtension struct {
	Create              func(ctx context.Context, input *DefaultsTestCreate, next DefaultsTestCreateQuery) (*DefaultsTest, error)
	CreateMany          func(ctx context.Context, inputs []*DefaultsTestCreate, next DefaultsTestCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, inputs []*DefaultsTestCreate, next DefaultsTestCreateManyAndReturnQuery) ([]*DefaultsTest, error)
	FindUnique          func(ctx context.Context, where UniquePredicate[DefaultsTest], additional []PredicateOf[DefaultsTest], selects *DefaultsTestSelect, omits *DefaultsTestOmit, next DefaultsTestFindUniqueQuery) (*DefaultsTest, error)
	FindFirst           func(ctx context.Context, params QueryParams[DefaultsTest], selects *DefaultsTestSelect, omits *DefaultsTestOmit, next DefaultsTestFindFirstQuery) (*DefaultsTest, error)
	FindMany            func(ctx context.Context, params QueryParams[DefaultsTest], selects *DefaultsTestSelect, omits *DefaultsTestOmit, next DefaultsTestFindManyQuery) ([]*DefaultsTest, error)
}

type DefaultsTestDelegate struct {
	client     *Queries
	extensions []DefaultsTestExtension
}

func (d *DefaultsTestDelegate) Use(exts ...DefaultsTestExtension) {
	d.extensions = append(d.extensions, exts...)
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

func selectDefaultsTestCols(selects *DefaultsTestSelect, omits *DefaultsTestOmit, forceCols ...string) []string {
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

func (s *DefaultsTestSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return false
}

type DefaultsTestCreateBuilder struct {
	*CreateBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit]
}

func (b *DefaultsTestCreateBuilder) OnConflict(target UniqueConstraintTarget) *DefaultsTestConflictBuilder[DefaultsTestCreateBuilder] {
	return &DefaultsTestConflictBuilder[DefaultsTestCreateBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *DefaultsTestCreateBuilder) SetUuid4(v string) *DefaultsTestCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "uuid4", Val: v})
	return b
}
func (b *DefaultsTestCreateBuilder) SetUuid7(v string) *DefaultsTestCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "uuid7", Val: v})
	return b
}
func (b *DefaultsTestCreateBuilder) SetUuidNoArgs(v string) *DefaultsTestCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "uuidNoArgs", Val: v})
	return b
}
func (b *DefaultsTestCreateBuilder) SetCuid1(v string) *DefaultsTestCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "cuid1", Val: v})
	return b
}
func (b *DefaultsTestCreateBuilder) SetCuid2(v string) *DefaultsTestCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "cuid2", Val: v})
	return b
}
func (b *DefaultsTestCreateBuilder) SetCuidNoArgs(v string) *DefaultsTestCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "cuidNoArgs", Val: v})
	return b
}
func (b *DefaultsTestCreateBuilder) SetUlid(v string) *DefaultsTestCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "ulid", Val: v})
	return b
}
func (b *DefaultsTestCreateBuilder) SetNanoid(v string) *DefaultsTestCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "nanoid", Val: v})
	return b
}
func (b *DefaultsTestCreateBuilder) SetNow(v time.Time) *DefaultsTestCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "now", Val: v})
	return b
}

func (d *DefaultsTestDelegate) Create(assignments ...FieldAssignment) *DefaultsTestCreateBuilder {
	return &DefaultsTestCreateBuilder{
		CreateBuilder: &CreateBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit]{
			assignments: assignments,
			execFunc:    d.executeCreate,
		},
	}
}

const (
	providedDefaultsTestUuid4      uint64 = 1 << 0
	providedDefaultsTestUuid7      uint64 = 1 << 1
	providedDefaultsTestUuidNoArgs uint64 = 1 << 2
	providedDefaultsTestCuid1      uint64 = 1 << 3
	providedDefaultsTestCuid2      uint64 = 1 << 4
	providedDefaultsTestCuidNoArgs uint64 = 1 << 5
	providedDefaultsTestUlid       uint64 = 1 << 6
	providedDefaultsTestNanoid     uint64 = 1 << 7
	providedDefaultsTestNow        uint64 = 1 << 8
)

func assignmentsToDefaultsTestCreate(assignments []FieldAssignment) (DefaultsTestCreate, error) {
	var input DefaultsTestCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "uuid4":
			provided |= providedDefaultsTestUuid4
			if v, ok := a.Val.(string); ok {
				input.Uuid4 = &v
				ValidateString(&errs, "uuid4", v, false, 0, false, false)
			} else {
				errs.Add("uuid4", a.Val, "type", "field uuid4 must be of type string")
			}
		case "uuid7":
			provided |= providedDefaultsTestUuid7
			if v, ok := a.Val.(string); ok {
				input.Uuid7 = &v
				ValidateString(&errs, "uuid7", v, false, 0, false, false)
			} else {
				errs.Add("uuid7", a.Val, "type", "field uuid7 must be of type string")
			}
		case "uuidNoArgs":
			provided |= providedDefaultsTestUuidNoArgs
			if v, ok := a.Val.(string); ok {
				input.UuidNoArgs = &v
				ValidateString(&errs, "uuidNoArgs", v, false, 0, false, false)
			} else {
				errs.Add("uuidNoArgs", a.Val, "type", "field uuidNoArgs must be of type string")
			}
		case "cuid1":
			provided |= providedDefaultsTestCuid1
			if v, ok := a.Val.(string); ok {
				input.Cuid1 = &v
				ValidateString(&errs, "cuid1", v, false, 0, false, false)
			} else {
				errs.Add("cuid1", a.Val, "type", "field cuid1 must be of type string")
			}
		case "cuid2":
			provided |= providedDefaultsTestCuid2
			if v, ok := a.Val.(string); ok {
				input.Cuid2 = &v
				ValidateString(&errs, "cuid2", v, false, 0, false, false)
			} else {
				errs.Add("cuid2", a.Val, "type", "field cuid2 must be of type string")
			}
		case "cuidNoArgs":
			provided |= providedDefaultsTestCuidNoArgs
			if v, ok := a.Val.(string); ok {
				input.CuidNoArgs = &v
				ValidateString(&errs, "cuidNoArgs", v, false, 0, false, false)
			} else {
				errs.Add("cuidNoArgs", a.Val, "type", "field cuidNoArgs must be of type string")
			}
		case "ulid":
			provided |= providedDefaultsTestUlid
			if v, ok := a.Val.(string); ok {
				input.Ulid = &v
				ValidateString(&errs, "ulid", v, false, 0, false, false)
			} else {
				errs.Add("ulid", a.Val, "type", "field ulid must be of type string")
			}
		case "nanoid":
			provided |= providedDefaultsTestNanoid
			if v, ok := a.Val.(string); ok {
				input.Nanoid = &v
				ValidateString(&errs, "nanoid", v, false, 0, false, false)
			} else {
				errs.Add("nanoid", a.Val, "type", "field nanoid must be of type string")
			}
		case "now":
			provided |= providedDefaultsTestNow
			if v, ok := a.Val.(time.Time); ok {
				input.Now = &v
			} else {
				errs.Add("now", a.Val, "type", "field now must be of type time.Time")
			}
		}
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *DefaultsTestCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 9)
	vals = make([]any, 0, 9)
	cols = append(cols, "uuid4")
	if s.Uuid4 != nil {
		vals = append(vals, *s.Uuid4)
	} else {
		vals = append(vals, generateUUID())
	}
	cols = append(cols, "uuid7")
	if s.Uuid7 != nil {
		vals = append(vals, *s.Uuid7)
	} else {
		vals = append(vals, generateUUID7())
	}
	cols = append(cols, "uuidNoArgs")
	if s.UuidNoArgs != nil {
		vals = append(vals, *s.UuidNoArgs)
	} else {
		vals = append(vals, generateUUID())
	}
	cols = append(cols, "cuid1")
	if s.Cuid1 != nil {
		vals = append(vals, *s.Cuid1)
	} else {
		vals = append(vals, generateCUID())
	}
	cols = append(cols, "cuid2")
	if s.Cuid2 != nil {
		vals = append(vals, *s.Cuid2)
	} else {
		vals = append(vals, generateCUID2())
	}
	cols = append(cols, "cuidNoArgs")
	if s.CuidNoArgs != nil {
		vals = append(vals, *s.CuidNoArgs)
	} else {
		vals = append(vals, generateCUID())
	}
	cols = append(cols, "ulid")
	if s.Ulid != nil {
		vals = append(vals, *s.Ulid)
	} else {
		vals = append(vals, generateULID())
	}
	cols = append(cols, "nanoid")
	if s.Nanoid != nil {
		vals = append(vals, *s.Nanoid)
	} else {
		vals = append(vals, generateNanoID())
	}
	cols = append(cols, "now")
	if s.Now != nil {
		vals = append(vals, *s.Now)
	} else {
		vals = append(vals, time.Now())
	}
	return
}

func (s *DefaultsTestCreate) ToRowMap() map[string]any {
	cols, vals := s.ToColsVals()
	m := make(map[string]any, len(cols))
	for i, c := range cols {
		m[c] = vals[i]
	}
	return m
}

func (d *DefaultsTestDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *DefaultsTestSelect, omits *DefaultsTestOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*DefaultsTest, error) {
	input, err := assignmentsToDefaultsTestCreate(assignments)
	if err != nil {
		return nil, err
	}

	curr := func(c context.Context, args *DefaultsTestCreate) (*DefaultsTest, error) {
		cols, vals := args.ToColsVals()

		returningCols := selectDefaultsTestCols(selects, omits)

		scanFunc := func(res *DefaultsTest, cols []string) []any {
			return res.ScanFields(cols)
		}

		pkCols := []string{
			"uuid4",
		}

		hasRelations := selects.hasAnyRelation()

		var res *DefaultsTest
		var err error
		if hasRelations {
			err = d.client.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = executeInsert(c, txQ, "DefaultsTest", cols, vals, returningCols, pkCols, scanFunc, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.DefaultsTest.loadRelations(c, []*DefaultsTest{res}, selects)
			})
		} else {
			res, err = executeInsert(c, d.client, "DefaultsTest", cols, vals, returningCols, pkCols, scanFunc, conflictTarget, conflictAction)
		}
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Create != nil {
			next, hook := curr, ext.Create
			curr = func(c context.Context, input *DefaultsTestCreate) (*DefaultsTest, error) {
				return hook(c, input, next)
			}
		}
	}

	return curr(ctx, &input)
}

type DefaultsTestCreateManyBuilder struct {
	*CreateManyBuilder[DefaultsTest]
}

func (b *DefaultsTestCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *DefaultsTestConflictBuilder[DefaultsTestCreateManyBuilder] {
	return &DefaultsTestConflictBuilder[DefaultsTestCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type DefaultsTestCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit]
}

func (b *DefaultsTestCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *DefaultsTestConflictBuilder[DefaultsTestCreateManyAndReturnBuilder] {
	return &DefaultsTestConflictBuilder[DefaultsTestCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (d *DefaultsTestDelegate) CreateMany(builders ...*DefaultsTestCreateBuilder) *DefaultsTestCreateManyBuilder {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &DefaultsTestCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[DefaultsTest]{
			records:  records,
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *DefaultsTestDelegate) CreateManyAndReturn(builders ...*DefaultsTestCreateBuilder) *DefaultsTestCreateManyAndReturnBuilder {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &DefaultsTestCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit]{
			records:  records,
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func (d *DefaultsTestDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs := make([]*DefaultsTestCreate, len(records))
	for i, rec := range records {
		input, err := assignmentsToDefaultsTestCreate(rec.Assignments)
		if err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &input
	}

	curr := func(c context.Context, args []*DefaultsTestCreate) (int64, error) {
		rowMaps := make([]map[string]any, len(args))
		for i, input := range args {
			rowMaps[i] = input.ToRowMap()
		}

		pkCols := []string{
			"uuid4",
		}

		return executeCreateMany(c, d.client, rowMaps, "DefaultsTest", defaultsTestDefaultCols, pkCols, conflictTarget, conflictAction)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.CreateMany != nil {
			next, hook := curr, ext.CreateMany
			curr = func(c context.Context, inputs []*DefaultsTestCreate) (int64, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
}

func (d *DefaultsTestDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *DefaultsTestSelect, omits *DefaultsTestOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*DefaultsTest, error) {
	inputs := make([]*DefaultsTestCreate, len(records))
	for i, rec := range records {
		input, err := assignmentsToDefaultsTestCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &input
	}

	curr := func(c context.Context, args []*DefaultsTestCreate) ([]*DefaultsTest, error) {
		rowMaps := make([]map[string]any, len(args))
		for i, input := range args {
			rowMaps[i] = input.ToRowMap()
		}

		pkCols := []string{
			"uuid4",
		}

		return executeCreateManyAndReturn(c, d.client, rowMaps, "DefaultsTest", defaultsTestDefaultCols, selects, omits,
			selectDefaultsTestCols,
			func(ctx context.Context, txQ *Queries, results []*DefaultsTest, sel *DefaultsTestSelect) error {
				return txQ.DefaultsTest.loadRelations(ctx, results, sel)
			},
			(*DefaultsTest).ScanFields,
			(*DefaultsTestSelect).hasAnyRelation,
			pkCols,
			conflictTarget,
			conflictAction,
		)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.CreateManyAndReturn != nil {
			next, hook := curr, ext.CreateManyAndReturn
			curr = func(c context.Context, inputs []*DefaultsTestCreate) ([]*DefaultsTest, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
}

type DefaultsTestConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *DefaultsTestConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *DefaultsTestConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *DefaultsTestConflictBuilder[B]) Update(fn func(u *DefaultsTestUpsert)) *B {
	var up ConflictUpdate
	u := newDefaultsTestUpsert(&up)
	fn(u)
	cb.setAction(ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}, cb.conflictTarget)
	return cb.builder
}

type DefaultsTestUpsert struct {
	Uuid4      fieldUpsert[string]
	Uuid7      fieldUpsert[string]
	UuidNoArgs fieldUpsert[string]
	Cuid1      fieldUpsert[string]
	Cuid2      fieldUpsert[string]
	CuidNoArgs fieldUpsert[string]
	Ulid       fieldUpsert[string]
	Nanoid     fieldUpsert[string]
	Now        fieldUpsert[time.Time]
}

func newDefaultsTestUpsert(up *ConflictUpdate) *DefaultsTestUpsert {
	return &DefaultsTestUpsert{
		Uuid4:      fieldUpsert[string]{column: "uuid4", update: up},
		Uuid7:      fieldUpsert[string]{column: "uuid7", update: up},
		UuidNoArgs: fieldUpsert[string]{column: "uuidNoArgs", update: up},
		Cuid1:      fieldUpsert[string]{column: "cuid1", update: up},
		Cuid2:      fieldUpsert[string]{column: "cuid2", update: up},
		CuidNoArgs: fieldUpsert[string]{column: "cuidNoArgs", update: up},
		Ulid:       fieldUpsert[string]{column: "ulid", update: up},
		Nanoid:     fieldUpsert[string]{column: "nanoid", update: up},
		Now:        fieldUpsert[time.Time]{column: "now", update: up},
	}
}
func (d *DefaultsTestDelegate) FindUnique(where UniquePredicate[DefaultsTest], additional ...PredicateOf[DefaultsTest]) *FindUniqueBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit] {
	return &FindUniqueBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *DefaultsTestDelegate) FindFirst(preds ...PredicateOf[DefaultsTest]) *FindFirstBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit] {
	return &FindFirstBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *DefaultsTestDelegate) FindMany(preds ...PredicateOf[DefaultsTest]) *FindManyBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit] {
	return &FindManyBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *DefaultsTestDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[DefaultsTest], additional []PredicateOf[DefaultsTest], selects *DefaultsTestSelect, omits *DefaultsTestOmit) (*DefaultsTest, error) {
	curr := func(c context.Context, w UniquePredicate[DefaultsTest], add []PredicateOf[DefaultsTest], sel *DefaultsTestSelect, o *DefaultsTestOmit) (*DefaultsTest, error) {
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
		allPreds := append([]PredicateOf[DefaultsTest]{w}, add...)
		whereClause, vals := CompilePredicates(d.client.dialect, allPreds)
		if whereClause != "" {
			whereClause = " WHERE " + whereClause
		}
		returningCols := selectDefaultsTestCols(sel, o)
		return executeSingleWithRelations(c, d.client, "DefaultsTest", whereClause, vals, returningCols,
			func(res *DefaultsTest, cols []string) []any { return res.ScanFields(cols) },
			sel.hasAnyRelation(),
			func(ctx context.Context, txQ *Queries, results []*DefaultsTest) error {
				return txQ.DefaultsTest.loadRelations(ctx, results, sel)
			},
			nil,
		)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindUnique != nil {
			next, hook := curr, ext.FindUnique
			curr = func(c context.Context, w UniquePredicate[DefaultsTest], add []PredicateOf[DefaultsTest], sel *DefaultsTestSelect, o *DefaultsTestOmit) (*DefaultsTest, error) {
				return hook(c, w, add, sel, o, next)
			}
		}
	}

	return curr(ctx, where, additional, selects, omits)
}

func (d *DefaultsTestDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[DefaultsTest],
	selects *DefaultsTestSelect,
	omits *DefaultsTestOmit,
) (*DefaultsTest, error) {
	curr := func(c context.Context, p QueryParams[DefaultsTest], sel *DefaultsTestSelect, o *DefaultsTestOmit) (*DefaultsTest, error) {
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
		returningCols := selectDefaultsTestCols(sel, o)
		return executeSingleWithRelations(c, d.client, "DefaultsTest", whereClause, vals, returningCols,
			func(res *DefaultsTest, cols []string) []any { return res.ScanFields(cols) },
			sel.hasAnyRelation(),
			func(ctx context.Context, txQ *Queries, results []*DefaultsTest) error {
				return txQ.DefaultsTest.loadRelations(ctx, results, sel)
			},
			p.Skip,
		)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindFirst != nil {
			next, hook := curr, ext.FindFirst
			curr = func(c context.Context, p QueryParams[DefaultsTest], sel *DefaultsTestSelect, o *DefaultsTestOmit) (*DefaultsTest, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}

func (d *DefaultsTestDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[DefaultsTest],
	selects *DefaultsTestSelect,
	omits *DefaultsTestOmit,
) ([]*DefaultsTest, error) {
	curr := func(c context.Context, p QueryParams[DefaultsTest], sel *DefaultsTestSelect, o *DefaultsTestOmit) ([]*DefaultsTest, error) {
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
		returningCols := selectDefaultsTestCols(sel, o)
		return executeManyWithRelations(c, d.client, "DefaultsTest", whereClause, vals, returningCols,
			func(res *DefaultsTest, cols []string) []any { return res.ScanFields(cols) },
			sel.hasAnyRelation(),
			func(ctx context.Context, txQ *Queries, results []*DefaultsTest) error {
				return txQ.DefaultsTest.loadRelations(ctx, results, sel)
			},
			p.Take,
			p.Skip,
		)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindMany != nil {
			next, hook := curr, ext.FindMany
			curr = func(c context.Context, p QueryParams[DefaultsTest], sel *DefaultsTestSelect, o *DefaultsTestOmit) ([]*DefaultsTest, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}
func (d *DefaultsTestDelegate) loadRelations(ctx context.Context, records []*DefaultsTest, selects *DefaultsTestSelect) error {
	_ = ctx
	if selects == nil || len(records) == 0 {
		return nil
	}

	return nil
}
