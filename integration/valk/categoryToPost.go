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

type CategoryToPostDelegate struct {
	client          *Queries
	beforeCreate    func(context.Context, *CategoryToPostCreate) error
	afterCreate     func(context.Context, []*CategoryToPost) error
	afterCreateMany func(context.Context, []CategoryToPostCreate, int64) error
}

func (d *CategoryToPostDelegate) BeforeCreate(hook func(context.Context, *CategoryToPostCreate) error) {
	d.beforeCreate = hook
}

func (d *CategoryToPostDelegate) AfterCreate(hook func(context.Context, []*CategoryToPost) error) {
	d.afterCreate = hook
}

func (d *CategoryToPostDelegate) AfterCreateMany(hook func(context.Context, []CategoryToPostCreate, int64) error) {
	d.afterCreateMany = hook
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

func (q *Queries) executeCategoryToPostCreate(ctx context.Context, assignments []FieldAssignment, selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error) {
	input := assignmentsToCategoryToPostCreate(assignments)

	if q.CategoryToPost.beforeCreate != nil {
		if err := q.CategoryToPost.beforeCreate(ctx, &input); err != nil {
			return nil, err
		}
	}

	if err := validateCategoryToPostCreate(assignments); err != nil {
		return nil, err
	}

	rowMap := input.ToRowMap()
	cols, vals := mapToColsVals(rowMap, CategoryToPostColOrder)

	returningCols := q.selectCategoryToPostCols(selects, omits)

	scanFunc := func(res *CategoryToPost, cols []string) []any {
		return res.ScanFields(cols)
	}

	idCol := ""

	hasRelations := selects.hasAnyRelation()

	var res *CategoryToPost
	var err error
	if hasRelations {
		err = q.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = executeInsert(ctx, txQ, "CategoryToPost", cols, vals, returningCols, idCol, scanFunc)
			if err != nil {
				return err
			}
			return txQ.loadCategoryToPostRelations(ctx, []*CategoryToPost{res}, selects)
		})
	} else {
		res, err = executeInsert(ctx, q, "CategoryToPost", cols, vals, returningCols, idCol, scanFunc)
	}
	if err != nil {
		return nil, err
	}

	if q.CategoryToPost.afterCreate != nil {
		if err := q.CategoryToPost.afterCreate(ctx, []*CategoryToPost{res}); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (d *CategoryToPostDelegate) CreateMany(builders ...*CategoryToPostCreateBuilder) *CreateManyBuilder[CategoryToPost] {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &CreateManyBuilder[CategoryToPost]{
		client:   d.client,
		records:  records,
		execFunc: d.client.executeCategoryToPostCreateMany,
	}
}

func (d *CategoryToPostDelegate) CreateManyAndReturn(builders ...*CategoryToPostCreateBuilder) *CreateManyAndReturnBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit] {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &CreateManyAndReturnBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
		client:   d.client,
		records:  records,
		execFunc: d.client.executeCategoryToPostCreateManyAndReturn,
	}
}

func (q *Queries) executeCategoryToPostCreateMany(ctx context.Context, records []RecordInput) (int64, error) {
	rowMaps := make([]map[string]any, len(records))
	inputs := make([]CategoryToPostCreate, len(records))
	for i, rec := range records {
		if err := validateCategoryToPostCreate(rec.Assignments); err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		input := assignmentsToCategoryToPostCreate(rec.Assignments)
		if q.CategoryToPost.beforeCreate != nil {
			if err := q.CategoryToPost.beforeCreate(ctx, &input); err != nil {
				return 0, err
			}
		}
		rowMaps[i] = input.ToRowMap()
		inputs[i] = input
	}
	count, err := executeCreateMany(ctx, q, rowMaps, "CategoryToPost", CategoryToPostColOrder)
	if err != nil {
		return 0, err
	}
	if q.CategoryToPost.afterCreateMany != nil {
		if err := q.CategoryToPost.afterCreateMany(ctx, inputs, count); err != nil {
			return 0, err
		}
	}
	return count, nil
}

func (q *Queries) executeCategoryToPostCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *CategoryToPostSelect, omits *CategoryToPostOmit) ([]*CategoryToPost, error) {
	rowMaps := make([]map[string]any, len(records))
	idCol := ""
	for i, rec := range records {
		if err := validateCategoryToPostCreate(rec.Assignments); err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		input := assignmentsToCategoryToPostCreate(rec.Assignments)
		if q.CategoryToPost.beforeCreate != nil {
			if err := q.CategoryToPost.beforeCreate(ctx, &input); err != nil {
				return nil, err
			}
		}
		rowMaps[i] = input.ToRowMap()
	}
	results, err := executeCreateManyAndReturn(ctx, q, rowMaps, "CategoryToPost", CategoryToPostColOrder, selects, omits,
		q.selectCategoryToPostCols,
		q.loadCategoryToPostRelations,
		(*CategoryToPost).ScanFields,
		(*CategoryToPostSelect).hasAnyRelation,
		idCol,
	)
	if err != nil {
		return nil, err
	}
	if q.CategoryToPost.afterCreate != nil {
		if err := q.CategoryToPost.afterCreate(ctx, results); err != nil {
			return nil, err
		}
	}
	return results, nil
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
	allPreds := append([]PredicateOf[CategoryToPost]{where}, additional...)
	whereClause, vals := CompilePredicates(q.dialect, allPreds)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	returningCols := q.selectCategoryToPostCols(selects, omits)
	return executeSingleWithRelations(ctx, q, "CategoryToPost", whereClause, vals, returningCols,
		func(res *CategoryToPost, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*CategoryToPost) error {
			return txQ.loadCategoryToPostRelations(ctx, results, selects)
		},
		nil,
	)
}

func (q *Queries) executeCategoryToPostFindFirst(
	ctx context.Context,
	params QueryParams[CategoryToPost],
	selects *CategoryToPostSelect,
	omits *CategoryToPostOmit,
) (*CategoryToPost, error) {
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
	returningCols := q.selectCategoryToPostCols(selects, omits)
	return executeSingleWithRelations(ctx, q, "CategoryToPost", whereClause, vals, returningCols,
		func(res *CategoryToPost, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*CategoryToPost) error {
			return txQ.loadCategoryToPostRelations(ctx, results, selects)
		},
		params.Skip,
	)
}

func (q *Queries) executeCategoryToPostFindMany(
	ctx context.Context,
	params QueryParams[CategoryToPost],
	selects *CategoryToPostSelect,
	omits *CategoryToPostOmit,
) ([]*CategoryToPost, error) {
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
	returningCols := q.selectCategoryToPostCols(selects, omits)
	return executeManyWithRelations(ctx, q, "CategoryToPost", whereClause, vals, returningCols,
		func(res *CategoryToPost, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*CategoryToPost) error {
			return txQ.loadCategoryToPostRelations(ctx, results, selects)
		},
		params.Take,
		params.Skip,
	)
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
