package valk

import (
	"context"
	"fmt"
	"slices"
)

// Category represents the database model
type Category struct {
	Id    int32             `db:"id" json:"id"`
	Name  string            `db:"name" json:"name"`
	Posts []*CategoryToPost `json:"posts,omitempty"`
}

// CategoryCreate is used for hooks only — the Create API uses FieldAssignment
type CategoryCreate struct {
	Id   *int32 `json:"id"`
	Name string `json:"name"`
}

// CategorySelect specifies which fields to include
type CategorySelect struct {
	Id    bool                      `json:"id"`
	Name  bool                      `json:"name"`
	Posts CategoryToPostSelectQuery `json:"posts,omitempty"`
}

// CategoryOmit specifies which fields to exclude
type CategoryOmit struct {
	Id   bool `json:"id"`
	Name bool `json:"name"`
}

type CategorySelectQuery interface {
	GetRelationParams() (*CategorySelect, *CategoryOmit, QueryParams[Category])
}

func (s *CategorySelect) GetRelationParams() (*CategorySelect, *CategoryOmit, QueryParams[Category]) {
	return s, nil, QueryParams[Category]{}
}

// CategoryQueryBuilder builds a query for the relation Category
type CategoryQueryBuilder struct {
	selects *CategorySelect
	omits   *CategoryOmit
	where   []PredicateOf[Category]
	take    *int
	skip    *int
	orderBy []OrderBy[Category]
}

func (b *CategoryQueryBuilder) Where(preds ...PredicateOf[Category]) *CategoryQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *CategoryQueryBuilder) Take(limit int) *CategoryQueryBuilder {
	b.take = &limit
	return b
}

func (b *CategoryQueryBuilder) Skip(offset int) *CategoryQueryBuilder {
	b.skip = &offset
	return b
}

func (b *CategoryQueryBuilder) OrderBy(orders ...OrderBy[Category]) *CategoryQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *CategoryQueryBuilder) Select(s CategorySelect) *CategoryQueryBuilder {
	b.selects = &s
	return b
}

func (b *CategoryQueryBuilder) Omit(o CategoryOmit) *CategoryQueryBuilder {
	b.omits = &o
	return b
}

func (b *CategoryQueryBuilder) GetRelationParams() (*CategorySelect, *CategoryOmit, QueryParams[Category]) {
	if b == nil {
		return nil, nil, QueryParams[Category]{}
	}
	return b.selects, b.omits, QueryParams[Category]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
	}
}

type CategoryCreateQuery = func(ctx context.Context, args *CategoryCreate) (*Category, error)
type CategoryCreateManyQuery = func(ctx context.Context, args []*CategoryCreate) (int64, error)
type CategoryCreateManyAndReturnQuery = func(ctx context.Context, args []*CategoryCreate) ([]*Category, error)
type CategoryFindUniqueQuery = func(ctx context.Context, where UniquePredicate[Category], additional []PredicateOf[Category], selects *CategorySelect, omits *CategoryOmit) (*Category, error)
type CategoryFindFirstQuery = func(ctx context.Context, params QueryParams[Category], selects *CategorySelect, omits *CategoryOmit) (*Category, error)
type CategoryFindManyQuery = func(ctx context.Context, params QueryParams[Category], selects *CategorySelect, omits *CategoryOmit) ([]*Category, error)

type CategoryExtension struct {
	Create              func(ctx context.Context, input *CategoryCreate, next CategoryCreateQuery) (*Category, error)
	CreateMany          func(ctx context.Context, inputs []*CategoryCreate, next CategoryCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, inputs []*CategoryCreate, next CategoryCreateManyAndReturnQuery) ([]*Category, error)
	FindUnique          func(ctx context.Context, where UniquePredicate[Category], additional []PredicateOf[Category], selects *CategorySelect, omits *CategoryOmit, next CategoryFindUniqueQuery) (*Category, error)
	FindFirst           func(ctx context.Context, params QueryParams[Category], selects *CategorySelect, omits *CategoryOmit, next CategoryFindFirstQuery) (*Category, error)
	FindMany            func(ctx context.Context, params QueryParams[Category], selects *CategorySelect, omits *CategoryOmit, next CategoryFindManyQuery) ([]*Category, error)
}

type CategoryDelegate struct {
	client     *Queries
	extensions []CategoryExtension
}

func (d *CategoryDelegate) Use(exts ...CategoryExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *Category) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "id":
			targets[i] = &m.Id
		case "name":
			targets[i] = &m.Name
		}
	}
	return targets
}

var categoryDefaultCols = []string{
	"id",
	"name",
}

func selectCategoryCols(selects *CategorySelect, omits *CategoryOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return categoryDefaultCols
	}

	anySelected := selects != nil && (selects.Id || selects.Name || selects.Posts != nil)

	specs := []colSpec{
		{"id", selects != nil && selects.Id, omits != nil && omits.Id, selects != nil && selects.hasAnyRelation()},
		{"name", selects != nil && selects.Name, omits != nil && omits.Name, false},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *CategorySelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Posts != nil
}

type CategoryCreateBuilder struct {
	*CreateBuilder[Category, CategorySelect, CategoryOmit]
}

func (b *CategoryCreateBuilder) OnConflict(target UniqueConstraintTarget) *CategoryConflictBuilder[CategoryCreateBuilder] {
	return &CategoryConflictBuilder[CategoryCreateBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *CategoryCreateBuilder) SetId(v int32) *CategoryCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}
func (b *CategoryCreateBuilder) SetName(v string) *CategoryCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}

func (d *CategoryDelegate) Create(assignments ...FieldAssignment) *CategoryCreateBuilder {
	return &CategoryCreateBuilder{
		CreateBuilder: &CreateBuilder[Category, CategorySelect, CategoryOmit]{
			assignments: assignments,
			execFunc:    d.executeCreate,
		},
	}
}

const (
	providedCategoryId   uint64 = 1 << 0
	providedCategoryName uint64 = 1 << 1
)

func assignmentsToCategoryCreate(assignments []FieldAssignment) (CategoryCreate, error) {
	var input CategoryCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "id":
			provided |= providedCategoryId
			if v, ok := a.Val.(int32); ok {
				input.Id = &v
				ValidateInt32(&errs, "id", v, "")
			} else {
				errs.Add("id", a.Val, "type", "field id must be of type int32")
			}
		case "name":
			provided |= providedCategoryName
			if v, ok := a.Val.(string); ok {
				input.Name = v
				ValidateString(&errs, "name", v, true, 0, false, false)
			} else {
				errs.Add("name", a.Val, "type", "field name must be of type string")
			}
		}
	}
	if provided&providedCategoryName == 0 {
		errs.Add("name", "", "required", "field Name is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *CategoryCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 2)
	vals = make([]any, 0, 2)
	if s.Id != nil {
		cols = append(cols, "id")
		vals = append(vals, *s.Id)
	}
	cols = append(cols, "name")
	vals = append(vals, s.Name)
	return
}

func (s *CategoryCreate) ToRowMap() map[string]any {
	cols, vals := s.ToColsVals()
	m := make(map[string]any, len(cols))
	for i, c := range cols {
		m[c] = vals[i]
	}
	return m
}

func (d *CategoryDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *CategorySelect, omits *CategoryOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*Category, error) {
	input, err := assignmentsToCategoryCreate(assignments)
	if err != nil {
		return nil, err
	}

	curr := func(c context.Context, args *CategoryCreate) (*Category, error) {
		cols, vals := args.ToColsVals()

		returningCols := selectCategoryCols(selects, omits)

		scanFunc := func(res *Category, cols []string) []any {
			return res.ScanFields(cols)
		}

		pkCols := []string{
			"id",
		}

		hasRelations := selects.hasAnyRelation()

		var res *Category
		var err error
		if hasRelations {
			err = d.client.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = executeInsert(c, txQ, "Category", cols, vals, returningCols, pkCols, scanFunc, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.Category.loadRelations(c, []*Category{res}, selects)
			})
		} else {
			res, err = executeInsert(c, d.client, "Category", cols, vals, returningCols, pkCols, scanFunc, conflictTarget, conflictAction)
		}
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Create != nil {
			next, hook := curr, ext.Create
			curr = func(c context.Context, input *CategoryCreate) (*Category, error) {
				return hook(c, input, next)
			}
		}
	}

	return curr(ctx, &input)
}

type CategoryCreateManyBuilder struct {
	*CreateManyBuilder[Category]
}

func (b *CategoryCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *CategoryConflictBuilder[CategoryCreateManyBuilder] {
	return &CategoryConflictBuilder[CategoryCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type CategoryCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[Category, CategorySelect, CategoryOmit]
}

func (b *CategoryCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *CategoryConflictBuilder[CategoryCreateManyAndReturnBuilder] {
	return &CategoryConflictBuilder[CategoryCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (d *CategoryDelegate) CreateMany(builders ...*CategoryCreateBuilder) *CategoryCreateManyBuilder {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &CategoryCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[Category]{
			records:  records,
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *CategoryDelegate) CreateManyAndReturn(builders ...*CategoryCreateBuilder) *CategoryCreateManyAndReturnBuilder {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &CategoryCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[Category, CategorySelect, CategoryOmit]{
			records:  records,
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func (d *CategoryDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs := make([]*CategoryCreate, len(records))
	for i, rec := range records {
		input, err := assignmentsToCategoryCreate(rec.Assignments)
		if err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &input
	}

	curr := func(c context.Context, args []*CategoryCreate) (int64, error) {
		rowMaps := make([]map[string]any, len(args))
		for i, input := range args {
			rowMaps[i] = input.ToRowMap()
		}

		pkCols := []string{
			"id",
		}

		return executeCreateMany(c, d.client, rowMaps, "Category", categoryDefaultCols, pkCols, conflictTarget, conflictAction)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.CreateMany != nil {
			next, hook := curr, ext.CreateMany
			curr = func(c context.Context, inputs []*CategoryCreate) (int64, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
}

func (d *CategoryDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *CategorySelect, omits *CategoryOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*Category, error) {
	inputs := make([]*CategoryCreate, len(records))
	for i, rec := range records {
		input, err := assignmentsToCategoryCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &input
	}

	curr := func(c context.Context, args []*CategoryCreate) ([]*Category, error) {
		rowMaps := make([]map[string]any, len(args))
		for i, input := range args {
			rowMaps[i] = input.ToRowMap()
		}

		pkCols := []string{
			"id",
		}

		return executeCreateManyAndReturn(c, d.client, rowMaps, "Category", categoryDefaultCols, selects, omits,
			selectCategoryCols,
			func(ctx context.Context, txQ *Queries, results []*Category, sel *CategorySelect) error {
				return txQ.Category.loadRelations(ctx, results, sel)
			},
			(*Category).ScanFields,
			(*CategorySelect).hasAnyRelation,
			pkCols,
			conflictTarget,
			conflictAction,
		)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.CreateManyAndReturn != nil {
			next, hook := curr, ext.CreateManyAndReturn
			curr = func(c context.Context, inputs []*CategoryCreate) ([]*Category, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
}

type CategoryConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *CategoryConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *CategoryConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *CategoryConflictBuilder[B]) Update(fn func(u *CategoryUpsert)) *B {
	var up ConflictUpdate
	u := newCategoryUpsert(&up)
	fn(u)
	cb.setAction(ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}, cb.conflictTarget)
	return cb.builder
}

type CategoryUpsert struct {
	Id   numericFieldUpsert[int32]
	Name fieldUpsert[string]
}

func newCategoryUpsert(up *ConflictUpdate) *CategoryUpsert {
	return &CategoryUpsert{
		Id: numericFieldUpsert[int32]{
			fieldUpsert: fieldUpsert[int32]{column: "id", update: up},
			tableName:   "Category",
		},
		Name: fieldUpsert[string]{column: "name", update: up},
	}
}
func (d *CategoryDelegate) FindUnique(where UniquePredicate[Category], additional ...PredicateOf[Category]) *FindUniqueBuilder[Category, CategorySelect, CategoryOmit] {
	return &FindUniqueBuilder[Category, CategorySelect, CategoryOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *CategoryDelegate) FindFirst(preds ...PredicateOf[Category]) *FindFirstBuilder[Category, CategorySelect, CategoryOmit] {
	return &FindFirstBuilder[Category, CategorySelect, CategoryOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *CategoryDelegate) FindMany(preds ...PredicateOf[Category]) *FindManyBuilder[Category, CategorySelect, CategoryOmit] {
	return &FindManyBuilder[Category, CategorySelect, CategoryOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *CategoryDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[Category], additional []PredicateOf[Category], selects *CategorySelect, omits *CategoryOmit) (*Category, error) {
	curr := func(c context.Context, w UniquePredicate[Category], add []PredicateOf[Category], sel *CategorySelect, o *CategoryOmit) (*Category, error) {
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
		allPreds := append([]PredicateOf[Category]{w}, add...)
		whereClause, vals := CompilePredicates(d.client.dialect, allPreds)
		if whereClause != "" {
			whereClause = " WHERE " + whereClause
		}
		returningCols := selectCategoryCols(sel, o)
		return executeSingleWithRelations(c, d.client, "Category", whereClause, vals, returningCols,
			func(res *Category, cols []string) []any { return res.ScanFields(cols) },
			sel.hasAnyRelation(),
			func(ctx context.Context, txQ *Queries, results []*Category) error {
				return txQ.Category.loadRelations(ctx, results, sel)
			},
			nil,
		)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindUnique != nil {
			next, hook := curr, ext.FindUnique
			curr = func(c context.Context, w UniquePredicate[Category], add []PredicateOf[Category], sel *CategorySelect, o *CategoryOmit) (*Category, error) {
				return hook(c, w, add, sel, o, next)
			}
		}
	}

	return curr(ctx, where, additional, selects, omits)
}

func (d *CategoryDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[Category],
	selects *CategorySelect,
	omits *CategoryOmit,
) (*Category, error) {
	curr := func(c context.Context, p QueryParams[Category], sel *CategorySelect, o *CategoryOmit) (*Category, error) {
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
		returningCols := selectCategoryCols(sel, o)
		return executeSingleWithRelations(c, d.client, "Category", whereClause, vals, returningCols,
			func(res *Category, cols []string) []any { return res.ScanFields(cols) },
			sel.hasAnyRelation(),
			func(ctx context.Context, txQ *Queries, results []*Category) error {
				return txQ.Category.loadRelations(ctx, results, sel)
			},
			p.Skip,
		)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindFirst != nil {
			next, hook := curr, ext.FindFirst
			curr = func(c context.Context, p QueryParams[Category], sel *CategorySelect, o *CategoryOmit) (*Category, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}

func (d *CategoryDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[Category],
	selects *CategorySelect,
	omits *CategoryOmit,
) ([]*Category, error) {
	curr := func(c context.Context, p QueryParams[Category], sel *CategorySelect, o *CategoryOmit) ([]*Category, error) {
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
		returningCols := selectCategoryCols(sel, o)
		return executeManyWithRelations(c, d.client, "Category", whereClause, vals, returningCols,
			func(res *Category, cols []string) []any { return res.ScanFields(cols) },
			sel.hasAnyRelation(),
			func(ctx context.Context, txQ *Queries, results []*Category) error {
				return txQ.Category.loadRelations(ctx, results, sel)
			},
			p.Take,
			p.Skip,
		)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindMany != nil {
			next, hook := curr, ext.FindMany
			curr = func(c context.Context, p QueryParams[Category], sel *CategorySelect, o *CategoryOmit) ([]*Category, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}
func (d *CategoryDelegate) loadRelations(ctx context.Context, records []*Category, selects *CategorySelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Posts != nil {
		relationSelects, relationOmits, relationParams := selects.Posts.GetRelationParams()
		returningCols := selectCategoryToPostCols(relationSelects, relationOmits, "categoryId")
		// Inverse holds the FK: CategoryToPost.categoryId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *Category) int32 { return p.Id }),
			"CategoryToPost",
			"categoryId",
			returningCols,
			scanInto(returningCols, (*CategoryToPost).ScanFields),
			directKey(func(c *CategoryToPost) int32 { return c.CategoryId }),
			appendMany(func(p *Category) *[]*CategoryToPost { return &p.Posts }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading posts: %w", err)
		}
		if err := d.client.CategoryToPost.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
