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

// Post represents the database model
type Post struct {
	Id         string            `db:"id" json:"id"`
	Title      string            `db:"title" json:"title"`
	Content    *string           `db:"content" json:"content"`
	Published  bool              `db:"published" json:"published"`
	AuthorId   string            `db:"authorId" json:"authorId"`
	Author     *User             `json:"author,omitempty"`
	Comments   []*Comment        `json:"comments,omitempty"`
	Categories []*CategoryToPost `json:"categories,omitempty"`
}

// PostCreateInput represents the input structure for creation
type PostCreateInput struct {
	Id        *string `json:"id"`
	Title     string  `json:"title"`
	Content   *string `json:"content"`
	Published *bool   `json:"published"`
	AuthorId  string  `json:"authorId"`
}

// PostSelect specifies which fields to include
type PostSelect struct {
	Id         bool                  `json:"id"`
	Title      bool                  `json:"title"`
	Content    bool                  `json:"content"`
	Published  bool                  `json:"published"`
	AuthorId   bool                  `json:"authorId"`
	Author     *UserSelect           `json:"author,omitempty"`
	Comments   *CommentSelect        `json:"comments,omitempty"`
	Categories *CategoryToPostSelect `json:"categories,omitempty"`
}

// PostOmit specifies which fields to exclude
type PostOmit struct {
	Id         bool                `json:"id"`
	Title      bool                `json:"title"`
	Content    bool                `json:"content"`
	Published  bool                `json:"published"`
	AuthorId   bool                `json:"authorId"`
	Author     *UserOmit           `json:"author,omitempty"`
	Comments   *CommentOmit        `json:"comments,omitempty"`
	Categories *CategoryToPostOmit `json:"categories,omitempty"`
}

type PostDelegate struct {
	client *Queries
}

func (m *Post) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "id":
			targets[i] = &m.Id
		case "title":
			targets[i] = &m.Title
		case "content":
			targets[i] = &m.Content
		case "published":
			targets[i] = &m.Published
		case "authorId":
			targets[i] = &m.AuthorId
		}
	}
	return targets
}

var postDefaultCols = []string{
	"id",
	"title",
	"content",
	"published",
	"authorId",
}

func (q *Queries) selectPostCols(selects *PostSelect, omits *PostOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return postDefaultCols
	}

	var cols []string
	var anySelected bool
	if selects != nil {
		if selects.Id {
			anySelected = true
		}
		if selects.Title {
			anySelected = true
		}
		if selects.Content {
			anySelected = true
		}
		if selects.Published {
			anySelected = true
		}
		if selects.AuthorId {
			anySelected = true
		}
		if selects.Author != nil {
			anySelected = true
		}
		if selects.Comments != nil {
			anySelected = true
		}
		if selects.Categories != nil {
			anySelected = true
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if !anySelected {
				include = true
			} else if selects.Id {
				include = true
			}
		} else if omits != nil {
			if omits.Id {
				include = false
			}
		}
		if include {
			cols = append(cols, "id")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if !anySelected {
				include = true
			} else if selects.Title {
				include = true
			}
		} else if omits != nil {
			if omits.Title {
				include = false
			}
		}
		if include {
			cols = append(cols, "title")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if !anySelected {
				include = true
			} else if selects.Content {
				include = true
			}
		} else if omits != nil {
			if omits.Content {
				include = false
			}
		}
		if include {
			cols = append(cols, "content")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if !anySelected {
				include = true
			} else if selects.Published {
				include = true
			}
		} else if omits != nil {
			if omits.Published {
				include = false
			}
		}
		if include {
			cols = append(cols, "published")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if !anySelected {
				include = true
			} else if selects.AuthorId {
				include = true
			}
			// Force-include FK when its relation is selected
			if selects.Author != nil {
				include = true
			}
		} else if omits != nil {
			if omits.AuthorId {
				include = false
			}
		}
		if include {
			cols = append(cols, "authorId")
		}
	}

	if len(cols) == 0 {
		cols = append(cols, "id")
		cols = append(cols, "title")
		cols = append(cols, "content")
		cols = append(cols, "published")
		cols = append(cols, "authorId")
	}

	// Force-include any requested columns
	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}
func (d *PostDelegate) Create(input PostCreateInput) *CreateBuilder[Post, PostCreateInput, PostSelect, PostOmit] {
	return &CreateBuilder[Post, PostCreateInput, PostSelect, PostOmit]{
		client:   d.client,
		input:    input,
		execFunc: d.client.executePostCreate,
	}
}

func (q *Queries) executePostCreate(ctx context.Context, input PostCreateInput, selects *PostSelect, omits *PostOmit) (*Post, error) {
	var cols []string
	var vals []any
	if input.Id != nil {
		cols = append(cols, q.dialect.Quote("id"))
		vals = append(vals, *input.Id)
	} else {
		cols = append(cols, q.dialect.Quote("id"))
		vals = append(vals, generateCUID())
	}
	cols = append(cols, q.dialect.Quote("title"))
	vals = append(vals, input.Title)
	if input.Content != nil {
		cols = append(cols, q.dialect.Quote("content"))
		vals = append(vals, *input.Content)
	}
	if input.Published != nil {
		cols = append(cols, q.dialect.Quote("published"))
		vals = append(vals, *input.Published)
	}
	cols = append(cols, q.dialect.Quote("authorId"))
	vals = append(vals, input.AuthorId)

	returningCols := q.selectPostCols(selects, omits)

	scanFunc := func(res *Post, cols []string) []any {
		return res.ScanFields(cols)
	}

	idCol := "id"

	res, err := executeInsert(ctx, q, "Post", cols, vals, returningCols, idCol, scanFunc)
	if err != nil {
		return nil, err
	}

	if err := q.loadPostRelations(ctx, []*Post{res}, selects); err != nil {
		return nil, err
	}

	return res, nil
}
func (q *Queries) loadPostRelations(ctx context.Context, records []*Post, selects *PostSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Author != nil {
		returningCols := q.selectUserCols(selects.Author, nil, "id")
		// Current model holds the FK: Post.authorId
		allChildren, err := loadRelation(
			ctx, q, records,
			func(p *Post) (string, bool) {
				return fmt.Sprint(p.AuthorId), true
			},
			"User",
			"id",
			returningCols,
			func(rows *sql.Rows, child *User) error {
				return rows.Scan(child.ScanFields(returningCols)...)
			},
			func(child *User) (string, bool) {
				return fmt.Sprint(child.Id), true
			},
			func(p *Post, children []*User) {
				if len(children) > 0 {
					p.Author = children[0]
				}
			},
		)
		if err != nil {
			return fmt.Errorf("loading author: %w", err)
		}
		if err := q.loadUserRelations(ctx, allChildren, selects.Author); err != nil {
			return err
		}
	}
	if selects.Comments != nil {
		returningCols := q.selectCommentCols(selects.Comments, nil, "postId")
		// Inverse holds the FK: Comment.postId
		allChildren, err := loadRelation(
			ctx, q, records,
			func(p *Post) (string, bool) {
				return fmt.Sprint(p.Id), true
			},
			"Comment",
			"postId",
			returningCols,
			func(rows *sql.Rows, child *Comment) error {
				return rows.Scan(child.ScanFields(returningCols)...)
			},
			func(child *Comment) (string, bool) {
				return fmt.Sprint(child.PostId), true
			},
			func(p *Post, children []*Comment) {
				p.Comments = append(p.Comments, children...)
			},
		)
		if err != nil {
			return fmt.Errorf("loading comments: %w", err)
		}
		if err := q.loadCommentRelations(ctx, allChildren, selects.Comments); err != nil {
			return err
		}
	}
	if selects.Categories != nil {
		returningCols := q.selectCategoryToPostCols(selects.Categories, nil, "postId")
		// Inverse holds the FK: CategoryToPost.postId
		allChildren, err := loadRelation(
			ctx, q, records,
			func(p *Post) (string, bool) {
				return fmt.Sprint(p.Id), true
			},
			"CategoryToPost",
			"postId",
			returningCols,
			func(rows *sql.Rows, child *CategoryToPost) error {
				return rows.Scan(child.ScanFields(returningCols)...)
			},
			func(child *CategoryToPost) (string, bool) {
				return fmt.Sprint(child.PostId), true
			},
			func(p *Post, children []*CategoryToPost) {
				p.Categories = append(p.Categories, children...)
			},
		)
		if err != nil {
			return fmt.Errorf("loading categories: %w", err)
		}
		if err := q.loadCategoryToPostRelations(ctx, allChildren, selects.Categories); err != nil {
			return err
		}
	}

	return nil
}
