package valk

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
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

// colMask returns a bit mask of columns that are set
func (s *DefaultsTestCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	mask |= 1 << 1
	mask |= 1 << 2
	mask |= 1 << 3
	mask |= 1 << 4
	mask |= 1 << 5
	mask |= 1 << 6
	mask |= 1 << 7
	mask |= 1 << 8
	return mask
}

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

type DefaultsTestQueryBuilder struct {
	selects *DefaultsTestSelect
	omits   *DefaultsTestOmit
	where   []PredicateOf[DefaultsTest]
	take    *int
	skip    *int
	orderBy []OrderBy[DefaultsTest]
	cursor  UniquePredicate[DefaultsTest]
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

func (b *DefaultsTestQueryBuilder) Cursor(where UniquePredicate[DefaultsTest]) *DefaultsTestQueryBuilder {
	b.cursor = where
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
		Cursor:  b.cursor,
	}
}

type DefaultsTestCreateQuery = func(ctx context.Context, args *DefaultsTestCreate) (*DefaultsTest, error)
type DefaultsTestCreateManyQuery = func(ctx context.Context, args []*DefaultsTestCreate) (int64, error)
type DefaultsTestCreateManyAndReturnQuery = func(ctx context.Context, args []*DefaultsTestCreate) ([]*DefaultsTest, error)
type DefaultsTestFindUniqueQuery = func(ctx context.Context, where UniquePredicate[DefaultsTest], additional []PredicateOf[DefaultsTest], selects *DefaultsTestSelect, omits *DefaultsTestOmit) (*DefaultsTest, error)
type DefaultsTestFindFirstQuery = func(ctx context.Context, params QueryParams[DefaultsTest], selects *DefaultsTestSelect, omits *DefaultsTestOmit) (*DefaultsTest, error)
type DefaultsTestFindManyQuery = func(ctx context.Context, params QueryParams[DefaultsTest], selects *DefaultsTestSelect, omits *DefaultsTestOmit) ([]*DefaultsTest, error)
type DefaultsTestDeleteManyQuery = func(ctx context.Context, preds []PredicateOf[DefaultsTest]) (int64, error)
type DefaultsTestDeleteQuery = func(ctx context.Context, where UniquePredicate[DefaultsTest], selects *DefaultsTestSelect, omits *DefaultsTestOmit) (*DefaultsTest, error)
type DefaultsTestCountQuery = func(ctx context.Context, params QueryParams[DefaultsTest]) (int64, error)
type DefaultsTestUpdateQuery = func(ctx context.Context, where UniquePredicate[DefaultsTest], additional []PredicateOf[DefaultsTest], assignments []FieldAssignment, selects *DefaultsTestSelect, omits *DefaultsTestOmit) (*DefaultsTest, error)
type DefaultsTestUpdateManyQuery = func(ctx context.Context, preds []PredicateOf[DefaultsTest], assignments []FieldAssignment) (int64, error)
type DefaultsTestUpdateManyAndReturnQuery = func(ctx context.Context, preds []PredicateOf[DefaultsTest], assignments []FieldAssignment, selects *DefaultsTestSelect, omits *DefaultsTestOmit) ([]*DefaultsTest, error)

type DefaultsTestExtension struct {
	Create              func(ctx context.Context, input *DefaultsTestCreate, next DefaultsTestCreateQuery) (*DefaultsTest, error)
	CreateMany          func(ctx context.Context, inputs []*DefaultsTestCreate, next DefaultsTestCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, inputs []*DefaultsTestCreate, next DefaultsTestCreateManyAndReturnQuery) ([]*DefaultsTest, error)
	FindUnique          func(ctx context.Context, where UniquePredicate[DefaultsTest], additional []PredicateOf[DefaultsTest], selects *DefaultsTestSelect, omits *DefaultsTestOmit, next DefaultsTestFindUniqueQuery) (*DefaultsTest, error)
	FindFirst           func(ctx context.Context, params QueryParams[DefaultsTest], selects *DefaultsTestSelect, omits *DefaultsTestOmit, next DefaultsTestFindFirstQuery) (*DefaultsTest, error)
	FindMany            func(ctx context.Context, params QueryParams[DefaultsTest], selects *DefaultsTestSelect, omits *DefaultsTestOmit, next DefaultsTestFindManyQuery) ([]*DefaultsTest, error)
	DeleteMany          func(ctx context.Context, preds []PredicateOf[DefaultsTest], next DefaultsTestDeleteManyQuery) (int64, error)
	Delete              func(ctx context.Context, where UniquePredicate[DefaultsTest], selects *DefaultsTestSelect, omits *DefaultsTestOmit, next DefaultsTestDeleteQuery) (*DefaultsTest, error)
	Count               func(ctx context.Context, params QueryParams[DefaultsTest], next DefaultsTestCountQuery) (int64, error)
	Update              func(ctx context.Context, where UniquePredicate[DefaultsTest], additional []PredicateOf[DefaultsTest], assignments []FieldAssignment, selects *DefaultsTestSelect, omits *DefaultsTestOmit, next DefaultsTestUpdateQuery) (*DefaultsTest, error)
	UpdateMany          func(ctx context.Context, preds []PredicateOf[DefaultsTest], assignments []FieldAssignment, next DefaultsTestUpdateManyQuery) (int64, error)
	UpdateManyAndReturn func(ctx context.Context, preds []PredicateOf[DefaultsTest], assignments []FieldAssignment, selects *DefaultsTestSelect, omits *DefaultsTestOmit, next DefaultsTestUpdateManyAndReturnQuery) ([]*DefaultsTest, error)
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

var defaultsTestPKCols = []string{
	"uuid4",
}

var defaultsTestUniqueCols = []string{
	"uuid4",
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

func (b *DefaultsTestCreateBuilder) Assignments(assignments ...FieldAssignmentOf[DefaultsTest]) *DefaultsTestCreateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment{Col: a.Col, Val: a.Val})
	}
	return b
}

func (d *DefaultsTestDelegate) Create() *DefaultsTestCreateBuilder {
	return &DefaultsTestCreateBuilder{
		CreateBuilder: &CreateBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit]{
			execFunc: d.executeCreate,
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

func partitionDefaultsTestInputs(dialect Dialect, inputs []*DefaultsTestCreate) [][]*DefaultsTestCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*DefaultsTestCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*DefaultsTestCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*DefaultsTestCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*DefaultsTestCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*DefaultsTestCreate{inputs}
}

func (d *DefaultsTestDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *DefaultsTestSelect, omits *DefaultsTestOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*DefaultsTest, error) {
	input, err := assignmentsToDefaultsTestCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectDefaultsTestCols(selects, omits)

	if len(d.extensions) == 0 {
		hasRelations := selects.hasAnyRelation()
		if hasRelations {
			var res *DefaultsTest
			err = d.client.transaction(ctx, func(txQ *Queries) error {
				var err error
				res, err = txQ.DefaultsTest.runCreate(ctx, cols, vals, returningCols, defaultsTestPKCols, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.DefaultsTest.loadRelations(ctx, []*DefaultsTest{res}, selects)
			})
			return res, err
		}
		return d.runCreate(ctx, cols, vals, returningCols, defaultsTestPKCols, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args *DefaultsTestCreate) (*DefaultsTest, error) {
		cols, vals := args.ToColsVals()
		returningCols := selectDefaultsTestCols(selects, omits)

		hasRelations := selects.hasAnyRelation()
		var res *DefaultsTest
		var err error
		if hasRelations {
			err = d.client.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = txQ.DefaultsTest.runCreate(c, cols, vals, returningCols, defaultsTestPKCols, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.DefaultsTest.loadRelations(c, []*DefaultsTest{res}, selects)
			})
		} else {
			res, err = d.runCreate(c, cols, vals, returningCols, defaultsTestPKCols, conflictTarget, conflictAction)
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

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args []*DefaultsTestCreate) (int64, error) {
		return d.runCreateMany(c, args, conflictTarget, conflictAction)
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

	if len(d.extensions) == 0 {
		hasRelations := selects != nil && selects.hasAnyRelation()
		if hasRelations {
			var res []*DefaultsTest
			err := d.client.transaction(ctx, func(txQ *Queries) error {
				var err error
				res, err = txQ.DefaultsTest.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.DefaultsTest.loadRelations(ctx, res, selects)
			})
			return res, err
		}
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args []*DefaultsTestCreate) ([]*DefaultsTest, error) {
		hasRelations := selects != nil && selects.hasAnyRelation()
		if hasRelations {
			var res []*DefaultsTest
			err := d.client.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = txQ.DefaultsTest.runCreateManyAndReturn(c, args, selects, omits, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.DefaultsTest.loadRelations(c, res, selects)
			})
			return res, err
		}
		return d.runCreateManyAndReturn(c, args, selects, omits, conflictTarget, conflictAction)
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

func (d *DefaultsTestDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	pkCols []string,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*DefaultsTest, error) {
	query, clauseArgs := buildSingleInsertSQL(d.client, "DefaultsTest", cols, returningCols, pkCols, conflictTarget, conflictAction, len(vals))
	if len(clauseArgs) > 0 {
		vals = append(vals, clauseArgs...)
	}

	var res DefaultsTest
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

func (d *DefaultsTestDelegate) runCreateFallback(ctx context.Context, query string, vals []any, cols []string, returningCols []string, pkCols []string) (*DefaultsTest, error) {
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
	selectSb.Grow(64 + len(returningCols)*15 + len("DefaultsTest") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "DefaultsTest")
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

	var res DefaultsTest
	if rows.Next() {
		if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
			return nil, err
		}
		return &res, nil
	}
	return nil, rows.Err()
}

func (d *DefaultsTestDelegate) buildBulkInsertSQL(q *Queries, batch []*DefaultsTestCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 9)
	for i, c := range defaultsTestDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "DefaultsTest")
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
			case "uuid4":
				if input.Uuid4 != nil {
					vals = append(vals, *input.Uuid4)
				} else {
					vals = append(vals, generateCUID())
				}
			case "uuid7":
				if input.Uuid7 != nil {
					vals = append(vals, *input.Uuid7)
				} else {
					vals = append(vals, generateUUID7())
				}
			case "uuidNoArgs":
				if input.UuidNoArgs != nil {
					vals = append(vals, *input.UuidNoArgs)
				} else {
					vals = append(vals, generateUUID())
				}
			case "cuid1":
				if input.Cuid1 != nil {
					vals = append(vals, *input.Cuid1)
				} else {
					vals = append(vals, generateCUID())
				}
			case "cuid2":
				if input.Cuid2 != nil {
					vals = append(vals, *input.Cuid2)
				} else {
					vals = append(vals, generateCUID2())
				}
			case "cuidNoArgs":
				if input.CuidNoArgs != nil {
					vals = append(vals, *input.CuidNoArgs)
				} else {
					vals = append(vals, generateCUID())
				}
			case "ulid":
				if input.Ulid != nil {
					vals = append(vals, *input.Ulid)
				} else {
					vals = append(vals, generateULID())
				}
			case "nanoid":
				if input.Nanoid != nil {
					vals = append(vals, *input.Nanoid)
				} else {
					vals = append(vals, generateNanoID())
				}
			case "now":
				if input.Now != nil {
					vals = append(vals, *input.Now)
				} else {
					vals = append(vals, time.Now())
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

func (d *DefaultsTestDelegate) runCreateMany(ctx context.Context, inputs []*DefaultsTestCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionDefaultsTestInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)

		var conflictCols []string
		if conflictTarget != nil {
			conflictCols = conflictTarget.UniqueColumns()
		}
		var nonConflictCols []string
		if conflictAction != nil && conflictAction.Type == ConflictActionUpdateNewValues {
			nonConflictCols = computeNonConflictCols(cols, conflictCols, defaultsTestPKCols)
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

func (d *DefaultsTestDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*DefaultsTestCreate,
	selects *DefaultsTestSelect,
	omits *DefaultsTestOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*DefaultsTest, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionDefaultsTestInputs(d.client.dialect, inputs)
	returningCols := selectDefaultsTestCols(selects, omits)
	hasRelations := selects != nil && selects.hasAnyRelation()

	recordsOut := make([]*DefaultsTest, 0, len(inputs))

	runBatch := func(txQ *Queries, batch []*DefaultsTestCreate) error {
		cols, vals, queryStr := d.buildBulkInsertSQL(txQ, batch, 1)

		var conflictCols []string
		if conflictTarget != nil {
			conflictCols = conflictTarget.UniqueColumns()
		}
		var nonConflictCols []string
		if conflictAction != nil && conflictAction.Type == ConflictActionUpdateNewValues {
			nonConflictCols = computeNonConflictCols(cols, conflictCols, defaultsTestPKCols)
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
				var res DefaultsTest
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
		selectSb.Grow(64 + len(returningCols)*15 + len("DefaultsTest") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			txQ.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		txQ.dialect.WriteQuotedIdent(&selectSb, "DefaultsTest")
		selectSb.WriteString(" WHERE ")
		txQ.dialect.WriteQuotedIdent(&selectSb, defaultsTestPKCols[0])
		selectSb.WriteString(" >= ")
		txQ.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		txQ.dialect.WriteQuotedIdent(&selectSb, defaultsTestPKCols[0])
		selectSb.WriteString(" < ")
		txQ.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := txQ.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var res DefaultsTest
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
				return txQ.DefaultsTest.loadRelations(ctx, recordsOut, selects)
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

type DefaultsTestUpdateBuilder struct {
	*UpdateBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit]
}

type DefaultsTestUpdateManyBuilder struct {
	*UpdateManyBuilder[DefaultsTest]
}

type DefaultsTestUpdateManyAndReturnBuilder struct {
	*UpdateManyAndReturnBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit]
}

func (b *DefaultsTestUpdateBuilder) SetUuid4(v string) *DefaultsTestUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "uuid4", Val: v})
	return b
}

func (b *DefaultsTestUpdateManyBuilder) SetUuid4(v string) *DefaultsTestUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "uuid4", Val: v})
	return b
}

func (b *DefaultsTestUpdateManyAndReturnBuilder) SetUuid4(v string) *DefaultsTestUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "uuid4", Val: v})
	return b
}
func (b *DefaultsTestUpdateBuilder) SetUuid7(v string) *DefaultsTestUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "uuid7", Val: v})
	return b
}

func (b *DefaultsTestUpdateManyBuilder) SetUuid7(v string) *DefaultsTestUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "uuid7", Val: v})
	return b
}

func (b *DefaultsTestUpdateManyAndReturnBuilder) SetUuid7(v string) *DefaultsTestUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "uuid7", Val: v})
	return b
}
func (b *DefaultsTestUpdateBuilder) SetUuidNoArgs(v string) *DefaultsTestUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "uuidNoArgs", Val: v})
	return b
}

func (b *DefaultsTestUpdateManyBuilder) SetUuidNoArgs(v string) *DefaultsTestUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "uuidNoArgs", Val: v})
	return b
}

func (b *DefaultsTestUpdateManyAndReturnBuilder) SetUuidNoArgs(v string) *DefaultsTestUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "uuidNoArgs", Val: v})
	return b
}
func (b *DefaultsTestUpdateBuilder) SetCuid1(v string) *DefaultsTestUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "cuid1", Val: v})
	return b
}

func (b *DefaultsTestUpdateManyBuilder) SetCuid1(v string) *DefaultsTestUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "cuid1", Val: v})
	return b
}

func (b *DefaultsTestUpdateManyAndReturnBuilder) SetCuid1(v string) *DefaultsTestUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "cuid1", Val: v})
	return b
}
func (b *DefaultsTestUpdateBuilder) SetCuid2(v string) *DefaultsTestUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "cuid2", Val: v})
	return b
}

func (b *DefaultsTestUpdateManyBuilder) SetCuid2(v string) *DefaultsTestUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "cuid2", Val: v})
	return b
}

func (b *DefaultsTestUpdateManyAndReturnBuilder) SetCuid2(v string) *DefaultsTestUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "cuid2", Val: v})
	return b
}
func (b *DefaultsTestUpdateBuilder) SetCuidNoArgs(v string) *DefaultsTestUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "cuidNoArgs", Val: v})
	return b
}

func (b *DefaultsTestUpdateManyBuilder) SetCuidNoArgs(v string) *DefaultsTestUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "cuidNoArgs", Val: v})
	return b
}

func (b *DefaultsTestUpdateManyAndReturnBuilder) SetCuidNoArgs(v string) *DefaultsTestUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "cuidNoArgs", Val: v})
	return b
}
func (b *DefaultsTestUpdateBuilder) SetUlid(v string) *DefaultsTestUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "ulid", Val: v})
	return b
}

func (b *DefaultsTestUpdateManyBuilder) SetUlid(v string) *DefaultsTestUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "ulid", Val: v})
	return b
}

func (b *DefaultsTestUpdateManyAndReturnBuilder) SetUlid(v string) *DefaultsTestUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "ulid", Val: v})
	return b
}
func (b *DefaultsTestUpdateBuilder) SetNanoid(v string) *DefaultsTestUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "nanoid", Val: v})
	return b
}

func (b *DefaultsTestUpdateManyBuilder) SetNanoid(v string) *DefaultsTestUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "nanoid", Val: v})
	return b
}

func (b *DefaultsTestUpdateManyAndReturnBuilder) SetNanoid(v string) *DefaultsTestUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "nanoid", Val: v})
	return b
}
func (b *DefaultsTestUpdateBuilder) SetNow(v time.Time) *DefaultsTestUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "now", Val: v})
	return b
}

func (b *DefaultsTestUpdateManyBuilder) SetNow(v time.Time) *DefaultsTestUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "now", Val: v})
	return b
}

func (b *DefaultsTestUpdateManyAndReturnBuilder) SetNow(v time.Time) *DefaultsTestUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "now", Val: v})
	return b
}

func (b *DefaultsTestUpdateBuilder) Assignments(assignments ...FieldAssignmentOf[DefaultsTest]) *DefaultsTestUpdateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment{Col: a.Col, Val: a.Val})
	}
	return b
}

func (b *DefaultsTestUpdateManyBuilder) Assignments(assignments ...FieldAssignmentOf[DefaultsTest]) *DefaultsTestUpdateManyBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment{Col: a.Col, Val: a.Val})
	}
	return b
}

func (b *DefaultsTestUpdateManyAndReturnBuilder) Assignments(assignments ...FieldAssignmentOf[DefaultsTest]) *DefaultsTestUpdateManyAndReturnBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment{Col: a.Col, Val: a.Val})
	}
	return b
}

func (d *DefaultsTestDelegate) Update(where UniquePredicate[DefaultsTest], additional ...PredicateOf[DefaultsTest]) *DefaultsTestUpdateBuilder {
	return &DefaultsTestUpdateBuilder{
		UpdateBuilder: &UpdateBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit]{
			where:      where,
			additional: additional,
			execFunc:   d.executeUpdate,
		},
	}
}

func (d *DefaultsTestDelegate) UpdateMany(preds ...PredicateOf[DefaultsTest]) *DefaultsTestUpdateManyBuilder {
	return &DefaultsTestUpdateManyBuilder{
		UpdateManyBuilder: &UpdateManyBuilder[DefaultsTest]{
			where:    preds,
			execFunc: d.executeUpdateMany,
		},
	}
}

func (d *DefaultsTestDelegate) UpdateManyAndReturn(preds ...PredicateOf[DefaultsTest]) *DefaultsTestUpdateManyAndReturnBuilder {
	return &DefaultsTestUpdateManyAndReturnBuilder{
		UpdateManyAndReturnBuilder: &UpdateManyAndReturnBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit]{
			where:    preds,
			execFunc: d.executeUpdateManyAndReturn,
		},
	}
}

func (d *DefaultsTestDelegate) buildUpdateSQL(preds []PredicateOf[DefaultsTest], assignments []FieldAssignment, returningCols []string) (string, []any) {
	whereClause, predVals, _ := CompilePredicates(d.client.dialect, preds, len(assignments)+1)

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	d.client.dialect.WriteQuotedIdent(&sb, "DefaultsTest")
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

func (d *DefaultsTestDelegate) executeUpdate(ctx context.Context, where UniquePredicate[DefaultsTest], additional []PredicateOf[DefaultsTest], assignments []FieldAssignment, selects *DefaultsTestSelect, omits *DefaultsTestOmit) (*DefaultsTest, error) {
	if len(d.extensions) == 0 {
		return d.runUpdate(ctx, where, additional, assignments, selects, omits)
	}

	curr := func(c context.Context, w UniquePredicate[DefaultsTest], add []PredicateOf[DefaultsTest], a []FieldAssignment, s *DefaultsTestSelect, o *DefaultsTestOmit) (*DefaultsTest, error) {
		return d.runUpdate(c, w, add, a, s, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Update != nil {
			next, hook := curr, ext.Update
			curr = func(c context.Context, w UniquePredicate[DefaultsTest], add []PredicateOf[DefaultsTest], a []FieldAssignment, s *DefaultsTestSelect, o *DefaultsTestOmit) (*DefaultsTest, error) {
				return hook(c, w, add, a, s, o, next)
			}
		}
	}

	return curr(ctx, where, additional, assignments, selects, omits)
}

func (d *DefaultsTestDelegate) runUpdate(ctx context.Context, where UniquePredicate[DefaultsTest], additional []PredicateOf[DefaultsTest], assignments []FieldAssignment, selects *DefaultsTestSelect, omits *DefaultsTestOmit) (*DefaultsTest, error) {
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
		var res *DefaultsTest
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.DefaultsTest.runUpdate(ctx, where, additional, assignments, selects, omits)
			} else {
				res, err = txQ.DefaultsTest.runUpdateFallback(ctx, where, additional, assignments, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectDefaultsTestCols(selects, omits, defaultsTestPKCols...)
	allPreds := append([]PredicateOf[DefaultsTest]{where}, additional...)
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

	var res DefaultsTest
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, scanErr
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*DefaultsTest{&res}, selects); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *DefaultsTestDelegate) execUpdateStmt(ctx context.Context, preds []PredicateOf[DefaultsTest], assignments []FieldAssignment) (int64, error) {
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

func (d *DefaultsTestDelegate) runUpdateFallback(ctx context.Context, where UniquePredicate[DefaultsTest], additional []PredicateOf[DefaultsTest], assignments []FieldAssignment, selects *DefaultsTestSelect, omits *DefaultsTestOmit) (*DefaultsTest, error) {
	allPreds := append([]PredicateOf[DefaultsTest]{where}, additional...)
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

func (d *DefaultsTestDelegate) executeUpdateMany(ctx context.Context, preds []PredicateOf[DefaultsTest], assignments []FieldAssignment) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runUpdateMany(ctx, preds, assignments)
	}

	curr := func(c context.Context, p []PredicateOf[DefaultsTest], a []FieldAssignment) (int64, error) {
		return d.runUpdateMany(c, p, a)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.UpdateMany != nil {
			next, hook := curr, ext.UpdateMany
			curr = func(c context.Context, p []PredicateOf[DefaultsTest], a []FieldAssignment) (int64, error) {
				return hook(c, p, a, next)
			}
		}
	}

	return curr(ctx, preds, assignments)
}

func (d *DefaultsTestDelegate) runUpdateMany(ctx context.Context, preds []PredicateOf[DefaultsTest], assignments []FieldAssignment) (int64, error) {
	return d.execUpdateStmt(ctx, preds, assignments)
}

// -----------------------------------------------------------------------------
// UpdateManyAndReturn
// -----------------------------------------------------------------------------

func (d *DefaultsTestDelegate) executeUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[DefaultsTest], assignments []FieldAssignment, selects *DefaultsTestSelect, omits *DefaultsTestOmit) ([]*DefaultsTest, error) {
	if len(d.extensions) == 0 {
		return d.runUpdateManyAndReturn(ctx, preds, assignments, selects, omits)
	}

	curr := func(c context.Context, p []PredicateOf[DefaultsTest], a []FieldAssignment, s *DefaultsTestSelect, o *DefaultsTestOmit) ([]*DefaultsTest, error) {
		return d.runUpdateManyAndReturn(c, p, a, s, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.UpdateManyAndReturn != nil {
			next, hook := curr, ext.UpdateManyAndReturn
			curr = func(c context.Context, p []PredicateOf[DefaultsTest], a []FieldAssignment, s *DefaultsTestSelect, o *DefaultsTestOmit) ([]*DefaultsTest, error) {
				return hook(c, p, a, s, o, next)
			}
		}
	}

	return curr(ctx, preds, assignments, selects, omits)
}

func (d *DefaultsTestDelegate) runUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[DefaultsTest], assignments []FieldAssignment, selects *DefaultsTestSelect, omits *DefaultsTestOmit) ([]*DefaultsTest, error) {
	if len(assignments) == 0 {
		return d.runFindMany(ctx, QueryParams[DefaultsTest]{Where: preds}, selects, omits)
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
		var res []*DefaultsTest
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.DefaultsTest.runUpdateManyAndReturn(ctx, preds, assignments, selects, omits)
			} else {
				res, err = txQ.DefaultsTest.runUpdateManyAndReturnFallback(ctx, preds, assignments, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectDefaultsTestCols(selects, omits, defaultsTestPKCols...)
	query, setVals := d.buildUpdateSQL(preds, assignments, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	results := make([]*DefaultsTest, 0)
	for rows.Next() {
		var res DefaultsTest
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

func (d *DefaultsTestDelegate) runUpdateManyAndReturnFallback(ctx context.Context, preds []PredicateOf[DefaultsTest], assignments []FieldAssignment, selects *DefaultsTestSelect, omits *DefaultsTestOmit) ([]*DefaultsTest, error) {
	affected, err := d.execUpdateStmt(ctx, preds, assignments)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return []*DefaultsTest{}, nil
	}
	return d.runFindMany(ctx, QueryParams[DefaultsTest]{Where: preds}, selects, omits)
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
	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, where, additional, selects, omits)
	}

	curr := func(c context.Context, w UniquePredicate[DefaultsTest], add []PredicateOf[DefaultsTest], sel *DefaultsTestSelect, o *DefaultsTestOmit) (*DefaultsTest, error) {
		return d.runFindUnique(c, w, add, sel, o)
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
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	curr := func(c context.Context, p QueryParams[DefaultsTest], sel *DefaultsTestSelect, o *DefaultsTestOmit) (*DefaultsTest, error) {
		return d.runFindFirst(c, p, sel, o)
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
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	curr := func(c context.Context, p QueryParams[DefaultsTest], sel *DefaultsTestSelect, o *DefaultsTestOmit) ([]*DefaultsTest, error) {
		return d.runFindMany(c, p, sel, o)
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

func (d *DefaultsTestDelegate) runFindUnique(ctx context.Context, where UniquePredicate[DefaultsTest], additional []PredicateOf[DefaultsTest], selects *DefaultsTestSelect, omits *DefaultsTestOmit) (*DefaultsTest, error) {
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
	allPreds := append([]PredicateOf[DefaultsTest]{where}, additional...)
	whereClause, vals, _ := CompilePredicates(d.client.dialect, allPreds)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	returningCols := selectDefaultsTestCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*DefaultsTest{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *DefaultsTestDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[DefaultsTest],
	selects *DefaultsTestSelect,
	omits *DefaultsTestOmit,
) (*DefaultsTest, error) {
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
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, defaultsTestPKCols, defaultsTestUniqueCols, "DefaultsTest", nextIdx, params.Take)
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
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, defaultsTestPKCols, defaultsTestUniqueCols, isCursorQuery, params.Take)
	returningCols := selectDefaultsTestCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, orderByClause, vals, returningCols, params.Skip)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*DefaultsTest{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *DefaultsTestDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[DefaultsTest],
	selects *DefaultsTestSelect,
	omits *DefaultsTestOmit,
) ([]*DefaultsTest, error) {
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
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, defaultsTestPKCols, defaultsTestUniqueCols, "DefaultsTest", nextIdx, params.Take)
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
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, defaultsTestPKCols, defaultsTestUniqueCols, isCursorQuery, params.Take)
	returningCols := selectDefaultsTestCols(selects, omits)

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

func (d *DefaultsTestDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*DefaultsTest, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "DefaultsTest", returningCols, whereClause, orderByClause, &limitOne, skip)
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

	var res DefaultsTest
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (d *DefaultsTestDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*DefaultsTest, error) {
	query := buildSelectSQL(d.client, "DefaultsTest", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*DefaultsTest, 0)
	for rows.Next() {
		var res DefaultsTest
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
func (d *DefaultsTestDelegate) DeleteMany(preds ...PredicateOf[DefaultsTest]) *DeleteManyBuilder[DefaultsTest] {
	return &DeleteManyBuilder[DefaultsTest]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *DefaultsTestDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[DefaultsTest]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	curr := func(c context.Context, p []PredicateOf[DefaultsTest]) (int64, error) {
		return d.runDeleteMany(c, p)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.DeleteMany != nil {
			next, hook := curr, ext.DeleteMany
			curr = func(c context.Context, p []PredicateOf[DefaultsTest]) (int64, error) {
				return hook(c, p, next)
			}
		}
	}

	return curr(ctx, preds)
}

func (d *DefaultsTestDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[DefaultsTest]) (int64, error) {
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
	d.client.dialect.WriteQuotedIdent(&sb, "DefaultsTest")
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

func (d *DefaultsTestDelegate) Delete(where UniquePredicate[DefaultsTest]) *DeleteBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit] {
	return &DeleteBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit]{
		where:    where,
		execFunc: d.executeDelete,
	}
}

func (d *DefaultsTestDelegate) executeDelete(ctx context.Context, where UniquePredicate[DefaultsTest], selects *DefaultsTestSelect, omits *DefaultsTestOmit) (*DefaultsTest, error) {
	if len(d.extensions) == 0 {
		return d.runDelete(ctx, where, selects, omits)
	}

	curr := func(c context.Context, w UniquePredicate[DefaultsTest], s *DefaultsTestSelect, o *DefaultsTestOmit) (*DefaultsTest, error) {
		return d.runDelete(c, w, s, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Delete != nil {
			next, hook := curr, ext.Delete
			curr = func(c context.Context, w UniquePredicate[DefaultsTest], s *DefaultsTestSelect, o *DefaultsTestOmit) (*DefaultsTest, error) {
				return hook(c, w, s, o, next)
			}
		}
	}

	return curr(ctx, where, selects, omits)
}

func (d *DefaultsTestDelegate) runDelete(ctx context.Context, where UniquePredicate[DefaultsTest], selects *DefaultsTestSelect, omits *DefaultsTestOmit) (*DefaultsTest, error) {
	if err := where.Validate(); err != nil {
		return nil, err
	}

	returningCols := selectDefaultsTestCols(selects, omits, defaultsTestPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *DefaultsTest
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.DefaultsTest.executeFindUnique(ctx, where, nil, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return sql.ErrNoRows
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "DefaultsTest")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[DefaultsTest]
			pkPreds = append(pkPreds, Predicate[DefaultsTest]{
				Data: PredicateData{
					Column:   "uuid4",
					Operator: "=",
					Value:    res.Uuid4,
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
	d.client.dialect.WriteQuotedIdent(&sb, "DefaultsTest")

	whereClause, vals, _ := CompilePredicates(d.client.dialect, []PredicateOf[DefaultsTest]{where})
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

	var row DefaultsTest
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, err
	}
	return &row, nil
}
func (d *DefaultsTestDelegate) Count(preds ...PredicateOf[DefaultsTest]) *CountBuilder[DefaultsTest] {
	return &CountBuilder[DefaultsTest]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *DefaultsTestDelegate) executeCount(ctx context.Context, params QueryParams[DefaultsTest]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	curr := func(c context.Context, p QueryParams[DefaultsTest]) (int64, error) {
		return d.runCount(c, p)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Count != nil {
			next, hook := curr, ext.Count
			curr = func(c context.Context, p QueryParams[DefaultsTest]) (int64, error) {
				return hook(c, p, next)
			}
		}
	}

	return curr(ctx, params)
}

func (d *DefaultsTestDelegate) runCount(ctx context.Context, params QueryParams[DefaultsTest]) (int64, error) {
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
		d.client.dialect.WriteQuotedIdent(&subQuery, "DefaultsTest")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "DefaultsTest")
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
func (d *DefaultsTestDelegate) loadRelations(ctx context.Context, records []*DefaultsTest, selects *DefaultsTestSelect) error {
	_ = ctx
	if selects == nil || len(records) == 0 {
		return nil
	}

	return nil
}
