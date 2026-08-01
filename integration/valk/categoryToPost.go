package valk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// CategoryToPost represents the database model
type CategoryToPost struct {
	PostId     string    `db:"postId" json:"postId"`
	CategoryId int32     `db:"categoryId" json:"categoryId"`
	Post       *Post     `json:"post,omitempty"`
	Category   *Category `json:"category,omitempty"`
}

// CategoryToPostCreate contains model input fields for CategoryToPost creation operations.
//
// Fields for CategoryToPost:
//
//	postId     string required
//	categoryId int32  required
type CategoryToPostCreate struct {
	PostId     string `json:"postId"`
	CategoryId int32  `json:"categoryId"`
}

// colMask returns a bit mask of columns that are set
func (s *CategoryToPostCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	mask |= 1 << 1
	return mask
}

// CategoryToPostUpdate contains model input fields for CategoryToPost update operations.
type CategoryToPostUpdate struct {
	PostId     *string `json:"postId"`
	CategoryId *int32  `json:"categoryId"`
}

func (u *CategoryToPostUpdate) ToColsVals() ([]string, []any) {
	var cols []string
	var vals []any
	if u.PostId != nil {
		cols = append(cols, "postId")
		vals = append(vals, u.PostId)
	}
	if u.CategoryId != nil {
		cols = append(cols, "categoryId")
		vals = append(vals, u.CategoryId)
	}
	return cols, vals
}

func assignmentsToCategoryToPostUpdate(assignments []FieldAssignment) (CategoryToPostUpdate, error) {
	var input CategoryToPostUpdate
	var errs ValidationError

	for _, a := range assignments {
		switch a.Col {
		case "postId":
			if v, ok := a.Val.(string); ok {
				input.PostId = &v
				errs.ValidateString("postId", v, false, 0, false, false)
			} else if v, ok := a.Val.(*string); ok {
				input.PostId = v
			} else {
				errs.Add("postId", a.Val, "type", "field postId must be of type string")
			}
		case "categoryId":
			if v, ok := a.Val.(int32); ok {
				input.CategoryId = &v
			} else if v, ok := a.Val.(*int32); ok {
				input.CategoryId = v
			} else {
				errs.Add("categoryId", a.Val, "type", "field categoryId must be of type int32")
			}
		}
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

// CategoryToPostSelect specifies which scalar and relation fields to select for CategoryToPost.
//
// Selectable fields:
//
//	-- Scalars --
//	postId     (bool)
//	categoryId (bool)
//	-- Relations --
//	post       (Post)
//	category   (Category)
type CategoryToPostSelect struct {
	PostId     bool            `json:"postId"`
	CategoryId bool            `json:"categoryId"`
	Post       *PostSelect     `json:"post,omitempty"`
	Category   *CategorySelect `json:"category,omitempty"`
}

var fullCategoryToPostSelectVal = &CategoryToPostSelect{
	PostId:     true,
	CategoryId: true,
}

func fullCategoryToPostSelect() *CategoryToPostSelect {
	return fullCategoryToPostSelectVal
}

func (s *CategoryToPostSelect) hasAnyScalar() bool {
	if s == nil {
		return false
	}
	return s.PostId || s.CategoryId
}

func (s *CategoryToPostSelect) hasAnySelected() bool {
	if s == nil {
		return false
	}
	return s.hasAnyScalar() || s.hasAnyRelation()
}

type CategoryToPostOmit struct {
	PostId     bool `json:"postId"`
	CategoryId bool `json:"categoryId"`
}

type CategoryToPostSelectQuery interface {
	GetRelationParams() (*CategoryToPostSelect, *CategoryToPostOmit, QueryParams[CategoryToPost])
}

func (s *CategoryToPostSelect) GetRelationParams() (*CategoryToPostSelect, *CategoryToPostOmit, QueryParams[CategoryToPost]) {
	return s, nil, QueryParams[CategoryToPost]{}
}

type CategoryToPostQueryBuilder struct {
	selects *CategoryToPostSelect
	omits   *CategoryToPostOmit
	where   []PredicateOf[CategoryToPost]
	take    *int
	skip    *int
	orderBy []OrderBy[CategoryToPost]
	cursor  UniquePredicate[CategoryToPost]
}

func (b *CategoryToPostQueryBuilder) Where(preds ...PredicateOf[CategoryToPost]) *CategoryToPostQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *CategoryToPostQueryBuilder) Take(limit int) *CategoryToPostQueryBuilder {
	b.take = &limit
	return b
}

func (b *CategoryToPostQueryBuilder) Skip(offset int) *CategoryToPostQueryBuilder {
	b.skip = &offset
	return b
}

func (b *CategoryToPostQueryBuilder) OrderBy(orders ...OrderBy[CategoryToPost]) *CategoryToPostQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *CategoryToPostQueryBuilder) Cursor(where UniquePredicate[CategoryToPost]) *CategoryToPostQueryBuilder {
	b.cursor = where
	return b
}

func (b *CategoryToPostQueryBuilder) Select(s CategoryToPostSelect) *CategoryToPostQueryBuilder {
	b.selects = &s
	return b
}

func (b *CategoryToPostQueryBuilder) Omit(o CategoryToPostOmit) *CategoryToPostQueryBuilder {
	b.omits = &o
	return b
}

func (b *CategoryToPostQueryBuilder) GetRelationParams() (*CategoryToPostSelect, *CategoryToPostOmit, QueryParams[CategoryToPost]) {
	if b == nil {
		return nil, nil, QueryParams[CategoryToPost]{}
	}
	return b.selects, b.omits, QueryParams[CategoryToPost]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
		Cursor:  b.cursor,
	}
}

// CategoryToPostCreateArgs is the input argument passed to CategoryToPost Create extension hooks.
//
// Fields for CategoryToPost:
//
//	postId     string required
//	categoryId int32  required
//
// Relations for CategoryToPost:
//
//	post       (Post)
//	category   (Category)
type CategoryToPostCreateArgs struct {
	// Data contains the model fields to insert.
	Data *CategoryToPostCreate
	// Select specifies which scalar and relation fields to select and return upon creation.
	Select *CategoryToPostSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

// CategoryToPostCreateManyArgs is the input argument passed to CategoryToPost CreateMany extension hooks.
//
// Fields for CategoryToPost:
//
//	postId     string required
//	categoryId int32  required
type CategoryToPostCreateManyArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*CategoryToPostCreate
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *CategoryToPostCreateManyArgs) AppendData(builders ...*CategoryToPostCreateBuilder) *CategoryToPostCreateManyArgs {
	for _, b := range builders {
		input, err := assignmentsToCategoryToPostCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// CategoryToPostCreateManyAndReturnArgs is the input argument passed to CategoryToPost CreateManyAndReturn extension hooks.
//
// Fields for CategoryToPost:
//
//	postId     string required
//	categoryId int32  required
//
// Relations for CategoryToPost:
//
//	post       (Post)
//	category   (Category)
type CategoryToPostCreateManyAndReturnArgs struct {
	// Data is the slice of model inputs to bulk insert.
	Data []*CategoryToPostCreate
	// Select specifies which scalar and relation fields to select and return for created records.
	Select *CategoryToPostSelect
	// ConflictTarget specifies the unique constraint target for ON CONFLICT clause.
	ConflictTarget UniqueConstraintTarget
	// ConflictAction specifies the resolution action for ON CONFLICT clause.
	ConflictAction *ConflictAction
}

func (a *CategoryToPostCreateManyAndReturnArgs) AppendData(builders ...*CategoryToPostCreateBuilder) *CategoryToPostCreateManyAndReturnArgs {
	for _, b := range builders {
		input, err := assignmentsToCategoryToPostCreate(b.assignments)
		if err != nil {
			panic(err)
		}
		a.Data = append(a.Data, &input)
	}
	return a
}

// CategoryToPostFindUniqueArgs is the input argument passed to CategoryToPost FindUnique extension hooks.
type CategoryToPostFindUniqueArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[CategoryToPost]
	// Select specifies which scalar and relation fields to select and return.
	Select *CategoryToPostSelect
}

func (a *CategoryToPostFindUniqueArgs) SetWhere(unique UniquePredicate[CategoryToPost], additional ...PredicateOf[CategoryToPost]) *CategoryToPostFindUniqueArgs {
	a.Where = make([]PredicateOf[CategoryToPost], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// CategoryToPostFindFirstArgs is the input argument passed to CategoryToPost FindFirst extension hooks.
type CategoryToPostFindFirstArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[CategoryToPost]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[CategoryToPost]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[CategoryToPost]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *CategoryToPostSelect
}

func (a *CategoryToPostFindFirstArgs) SetWhere(preds ...PredicateOf[CategoryToPost]) *CategoryToPostFindFirstArgs {
	a.Where = preds
	return a
}

func (a *CategoryToPostFindFirstArgs) SetOrderBy(orders ...OrderBy[CategoryToPost]) *CategoryToPostFindFirstArgs {
	a.OrderBy = orders
	return a
}

func (a *CategoryToPostFindFirstArgs) SetCursor(cursor UniquePredicate[CategoryToPost]) *CategoryToPostFindFirstArgs {
	a.Cursor = cursor
	return a
}

func (a *CategoryToPostFindFirstArgs) SetSkip(n int) *CategoryToPostFindFirstArgs {
	a.Skip = &n
	return a
}

func (a *CategoryToPostFindFirstArgs) SetTake(n int) *CategoryToPostFindFirstArgs {
	a.Take = &n
	return a
}

// CategoryToPostFindManyArgs is the input argument passed to CategoryToPost FindMany extension hooks.
type CategoryToPostFindManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[CategoryToPost]
	// OrderBy specifies sorting definitions.
	OrderBy []OrderBy[CategoryToPost]
	// Cursor specifies cursor-based pagination parameters.
	Cursor UniquePredicate[CategoryToPost]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
	// Select specifies which scalar and relation fields to select and return.
	Select *CategoryToPostSelect
}

func (a *CategoryToPostFindManyArgs) SetWhere(preds ...PredicateOf[CategoryToPost]) *CategoryToPostFindManyArgs {
	a.Where = preds
	return a
}

func (a *CategoryToPostFindManyArgs) SetOrderBy(orders ...OrderBy[CategoryToPost]) *CategoryToPostFindManyArgs {
	a.OrderBy = orders
	return a
}

func (a *CategoryToPostFindManyArgs) SetCursor(cursor UniquePredicate[CategoryToPost]) *CategoryToPostFindManyArgs {
	a.Cursor = cursor
	return a
}

func (a *CategoryToPostFindManyArgs) SetSkip(n int) *CategoryToPostFindManyArgs {
	a.Skip = &n
	return a
}

func (a *CategoryToPostFindManyArgs) SetTake(n int) *CategoryToPostFindManyArgs {
	a.Take = &n
	return a
}

// CategoryToPostCountArgs is the input argument passed to CategoryToPost Count extension hooks.
type CategoryToPostCountArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[CategoryToPost]
	// Skip specifies offset for pagination.
	Skip *int
	// Take specifies limit for pagination.
	Take *int
}

func (a *CategoryToPostCountArgs) SetWhere(preds ...PredicateOf[CategoryToPost]) *CategoryToPostCountArgs {
	a.Where = preds
	return a
}

func (a *CategoryToPostCountArgs) SetSkip(n int) *CategoryToPostCountArgs {
	a.Skip = &n
	return a
}

func (a *CategoryToPostCountArgs) SetTake(n int) *CategoryToPostCountArgs {
	a.Take = &n
	return a
}

// CategoryToPostDeleteArgs is the input argument passed to CategoryToPost Delete extension hooks.
type CategoryToPostDeleteArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[CategoryToPost]
	// Select specifies which scalar and relation fields to select and return for deleted record.
	Select *CategoryToPostSelect
}

func (a *CategoryToPostDeleteArgs) SetWhere(unique UniquePredicate[CategoryToPost], additional ...PredicateOf[CategoryToPost]) *CategoryToPostDeleteArgs {
	a.Where = make([]PredicateOf[CategoryToPost], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// CategoryToPostDeleteManyArgs is the input argument passed to CategoryToPost DeleteMany extension hooks.
type CategoryToPostDeleteManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[CategoryToPost]
}

func (a *CategoryToPostDeleteManyArgs) SetWhere(preds ...PredicateOf[CategoryToPost]) *CategoryToPostDeleteManyArgs {
	a.Where = preds
	return a
}

// CategoryToPostUpdateArgs is the input argument passed to CategoryToPost Update extension hooks.
type CategoryToPostUpdateArgs struct {
	// Where contains all query filter predicates (merged primary unique constraint and additional predicates).
	Where []PredicateOf[CategoryToPost]
	// Data contains the model fields to update.
	Data *CategoryToPostUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *CategoryToPostSelect
}

func (a *CategoryToPostUpdateArgs) SetWhere(unique UniquePredicate[CategoryToPost], additional ...PredicateOf[CategoryToPost]) *CategoryToPostUpdateArgs {
	a.Where = make([]PredicateOf[CategoryToPost], 0, 1+len(additional))
	a.Where = append(a.Where, unique)
	a.Where = append(a.Where, additional...)
	return a
}

// CategoryToPostUpdateManyArgs is the input argument passed to CategoryToPost UpdateMany extension hooks.
type CategoryToPostUpdateManyArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[CategoryToPost]
	// Data contains the model fields to update.
	Data *CategoryToPostUpdate
}

func (a *CategoryToPostUpdateManyArgs) SetWhere(preds ...PredicateOf[CategoryToPost]) *CategoryToPostUpdateManyArgs {
	a.Where = preds
	return a
}

// CategoryToPostUpdateManyAndReturnArgs is the input argument passed to CategoryToPost UpdateManyAndReturn extension hooks.
type CategoryToPostUpdateManyAndReturnArgs struct {
	// Where contains all query filter predicates.
	Where []PredicateOf[CategoryToPost]
	// Data contains the model fields to update.
	Data *CategoryToPostUpdate
	// Select specifies which scalar and relation fields to select and return upon update.
	Select *CategoryToPostSelect
}

func (a *CategoryToPostUpdateManyAndReturnArgs) SetWhere(preds ...PredicateOf[CategoryToPost]) *CategoryToPostUpdateManyAndReturnArgs {
	a.Where = preds
	return a
}

type CategoryToPostCreateQuery = func(ctx context.Context, args *CategoryToPostCreateArgs) (*CategoryToPost, error)
type CategoryToPostCreateManyQuery = func(ctx context.Context, args *CategoryToPostCreateManyArgs) (int64, error)
type CategoryToPostCreateManyAndReturnQuery = func(ctx context.Context, args *CategoryToPostCreateManyAndReturnArgs) ([]*CategoryToPost, error)
type CategoryToPostFindUniqueQuery = func(ctx context.Context, args *CategoryToPostFindUniqueArgs) (*CategoryToPost, error)
type CategoryToPostFindFirstQuery = func(ctx context.Context, args *CategoryToPostFindFirstArgs) (*CategoryToPost, error)
type CategoryToPostFindManyQuery = func(ctx context.Context, args *CategoryToPostFindManyArgs) ([]*CategoryToPost, error)
type CategoryToPostDeleteQuery = func(ctx context.Context, args *CategoryToPostDeleteArgs) (*CategoryToPost, error)
type CategoryToPostDeleteManyQuery = func(ctx context.Context, args *CategoryToPostDeleteManyArgs) (int64, error)
type CategoryToPostCountQuery = func(ctx context.Context, args *CategoryToPostCountArgs) (int64, error)
type CategoryToPostUpdateQuery = func(ctx context.Context, args *CategoryToPostUpdateArgs) (*CategoryToPost, error)
type CategoryToPostUpdateManyQuery = func(ctx context.Context, args *CategoryToPostUpdateManyArgs) (int64, error)
type CategoryToPostUpdateManyAndReturnQuery = func(ctx context.Context, args *CategoryToPostUpdateManyAndReturnArgs) ([]*CategoryToPost, error)

type CategoryToPostExtension struct {
	Create              func(ctx context.Context, args *CategoryToPostCreateArgs, next CategoryToPostCreateQuery) (*CategoryToPost, error)
	CreateMany          func(ctx context.Context, args *CategoryToPostCreateManyArgs, next CategoryToPostCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, args *CategoryToPostCreateManyAndReturnArgs, next CategoryToPostCreateManyAndReturnQuery) ([]*CategoryToPost, error)
	FindUnique          func(ctx context.Context, args *CategoryToPostFindUniqueArgs, next CategoryToPostFindUniqueQuery) (*CategoryToPost, error)
	FindFirst           func(ctx context.Context, args *CategoryToPostFindFirstArgs, next CategoryToPostFindFirstQuery) (*CategoryToPost, error)
	FindMany            func(ctx context.Context, args *CategoryToPostFindManyArgs, next CategoryToPostFindManyQuery) ([]*CategoryToPost, error)
	Delete              func(ctx context.Context, args *CategoryToPostDeleteArgs, next CategoryToPostDeleteQuery) (*CategoryToPost, error)
	DeleteMany          func(ctx context.Context, args *CategoryToPostDeleteManyArgs, next CategoryToPostDeleteManyQuery) (int64, error)
	Count               func(ctx context.Context, args *CategoryToPostCountArgs, next CategoryToPostCountQuery) (int64, error)
	Update              func(ctx context.Context, args *CategoryToPostUpdateArgs, next CategoryToPostUpdateQuery) (*CategoryToPost, error)
	UpdateMany          func(ctx context.Context, args *CategoryToPostUpdateManyArgs, next CategoryToPostUpdateManyQuery) (int64, error)
	UpdateManyAndReturn func(ctx context.Context, args *CategoryToPostUpdateManyAndReturnArgs, next CategoryToPostUpdateManyAndReturnQuery) ([]*CategoryToPost, error)
}

type CategoryToPostDelegate struct {
	client     *Queries
	extensions []CategoryToPostExtension
}

func (d *CategoryToPostDelegate) Use(exts ...CategoryToPostExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *CategoryToPost) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "postId":
			targets[i] = &m.PostId
		case "categoryId":
			targets[i] = &m.CategoryId
		}
	}
	return targets
}

var categoryToPostDefaultCols = []string{
	"postId",
	"categoryId",
}

var categoryToPostPKCols = []string{
	"postId",
	"categoryId",
}

var categoryToPostUniqueCols = []string{}

func selectCategoryToPostCols(selects *CategoryToPostSelect, omits *CategoryToPostOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return categoryToPostDefaultCols
	}

	anySelected := selects != nil && (selects.PostId || selects.CategoryId || selects.Post != nil || selects.Category != nil)

	specs := []colSpec{
		{"postId", selects != nil && selects.PostId, omits != nil && omits.PostId, selects != nil && selects.Post != nil},
		{"categoryId", selects != nil && selects.CategoryId, omits != nil && omits.CategoryId, selects != nil && selects.Category != nil},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *CategoryToPostSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Post != nil || s.Category != nil
}

type CategoryToPostCreateBuilder struct {
	*CreateBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]
}

func (b *CategoryToPostCreateBuilder) Select(s CategoryToPostSelect) *CategoryToPostCreateBuilder {
	b.selects = &s
	return b
}

func (b *CategoryToPostCreateBuilder) Omit(o CategoryToPostOmit) *CategoryToPostCreateBuilder {
	b.omits = &o
	return b
}

func (b *CategoryToPostCreateBuilder) OnConflict(target UniqueConstraintTarget) *CategoryToPostConflictBuilder[CategoryToPostCreateBuilder] {
	return &CategoryToPostConflictBuilder[CategoryToPostCreateBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *CategoryToPostCreateBuilder) SetPostId(v string) *CategoryToPostCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "postId", Val: v})
	return b
}
func (b *CategoryToPostCreateBuilder) SetCategoryId(v int32) *CategoryToPostCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "categoryId", Val: v})
	return b
}

func (b *CategoryToPostCreateBuilder) Assignments(assignments ...FieldAssignmentOf[CategoryToPost]) *CategoryToPostCreateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment{Col: a.Col, Val: a.Val})
	}
	return b
}

func (d *CategoryToPostDelegate) Create() *CategoryToPostCreateBuilder {
	return &CategoryToPostCreateBuilder{
		CreateBuilder: &CreateBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
			execFunc: d.executeCreate,
		},
	}
}

const (
	providedCategoryToPostPostId     uint64 = 1 << 0
	providedCategoryToPostCategoryId uint64 = 1 << 1
)

func assignmentsToCategoryToPostCreate(assignments []FieldAssignment) (CategoryToPostCreate, error) {
	var input CategoryToPostCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "postId":
			provided |= providedCategoryToPostPostId
			if v, ok := a.Val.(string); ok {
				input.PostId = v
				errs.ValidateString("postId", v, true, 0, false, false)
			} else {
				errs.Add("postId", a.Val, "type", "field postId must be of type string")
			}
		case "categoryId":
			provided |= providedCategoryToPostCategoryId
			if v, ok := a.Val.(int32); ok {
				input.CategoryId = v
				errs.ValidateInt32("categoryId", v, "")
			} else {
				errs.Add("categoryId", a.Val, "type", "field categoryId must be of type int32")
			}
		}
	}
	if provided&providedCategoryToPostPostId == 0 {
		errs.Add("postId", "", "required", "field PostId is required")
	}
	if provided&providedCategoryToPostCategoryId == 0 {
		errs.Add("categoryId", nil, "required", "field CategoryId is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *CategoryToPostCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 2)
	vals = make([]any, 0, 2)
	cols = append(cols, "postId")
	vals = append(vals, s.PostId)
	cols = append(cols, "categoryId")
	vals = append(vals, s.CategoryId)
	return
}

func partitionCategoryToPostInputs(dialect Dialect, inputs []*CategoryToPostCreate) [][]*CategoryToPostCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*CategoryToPostCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*CategoryToPostCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*CategoryToPostCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*CategoryToPostCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*CategoryToPostCreate{inputs}
}

func (d *CategoryToPostDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *CategoryToPostSelect, omits *CategoryToPostOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*CategoryToPost, error) {
	input, err := assignmentsToCategoryToPostCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectCategoryToPostCols(selects, omits)

	if len(d.extensions) == 0 {
		return d.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullCategoryToPostSelect()
	}

	args := &CategoryToPostCreateArgs{
		Data:           &input,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *CategoryToPostCreateArgs) (*CategoryToPost, error) {
		cCols, cVals := a.Data.ToColsVals()
		cReturningCols := selectCategoryToPostCols(a.Select, omits)
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
			curr = func(c context.Context, a *CategoryToPostCreateArgs) (*CategoryToPost, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

type CategoryToPostCreateManyBuilder struct {
	*CreateManyBuilder[CategoryToPost]
}

func (b *CategoryToPostCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *CategoryToPostConflictBuilder[CategoryToPostCreateManyBuilder] {
	return &CategoryToPostConflictBuilder[CategoryToPostCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type CategoryToPostCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]
}

func (b *CategoryToPostCreateManyAndReturnBuilder) Select(s CategoryToPostSelect) *CategoryToPostCreateManyAndReturnBuilder {
	b.selects = &s
	return b
}

func (b *CategoryToPostCreateManyAndReturnBuilder) Omit(o CategoryToPostOmit) *CategoryToPostCreateManyAndReturnBuilder {
	b.omits = &o
	return b
}

func (b *CategoryToPostCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *CategoryToPostConflictBuilder[CategoryToPostCreateManyAndReturnBuilder] {
	return &CategoryToPostConflictBuilder[CategoryToPostCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func createBuildersToCategoryToPostRecordInputs(builders []*CategoryToPostCreateBuilder) []RecordInput {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return records
}

func (d *CategoryToPostDelegate) CreateMany(builders ...*CategoryToPostCreateBuilder) *CategoryToPostCreateManyBuilder {
	return &CategoryToPostCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[CategoryToPost]{
			records:  createBuildersToCategoryToPostRecordInputs(builders),
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *CategoryToPostDelegate) CreateManyAndReturn(builders ...*CategoryToPostCreateBuilder) *CategoryToPostCreateManyAndReturnBuilder {
	return &CategoryToPostCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
			records:  createBuildersToCategoryToPostRecordInputs(builders),
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func recordsToCategoryToPostCreateInputs(records []RecordInput) ([]*CategoryToPostCreate, error) {
	structs := make([]CategoryToPostCreate, len(records))
	inputs := make([]*CategoryToPostCreate, len(records))
	for i, rec := range records {
		var err error
		structs[i], err = assignmentsToCategoryToPostCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &structs[i]
	}
	return inputs, nil
}

func (d *CategoryToPostDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs, err := recordsToCategoryToPostCreateInputs(records)
	if err != nil {
		return 0, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	args := &CategoryToPostCreateManyArgs{
		Data:           inputs,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *CategoryToPostCreateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *CategoryToPostCreateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *CategoryToPostDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *CategoryToPostSelect, omits *CategoryToPostOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*CategoryToPost, error) {
	inputs, err := recordsToCategoryToPostCreateInputs(records)
	if err != nil {
		return nil, err
	}

	if len(d.extensions) == 0 {
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullCategoryToPostSelect()
	}

	args := &CategoryToPostCreateManyAndReturnArgs{
		Data:           inputs,
		Select:         selects,
		ConflictTarget: conflictTarget,
		ConflictAction: conflictAction,
	}

	curr := func(c context.Context, a *CategoryToPostCreateManyAndReturnArgs) ([]*CategoryToPost, error) {
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
			curr = func(c context.Context, a *CategoryToPostCreateManyAndReturnArgs) ([]*CategoryToPost, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *CategoryToPostDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	selects *CategoryToPostSelect,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*CategoryToPost, error) {
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := hasRelations && !d.client.inTx()

	if useTx {
		var res *CategoryToPost
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.CategoryToPost.runCreate(ctx, cols, vals, returningCols, selects, conflictTarget, conflictAction)
			if err != nil {
				return err
			}
			return txQ.CategoryToPost.loadRelations(ctx, []*CategoryToPost{res}, selects)
		})
		return res, err
	}

	query, clauseArgs := buildSingleInsertSQL(d.client, "CategoryToPost", cols, returningCols, categoryToPostPKCols, conflictTarget, conflictAction, len(vals))
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
				return nil, TranslateDBError(err)
			}
			return nil, nil
		}

		var res CategoryToPost
		scanErr := rows.Scan(res.ScanFields(returningCols)...)
		rows.Close()
		if scanErr != nil {
			return nil, TranslateDBError(scanErr)
		}

		return &res, nil
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, categoryToPostPKCols)
}

func (d *CategoryToPostDelegate) runCreateFallback(
	ctx context.Context,
	query string,
	vals []any,
	cols []string,
	returningCols []string,
	pkCols []string,
) (*CategoryToPost, error) {
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
				return nil, TranslateDBError(err)
			}
			val = lastID
		}
		pkVals = append(pkVals, val)
	}

	var selectSb strings.Builder
	selectSb.Grow(64 + len(returningCols)*15 + len("CategoryToPost") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "CategoryToPost")
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
			return nil, TranslateDBError(err)
		}
		return nil, nil
	}

	var res CategoryToPost
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	return &res, nil
}

func (d *CategoryToPostDelegate) buildBulkInsertSQL(q *Queries, batch []*CategoryToPostCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 2)
	for i, c := range categoryToPostDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "CategoryToPost")
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
			case "postId":
				vals = append(vals, input.PostId)
			case "categoryId":
				vals = append(vals, input.CategoryId)
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

func applyCategoryToPostConflictClause(dialect Dialect, queryStr string, vals []any, cols []string, pkCols []string, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (string, []any) {
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

func scanCategoryToPostRows(rows *sql.Rows, returningCols []string) ([]*CategoryToPost, error) {
	var records []*CategoryToPost
	for rows.Next() {
		var res CategoryToPost
		if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
			rows.Close()
			return nil, TranslateDBError(err)
		}
		records = append(records, &res)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, TranslateDBError(err)
	}
	rows.Close()
	return records, nil
}

func (d *CategoryToPostDelegate) runCreateMany(ctx context.Context, inputs []*CategoryToPostCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionCategoryToPostInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyCategoryToPostConflictClause(d.client.dialect, queryStr, vals, cols, categoryToPostPKCols, conflictTarget, conflictAction)

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

func (d *CategoryToPostDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*CategoryToPostCreate,
	selects *CategoryToPostSelect,
	omits *CategoryToPostOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*CategoryToPost, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionCategoryToPostInputs(d.client.dialect, inputs)
	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := (len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning) && !d.client.inTx()

	if useTx {
		var res []*CategoryToPost
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if txQ.dialect.SupportsInsertReturning {
				res, err = txQ.CategoryToPost.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			} else {
				res, err = txQ.CategoryToPost.runCreateManyAndReturnFallback(ctx, inputs, selects, omits, conflictTarget, conflictAction)
			}
			if err != nil {
				return err
			}
			if hasRelations {
				return txQ.CategoryToPost.loadRelations(ctx, res, selects)
			}
			return nil
		})
		return res, err
	}

	returningCols := selectCategoryToPostCols(selects, omits, categoryToPostPKCols...)
	recordsOut := make([]*CategoryToPost, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyCategoryToPostConflictClause(d.client.dialect, queryStr, vals, cols, categoryToPostPKCols, conflictTarget, conflictAction)

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

		scanned, err := scanCategoryToPostRows(rows, returningCols)
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

func (d *CategoryToPostDelegate) runCreateManyAndReturnFallback(
	ctx context.Context,
	inputs []*CategoryToPostCreate,
	selects *CategoryToPostSelect,
	omits *CategoryToPostOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*CategoryToPost, error) {
	batches := partitionCategoryToPostInputs(d.client.dialect, inputs)
	returningCols := selectCategoryToPostCols(selects, omits, categoryToPostPKCols...)
	recordsOut := make([]*CategoryToPost, 0, len(inputs))

	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)
		queryStr, vals = applyCategoryToPostConflictClause(d.client.dialect, queryStr, vals, cols, categoryToPostPKCols, conflictTarget, conflictAction)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return nil, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("CategoryToPost") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			d.client.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		d.client.dialect.WriteQuotedIdent(&selectSb, "CategoryToPost")
		selectSb.WriteString(" WHERE ")
		d.client.dialect.WriteQuotedIdent(&selectSb, categoryToPostPKCols[0])
		selectSb.WriteString(" >= ")
		d.client.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		d.client.dialect.WriteQuotedIdent(&selectSb, categoryToPostPKCols[0])
		selectSb.WriteString(" < ")
		d.client.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := d.client.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return nil, err
		}

		scanned, err := scanCategoryToPostRows(rows, returningCols)
		if err != nil {
			return nil, err
		}
		recordsOut = append(recordsOut, scanned...)
	}

	return recordsOut, nil
}

type CategoryToPostConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *CategoryToPostConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *CategoryToPostConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *CategoryToPostConflictBuilder[B]) Update(fn func(u *CategoryToPostUpsert)) *B {
	cb.setAction(*CategoryToPostConflictUpdate(fn), cb.conflictTarget)
	return cb.builder
}

// CategoryToPostConflictUpdate creates a custom ConflictAction for CategoryToPost upsert conflicts.
func CategoryToPostConflictUpdate(fn func(u *CategoryToPostUpsert)) *ConflictAction {
	var up ConflictUpdate
	u := newCategoryToPostUpsert(&up)
	fn(u)
	return &ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}
}

type CategoryToPostUpsert struct {
	PostId     fieldUpsert[string]
	CategoryId numericFieldUpsert[int32]
}

func newCategoryToPostUpsert(up *ConflictUpdate) *CategoryToPostUpsert {
	return &CategoryToPostUpsert{
		PostId: fieldUpsert[string]{column: "postId", update: up},
		CategoryId: numericFieldUpsert[int32]{
			fieldUpsert: fieldUpsert[int32]{column: "categoryId", update: up},
			tableName:   "CategoryToPost",
		},
	}
}

type CategoryToPostUpdateBuilder struct {
	*UpdateBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]
}

type CategoryToPostUpdateManyBuilder struct {
	*UpdateManyBuilder[CategoryToPost]
}

type CategoryToPostUpdateManyAndReturnBuilder struct {
	*UpdateManyAndReturnBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]
}

func (b *CategoryToPostUpdateBuilder) SetPostId(v string) *CategoryToPostUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "postId", Val: v})
	return b
}

func (b *CategoryToPostUpdateManyBuilder) SetPostId(v string) *CategoryToPostUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "postId", Val: v})
	return b
}

func (b *CategoryToPostUpdateManyAndReturnBuilder) SetPostId(v string) *CategoryToPostUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "postId", Val: v})
	return b
}
func (b *CategoryToPostUpdateBuilder) SetCategoryId(v int32) *CategoryToPostUpdateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "categoryId", Val: v})
	return b
}

func (b *CategoryToPostUpdateManyBuilder) SetCategoryId(v int32) *CategoryToPostUpdateManyBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "categoryId", Val: v})
	return b
}

func (b *CategoryToPostUpdateManyAndReturnBuilder) SetCategoryId(v int32) *CategoryToPostUpdateManyAndReturnBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "categoryId", Val: v})
	return b
}

func (b *CategoryToPostUpdateBuilder) Assignments(assignments ...FieldAssignmentOf[CategoryToPost]) *CategoryToPostUpdateBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment{Col: a.Col, Val: a.Val})
	}
	return b
}

func (b *CategoryToPostUpdateManyBuilder) Assignments(assignments ...FieldAssignmentOf[CategoryToPost]) *CategoryToPostUpdateManyBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment{Col: a.Col, Val: a.Val})
	}
	return b
}

func (b *CategoryToPostUpdateManyAndReturnBuilder) Assignments(assignments ...FieldAssignmentOf[CategoryToPost]) *CategoryToPostUpdateManyAndReturnBuilder {
	for _, a := range assignments {
		b.assignments = append(b.assignments, FieldAssignment{Col: a.Col, Val: a.Val})
	}
	return b
}

func (d *CategoryToPostDelegate) Update(where UniquePredicate[CategoryToPost], additional ...PredicateOf[CategoryToPost]) *CategoryToPostUpdateBuilder {
	return &CategoryToPostUpdateBuilder{
		UpdateBuilder: &UpdateBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
			where:      where,
			additional: additional,
			execFunc:   d.executeUpdate,
		},
	}
}

func (d *CategoryToPostDelegate) UpdateMany(preds ...PredicateOf[CategoryToPost]) *CategoryToPostUpdateManyBuilder {
	return &CategoryToPostUpdateManyBuilder{
		UpdateManyBuilder: &UpdateManyBuilder[CategoryToPost]{
			where:    preds,
			execFunc: d.executeUpdateMany,
		},
	}
}

func (d *CategoryToPostDelegate) UpdateManyAndReturn(preds ...PredicateOf[CategoryToPost]) *CategoryToPostUpdateManyAndReturnBuilder {
	return &CategoryToPostUpdateManyAndReturnBuilder{
		UpdateManyAndReturnBuilder: &UpdateManyAndReturnBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
			where:    preds,
			execFunc: d.executeUpdateManyAndReturn,
		},
	}
}

func (d *CategoryToPostDelegate) buildUpdateSQL(preds []PredicateOf[CategoryToPost], cols []string, vals []any, returningCols []string) (string, []any) {
	whereClause, predVals, _ := CompilePredicates(d.client.dialect, preds, len(cols)+1)

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	d.client.dialect.WriteQuotedIdent(&sb, "CategoryToPost")
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

func (d *CategoryToPostDelegate) executeUpdate(ctx context.Context, where UniquePredicate[CategoryToPost], additional []PredicateOf[CategoryToPost], assignments []FieldAssignment, selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error) {
	allWhere := make([]PredicateOf[CategoryToPost], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	input, err := assignmentsToCategoryToPostUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdate(ctx, allWhere, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullCategoryToPostSelect()
	}

	args := &CategoryToPostUpdateArgs{
		Where:  allWhere,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *CategoryToPostUpdateArgs) (*CategoryToPost, error) {
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
			curr = func(c context.Context, a *CategoryToPostUpdateArgs) (*CategoryToPost, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *CategoryToPostDelegate) runUpdate(ctx context.Context, preds []PredicateOf[CategoryToPost], cols []string, vals []any, selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error) {
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
		var res *CategoryToPost
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.CategoryToPost.runUpdate(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.CategoryToPost.runUpdateFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectCategoryToPostCols(selects, omits, categoryToPostPKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		err := rows.Err()
		rows.Close()
		if err != nil {
			return nil, TranslateDBError(err)
		}
		return nil, &NotFoundError{Model: "CategoryToPost"}
	}

	var res CategoryToPost
	scanErr := rows.Scan(res.ScanFields(returningCols)...)
	rows.Close()
	if scanErr != nil {
		return nil, TranslateDBError(scanErr)
	}

	if selects != nil && selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*CategoryToPost{&res}, selects); err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (d *CategoryToPostDelegate) execUpdateStmt(ctx context.Context, preds []PredicateOf[CategoryToPost], cols []string, vals []any) (int64, error) {
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

func (d *CategoryToPostDelegate) runUpdateFallback(ctx context.Context, preds []PredicateOf[CategoryToPost], cols []string, vals []any, selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, &NotFoundError{Model: "CategoryToPost"}
	}
	return d.runFindUnique(ctx, preds, selects, omits)
}

// -----------------------------------------------------------------------------
// UpdateMany
// -----------------------------------------------------------------------------

func (d *CategoryToPostDelegate) executeUpdateMany(ctx context.Context, preds []PredicateOf[CategoryToPost], assignments []FieldAssignment) (int64, error) {
	input, err := assignmentsToCategoryToPostUpdate(assignments)
	if err != nil {
		return 0, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.execUpdateStmt(ctx, preds, cols, vals)
	}

	args := &CategoryToPostUpdateManyArgs{
		Where: preds,
		Data:  &input,
	}

	curr := func(c context.Context, a *CategoryToPostUpdateManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *CategoryToPostUpdateManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

// -----------------------------------------------------------------------------
// UpdateManyAndReturn
// -----------------------------------------------------------------------------

func (d *CategoryToPostDelegate) executeUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[CategoryToPost], assignments []FieldAssignment, selects *CategoryToPostSelect, omits *CategoryToPostOmit) ([]*CategoryToPost, error) {
	input, err := assignmentsToCategoryToPostUpdate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()

	if len(d.extensions) == 0 {
		return d.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullCategoryToPostSelect()
	}

	args := &CategoryToPostUpdateManyAndReturnArgs{
		Where:  preds,
		Data:   &input,
		Select: selects,
	}

	curr := func(c context.Context, a *CategoryToPostUpdateManyAndReturnArgs) ([]*CategoryToPost, error) {
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
			curr = func(c context.Context, a *CategoryToPostUpdateManyAndReturnArgs) ([]*CategoryToPost, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *CategoryToPostDelegate) runUpdateManyAndReturn(ctx context.Context, preds []PredicateOf[CategoryToPost], cols []string, vals []any, selects *CategoryToPostSelect, omits *CategoryToPostOmit) ([]*CategoryToPost, error) {
	if len(cols) == 0 {
		return d.runFindMany(ctx, QueryParams[CategoryToPost]{Where: preds}, selects, omits)
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
		var res []*CategoryToPost
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			if d.client.dialect.SupportsUpdateReturning {
				res, err = txQ.CategoryToPost.runUpdateManyAndReturn(ctx, preds, cols, vals, selects, omits)
			} else {
				res, err = txQ.CategoryToPost.runUpdateManyAndReturnFallback(ctx, preds, cols, vals, selects, omits)
			}
			return err
		})
		return res, err
	}

	returningCols := selectCategoryToPostCols(selects, omits, categoryToPostPKCols...)
	query, setVals := d.buildUpdateSQL(preds, cols, vals, returningCols)

	rows, err := d.client.query(ctx, query, setVals...)
	if err != nil {
		return nil, err
	}

	scanned, err := scanCategoryToPostRows(rows, returningCols)
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

func (d *CategoryToPostDelegate) runUpdateManyAndReturnFallback(ctx context.Context, preds []PredicateOf[CategoryToPost], cols []string, vals []any, selects *CategoryToPostSelect, omits *CategoryToPostOmit) ([]*CategoryToPost, error) {
	affected, err := d.execUpdateStmt(ctx, preds, cols, vals)
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return []*CategoryToPost{}, nil
	}
	return d.runFindMany(ctx, QueryParams[CategoryToPost]{Where: preds}, selects, omits)
}
func (d *CategoryToPostDelegate) FindUnique(where UniquePredicate[CategoryToPost], additional ...PredicateOf[CategoryToPost]) *FindUniqueBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit] {
	return &FindUniqueBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *CategoryToPostDelegate) FindFirst(preds ...PredicateOf[CategoryToPost]) *FindFirstBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit] {
	return &FindFirstBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *CategoryToPostDelegate) FindMany(preds ...PredicateOf[CategoryToPost]) *FindManyBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit] {
	return &FindManyBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *CategoryToPostDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[CategoryToPost], additional []PredicateOf[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error) {
	allWhere := make([]PredicateOf[CategoryToPost], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullCategoryToPostSelect()
	}

	args := &CategoryToPostFindUniqueArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *CategoryToPostFindUniqueArgs) (*CategoryToPost, error) {
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
			curr = func(c context.Context, a *CategoryToPostFindUniqueArgs) (*CategoryToPost, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *CategoryToPostDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[CategoryToPost],
	selects *CategoryToPostSelect,
	omits *CategoryToPostOmit,
) (*CategoryToPost, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullCategoryToPostSelect()
	}

	args := &CategoryToPostFindFirstArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *CategoryToPostFindFirstArgs) (*CategoryToPost, error) {
		return d.runFindFirst(c, QueryParams[CategoryToPost]{
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
			curr = func(c context.Context, a *CategoryToPostFindFirstArgs) (*CategoryToPost, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *CategoryToPostDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[CategoryToPost],
	selects *CategoryToPostSelect,
	omits *CategoryToPostOmit,
) ([]*CategoryToPost, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullCategoryToPostSelect()
	}

	args := &CategoryToPostFindManyArgs{
		Where:   params.Where,
		OrderBy: params.OrderBy,
		Cursor:  params.Cursor,
		Skip:    params.Skip,
		Take:    params.Take,
		Select:  selects,
	}

	curr := func(c context.Context, a *CategoryToPostFindManyArgs) ([]*CategoryToPost, error) {
		return d.runFindMany(c, QueryParams[CategoryToPost]{
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
			curr = func(c context.Context, a *CategoryToPostFindManyArgs) ([]*CategoryToPost, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *CategoryToPostDelegate) runFindUnique(ctx context.Context, where []PredicateOf[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error) {
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
	returningCols := selectCategoryToPostCols(selects, omits)

	res, err := d.queryOne(ctx, whereClause, "", vals, returningCols, nil)
	if err != nil || res == nil {
		return res, err
	}
	if selects.hasAnyRelation() {
		if err := d.loadRelations(ctx, []*CategoryToPost{res}, selects); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (d *CategoryToPostDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[CategoryToPost],
	selects *CategoryToPostSelect,
	omits *CategoryToPostOmit,
) (*CategoryToPost, error) {
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

func (d *CategoryToPostDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[CategoryToPost],
	selects *CategoryToPostSelect,
	omits *CategoryToPostOmit,
) ([]*CategoryToPost, error) {
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
		cClause, cVals, err := compileCursorClause(d.client.dialect, params.Cursor, params.OrderBy, categoryToPostPKCols, categoryToPostUniqueCols, "CategoryToPost", nextIdx, params.Take)
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
	orderByClause := formatOrderBySQL(d.client.dialect, params.OrderBy, categoryToPostPKCols, categoryToPostUniqueCols, isCursorQuery, params.Take)
	returningCols := selectCategoryToPostCols(selects, omits)

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

func (d *CategoryToPostDelegate) queryOne(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, skip *int) (*CategoryToPost, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "CategoryToPost", returningCols, whereClause, orderByClause, &limitOne, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
			return nil, nil
		}
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
				return nil, nil
			}
			return nil, TranslateDBError(err)
		}
		return nil, nil
	}

	var res CategoryToPost
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if errors.Is(err, sql.ErrNoRows) || IsNotFound(err) {
			return nil, nil
		}
		return nil, TranslateDBError(err)
	}

	return &res, nil
}

func (d *CategoryToPostDelegate) queryMany(ctx context.Context, whereClause string, orderByClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*CategoryToPost, error) {
	query := buildSelectSQL(d.client, "CategoryToPost", returningCols, whereClause, orderByClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, TranslateDBError(err)
	}
	defer rows.Close()

	results := make([]*CategoryToPost, 0)
	for rows.Next() {
		var res CategoryToPost
		if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
			return nil, TranslateDBError(err)
		}
		results = append(results, &res)
	}
	if err := rows.Err(); err != nil {
		return nil, TranslateDBError(err)
	}
	if take != nil && *take < 0 {
		reverseSlice(results)
	}
	return results, nil
}
func (d *CategoryToPostDelegate) DeleteMany(preds ...PredicateOf[CategoryToPost]) *DeleteManyBuilder[CategoryToPost] {
	return &DeleteManyBuilder[CategoryToPost]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *CategoryToPostDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[CategoryToPost]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	args := &CategoryToPostDeleteManyArgs{
		Where: preds,
	}

	curr := func(c context.Context, a *CategoryToPostDeleteManyArgs) (int64, error) {
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
			curr = func(c context.Context, a *CategoryToPostDeleteManyArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *CategoryToPostDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[CategoryToPost]) (int64, error) {
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
	d.client.dialect.WriteQuotedIdent(&sb, "CategoryToPost")
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

func (d *CategoryToPostDelegate) Delete(where UniquePredicate[CategoryToPost], additional ...PredicateOf[CategoryToPost]) *DeleteBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit] {
	return &DeleteBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeDelete,
	}
}

func (d *CategoryToPostDelegate) executeDelete(ctx context.Context, where UniquePredicate[CategoryToPost], additional []PredicateOf[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error) {
	allWhere := make([]PredicateOf[CategoryToPost], 0, 1+len(additional))
	allWhere = append(allWhere, where)
	allWhere = append(allWhere, additional...)

	if len(d.extensions) == 0 {
		return d.runDelete(ctx, allWhere, selects, omits)
	}

	if selects == nil || !selects.hasAnySelected() {
		selects = fullCategoryToPostSelect()
	}

	args := &CategoryToPostDeleteArgs{
		Where:  allWhere,
		Select: selects,
	}

	curr := func(c context.Context, a *CategoryToPostDeleteArgs) (*CategoryToPost, error) {
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
			curr = func(c context.Context, a *CategoryToPostDeleteArgs) (*CategoryToPost, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *CategoryToPostDelegate) runDelete(ctx context.Context, where []PredicateOf[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}

	returningCols := selectCategoryToPostCols(selects, omits, categoryToPostPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *CategoryToPost
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.CategoryToPost.runFindUnique(ctx, where, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return &NotFoundError{Model: "CategoryToPost"}
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "CategoryToPost")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[CategoryToPost]
			pkPreds = append(pkPreds, Predicate[CategoryToPost]{
				Data: PredicateData{
					Column:   "postId",
					Operator: "=",
					Value:    res.PostId,
				},
			})
			pkPreds = append(pkPreds, Predicate[CategoryToPost]{
				Data: PredicateData{
					Column:   "categoryId",
					Operator: "=",
					Value:    res.CategoryId,
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
	d.client.dialect.WriteQuotedIdent(&sb, "CategoryToPost")

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
			return nil, TranslateDBError(err)
		}
		return nil, &NotFoundError{Model: "CategoryToPost"}
	}

	var row CategoryToPost
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, TranslateDBError(err)
	}
	return &row, nil
}
func (d *CategoryToPostDelegate) Count(preds ...PredicateOf[CategoryToPost]) *CountBuilder[CategoryToPost] {
	return &CountBuilder[CategoryToPost]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *CategoryToPostDelegate) executeCount(ctx context.Context, params QueryParams[CategoryToPost]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	args := &CategoryToPostCountArgs{
		Where: params.Where,
		Skip:  params.Skip,
		Take:  params.Take,
	}

	curr := func(c context.Context, a *CategoryToPostCountArgs) (int64, error) {
		return d.runCount(c, QueryParams[CategoryToPost]{
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
			curr = func(c context.Context, a *CategoryToPostCountArgs) (int64, error) {
				return hook(c, a, next)
			}
		}
	}

	return curr(ctx, args)
}

func (d *CategoryToPostDelegate) runCount(ctx context.Context, params QueryParams[CategoryToPost]) (int64, error) {
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
		d.client.dialect.WriteQuotedIdent(&subQuery, "CategoryToPost")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "CategoryToPost")
		if whereClause != "" {
			sb.WriteString(whereClause)
		}
		query = sb.String()
	}

	rows, err := d.client.query(ctx, query, vals...)
	if err != nil {
		return 0, TranslateDBError(err)
	}
	defer rows.Close()

	var count int64
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, TranslateDBError(err)
		}
	}
	if err := rows.Err(); err != nil {
		return 0, TranslateDBError(err)
	}
	return count, nil
}
func (d *CategoryToPostDelegate) loadRelations(ctx context.Context, records []*CategoryToPost, selects *CategoryToPostSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Post != nil {
		relationSelects, relationOmits, relationParams := selects.Post.GetRelationParams()
		returningCols := selectPostCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: CategoryToPost.postId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *CategoryToPost) string { return p.PostId }),
			"Post",
			"id",
			returningCols,
			scanInto(returningCols, (*Post).ScanFields),
			directKey(func(c *Post) string { return c.Id }),
			setOne(func(p *CategoryToPost, c *Post) { p.Post = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading post: %w", err)
		}
		if err := d.client.Post.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Category != nil {
		relationSelects, relationOmits, relationParams := selects.Category.GetRelationParams()
		returningCols := selectCategoryCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: CategoryToPost.categoryId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *CategoryToPost) int32 { return p.CategoryId }),
			"Category",
			"id",
			returningCols,
			scanInto(returningCols, (*Category).ScanFields),
			directKey(func(c *Category) int32 { return c.Id }),
			setOne(func(p *CategoryToPost, c *Category) { p.Category = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading category: %w", err)
		}
		if err := d.client.Category.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
