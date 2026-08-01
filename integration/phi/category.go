package phi

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
)

// Category represents the database model
type Category struct {
	Id    int32             `db:"id" json:"id"`
	Name  string            `db:"name" json:"name"`
	Posts []*CategoryToPost `json:"posts,omitempty"`
}

// CategoryCreate contains model input fields for Category creation operations.
//
// Fields for Category:
//
//	id   int32  default: autoincrement()
//	name string required
type CategoryCreate struct {
	Id   *int32 `json:"id"`
	Name string `json:"name"`
}

// colMask returns a bit mask of columns that are set
func (s *CategoryCreate) colMask() uint64 {
	var mask uint64
	if s.Id != nil {
		mask |= 1 << 0
	}
	mask |= 1 << 1
	return mask
}

// CategoryUpdate contains model input fields for Category update operations.
type CategoryUpdate struct {
	Id   *int32  `json:"id"`
	Name *string `json:"name"`
}

func (u *CategoryUpdate) ToColsVals() ([]string, []any) {
	var cols []string
	var vals []any
	if u.Id != nil {
		cols = append(cols, "id")
		vals = append(vals, u.Id)
	}
	if u.Name != nil {
		cols = append(cols, "name")
		vals = append(vals, u.Name)
	}
	return cols, vals
}

func assignmentsToCategoryUpdate(assignments []FieldAssignment) (CategoryUpdate, error) {
	var input CategoryUpdate
	var errs ValidationError

	for _, a := range assignments {
		switch a.Col {
		case "id":
			if v, ok := a.Val.(int32); ok {
				input.Id = &v
			} else if v, ok := a.Val.(*int32); ok {
				input.Id = v
			} else {
				errs.Add("id", a.Val, "type", "field id must be of type int32")
			}
		case "name":
			if v, ok := a.Val.(string); ok {
				input.Name = &v
				errs.ValidateString("name", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.Name = v
			} else {
				errs.Add("name", a.Val, "type", "field name must be of type string")
			}
		}
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

// CategorySelect specifies which scalar and relation fields to select for Category.
//
// Selectable fields:
//
//	-- Scalars --
//	id   (bool)
//	name (bool)
//	-- Relations --
//	posts ([]CategoryToPost)
type CategorySelect struct {
	Id    bool                      `json:"id"`
	Name  bool                      `json:"name"`
	Posts CategoryToPostSelectQuery `json:"posts,omitempty"`
}

var fullCategorySelectVal = &CategorySelect{
	Id:   true,
	Name: true,
}

func fullCategorySelect() *CategorySelect {
	return fullCategorySelectVal
}

func (s *CategorySelect) hasAnyScalar() bool {
	if s == nil {
		return false
	}
	return s.Id || s.Name
}

func (s *CategorySelect) hasAnySelected() bool {
	if s == nil {
		return false
	}
	return s.hasAnyScalar() || s.hasAnyRelation()
}

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

type CategoryQueryBuilder struct {
	selects *CategorySelect
	omits   *CategoryOmit
	where   []PredicateOf[Category]
	take    *int
	skip    *int
	orderBy []OrderBy[Category]
	cursor  UniquePredicate[Category]
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

func (b *CategoryQueryBuilder) Cursor(where UniquePredicate[Category]) *CategoryQueryBuilder {
	b.cursor = where
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
		Cursor:  b.cursor,
	}
}

// CategoryCreateArgs is the input argument passed to Category Create extension hooks.
//
// Fields for Category:
//
//	id   int32  default: autoincrement()
//	name string required
//
// Relations for Category:
//
//	posts ([]CategoryToPost)
type CategoryCreateArgs struct {
	// Data contains the model fields to insert.
	Data *CategoryCreate
	// Select specifies which scalar and relation fields to select and return upon creation.
	Select *CategorySelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

// CategoryCreateManyArgs is the input argument passed to Category CreateMany extension hooks.
//
// Fields for Category:
//
//	id   int32  default: autoincrement()
//	name string required
type CategoryCreateManyArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*CategoryCreate
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *CategoryCreateManyArgs) AppendData(builders ...*CategoryCreateBuilder) *CategoryCreateManyArgs {
	for _, b := range builders {
		input, err := assignmentsToCategoryCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// CategoryCreateManyAndReturnArgs is the input argument passed to Category CreateManyAndReturn extension hooks.
//
// Fields for Category:
//
//	id   int32  default: autoincrement()
//	name string required
//
// Relations for Category:
//
//	posts ([]CategoryToPost)
type CategoryCreateManyAndReturnArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*CategoryCreate
	// Select specifies which scalar and relation fields to select and return for created records.
	Select *CategorySelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *CategoryCreateManyAndReturnArgs) AppendData(builders ...*CategoryCreateBuilder) *CategoryCreateManyAndReturnArgs {
	for _, b := range builders {
		input, err := assignmentsToCategoryCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// CategoryFindUniqueArgs is the input argument passed to Category FindUnique extension hooks.
type CategoryFindUniqueArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[Category]
	// Select specifies which scalar and relation fields to select and return.
	Select *CategorySelect
}

func (a *CategoryFindUniqueArgs) SetWhere(unique UniquePredicate[Category], additional ...PredicateOf[Category]) *CategoryFindUniqueArgs {
	a.Where = make([]PredicateOf[Category], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// CategoryFindFirstArgs is the input argument passed to Category FindFirst extension hooks.
type CategoryFindFirstArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Category]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[Category]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[Category]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *CategorySelect
}

func (a *CategoryFindFirstArgs) SetWhere(preds ...PredicateOf[Category]) *CategoryFindFirstArgs {
	a.Where = preds
	return a
}

func (a *CategoryFindFirstArgs) SetOrderBy(orders ...OrderBy[Category]) *CategoryFindFirstArgs {
	a.OrderBy = orders
	return a
}

func (a *CategoryFindFirstArgs) SetCursor(cursor UniquePredicate[Category]) *CategoryFindFirstArgs {
	a.Cursor = cursor
	return a
}

func (a *CategoryFindFirstArgs) SetSkip(n int) *CategoryFindFirstArgs {
	a.Skip = &n
	return a
}

func (a *CategoryFindFirstArgs) SetTake(n int) *CategoryFindFirstArgs {
	a.Take = &n
	return a
}

// CategoryFindManyArgs is the input argument passed to Category FindMany extension hooks.
type CategoryFindManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Category]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[Category]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[Category]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *CategorySelect
}

func (a *CategoryFindManyArgs) SetWhere(preds ...PredicateOf[Category]) *CategoryFindManyArgs {
	a.Where = preds
	return a
}

func (a *CategoryFindManyArgs) SetOrderBy(orders ...OrderBy[Category]) *CategoryFindManyArgs {
	a.OrderBy = orders
	return a
}

func (a *CategoryFindManyArgs) SetCursor(cursor UniquePredicate[Category]) *CategoryFindManyArgs {
	a.Cursor = cursor
	return a
}

func (a *CategoryFindManyArgs) SetSkip(n int) *CategoryFindManyArgs {
	a.Skip = &n
	return a
}

func (a *CategoryFindManyArgs) SetTake(n int) *CategoryFindManyArgs {
	a.Take = &n
	return a
}

// CategoryCountArgs is the input argument passed to Category Count extension hooks.
type CategoryCountArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Category]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
}

func (a *CategoryCountArgs) SetWhere(preds ...PredicateOf[Category]) *CategoryCountArgs {
	a.Where = preds
	return a
}

func (a *CategoryCountArgs) SetSkip(n int) *CategoryCountArgs {
	a.Skip = &n
	return a
}

func (a *CategoryCountArgs) SetTake(n int) *CategoryCountArgs {
	a.Take = &n
	return a
}

// CategoryDeleteArgs is the input argument passed to Category Delete extension hooks.
type CategoryDeleteArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[Category]
	// Select specifies which scalar and relation fields to select and return for deleted record.
	Select *CategorySelect
}

func (a *CategoryDeleteArgs) SetWhere(unique UniquePredicate[Category], additional ...PredicateOf[Category]) *CategoryDeleteArgs {
	a.Where = make([]PredicateOf[Category], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// CategoryDeleteManyArgs is the input argument passed to Category DeleteMany extension hooks.
type CategoryDeleteManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Category]
}

func (a *CategoryDeleteManyArgs) SetWhere(preds ...PredicateOf[Category]) *CategoryDeleteManyArgs {
	a.Where = preds
	return a
}

// CategoryUpdateArgs is the input argument passed to Category Update extension hooks.
type CategoryUpdateArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[Category]
	// Data contains the model fields to update.
	Data *CategoryUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *CategorySelect
}

func (a *CategoryUpdateArgs) SetWhere(unique UniquePredicate[Category], additional ...PredicateOf[Category]) *CategoryUpdateArgs {
	a.Where = make([]PredicateOf[Category], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// CategoryUpdateManyArgs is the input argument passed to Category UpdateMany extension hooks.
type CategoryUpdateManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Category]
	// Data contains the model fields to update.
	Data *CategoryUpdate
}

func (a *CategoryUpdateManyArgs) SetWhere(preds ...PredicateOf[Category]) *CategoryUpdateManyArgs {
	a.Where = preds
	return a
}

// CategoryUpdateManyAndReturnArgs is the input argument passed to Category UpdateManyAndReturn extension hooks.
type CategoryUpdateManyAndReturnArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[Category]
	// Data contains the model fields to update.
	Data *CategoryUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *CategorySelect
}

func (a *CategoryUpdateManyAndReturnArgs) SetWhere(preds ...PredicateOf[Category]) *CategoryUpdateManyAndReturnArgs {
	a.Where = preds
	return a
}

type CategoryCreateQuery = func(ctx context.Context, args *CategoryCreateArgs) (*Category, error)
type CategoryCreateManyQuery = func(ctx context.Context, args *CategoryCreateManyArgs) (int64, error)
type CategoryCreateManyAndReturnQuery = func(ctx context.Context, args *CategoryCreateManyAndReturnArgs) ([]*Category, error)
type CategoryFindUniqueQuery = func(ctx context.Context, args *CategoryFindUniqueArgs) (*Category, error)
type CategoryFindFirstQuery = func(ctx context.Context, args *CategoryFindFirstArgs) (*Category, error)
type CategoryFindManyQuery = func(ctx context.Context, args *CategoryFindManyArgs) ([]*Category, error)
type CategoryDeleteQuery = func(ctx context.Context, args *CategoryDeleteArgs) (*Category, error)
type CategoryDeleteManyQuery = func(ctx context.Context, args *CategoryDeleteManyArgs) (int64, error)
type CategoryCountQuery = func(ctx context.Context, args *CategoryCountArgs) (int64, error)
type CategoryUpdateQuery = func(ctx context.Context, args *CategoryUpdateArgs) (*Category, error)
type CategoryUpdateManyQuery = func(ctx context.Context, args *CategoryUpdateManyArgs) (int64, error)
type CategoryUpdateManyAndReturnQuery = func(ctx context.Context, args *CategoryUpdateManyAndReturnArgs) ([]*Category, error)

type CategoryExtension struct {
	Create              func(ctx context.Context, args *CategoryCreateArgs, next CategoryCreateQuery) (*Category, error)
	CreateMany          func(ctx context.Context, args *CategoryCreateManyArgs, next CategoryCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, args *CategoryCreateManyAndReturnArgs, next CategoryCreateManyAndReturnQuery) ([]*Category, error)
	FindUnique          func(ctx context.Context, args *CategoryFindUniqueArgs, next CategoryFindUniqueQuery) (*Category, error)
	FindFirst           func(ctx context.Context, args *CategoryFindFirstArgs, next CategoryFindFirstQuery) (*Category, error)
	FindMany            func(ctx context.Context, args *CategoryFindManyArgs, next CategoryFindManyQuery) ([]*Category, error)
	Delete              func(ctx context.Context, args *CategoryDeleteArgs, next CategoryDeleteQuery) (*Category, error)
	DeleteMany          func(ctx context.Context, args *CategoryDeleteManyArgs, next CategoryDeleteManyQuery) (int64, error)
	Count               func(ctx context.Context, args *CategoryCountArgs, next CategoryCountQuery) (int64, error)
	Update              func(ctx context.Context, args *CategoryUpdateArgs, next CategoryUpdateQuery) (*Category, error)
	UpdateMany          func(ctx context.Context, args *CategoryUpdateManyArgs, next CategoryUpdateManyQuery) (int64, error)
	UpdateManyAndReturn func(ctx context.Context, args *CategoryUpdateManyAndReturnArgs, next CategoryUpdateManyAndReturnQuery) ([]*Category, error)
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

var categoryPKCols = []string{
	"id",
}

var categoryUniqueCols = []string{
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

func (b *CategoryCreateBuilder) Select(s CategorySelect) *CategoryCreateBuilder {
	b.selects = &s
	return b
}

func (b *CategoryCreateBuilder) Omit(o CategoryOmit) *CategoryCreateBuilder {
	b.omits = &o
	return b
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

func (b *CategoryCreateBuilder) Assignments(assignments ...FieldAssignmentOf[Category]) *CategoryCreateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *CategoryDelegate) Create() *CategoryCreateBuilder {
	return &CategoryCreateBuilder{
		CreateBuilder: &CreateBuilder[Category, CategorySelect, CategoryOmit]{
			execFunc: d.executeCreate,
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
				errs.ValidateInt32("id", v, "")
			} else {
				errs.Add("id", a.Val, "type", "field id must be of type int32")
			}
		case "name":
			provided |= providedCategoryName
			if v, ok := a.Val.(string); ok {
				input.Name = v
				errs.ValidateString("name", v, true, 0, false, false)
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

func partitionCategoryInputs(dialect Dialect, inputs []*CategoryCreate) [][]*CategoryCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*CategoryCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*CategoryCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*CategoryCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*CategoryCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*CategoryCreate{inputs}
}

func (d *CategoryDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *CategorySelect, omits *CategoryOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*Category, error) {
	input, err := assignmentsToCategoryCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectCategoryCols(selects, omits)

	if len(d.extensions) == 0 {
		return d.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullCategorySelect()
	}

	args := &CategoryCreateArgs{
		Data:           &input,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *CategoryCreateArgs) (*Category, error) {
		cCols, cVals := a.Data.ToColsVals()
		cReturningCols := selectCategoryCols(a.Select, omits)
		return d.runCreate(c, cCols, cVals, cReturningCols, a.Select, a.ConflictTarget, a.ConflictAction)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.Create != nil {
			return ext.Create(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.Create != nil {
			next, hook := curr, ext.Create
			curr = func(c context.Context, a *CategoryCreateArgs) (*Category, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
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

func (b *CategoryCreateManyAndReturnBuilder) Select(s CategorySelect) *CategoryCreateManyAndReturnBuilder {
	b.selects = &s
	return b
}

func (b *CategoryCreateManyAndReturnBuilder) Omit(o CategoryOmit) *CategoryCreateManyAndReturnBuilder {
	b.omits = &o
	return b
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

func createBuildersToCategoryRecordInputs(builders []*CategoryCreateBuilder) []RecordInput {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return records
}

func (d *CategoryDelegate) CreateMany(builders ...*CategoryCreateBuilder) *CategoryCreateManyBuilder {
	return &CategoryCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[Category]{
			records:  createBuildersToCategoryRecordInputs(builders),
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *CategoryDelegate) CreateManyAndReturn(builders ...*CategoryCreateBuilder) *CategoryCreateManyAndReturnBuilder {
	return &CategoryCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[Category, CategorySelect, CategoryOmit]{
			records:  createBuildersToCategoryRecordInputs(builders),
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func recordsToCategoryCreateInputs(records []RecordInput) ([]*CategoryCreate, error) {
	structs := make([]CategoryCreate, len(records))
	inputs := make([]*CategoryCreate, len(records))
	for i, rec := range records {
		var err error
		structs[i], err = assignmentsToCategoryCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &structs[i]
	}
	return inputs, nil
}

func (d *CategoryDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs, err := recordsToCategoryCreateInputs(records)
	if err != nil {
		return 0, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	args := &CategoryCreateManyArgs{
		Data:           inputs,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *CategoryCreateManyArgs) (int64, error) {
		return d.runCreateMany(c, a.Data, a.ConflictTarget, a.ConflictAction)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.CreateMany != nil {
			return ext.CreateMany(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.CreateMany != nil {
			next, hook := curr, ext.CreateMany
			curr = func(c context.Context, a *CategoryCreateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *CategoryDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *CategorySelect, omits *CategoryOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*Category, error) {
	inputs, err := recordsToCategoryCreateInputs(records)
	if err != nil {
		return nil, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullCategorySelect()
	}

	args := &CategoryCreateManyAndReturnArgs{
		Data:           inputs,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *CategoryCreateManyAndReturnArgs) ([]*Category, error) {
		return d.runCreateManyAndReturn(c, a.Data, a.Select, omits, a.ConflictTarget, a.ConflictAction)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.CreateManyAndReturn != nil {
			return ext.CreateManyAndReturn(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.CreateManyAndReturn != nil {
			next, hook := curr, ext.CreateManyAndReturn
			curr = func(c context.Context, a *CategoryCreateManyAndReturnArgs) ([]*Category, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *CategoryDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	selects *CategorySelect,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*Category, error) {
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := hasRelations && !d.client.inTx()

	if useTx {
		var res *Category
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Category.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
			if err != nil {
				return err
			}
			return txQ.Category.loadRelations(ctx, []*Category{res}, selects)
		})
		return res, err
	}

	query, clauseArgs := buildSingleInsertSQL(d.client, "Category", cols, returningCols, categoryPKCols, conflictTarget, conflictAction, len(vals))
	if len(clauseArgs) > 0 {
		vals = append(vals, clauseArgs...)
	}

	if d.client.dialect.SupportsInsertReturning {
		rows, err := d.client.query(ctx, query, vals...)
		if err != nil {
			return nil, err
		}

		if !rows.Next() {
			err := rows.Err()
			rows.Close()
			if err != nil {
				return nil, err
			}
			return nil, nil
		}

		var res Category
		scanErr := rows.Scan(res.ScanFields(returningCols)...)
		rows.Close()
		if scanErr != nil {
			return nil, scanErr
		}

		return &res, nil
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, categoryPKCols)
}

func (d *CategoryDelegate) runCreateFallback(
	ctx context.Context,
	query string,
	vals []any,
	cols []string,
	returningCols []string,
	pkCols []string,
) (*Category, error) {
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
	selectSb.Grow(64 + len(returningCols)*15 + len("Category") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "Category")
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

	if !rows.Next() {
		err := rows.Err()
		rows.Close()
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	var res Category
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, scanErr
	}

	return &res, nil
}

func (d *CategoryDelegate) buildBulkInsertSQL(q *Queries, batch []*CategoryCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 2)
	for i, c := range categoryDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "Category")
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
			case "id":
				if input.Id != nil {
					vals = append(vals, *input.Id)
				} else {
					writeDefault = true
				}
			case "name":
				vals = append(vals, input.Name)
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

func applyCategoryConflictClause(dialect Dialect, queryStr string, vals []any, cols []string, pkCols []string, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (string, []any) {
	var conflictCols []string
	if conflictTarget != nil {
		conflictCols = conflictTarget.UniqueColumns()
	}
	var nonConflictCols []string
	if conflictAction != nil && conflictAction.Type == ConflictActionUpdateNewValues {
		nonConflictCols = computeNonConflictCols(cols, conflictCols, pkCols)
	}
	clause, clauseArgs := dialect.BuildConflictClause(conflictCols, conflictAction, nonConflictCols, len(vals)+1)
	queryStr += clause
	if len(clauseArgs) > 0 {
		vals = append(vals, clauseArgs...)
	}
	return queryStr, vals
}

func scanCategoryRows(rows *sql.Rows, returningCols []string) ([]*Category, error) {
	var records []*Category
	for rows.Next() {
		var res Category
		if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
			rows.Close()
			return nil, err
		}
		records = append(records, &res)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()
	return records, nil
}

func (d *CategoryDelegate) runCreateMany(ctx context.Context, inputs []*CategoryCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionCategoryInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyCategoryConflictClause(d.client.dialect, queryStr, vals, cols, categoryPKCols, conflictTarget, conflictAction)

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

func (d *CategoryDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*CategoryCreate,
	selects *CategorySelect,
	omits *CategoryOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*Category, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionCategoryInputs(d.client.dialect, inputs)
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning) && !d.client.inTx()

	if useTx {
		var res []*Category
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if txQ.dialect.SupportsInsertReturning {
				res, err = txQ.Category.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			} else {
				res, err = txQ.Category.runCreateManyAndReturnFallback(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			}
			if err != nil {
				return err
			}
			if hasRelations {
				return txQ.Category.loadRelations(ctx, res, selects)
			}
			return nil
		})
		return res, err
	}

	returningCols := selectCategoryCols(selects, omits, categoryPKCols...)
	recordsOut := make([]*Category, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyCategoryConflictClause(d.client.dialect, queryStr, vals, cols, categoryPKCols, conflictTarget, conflictAction)

		if len(returningCols) > 0 {
			var retSb strings.Builder
			retSb.Grow(12 + len(returningCols)*15)
			retSb.WriteString(" RETURNING ")
			for i, col := range returningCols {
				if i > 0 {
					retSb.WriteString(", ")
				}
				d.client.dialect.WriteQuotedIdent(&retSb, col)
			}
			queryStr += retSb.String()
		}

		rows, err := d.client.query(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		scanned, err := scanCategoryRows(rows, returningCols)
		if err != nil {
			return nil, err
		}
		recordsOut = append(recordsOut, scanned...)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, recordsOut, selects); err != nil {
			return nil, err
		}
	}

	return recordsOut, nil
}

func (d *CategoryDelegate) runCreateManyAndReturnFallback(
	ctx context.Context,
	inputs []*CategoryCreate,
	selects *CategorySelect,
	omits *CategoryOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*Category, error) {
	batches := partitionCategoryInputs(d.client.dialect, inputs)
	returningCols := selectCategoryCols(selects, omits, categoryPKCols...)
	recordsOut := make([]*Category, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyCategoryConflictClause(d.client.dialect, queryStr, vals, cols, categoryPKCols, conflictTarget, conflictAction)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("Category") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		d.client.dialect.WriteQuotedIdent(&selectSb, "Category")
		selectSb.WriteString(" WHERE ")
		d.client.dialect.WriteQuotedIdent(&selectSb, categoryPKCols[0])
		selectSb.WriteString(" >= ")
		d.client.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		d.client.dialect.WriteQuotedIdent(&selectSb, categoryPKCols[0])
		selectSb.WriteString(" < ")
		d.client.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := d.client.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return nil, err
		}

		scanned, err := scanCategoryRows(rows, returningCols)
		if err != nil {
			return nil, err
		}
		recordsOut = append(recordsOut, scanned...)
	}

	return recordsOut, nil
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
	cb.setAction(*CategoryConflictUpdate(fn), cb.conflictTarget)
	return cb.builder
}

// CategoryConflictUpdate creates a custom ConflictAction for Category upsert conflicts.
func CategoryConflictUpdate(fn func(u *CategoryUpsert)) *ConflictAction {
	var up ConflictUpdate
	u := newCategoryUpsert(&up)
	fn(u)
	return &ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}
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

type CategoryUpdateBuilder struct {
	*UpdateBuilder[Category, CategorySelect, CategoryOmit]
}

type CategoryUpdateManyBuilder struct {
	*UpdateManyBuilder[Category]
}

type CategoryUpdateManyAndReturnBuilder struct {
	*UpdateManyAndReturnBuilder[Category, CategorySelect, CategoryOmit]
}

func (b *CategoryUpdateBuilder) SetId(v int32) *CategoryUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}

func (b *CategoryUpdateManyBuilder) SetId(v int32) *CategoryUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}

func (b *CategoryUpdateManyAndReturnBuilder) SetId(v int32) *CategoryUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}
func (b *CategoryUpdateBuilder) SetName(v string) *CategoryUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}

func (b *CategoryUpdateManyBuilder) SetName(v string) *CategoryUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}

func (b *CategoryUpdateManyAndReturnBuilder) SetName(v string) *CategoryUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}

func (b *CategoryUpdateBuilder) Assignments(assignments ...FieldAssignmentOf[Category]) *CategoryUpdateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *CategoryUpdateManyBuilder) Assignments(assignments ...FieldAssignmentOf[Category]) *CategoryUpdateManyBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (b *CategoryUpdateManyAndReturnBuilder) Assignments(assignments ...FieldAssignmentOf[Category]) *CategoryUpdateManyAndReturnBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment(a))
	}
	return b
}

func (d *CategoryDelegate) Update(where UniquePredicate[Category], additional ...PredicateOf[Category]) *CategoryUpdateBuilder {
	return &CategoryUpdateBuilder{
		UpdateBuilder: &UpdateBuilder[Category, CategorySelect, CategoryOmit]{
			where:      where,
			additional: additional,
			execFunc:   d.executeUpdate,
		},
	}
}

func (d *CategoryDelegate) UpdateMany(preds ...PredicateOf[Category]) *CategoryUpdateManyBuilder {
	return &CategoryUpdateManyBuilder{
		UpdateManyBuilder: &UpdateManyBuilder[Category]{
			where:    preds,
			execFunc: d.executeUpdateMany,
		},
	}
}

func (d *CategoryDelegate) UpdateManyAndReturn(preds ...PredicateOf[Category]) *CategoryUpdateManyAndReturnBuilder {
	return &CategoryUpdateManyAndReturnBuilder{
		UpdateManyAndReturnBuilder: &UpdateManyAndReturnBuilder[Category, CategorySelect, CategoryOmit]{
			where:    preds,
			execFunc: d.executeUpdateManyAndReturn,
		},
	}
}

func (d *CategoryDelegate) buildUpdateSQL(preds []PredicateOf[Category], cols []string, vals []any, returningCols []string) (string, []any) {
	whereClause, predVals, _ := CompilePredicates(d.client.dialect, preds, len(cols)+1)

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	d.client.dialect.WriteQuotedIdent(&sb, "Category")
	sb.WriteString(" SET ")

	setVals := make([]any, 0, len(cols)+len(predVals))
	for i, col := range cols {
		if i > 0 {
			sb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&sb, col)
		sb.WriteString(" = ")
		d.client.dialect.WritePlaceholder(&sb, i+1)
		setVals = append(setVals, vals[i])
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

func (d *CategoryDelegate) executeUpdate(ctx context.Context, where UniquePredicate[Category], additional []PredicateOf[Category], assignments []FieldAssignment, selects *CategorySelect, omits *CategoryOmit) (*Category, error) {
	allWhere := make([]PredicateOf[Category], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	input, err := assignmentsToCategoryUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdate(ctx, allWhere, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullCategorySelect()
	}

	args := &CategoryUpdateArgs{
		Where:  allWhere,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *CategoryUpdateArgs) (*Category, error) {
		extCols, extVals := a.Data.ToColsVals()
		return d.runUpdate(c, a.Where, extCols, extVals, a.Select, omits)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.Update != nil {
			return ext.Update(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.Update != nil {
			next, hook := curr, ext.Update
			curr = func(c context.Context, a *CategoryUpdateArgs) (*Category, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *CategoryDelegate) runUpdate(ctx context.Context, preds []PredicateOf[Category], cols []string, vals []any, selects *CategorySelect, omits *CategoryOmit) (*Category, error) {
	if len(cols) == 0 {
		return d.runFindUnique(ctx, preds, selects, omits)
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
		var res *Category
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.Category.runUpdate(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.Category.runUpdateFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectCategoryCols(selects, omits, categoryPKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

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

	var res Category
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, scanErr
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*Category{&res}, selects); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *CategoryDelegate) execUpdateStmt(ctx context.Context, preds []PredicateOf[Category], cols []string, vals []any) (int64, error) {
	if len(cols) == 0 {
		return 0, nil
	}

	for _, pr := range preds {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return 0, err
			}
		}
	}

	query, setVals := d.buildUpdateSQL(preds, cols, vals, nil)
	result, err := d.client.exec(ctx, query, setVals...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (d *CategoryDelegate) runUpdateFallback(ctx context.Context, preds []PredicateOf[Category], cols []string, vals []any, selects *CategorySelect, omits *CategoryOmit) (*Category, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, sql.ErrNoRows
	}
	return d.runFindUnique(ctx, preds, selects, omits)
}

// -----------------------------------------------------------------------------
// UpdateMany
// -----------------------------------------------------------------------------

func (d *CategoryDelegate) executeUpdateMany(ctx context.Context, preds []PredicateOf[Category], assignments []FieldAssignment) (int64, error) {
	input, err := assignmentsToCategoryUpdate(assignments)
	if err != nil {
		return 0, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.execUpdateStmt(ctx, preds, cols, vals)
	}

	args := &CategoryUpdateManyArgs{
		Where: preds,
		Data:  &input,
	}

	curr := func(c context.Context, a *CategoryUpdateManyArgs) (int64, error) {
		extCols, extVals := a.Data.ToColsVals()
		return d.execUpdateStmt(c, a.Where, extCols, extVals)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.UpdateMany != nil {
			return ext.UpdateMany(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.UpdateMany != nil {
			next, hook := curr, ext.UpdateMany
			curr = func(c context.Context, a *CategoryUpdateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

// -----------------------------------------------------------------------------
// UpdateManyAndReturn
// -----------------------------------------------------------------------------

func (d *CategoryDelegate) executeUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[Category], assignments []FieldAssignment, selects *CategorySelect, omits *CategoryOmit) ([]*Category, error) {
	input, err := assignmentsToCategoryUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullCategorySelect()
	}

	args := &CategoryUpdateManyAndReturnArgs{
		Where:  preds,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *CategoryUpdateManyAndReturnArgs) ([]*Category, error) {
		extCols, extVals := a.Data.ToColsVals()
		return d.runUpdateManyAndReturn(c, a.Where, extCols, extVals, a.Select, omits)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.UpdateManyAndReturn != nil {
			return ext.UpdateManyAndReturn(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.UpdateManyAndReturn != nil {
			next, hook := curr, ext.UpdateManyAndReturn
			curr = func(c context.Context, a *CategoryUpdateManyAndReturnArgs) ([]*Category, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *CategoryDelegate) runUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[Category], cols []string, vals []any, selects *CategorySelect, omits *CategoryOmit) ([]*Category, error) {
	if len(cols) == 0 {
		return d.runFindMany(ctx, QueryParams[Category]{Where: preds}, selects, omits)
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
		var res []*Category
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.Category.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.Category.runUpdateManyAndReturnFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectCategoryCols(selects, omits, categoryPKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	scanned, err := scanCategoryRows(rows, returningCols)
	if err != nil {
		return nil, err
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, scanned, selects); err != nil {
			return nil, err
		}
	}

	return scanned, nil
}

func (d *CategoryDelegate) runUpdateManyAndReturnFallback(ctx context.Context, preds []PredicateOf[Category], cols []string, vals []any, selects *CategorySelect, omits *CategoryOmit) ([]*Category, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return []*Category{}, nil
	}
	return d.runFindMany(ctx, QueryParams[Category]{Where: preds}, selects, omits)
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
	allWhere := make([]PredicateOf[Category], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullCategorySelect()
	}

	args := &CategoryFindUniqueArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *CategoryFindUniqueArgs) (*Category, error) {
		return d.runFindUnique(c, a.Where, a.Select, omits)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.FindUnique != nil {
			return ext.FindUnique(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.FindUnique != nil {
			next, hook := curr, ext.FindUnique
			curr = func(c context.Context, a *CategoryFindUniqueArgs) (*Category, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *CategoryDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[Category],
	selects *CategorySelect,
	omits *CategoryOmit,
) (*Category, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullCategorySelect()
	}

	args := &CategoryFindFirstArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *CategoryFindFirstArgs) (*Category, error) {
		return d.runFindFirst(c, QueryParams[Category]{
			Where:   a.Where,
			OrderBy: a.OrderBy,
			Cursor:  a.Cursor,
			Skip:    a.Skip,
			Take:    a.Take,
		}, a.Select, omits)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.FindFirst != nil {
			return ext.FindFirst(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.FindFirst != nil {
			next, hook := curr, ext.FindFirst
			curr = func(c context.Context, a *CategoryFindFirstArgs) (*Category, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *CategoryDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[Category],
	selects *CategorySelect,
	omits *CategoryOmit,
) ([]*Category, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullCategorySelect()
	}

	args := &CategoryFindManyArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *CategoryFindManyArgs) ([]*Category, error) {
		return d.runFindMany(c, QueryParams[Category]{
			Where:   a.Where,
			OrderBy: a.OrderBy,
			Cursor:  a.Cursor,
			Skip:    a.Skip,
			Take:    a.Take,
		}, a.Select, omits)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.FindMany != nil {
			return ext.FindMany(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.FindMany != nil {
			next, hook := curr, ext.FindMany
			curr = func(c context.Context, a *CategoryFindManyArgs) ([]*Category, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *CategoryDelegate) runFindUnique(ctx context.Context, where []PredicateOf[Category], selects *CategorySelect, omits *CategoryOmit) (*Category, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}
	whereClause, vals, _ := CompilePredicates(d.client.dialect, where)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	returningCols := selectCategoryCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*Category{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *CategoryDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[Category],
	selects *CategorySelect,
	omits *CategoryOmit,
) (*Category, error) {
	one := 1
	negOne := -1
	if params.Take != nil && *params.Take < 0 {
		params.Take = &negOne
	} else {
		params.Take = &one
	}

	results, err := d.runFindMany(ctx, params, selects, omits)
	if err != nil || len(results) == 0 {
		return nil, err
	}
	return results[0], nil
}

func (d *CategoryDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[Category],
	selects *CategorySelect,
	omits *CategoryOmit,
) ([]*Category, error) {
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
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, categoryPKCols, categoryUniqueCols, "Category", nextIdx, params.Take)
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
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, categoryPKCols, categoryUniqueCols, isCursorQuery, params.Take)
	returningCols := selectCategoryCols(selects, omits)

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

func (d *CategoryDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*Category, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "Category", returningCols, whereClause, orderByClause, &limitOne, skip)
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

	var res Category
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (d *CategoryDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*Category, error) {
	query := buildSelectSQL(d.client, "Category", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*Category, 0)
	for rows.Next() {
		var res Category
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
func (d *CategoryDelegate) DeleteMany(preds ...PredicateOf[Category]) *DeleteManyBuilder[Category] {
	return &DeleteManyBuilder[Category]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *CategoryDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[Category]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	args := &CategoryDeleteManyArgs{
		Where: preds,
	}

	curr := func(c context.Context, a *CategoryDeleteManyArgs) (int64, error) {
		return d.runDeleteMany(c, a.Where)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.DeleteMany != nil {
			return ext.DeleteMany(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.DeleteMany != nil {
			next, hook := curr, ext.DeleteMany
			curr = func(c context.Context, a *CategoryDeleteManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *CategoryDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[Category]) (int64, error) {
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
	d.client.dialect.WriteQuotedIdent(&sb, "Category")
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

func (d *CategoryDelegate) Delete(where UniquePredicate[Category], additional ...PredicateOf[Category]) *DeleteBuilder[Category, CategorySelect, CategoryOmit] {
	return &DeleteBuilder[Category, CategorySelect, CategoryOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeDelete,
	}
}

func (d *CategoryDelegate) executeDelete(ctx context.Context, where UniquePredicate[Category], additional []PredicateOf[Category], selects *CategorySelect, omits *CategoryOmit) (*Category, error) {
	allWhere := make([]PredicateOf[Category], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runDelete(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullCategorySelect()
	}

	args := &CategoryDeleteArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *CategoryDeleteArgs) (*Category, error) {
		return d.runDelete(c, a.Where, a.Select, omits)
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.Delete != nil {
			return ext.Delete(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.Delete != nil {
			next, hook := curr, ext.Delete
			curr = func(c context.Context, a *CategoryDeleteArgs) (*Category, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *CategoryDelegate) runDelete(ctx context.Context, where []PredicateOf[Category], selects *CategorySelect, omits *CategoryOmit) (*Category, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}

	returningCols := selectCategoryCols(selects, omits, categoryPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *Category
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Category.runFindUnique(ctx, where, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return sql.ErrNoRows
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "Category")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[Category]
			pkPreds = append(pkPreds, Predicate[Category]{
				Data: PredicateData{
					Column:   "id",
					Operator: "=",
					Value:    res.Id,
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
	d.client.dialect.WriteQuotedIdent(&sb, "Category")

	whereClause, vals, _ := CompilePredicates(d.client.dialect, where)
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

	var row Category
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, err
	}
	return &row, nil
}
func (d *CategoryDelegate) Count(preds ...PredicateOf[Category]) *CountBuilder[Category] {
	return &CountBuilder[Category]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *CategoryDelegate) executeCount(ctx context.Context, params QueryParams[Category]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	args := &CategoryCountArgs{
		Where: params.Where,
		Skip:  params.Skip,
		Take:  params.Take,
	}

	curr := func(c context.Context, a *CategoryCountArgs) (int64, error) {
		return d.runCount(c, QueryParams[Category]{
			Where: a.Where,
			Skip:  a.Skip,
			Take:  a.Take,
		})
	}

	if len(d.extensions) == 1 {
		if ext := d.extensions[0]; ext.Count != nil {
			return ext.Count(ctx, args, curr)
		}
	}

	for i := len(d.extensions) - 1; i >= 0; i-- {
		ext := d.extensions[i]
		if ext.Count != nil {
			next, hook := curr, ext.Count
			curr = func(c context.Context, a *CategoryCountArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *CategoryDelegate) runCount(ctx context.Context, params QueryParams[Category]) (int64, error) {
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
		d.client.dialect.WriteQuotedIdent(&subQuery, "Category")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "Category")
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
