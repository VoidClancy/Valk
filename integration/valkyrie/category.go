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
func (d *CategoryDelegate) Create(input CategoryCreateInput) *CreateBuilder[Category, CategoryCreateInput, CategorySelect, CategoryOmit] {
	return &CreateBuilder[Category, CategoryCreateInput, CategorySelect, CategoryOmit]{
		client:   d.client,
		input:    input,
		execFunc: d.client.executeCategoryCreate,
	}
}

func (q *Queries) executeCategoryCreate(ctx context.Context, input CategoryCreateInput, selects *CategorySelect, omits *CategoryOmit) (*Category, error) {
	var cols []string
	var vals []any
	if input.Id != nil {
		cols = append(cols, q.dialect.Quote("id"))
		vals = append(vals, *input.Id)
	} else {
		cols = append(cols, q.dialect.Quote("id"))
	}
	cols = append(cols, q.dialect.Quote("name"))
	vals = append(vals, input.Name)

	returningCols := q.selectCategoryCols(selects, omits)

	scanFunc := func(res *Category, cols []string) []any {
		return res.ScanFields(cols)
	}

	idCol := "id"

	res, err := executeInsert(ctx, q, "Category", cols, vals, returningCols, idCol, scanFunc)
	if err != nil {
		return nil, err
	}

	if err := q.loadCategoryRelations(ctx, []*Category{res}, selects); err != nil {
		return nil, err
	}

	return res, nil
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
