package valk

import (
	"context"
	"fmt"
	"slices"
)

// CategoryToPost represents the database model
type CategoryToPost struct {
	PostId     string    `db:"postId" json:"postId"`
	CategoryId int32     `db:"categoryId" json:"categoryId"`
	Post       *Post     `json:"post,omitempty"`
	Category   *Category `json:"category,omitempty"`
}

// CategoryToPostCreate is used for hooks only — the Create API uses FieldAssignment
type CategoryToPostCreate struct {
	PostId     string `json:"postId"`
	CategoryId int32  `json:"categoryId"`
}

// CategoryToPostSelect specifies which fields to include
type CategoryToPostSelect struct {
	PostId     bool                `json:"postId"`
	CategoryId bool                `json:"categoryId"`
	Post       PostSelectQuery     `json:"post,omitempty"`
	Category   CategorySelectQuery `json:"category,omitempty"`
}

// CategoryToPostOmit specifies which fields to exclude
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

// CategoryToPostQueryBuilder builds a query for the relation CategoryToPost
type CategoryToPostQueryBuilder struct {
	selects *CategoryToPostSelect
	omits   *CategoryToPostOmit
	where   []PredicateOf[CategoryToPost]
	take    *int
	skip    *int
	orderBy []OrderBy[CategoryToPost]
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
	}
}

type CategoryToPostCreateQuery = func(ctx context.Context, args *CategoryToPostCreate) (*CategoryToPost, error)
type CategoryToPostCreateManyQuery = func(ctx context.Context, args []*CategoryToPostCreate) (int64, error)
type CategoryToPostCreateManyAndReturnQuery = func(ctx context.Context, args []*CategoryToPostCreate) ([]*CategoryToPost, error)
type CategoryToPostFindUniqueQuery = func(ctx context.Context, where UniquePredicate[CategoryToPost], additional []PredicateOf[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error)
type CategoryToPostFindFirstQuery = func(ctx context.Context, params QueryParams[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error)
type CategoryToPostFindManyQuery = func(ctx context.Context, params QueryParams[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit) ([]*CategoryToPost, error)

type CategoryToPostExtension struct {
	Create              func(ctx context.Context, input *CategoryToPostCreate, next CategoryToPostCreateQuery) (*CategoryToPost, error)
	CreateMany          func(ctx context.Context, inputs []*CategoryToPostCreate, next CategoryToPostCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, inputs []*CategoryToPostCreate, next CategoryToPostCreateManyAndReturnQuery) ([]*CategoryToPost, error)
	FindUnique          func(ctx context.Context, where UniquePredicate[CategoryToPost], additional []PredicateOf[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit, next CategoryToPostFindUniqueQuery) (*CategoryToPost, error)
	FindFirst           func(ctx context.Context, params QueryParams[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit, next CategoryToPostFindFirstQuery) (*CategoryToPost, error)
	FindMany            func(ctx context.Context, params QueryParams[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit, next CategoryToPostFindManyQuery) ([]*CategoryToPost, error)
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

func (q *Queries) selectCategoryToPostCols(selects *CategoryToPostSelect, omits *CategoryToPostOmit, forceCols ...string) []string {
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

var CategoryToPostColOrder = []string{
	"postId",
	"categoryId",
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

func (d *CategoryToPostDelegate) Create(assignments ...FieldAssignment) *CategoryToPostCreateBuilder {
	return &CategoryToPostCreateBuilder{
		CreateBuilder: &CreateBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
			client:      d.client,
			assignments: assignments,
			execFunc:    d.client.executeCategoryToPostCreate,
		},
	}
}

func validateCategoryToPostCreate(assignments []FieldAssignment) error {
	errs := &ValidationError{}

	provided := make(map[string]bool)
	for _, a := range assignments {
		provided[a.Col] = true
		switch a.Col {
		case "postId":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "postId", v, true, 0, false, false)
			} else {
				errs.Add("postId", a.Val, "type", "field postId must be of type string")
			}
		case "categoryId":
			if v, ok := a.Val.(int32); ok {
				ValidateInt32(errs, "categoryId", v, "")
			} else {
				errs.Add("categoryId", a.Val, "type", "field categoryId must be of type int32")
			}
		}
	}
	if !provided["postId"] {
		errs.Add("postId", "", "required", "field PostId is required")
	}
	if !provided["categoryId"] {
		errs.Add("categoryId", nil, "required", "field CategoryId is required")
	}

	if errs.HasErrors() {
		return *errs
	}
	return nil
}

func assignmentsToCategoryToPostCreate(assignments []FieldAssignment) CategoryToPostCreate {
	var input CategoryToPostCreate
	for _, a := range assignments {
		switch a.Col {
		case "postId":
			if v, ok := a.Val.(string); ok {
				input.PostId = v
			}
		case "categoryId":
			if v, ok := a.Val.(int32); ok {
				input.CategoryId = v
			}
		}
	}
	return input
}

func (s *CategoryToPostCreate) ToRowMap() map[string]any {
	m := make(map[string]any, 2)
	m["postId"] = s.PostId
	m["categoryId"] = s.CategoryId
	return m
}

func (q *Queries) executeCategoryToPostCreate(ctx context.Context, assignments []FieldAssignment, selects *CategoryToPostSelect, omits *CategoryToPostOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*CategoryToPost, error) {
	input := assignmentsToCategoryToPostCreate(assignments)

	curr := func(c context.Context, args *CategoryToPostCreate) (*CategoryToPost, error) {
		if err := validateCategoryToPostCreate(assignments); err != nil {
			return nil, err
		}

		rowMap := args.ToRowMap()
		cols, vals := mapToColsVals(rowMap, CategoryToPostColOrder)

		returningCols := q.selectCategoryToPostCols(selects, omits)

		scanFunc := func(res *CategoryToPost, cols []string) []any {
			return res.ScanFields(cols)
		}

		pkCols := []string{
			"postId",
			"categoryId",
		}

		hasRelations := selects.hasAnyRelation()

		var res *CategoryToPost
		var err error
		if hasRelations {
			err = q.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = executeInsert(c, txQ, "CategoryToPost", cols, vals, returningCols, pkCols, scanFunc, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.loadCategoryToPostRelations(c, []*CategoryToPost{res}, selects)
			})
		} else {
			res, err = executeInsert(c, q, "CategoryToPost", cols, vals, returningCols, pkCols, scanFunc, conflictTarget, conflictAction)
		}
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	for _, ext := range slices.Backward(q.CategoryToPost.extensions) {
		if ext.Create != nil {
			next, hook := curr, ext.Create
			curr = func(c context.Context, input *CategoryToPostCreate) (*CategoryToPost, error) {
				return hook(c, input, next)
			}
		}
	}

	return curr(ctx, &input)
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

func (d *CategoryToPostDelegate) CreateMany(builders ...*CategoryToPostCreateBuilder) *CategoryToPostCreateManyBuilder {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &CategoryToPostCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[CategoryToPost]{
			client:   d.client,
			records:  records,
			execFunc: d.client.executeCategoryToPostCreateMany,
		},
	}
}

func (d *CategoryToPostDelegate) CreateManyAndReturn(builders ...*CategoryToPostCreateBuilder) *CategoryToPostCreateManyAndReturnBuilder {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &CategoryToPostCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
			client:   d.client,
			records:  records,
			execFunc: d.client.executeCategoryToPostCreateManyAndReturn,
		},
	}
}

func (q *Queries) executeCategoryToPostCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs := make([]*CategoryToPostCreate, len(records))
	for i, rec := range records {
		if err := validateCategoryToPostCreate(rec.Assignments); err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		input := assignmentsToCategoryToPostCreate(rec.Assignments)
		inputs[i] = &input
	}

	curr := func(c context.Context, args []*CategoryToPostCreate) (int64, error) {
		rowMaps := make([]map[string]any, len(args))
		for i, input := range args {
			rowMaps[i] = input.ToRowMap()
		}

		pkCols := []string{
			"postId",
			"categoryId",
		}

		return executeCreateMany(c, q, rowMaps, "CategoryToPost", CategoryToPostColOrder, pkCols, conflictTarget, conflictAction)
	}

	for _, ext := range slices.Backward(q.CategoryToPost.extensions) {
		if ext.CreateMany != nil {
			next, hook := curr, ext.CreateMany
			curr = func(c context.Context, inputs []*CategoryToPostCreate) (int64, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
}

func (q *Queries) executeCategoryToPostCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *CategoryToPostSelect, omits *CategoryToPostOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*CategoryToPost, error) {
	inputs := make([]*CategoryToPostCreate, len(records))
	for i, rec := range records {
		if err := validateCategoryToPostCreate(rec.Assignments); err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		input := assignmentsToCategoryToPostCreate(rec.Assignments)
		inputs[i] = &input
	}

	curr := func(c context.Context, args []*CategoryToPostCreate) ([]*CategoryToPost, error) {
		rowMaps := make([]map[string]any, len(args))
		for i, input := range args {
			rowMaps[i] = input.ToRowMap()
		}

		pkCols := []string{
			"postId",
			"categoryId",
		}

		return executeCreateManyAndReturn(c, q, rowMaps, "CategoryToPost", CategoryToPostColOrder, selects, omits,
			q.selectCategoryToPostCols,
			q.loadCategoryToPostRelations,
			(*CategoryToPost).ScanFields,
			(*CategoryToPostSelect).hasAnyRelation,
			pkCols,
			conflictTarget,
			conflictAction,
		)
	}

	for _, ext := range slices.Backward(q.CategoryToPost.extensions) {
		if ext.CreateManyAndReturn != nil {
			next, hook := curr, ext.CreateManyAndReturn
			curr = func(c context.Context, inputs []*CategoryToPostCreate) ([]*CategoryToPost, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
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
	var up ConflictUpdate
	u := newCategoryToPostUpsert(&up)
	fn(u)
	cb.setAction(ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}, cb.conflictTarget)
	return cb.builder
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
func (d *CategoryToPostDelegate) FindUnique(where UniquePredicate[CategoryToPost], additional ...PredicateOf[CategoryToPost]) *FindUniqueBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit] {
	return &FindUniqueBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
		client:     d.client,
		where:      where,
		additional: additional,
		execFunc:   d.client.executeCategoryToPostFindUnique,
	}
}

func (d *CategoryToPostDelegate) FindFirst(preds ...PredicateOf[CategoryToPost]) *FindFirstBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit] {
	return &FindFirstBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeCategoryToPostFindFirst,
	}
}

func (d *CategoryToPostDelegate) FindMany(preds ...PredicateOf[CategoryToPost]) *FindManyBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit] {
	return &FindManyBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeCategoryToPostFindMany,
	}
}

func (q *Queries) executeCategoryToPostFindUnique(ctx context.Context, where UniquePredicate[CategoryToPost], additional []PredicateOf[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error) {
	curr := func(c context.Context, w UniquePredicate[CategoryToPost], add []PredicateOf[CategoryToPost], sel *CategoryToPostSelect, o *CategoryToPostOmit) (*CategoryToPost, error) {
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
		allPreds := append([]PredicateOf[CategoryToPost]{w}, add...)
		whereClause, vals := CompilePredicates(q.dialect, allPreds)
		if whereClause != "" {
			whereClause = " WHERE " + whereClause
		}
		returningCols := q.selectCategoryToPostCols(sel, o)
		return executeSingleWithRelations(c, q, "CategoryToPost", whereClause, vals, returningCols,
			func(res *CategoryToPost, cols []string) []any { return res.ScanFields(cols) },
			sel.hasAnyRelation(),
			func(ctx context.Context, txQ *Queries, results []*CategoryToPost) error {
				return txQ.loadCategoryToPostRelations(ctx, results, sel)
			},
			nil,
		)
	}

	for _, ext := range slices.Backward(q.CategoryToPost.extensions) {
		if ext.FindUnique != nil {
			next, hook := curr, ext.FindUnique
			curr = func(c context.Context, w UniquePredicate[CategoryToPost], add []PredicateOf[CategoryToPost], sel *CategoryToPostSelect, o *CategoryToPostOmit) (*CategoryToPost, error) {
				return hook(c, w, add, sel, o, next)
			}
		}
	}

	return curr(ctx, where, additional, selects, omits)
}

func (q *Queries) executeCategoryToPostFindFirst(
	ctx context.Context,
	params QueryParams[CategoryToPost],
	selects *CategoryToPostSelect,
	omits *CategoryToPostOmit,
) (*CategoryToPost, error) {
	curr := func(c context.Context, p QueryParams[CategoryToPost], sel *CategoryToPostSelect, o *CategoryToPostOmit) (*CategoryToPost, error) {
		for _, pr := range p.Where {
			if pr != nil {
				if err := pr.Validate(); err != nil {
					return nil, err
				}
			}
		}
		whereClause, vals := CompilePredicates(q.dialect, p.Where)
		if whereClause != "" {
			whereClause = " WHERE " + whereClause
		}
		returningCols := q.selectCategoryToPostCols(sel, o)
		return executeSingleWithRelations(c, q, "CategoryToPost", whereClause, vals, returningCols,
			func(res *CategoryToPost, cols []string) []any { return res.ScanFields(cols) },
			sel.hasAnyRelation(),
			func(ctx context.Context, txQ *Queries, results []*CategoryToPost) error {
				return txQ.loadCategoryToPostRelations(ctx, results, sel)
			},
			p.Skip,
		)
	}

	for _, ext := range slices.Backward(q.CategoryToPost.extensions) {
		if ext.FindFirst != nil {
			next, hook := curr, ext.FindFirst
			curr = func(c context.Context, p QueryParams[CategoryToPost], sel *CategoryToPostSelect, o *CategoryToPostOmit) (*CategoryToPost, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}

func (q *Queries) executeCategoryToPostFindMany(
	ctx context.Context,
	params QueryParams[CategoryToPost],
	selects *CategoryToPostSelect,
	omits *CategoryToPostOmit,
) ([]*CategoryToPost, error) {
	curr := func(c context.Context, p QueryParams[CategoryToPost], sel *CategoryToPostSelect, o *CategoryToPostOmit) ([]*CategoryToPost, error) {
		for _, pr := range p.Where {
			if pr != nil {
				if err := pr.Validate(); err != nil {
					return nil, err
				}
			}
		}
		whereClause, vals := CompilePredicates(q.dialect, p.Where)
		if whereClause != "" {
			whereClause = " WHERE " + whereClause
		}
		returningCols := q.selectCategoryToPostCols(sel, o)
		return executeManyWithRelations(c, q, "CategoryToPost", whereClause, vals, returningCols,
			func(res *CategoryToPost, cols []string) []any { return res.ScanFields(cols) },
			sel.hasAnyRelation(),
			func(ctx context.Context, txQ *Queries, results []*CategoryToPost) error {
				return txQ.loadCategoryToPostRelations(ctx, results, sel)
			},
			p.Take,
			p.Skip,
		)
	}

	for _, ext := range slices.Backward(q.CategoryToPost.extensions) {
		if ext.FindMany != nil {
			next, hook := curr, ext.FindMany
			curr = func(c context.Context, p QueryParams[CategoryToPost], sel *CategoryToPostSelect, o *CategoryToPostOmit) ([]*CategoryToPost, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}
func (q *Queries) loadCategoryToPostRelations(ctx context.Context, records []*CategoryToPost, selects *CategoryToPostSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Post != nil {
		relationSelects, relationOmits, relationParams := selects.Post.GetRelationParams()
		returningCols := q.selectPostCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: CategoryToPost.postId
		allChildren, err := loadRelation(
			ctx, q, records,
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
		if err := q.loadPostRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Category != nil {
		relationSelects, relationOmits, relationParams := selects.Category.GetRelationParams()
		returningCols := q.selectCategoryCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: CategoryToPost.categoryId
		allChildren, err := loadRelation(
			ctx, q, records,
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
		if err := q.loadCategoryRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
