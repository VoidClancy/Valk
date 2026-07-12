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
	Id    bool                  `json:"id"`
	Name  bool                  `json:"name"`
	Posts *CategoryToPostSelect `json:"posts,omitempty"`
}

// CategoryOmit specifies which fields to exclude
type CategoryOmit struct {
	Id    bool                `json:"id"`
	Name  bool                `json:"name"`
	Posts *CategoryToPostOmit `json:"posts,omitempty"`
}

type CategoryDelegate struct {
	client          *Queries
	beforeCreate    func(context.Context, *CategoryCreate) error
	afterCreate     func(context.Context, []*Category) error
	afterCreateMany func(context.Context, []CategoryCreate, int64) error
}

func (d *CategoryDelegate) BeforeCreate(hook func(context.Context, *CategoryCreate) error) {
	d.beforeCreate = hook
}

func (d *CategoryDelegate) AfterCreate(hook func(context.Context, []*Category) error) {
	d.afterCreate = hook
}

func (d *CategoryDelegate) AfterCreateMany(hook func(context.Context, []CategoryCreate, int64) error) {
	d.afterCreateMany = hook
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

func (q *Queries) selectCategoryCols(selects *CategorySelect, omits *CategoryOmit, forceCols ...string) []string {
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

var CategoryColOrder = []string{
	"id",
	"name",
}

func (s *CategorySelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Posts != nil
}

func (d *CategoryDelegate) Create(assignments ...FieldAssignment) *CreateBuilder[Category, CategorySelect, CategoryOmit] {
	return &CreateBuilder[Category, CategorySelect, CategoryOmit]{
		client:      d.client,
		assignments: assignments,
		execFunc:    d.client.executeCategoryCreate,
	}
}

func validateCategoryCreate(assignments []FieldAssignment) error {
	errs := &ValidationError{}

	provided := make(map[string]bool)
	for _, a := range assignments {
		provided[a.Col] = true
		switch a.Col {
		case "id":
			if v, ok := a.Val.(int32); ok {
				ValidateInt32(errs, "id", v, "")
			} else {
				errs.Add("id", a.Val, "type", "field id must be of type int32")
			}
		case "name":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "name", v, true, 0, false, false)
			} else {
				errs.Add("name", a.Val, "type", "field name must be of type string")
			}
		}
	}
	if !provided["name"] {
		errs.Add("name", "", "required", "field Name is required")
	}

	if errs.HasErrors() {
		return *errs
	}
	return nil
}

func assignmentsToCategoryCreate(assignments []FieldAssignment) CategoryCreate {
	var input CategoryCreate
	for _, a := range assignments {
		switch a.Col {
		case "id":
			if v, ok := a.Val.(int32); ok {
				input.Id = &v
			}
		case "name":
			if v, ok := a.Val.(string); ok {
				input.Name = v
			}
		}
	}
	return input
}

func (s *CategoryCreate) ToRowMap() map[string]any {
	m := make(map[string]any, 2)
	if s.Id != nil {
		m["id"] = *s.Id
	}
	m["name"] = s.Name
	return m
}

func (q *Queries) executeCategoryCreate(ctx context.Context, assignments []FieldAssignment, selects *CategorySelect, omits *CategoryOmit) (*Category, error) {
	input := assignmentsToCategoryCreate(assignments)

	if q.Category.beforeCreate != nil {
		if err := q.Category.beforeCreate(ctx, &input); err != nil {
			return nil, err
		}
	}

	if err := validateCategoryCreate(assignments); err != nil {
		return nil, err
	}

	rowMap := input.ToRowMap()
	var cols []string
	var vals []any
	for _, col := range CategoryColOrder {
		if val, ok := rowMap[col]; ok {
			cols = append(cols, col)
			vals = append(vals, val)
		}
	}

	returningCols := q.selectCategoryCols(selects, omits)

	scanFunc := func(res *Category, cols []string) []any {
		return res.ScanFields(cols)
	}

	idCol := "id"

	hasRelations := selects.hasAnyRelation()

	var res *Category
	var err error
	if hasRelations {
		err = q.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = executeInsert(ctx, txQ, "Category", cols, vals, returningCols, idCol, scanFunc)
			if err != nil {
				return err
			}
			return txQ.loadCategoryRelations(ctx, []*Category{res}, selects)
		})
	} else {
		res, err = executeInsert(ctx, q, "Category", cols, vals, returningCols, idCol, scanFunc)
	}
	if err != nil {
		return nil, err
	}

	if q.Category.afterCreate != nil {
		if err := q.Category.afterCreate(ctx, []*Category{res}); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (d *CategoryDelegate) CreateMany(records ...RecordInput) *CreateManyBuilder[Category] {
	return &CreateManyBuilder[Category]{
		client:   d.client,
		records:  records,
		execFunc: d.client.executeCategoryCreateMany,
	}
}

func (d *CategoryDelegate) CreateManyAndReturn(records ...RecordInput) *CreateManyAndReturnBuilder[Category, CategorySelect, CategoryOmit] {
	return &CreateManyAndReturnBuilder[Category, CategorySelect, CategoryOmit]{
		client:   d.client,
		records:  records,
		execFunc: d.client.executeCategoryCreateManyAndReturn,
	}
}

func (q *Queries) executeCategoryCreateMany(ctx context.Context, records []RecordInput) (int64, error) {
	rowMaps := make([]map[string]any, len(records))
	inputs := make([]CategoryCreate, len(records))
	for i, rec := range records {
		if err := validateCategoryCreate(rec.Assignments); err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		input := assignmentsToCategoryCreate(rec.Assignments)
		if q.Category.beforeCreate != nil {
			if err := q.Category.beforeCreate(ctx, &input); err != nil {
				return 0, err
			}
		}
		rowMaps[i] = input.ToRowMap()
		inputs[i] = input
	}
	count, err := executeCreateMany(ctx, q, rowMaps, "Category", CategoryColOrder)
	if err != nil {
		return 0, err
	}
	if q.Category.afterCreateMany != nil {
		if err := q.Category.afterCreateMany(ctx, inputs, count); err != nil {
			return 0, err
		}
	}
	return count, nil
}

func (q *Queries) executeCategoryCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *CategorySelect, omits *CategoryOmit) ([]*Category, error) {
	rowMaps := make([]map[string]any, len(records))
	idCol := "id"
	for i, rec := range records {
		if err := validateCategoryCreate(rec.Assignments); err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		input := assignmentsToCategoryCreate(rec.Assignments)
		if q.Category.beforeCreate != nil {
			if err := q.Category.beforeCreate(ctx, &input); err != nil {
				return nil, err
			}
		}
		rowMaps[i] = input.ToRowMap()
	}
	results, err := executeCreateManyAndReturn(ctx, q, rowMaps, "Category", CategoryColOrder, selects, omits,
		q.selectCategoryCols,
		q.loadCategoryRelations,
		(*Category).ScanFields,
		(*CategorySelect).hasAnyRelation,
		idCol,
	)
	if err != nil {
		return nil, err
	}
	if q.Category.afterCreate != nil {
		if err := q.Category.afterCreate(ctx, results); err != nil {
			return nil, err
		}
	}
	return results, nil
}
func (d *CategoryDelegate) FindUnique(where UniquePredicate) *FindUniqueBuilder[Category, CategorySelect, CategoryOmit] {
	return &FindUniqueBuilder[Category, CategorySelect, CategoryOmit]{
		client:   d.client,
		where:    where,
		execFunc: d.client.executeCategoryFindUnique,
	}
}

func (d *CategoryDelegate) FindFirst(preds ...Predicate) *FindFirstBuilder[Category, CategorySelect, CategoryOmit] {
	return &FindFirstBuilder[Category, CategorySelect, CategoryOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeCategoryFindFirst,
	}
}

func (d *CategoryDelegate) FindMany(preds ...Predicate) *FindManyBuilder[Category, CategorySelect, CategoryOmit] {
	return &FindManyBuilder[Category, CategorySelect, CategoryOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeCategoryFindMany,
	}
}

func (q *Queries) executeCategoryFindUnique(ctx context.Context, where UniquePredicate, selects *CategorySelect, omits *CategoryOmit) (*Category, error) {
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
	returningCols := q.selectCategoryCols(selects, omits)
	return executeSingleWithRelations(ctx, q, "Category", whereClause, vals, returningCols,
		func(res *Category, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*Category) error {
			return txQ.loadCategoryRelations(ctx, results, selects)
		},
		nil,
	)
}

func (q *Queries) executeCategoryFindFirst(
	ctx context.Context,
	params QueryParams,
	selects *CategorySelect,
	omits *CategoryOmit,
) (*Category, error) {
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
	returningCols := q.selectCategoryCols(selects, omits)
	return executeSingleWithRelations(ctx, q, "Category", whereClause, vals, returningCols,
		func(res *Category, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*Category) error {
			return txQ.loadCategoryRelations(ctx, results, selects)
		},
		params.Skip,
	)
}

func (q *Queries) executeCategoryFindMany(
	ctx context.Context,
	params QueryParams,
	selects *CategorySelect,
	omits *CategoryOmit,
) ([]*Category, error) {
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
	returningCols := q.selectCategoryCols(selects, omits)
	return executeManyWithRelations(ctx, q, "Category", whereClause, vals, returningCols,
		func(res *Category, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*Category) error {
			return txQ.loadCategoryRelations(ctx, results, selects)
		},
		params.Take,
		params.Skip,
	)
}
func (q *Queries) loadCategoryRelations(ctx context.Context, records []*Category, selects *CategorySelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Posts != nil {
		returningCols := q.selectCategoryToPostCols(selects.Posts, nil, "categoryId")
		// Inverse holds the FK: CategoryToPost.categoryId
		allChildren, err := loadRelation(
			ctx, q, records,
			directKey(func(p *Category) int32 { return p.Id }),
			"CategoryToPost",
			"categoryId",
			returningCols,
			scanInto(returningCols, (*CategoryToPost).ScanFields),
			directKey(func(c *CategoryToPost) int32 { return c.CategoryId }),
			appendMany(func(p *Category) *[]*CategoryToPost { return &p.Posts }),
		)
		if err != nil {
			return fmt.Errorf("loading posts: %w", err)
		}
		if err := q.loadCategoryToPostRelations(ctx, allChildren, selects.Posts); err != nil {
			return err
		}
	}

	return nil
}
