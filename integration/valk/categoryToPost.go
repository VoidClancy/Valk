package valk

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"
)

// CategoryToPost represents the database model
type CategoryToPost struct {
	PostId     string    `db:"postId" json:"postId"`
	CategoryId int32     `db:"categoryId" json:"categoryId"`
	Post       *Post     `json:"post,omitempty"`
	Category   *Category `json:"category,omitempty"`
}

// CategoryToPostCreate represents the input structure for creation
type CategoryToPostCreate struct {
	PostId     string `json:"postId"`
	CategoryId int32  `json:"categoryId"`
}

// CategoryToPostSelect specifies which fields to include
type CategoryToPostSelect struct {
	PostId     bool            `json:"postId"`
	CategoryId bool            `json:"categoryId"`
	Post       *PostSelect     `json:"post,omitempty"`
	Category   *CategorySelect `json:"category,omitempty"`
}

// CategoryToPostOmit specifies which fields to exclude
type CategoryToPostOmit struct {
	PostId     bool          `json:"postId"`
	CategoryId bool          `json:"categoryId"`
	Post       *PostOmit     `json:"post,omitempty"`
	Category   *CategoryOmit `json:"category,omitempty"`
}

type CategoryToPostDelegate struct {
	client       *Queries
	beforeCreate func(context.Context, *CategoryToPostCreate) error
	afterCreate  func(context.Context, *CategoryToPost) error
}

func (d *CategoryToPostDelegate) BeforeCreate(hook func(context.Context, *CategoryToPostCreate) error) {
	d.beforeCreate = hook
}

func (d *CategoryToPostDelegate) AfterCreate(hook func(context.Context, *CategoryToPost) error) {
	d.afterCreate = hook
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

func (input CategoryToPostCreate) Validate() error {
	errs := &ValidationError{}
	if input.PostId == "" {
		errs.Add("postId", input.PostId, "required", "field PostId is required")
	}
	if strings.Contains(input.PostId, "\x00") {
		errs.Add("postId", input.PostId, "safety", "string cannot contain null bytes")
	}
	if !utf8.ValidString(input.PostId) {
		errs.Add("postId", input.PostId, "safety", "string must be valid UTF-8")
	}

	if errs.HasErrors() {
		return *errs
	}
	return nil
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

func (d *CategoryToPostDelegate) Create(input CategoryToPostCreate) *CreateBuilder[CategoryToPost, CategoryToPostCreate, CategoryToPostSelect, CategoryToPostOmit] {
	return &CreateBuilder[CategoryToPost, CategoryToPostCreate, CategoryToPostSelect, CategoryToPostOmit]{
		client:   d.client,
		input:    input,
		execFunc: d.client.executeCategoryToPostCreate,
	}
}

func (q *Queries) executeCategoryToPostCreate(ctx context.Context, input CategoryToPostCreate, selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error) {
	if q.CategoryToPost.beforeCreate != nil {
		if err := q.CategoryToPost.beforeCreate(ctx, &input); err != nil {
			return nil, err
		}
	}

	if err := input.Validate(); err != nil {
		return nil, err
	}
	var cols []string
	var vals []any
	cols = append(cols, "postId")
	vals = append(vals, input.PostId)
	cols = append(cols, "categoryId")
	vals = append(vals, input.CategoryId)

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
		if err := q.CategoryToPost.afterCreate(ctx, res); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (q *Queries) CategoryToPostInputToMap(input CategoryToPostCreate) map[string]any {
	m := make(map[string]any)
	m["postId"] = input.PostId
	m["categoryId"] = input.CategoryId
	return m
}

func (d *CategoryToPostDelegate) CreateMany(inputs []CategoryToPostCreate) *CreateManyBuilder[CategoryToPost, CategoryToPostCreate] {
	return &CreateManyBuilder[CategoryToPost, CategoryToPostCreate]{
		client:   d.client,
		inputs:   inputs,
		execFunc: d.client.executeCategoryToPostCreateMany,
	}
}

func (d *CategoryToPostDelegate) CreateManyAndReturn(inputs []CategoryToPostCreate) *CreateManyAndReturnBuilder[CategoryToPost, CategoryToPostCreate, CategoryToPostSelect, CategoryToPostOmit] {
	return &CreateManyAndReturnBuilder[CategoryToPost, CategoryToPostCreate, CategoryToPostSelect, CategoryToPostOmit]{
		client:   d.client,
		inputs:   inputs,
		execFunc: d.client.executeCategoryToPostCreateManyAndReturn,
	}
}

func (q *Queries) executeCategoryToPostCreateMany(ctx context.Context, inputs []CategoryToPostCreate) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}
	for i, input := range inputs {
		if err := input.Validate(); err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
	}

	if q.dialect.SupportsBulkInsert() {
		rowMaps := make([]map[string]any, len(inputs))
		for i, input := range inputs {
			rowMaps[i] = q.CategoryToPostInputToMap(input)
		}
		query, vals := buildBulkInsertSQL(q.dialect, "CategoryToPost", rowMaps, CategoryToPostColOrder, nil)
		res, err := q.exec(ctx, query, vals...)
		if err != nil {
			return 0, err
		}
		return res.RowsAffected()
	}

	var count int64
	err := q.transaction(ctx, func(txQ *Queries) error {
		for _, input := range inputs {
			_, err := txQ.executeCategoryToPostCreate(ctx, input, nil, nil)
			if err != nil {
				return err
			}
			count++
		}
		return nil
	})
	return count, err
}

func (q *Queries) executeCategoryToPostCreateManyAndReturn(ctx context.Context, inputs []CategoryToPostCreate, selects *CategoryToPostSelect, omits *CategoryToPostOmit) ([]*CategoryToPost, error) {
	if len(inputs) == 0 {
		return nil, nil
	}
	for i, input := range inputs {
		if err := input.Validate(); err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
	}

	hasRelations := selects.hasAnyRelation()
	returningCols := q.selectCategoryToPostCols(selects, omits)

	if q.dialect.SupportsBulkInsert() {
		rowMaps := make([]map[string]any, len(inputs))
		for i, input := range inputs {
			rowMaps[i] = q.CategoryToPostInputToMap(input)
		}
		query, vals := buildBulkInsertSQL(q.dialect, "CategoryToPost", rowMaps, CategoryToPostColOrder, returningCols)
		records := make([]*CategoryToPost, 0)
		err := q.transaction(ctx, func(txQ *Queries) error {
			rows, err := txQ.query(ctx, query, vals...)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var record CategoryToPost
				if err := rows.Scan(record.ScanFields(returningCols)...); err != nil {
					return err
				}
				records = append(records, &record)
			}
			if err := rows.Err(); err != nil {
				return err
			}
			if hasRelations {
				return txQ.loadCategoryToPostRelations(ctx, records, selects)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return records, nil
	}

	// Fallback to loop inside transaction
	records := make([]*CategoryToPost, 0)
	err := q.transaction(ctx, func(txQ *Queries) error {
		for _, input := range inputs {
			res, err := txQ.executeCategoryToPostCreate(ctx, input, nil, nil)
			if err != nil {
				return err
			}
			records = append(records, res)
		}

		if hasRelations {
			return txQ.loadCategoryToPostRelations(ctx, records, selects)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}
func (d *CategoryToPostDelegate) FindUnique(where UniquePredicate) *FindUniqueBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit] {
	return &FindUniqueBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
		client:   d.client,
		where:    where,
		execFunc: d.client.executeCategoryToPostFindUnique,
	}
}

func (d *CategoryToPostDelegate) FindFirst(preds ...Predicate) *FindFirstBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit] {
	return &FindFirstBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeCategoryToPostFindFirst,
	}
}

func (d *CategoryToPostDelegate) FindMany(preds ...Predicate) *FindManyBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit] {
	return &FindManyBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeCategoryToPostFindMany,
	}
}

func (q *Queries) executeCategoryToPostFindUnique(ctx context.Context, where UniquePredicate, selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error) {
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
	returningCols := q.selectCategoryToPostCols(selects, omits)
	return executeSingleWithRelations(ctx, q, "CategoryToPost", whereClause, vals, returningCols,
		func(res *CategoryToPost, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*CategoryToPost) error {
			return txQ.loadCategoryToPostRelations(ctx, results, selects)
		},
	)
}

func (q *Queries) executeCategoryToPostFindFirst(ctx context.Context, where []Predicate, selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error) {
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
	returningCols := q.selectCategoryToPostCols(selects, omits)
	return executeSingleWithRelations(ctx, q, "CategoryToPost", whereClause, vals, returningCols,
		func(res *CategoryToPost, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*CategoryToPost) error {
			return txQ.loadCategoryToPostRelations(ctx, results, selects)
		},
	)
}

func (q *Queries) executeCategoryToPostFindMany(ctx context.Context, where []Predicate, selects *CategoryToPostSelect, omits *CategoryToPostOmit) ([]*CategoryToPost, error) {
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
	returningCols := q.selectCategoryToPostCols(selects, omits)
	return executeManyWithRelations(ctx, q, "CategoryToPost", whereClause, vals, returningCols,
		func(res *CategoryToPost, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*CategoryToPost) error {
			return txQ.loadCategoryToPostRelations(ctx, results, selects)
		},
	)
}
func (q *Queries) loadCategoryToPostRelations(ctx context.Context, records []*CategoryToPost, selects *CategoryToPostSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Post != nil {
		returningCols := q.selectPostCols(selects.Post, nil, "id")
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
		)
		if err != nil {
			return fmt.Errorf("loading post: %w", err)
		}
		if err := q.loadPostRelations(ctx, allChildren, selects.Post); err != nil {
			return err
		}
	}
	if selects.Category != nil {
		returningCols := q.selectCategoryCols(selects.Category, nil, "id")
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
		)
		if err != nil {
			return fmt.Errorf("loading category: %w", err)
		}
		if err := q.loadCategoryRelations(ctx, allChildren, selects.Category); err != nil {
			return err
		}
	}

	return nil
}
