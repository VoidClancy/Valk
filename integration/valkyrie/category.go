package valkyrie

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"time"
)

var _ = time.Time{}
var _ = fmt.Sprintf
var _ = strings.Join
var _ = context.Background
var _ = sql.LevelDefault
var _ = slices.Contains[[]string, string]

// Category represents the database model
type Category struct {
	Id    int32             `db:"id" json:"id"`
	Name  string            `db:"name" json:"name"`
	Posts []*CategoryToPost `json:"posts,omitempty"`
}

// CategoryCreateInput represents the input structure for creation
type CategoryCreateInput struct {
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
	client *Queries
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
		{"id", selects != nil && selects.Id, omits != nil && omits.Id, false},
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

func (input CategoryCreateInput) Validate() error {
	if input.Name == "" {
		return fmt.Errorf("field Name is required")
	}
	return nil
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

func (d *CategoryDelegate) Create(input CategoryCreateInput) *CreateBuilder[Category, CategoryCreateInput, CategorySelect, CategoryOmit] {
	return &CreateBuilder[Category, CategoryCreateInput, CategorySelect, CategoryOmit]{
		client:   d.client,
		input:    input,
		execFunc: d.client.executeCategoryCreate,
	}
}

func (q *Queries) executeCategoryCreate(ctx context.Context, input CategoryCreateInput, selects *CategorySelect, omits *CategoryOmit) (*Category, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	m := q.CategoryInputToMap(input)
	cols, vals := mapToColsVals(m, CategoryColOrder)

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

	return res, nil
}

func (q *Queries) CategoryInputToMap(input CategoryCreateInput) map[string]any {
	m := make(map[string]any)
	if input.Id != nil {
		m["id"] = *input.Id
	} else {
	}
	m["name"] = input.Name
	return m
}

func (d *CategoryDelegate) CreateMany(inputs []CategoryCreateInput) *CreateManyBuilder[Category, CategoryCreateInput] {
	return &CreateManyBuilder[Category, CategoryCreateInput]{
		client:   d.client,
		inputs:   inputs,
		execFunc: d.client.executeCategoryCreateMany,
	}
}

func (d *CategoryDelegate) CreateManyAndReturn(inputs []CategoryCreateInput) *CreateManyAndReturnBuilder[Category, CategoryCreateInput, CategorySelect, CategoryOmit] {
	return &CreateManyAndReturnBuilder[Category, CategoryCreateInput, CategorySelect, CategoryOmit]{
		client:   d.client,
		inputs:   inputs,
		execFunc: d.client.executeCategoryCreateManyAndReturn,
	}
}

func (q *Queries) executeCategoryCreateMany(ctx context.Context, inputs []CategoryCreateInput) (int64, error) {
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
			rowMaps[i] = q.CategoryInputToMap(input)
		}
		query, vals := buildBulkInsertSQL(q.dialect, "Category", rowMaps, CategoryColOrder, nil)
		res, err := q.exec(ctx, query, vals...)
		if err != nil {
			return 0, err
		}
		return res.RowsAffected()
	}

	var count int64
	err := q.transaction(ctx, func(txQ *Queries) error {
		for _, input := range inputs {
			_, err := txQ.executeCategoryCreate(ctx, input, nil, nil)
			if err != nil {
				return err
			}
			count++
		}
		return nil
	})
	return count, err
}

func (q *Queries) executeCategoryCreateManyAndReturn(ctx context.Context, inputs []CategoryCreateInput, selects *CategorySelect, omits *CategoryOmit) ([]*Category, error) {
	if len(inputs) == 0 {
		return nil, nil
	}
	for i, input := range inputs {
		if err := input.Validate(); err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
	}

	hasRelations := selects.hasAnyRelation()
	returningCols := q.selectCategoryCols(selects, omits)

	if q.dialect.SupportsBulkInsert() {
		rowMaps := make([]map[string]any, len(inputs))
		for i, input := range inputs {
			rowMaps[i] = q.CategoryInputToMap(input)
		}
		query, vals := buildBulkInsertSQL(q.dialect, "Category", rowMaps, CategoryColOrder, returningCols)
		var records []*Category
		err := q.transaction(ctx, func(txQ *Queries) error {
			rows, err := txQ.query(ctx, query, vals...)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var record Category
				if err := rows.Scan(record.ScanFields(returningCols)...); err != nil {
					return err
				}
				records = append(records, &record)
			}
			if err := rows.Err(); err != nil {
				return err
			}
			if hasRelations {
				return txQ.loadCategoryRelations(ctx, records, selects)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return records, nil
	}

	// Fallback to loop inside transaction
	var records []*Category
	err := q.transaction(ctx, func(txQ *Queries) error {
		for _, input := range inputs {
			res, err := txQ.executeCategoryCreate(ctx, input, nil, nil)
			if err != nil {
				return err
			}
			records = append(records, res)
		}

		if hasRelations {
			return txQ.loadCategoryRelations(ctx, records, selects)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return records, nil
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
