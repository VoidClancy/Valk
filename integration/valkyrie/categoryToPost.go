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

	var cols []string
	var anySelected bool
	if selects != nil {
		if selects.PostId {
			anySelected = true
		}
		if selects.CategoryId {
			anySelected = true
		}
		if selects.Post != nil {
			anySelected = true
		}
		if selects.Category != nil {
			anySelected = true
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if !anySelected {
				include = true
			} else if selects.PostId {
				include = true
			}
			// Force-include FK when its relation is selected
			if selects.Post != nil {
				include = true
			}
		} else if omits != nil {
			if omits.PostId {
				include = false
			}
		}
		if include {
			cols = append(cols, "postId")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if !anySelected {
				include = true
			} else if selects.CategoryId {
				include = true
			}
			// Force-include FK when its relation is selected
			if selects.Category != nil {
				include = true
			}
		} else if omits != nil {
			if omits.CategoryId {
				include = false
			}
		}
		if include {
			cols = append(cols, "categoryId")
		}
	}

	if len(cols) == 0 {
		cols = append(cols, "postId")
		cols = append(cols, "categoryId")
	}

	// Force-include any requested columns
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

	res, err := executeInsert(ctx, q, "CategoryToPost", cols, vals, returningCols, idCol, scanFunc)
	if err != nil {
		return nil, err
	}

	if err := q.loadCategoryToPostRelations(ctx, []*CategoryToPost{res}, selects); err != nil {
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
			func(p *CategoryToPost) (string, bool) {
				return fmt.Sprint(p.PostId), true
			},
			"Post",
			"id",
			returningCols,
			func(rows *sql.Rows, child *Post) error {
				return rows.Scan(child.ScanFields(returningCols)...)
			},
			func(child *Post) (string, bool) {
				return fmt.Sprint(child.Id), true
			},
			func(p *CategoryToPost, children []*Post) {
				if len(children) > 0 {
					p.Post = children[0]
				}
			},
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
			func(p *CategoryToPost) (string, bool) {
				return fmt.Sprint(p.CategoryId), true
			},
			"Category",
			"id",
			returningCols,
			func(rows *sql.Rows, child *Category) error {
				return rows.Scan(child.ScanFields(returningCols)...)
			},
			func(child *Category) (string, bool) {
				return fmt.Sprint(child.Id), true
			},
			func(p *CategoryToPost, children []*Category) {
				if len(children) > 0 {
					p.Category = children[0]
				}
			},
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
