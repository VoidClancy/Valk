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

// CategoryToPost represents the database model
type CategoryToPost struct {
	PostId     string    `db:"postId" json:"postId"`
	CategoryId int32     `db:"categoryId" json:"categoryId"`
	Post       *Post     `json:"post,omitempty"`
	Category   *Category `json:"category,omitempty"`
}

// CategoryToPostCreateInput represents the input structure for creation
type CategoryToPostCreateInput struct {
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
	client *Queries
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
func (d *CategoryToPostDelegate) Create(input CategoryToPostCreateInput) *CreateBuilder[CategoryToPost, CategoryToPostCreateInput, CategoryToPostSelect, CategoryToPostOmit] {
	return &CreateBuilder[CategoryToPost, CategoryToPostCreateInput, CategoryToPostSelect, CategoryToPostOmit]{
		client:   d.client,
		input:    input,
		execFunc: d.client.executeCategoryToPostCreate,
	}
}

func (q *Queries) executeCategoryToPostCreate(ctx context.Context, input CategoryToPostCreateInput, selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error) {
	var cols []string
	var vals []any
	cols = append(cols, q.dialect.Quote("postId"))
	vals = append(vals, input.PostId)
	cols = append(cols, q.dialect.Quote("categoryId"))
	vals = append(vals, input.CategoryId)

	returningCols := q.selectCategoryToPostCols(selects, omits)

	scanFunc := func(res *CategoryToPost, cols []string) []any {
		return res.ScanFields(cols)
	}

	idCol := ""
	hasRelations := selects != nil && (selects.Post != nil || selects.Category != nil)

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

	return res, nil
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
