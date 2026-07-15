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

type DefaultsTestCreateBuilder struct {
	*CreateBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit]
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
			client:      d.client,
			assignments: assignments,
			execFunc:    d.client.executeDefaultsTestCreate,
		},
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
				ValidateString(errs, "uuid4", v, false, 0, false, false)
			} else {
				errs.Add("uuid4", a.Val, "type", "field uuid4 must be of type string")
			}
		case "uuid7":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "uuid7", v, false, 0, false, false)
			} else {
				errs.Add("uuid7", a.Val, "type", "field uuid7 must be of type string")
			}
		case "uuidNoArgs":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "uuidNoArgs", v, false, 0, false, false)
			} else {
				errs.Add("uuidNoArgs", a.Val, "type", "field uuidNoArgs must be of type string")
			}
		case "cuid1":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "cuid1", v, false, 0, false, false)
			} else {
				errs.Add("cuid1", a.Val, "type", "field cuid1 must be of type string")
			}
		case "cuid2":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "cuid2", v, false, 0, false, false)
			} else {
				errs.Add("cuid2", a.Val, "type", "field cuid2 must be of type string")
			}
		case "cuidNoArgs":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "cuidNoArgs", v, false, 0, false, false)
			} else {
				errs.Add("cuidNoArgs", a.Val, "type", "field cuidNoArgs must be of type string")
			}
		case "ulid":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "ulid", v, false, 0, false, false)
			} else {
				errs.Add("ulid", a.Val, "type", "field ulid must be of type string")
			}
		case "nanoid":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "nanoid", v, false, 0, false, false)
			} else {
				errs.Add("nanoid", a.Val, "type", "field nanoid must be of type string")
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

	rowMap := input.ToRowMap()
	cols, vals := mapToColsVals(rowMap, DefaultsTestColOrder)

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

func (d *DefaultsTestDelegate) CreateMany(builders ...*DefaultsTestCreateBuilder) *CreateManyBuilder[DefaultsTest] {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &CreateManyBuilder[DefaultsTest]{
		client:   d.client,
		records:  records,
		execFunc: d.client.executeDefaultsTestCreateMany,
	}
}

func (d *DefaultsTestDelegate) CreateManyAndReturn(builders ...*DefaultsTestCreateBuilder) *CreateManyAndReturnBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit] {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
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
func (d *DefaultsTestDelegate) FindUnique(where UniquePredicate[DefaultsTest], additional ...PredicateOf[DefaultsTest]) *FindUniqueBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit] {
	return &FindUniqueBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit]{
		client:     d.client,
		where:      where,
		additional: additional,
		execFunc:   d.client.executeDefaultsTestFindUnique,
	}
}

func (d *DefaultsTestDelegate) FindFirst(preds ...PredicateOf[DefaultsTest]) *FindFirstBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit] {
	return &FindFirstBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeDefaultsTestFindFirst,
	}
}

func (d *DefaultsTestDelegate) FindMany(preds ...PredicateOf[DefaultsTest]) *FindManyBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit] {
	return &FindManyBuilder[DefaultsTest, DefaultsTestSelect, DefaultsTestOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeDefaultsTestFindMany,
	}
}

func (q *Queries) executeDefaultsTestFindUnique(ctx context.Context, where UniquePredicate[DefaultsTest], additional []PredicateOf[DefaultsTest], selects *DefaultsTestSelect, omits *DefaultsTestOmit) (*DefaultsTest, error) {
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
	whereClause, vals := CompilePredicates(q.dialect, allPreds)
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
		nil,
	)
}

func (q *Queries) executeDefaultsTestFindFirst(
	ctx context.Context,
	params QueryParams[DefaultsTest],
	selects *DefaultsTestSelect,
	omits *DefaultsTestOmit,
) (*DefaultsTest, error) {
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
	returningCols := q.selectDefaultsTestCols(selects, omits)
	return executeSingleWithRelations(ctx, q, "DefaultsTest", whereClause, vals, returningCols,
		func(res *DefaultsTest, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*DefaultsTest) error {
			return txQ.loadDefaultsTestRelations(ctx, results, selects)
		},
		params.Skip,
	)
}

func (q *Queries) executeDefaultsTestFindMany(
	ctx context.Context,
	params QueryParams[DefaultsTest],
	selects *DefaultsTestSelect,
	omits *DefaultsTestOmit,
) ([]*DefaultsTest, error) {
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
	returningCols := q.selectDefaultsTestCols(selects, omits)
	return executeManyWithRelations(ctx, q, "DefaultsTest", whereClause, vals, returningCols,
		func(res *DefaultsTest, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*DefaultsTest) error {
			return txQ.loadDefaultsTestRelations(ctx, results, selects)
		},
		params.Take,
		params.Skip,
	)
}
func (q *Queries) loadDefaultsTestRelations(ctx context.Context, records []*DefaultsTest, selects *DefaultsTestSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}

	return nil
}
