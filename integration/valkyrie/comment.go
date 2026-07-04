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

// Comment represents the database model
type Comment struct {
	Id       string `db:"id" json:"id"`
	Textify  int32  `db:"textify" json:"textify"`
	Dummy3   string `db:"dummy3" json:"dummy3"`
	Dummy1   int32  `db:"dummy1" json:"dummy1"`
	Dummy2   string `db:"dummy2" json:"dummy2"`
	PostId   string `db:"postId" json:"postId"`
	AuthorId string `db:"authorId" json:"authorId"`
	Post     *Post  `json:"post,omitempty"`
	Author   *User  `json:"author,omitempty"`
}

// CommentCreateInput represents the input structure for creation
type CommentCreateInput struct {
	Id       *string `json:"id"`
	Textify  int32   `json:"textify"`
	Dummy3   string  `json:"dummy3"`
	Dummy1   int32   `json:"dummy1"`
	Dummy2   string  `json:"dummy2"`
	PostId   string  `json:"postId"`
	AuthorId string  `json:"authorId"`
}

// CommentSelect specifies which fields to include
type CommentSelect struct {
	Id       bool        `json:"id"`
	Textify  bool        `json:"textify"`
	Dummy3   bool        `json:"dummy3"`
	Dummy1   bool        `json:"dummy1"`
	Dummy2   bool        `json:"dummy2"`
	PostId   bool        `json:"postId"`
	AuthorId bool        `json:"authorId"`
	Post     *PostSelect `json:"post,omitempty"`
	Author   *UserSelect `json:"author,omitempty"`
}

// CommentOmit specifies which fields to exclude
type CommentOmit struct {
	Id       bool      `json:"id"`
	Textify  bool      `json:"textify"`
	Dummy3   bool      `json:"dummy3"`
	Dummy1   bool      `json:"dummy1"`
	Dummy2   bool      `json:"dummy2"`
	PostId   bool      `json:"postId"`
	AuthorId bool      `json:"authorId"`
	Post     *PostOmit `json:"post,omitempty"`
	Author   *UserOmit `json:"author,omitempty"`
}

type CommentDelegate struct {
	client *Queries
}

func (m *Comment) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "id":
			targets[i] = &m.Id
		case "textify":
			targets[i] = &m.Textify
		case "dummy3":
			targets[i] = &m.Dummy3
		case "dummy1":
			targets[i] = &m.Dummy1
		case "dummy2":
			targets[i] = &m.Dummy2
		case "postId":
			targets[i] = &m.PostId
		case "authorId":
			targets[i] = &m.AuthorId
		}
	}
	return targets
}

var commentDefaultCols = []string{
	"id",
	"textify",
	"dummy3",
	"dummy1",
	"dummy2",
	"postId",
	"authorId",
}

func (q *Queries) selectCommentCols(selects *CommentSelect, omits *CommentOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return commentDefaultCols
	}

	var cols []string
	var anySelected bool
	if selects != nil {
		if selects.Id {
			anySelected = true
		}
		if selects.Textify {
			anySelected = true
		}
		if selects.Dummy3 {
			anySelected = true
		}
		if selects.Dummy1 {
			anySelected = true
		}
		if selects.Dummy2 {
			anySelected = true
		}
		if selects.PostId {
			anySelected = true
		}
		if selects.AuthorId {
			anySelected = true
		}
		if selects.Post != nil {
			anySelected = true
		}
		if selects.Author != nil {
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
			} else if selects.Textify {
				include = true
			}
		} else if omits != nil {
			if omits.Textify {
				include = false
			}
		}
		if include {
			cols = append(cols, "textify")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if !anySelected {
				include = true
			} else if selects.Dummy3 {
				include = true
			}
		} else if omits != nil {
			if omits.Dummy3 {
				include = false
			}
		}
		if include {
			cols = append(cols, "dummy3")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if !anySelected {
				include = true
			} else if selects.Dummy1 {
				include = true
			}
		} else if omits != nil {
			if omits.Dummy1 {
				include = false
			}
		}
		if include {
			cols = append(cols, "dummy1")
		}
	}
	{
		include := true
		if selects != nil {
			include = false
			if !anySelected {
				include = true
			} else if selects.Dummy2 {
				include = true
			}
		} else if omits != nil {
			if omits.Dummy2 {
				include = false
			}
		}
		if include {
			cols = append(cols, "dummy2")
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
		cols = append(cols, "textify")
		cols = append(cols, "dummy3")
		cols = append(cols, "dummy1")
		cols = append(cols, "dummy2")
		cols = append(cols, "postId")
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
func (d *CommentDelegate) Create(input CommentCreateInput) *CreateBuilder[Comment, CommentCreateInput, CommentSelect, CommentOmit] {
	return &CreateBuilder[Comment, CommentCreateInput, CommentSelect, CommentOmit]{
		client:   d.client,
		input:    input,
		execFunc: d.client.executeCommentCreate,
	}
}

func (q *Queries) executeCommentCreate(ctx context.Context, input CommentCreateInput, selects *CommentSelect, omits *CommentOmit) (*Comment, error) {
	var cols []string
	var vals []any
	if input.Id != nil {
		cols = append(cols, q.dialect.Quote("id"))
		vals = append(vals, *input.Id)
	} else {
		cols = append(cols, q.dialect.Quote("id"))
		vals = append(vals, generateCUID())
	}
	cols = append(cols, q.dialect.Quote("textify"))
	vals = append(vals, input.Textify)
	cols = append(cols, q.dialect.Quote("dummy3"))
	vals = append(vals, input.Dummy3)
	cols = append(cols, q.dialect.Quote("dummy1"))
	vals = append(vals, input.Dummy1)
	cols = append(cols, q.dialect.Quote("dummy2"))
	vals = append(vals, input.Dummy2)
	cols = append(cols, q.dialect.Quote("postId"))
	vals = append(vals, input.PostId)
	cols = append(cols, q.dialect.Quote("authorId"))
	vals = append(vals, input.AuthorId)

	returningCols := q.selectCommentCols(selects, omits)

	scanFunc := func(res *Comment, cols []string) []any {
		return res.ScanFields(cols)
	}

	idCol := "id"

	res, err := executeInsert(ctx, q, "Comment", cols, vals, returningCols, idCol, scanFunc)
	if err != nil {
		return nil, err
	}

	if err := q.loadCommentRelations(ctx, []*Comment{res}, selects); err != nil {
		return nil, err
	}

	return res, nil
}
func (q *Queries) loadCommentRelations(ctx context.Context, records []*Comment, selects *CommentSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Post != nil {
		returningCols := q.selectPostCols(selects.Post, nil, "id")
		// Current model holds the FK: Comment.postId
		allChildren, err := loadRelation(
			ctx, q, records,
			func(p *Comment) (string, bool) {
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
			func(p *Comment, children []*Post) {
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
	if selects.Author != nil {
		returningCols := q.selectUserCols(selects.Author, nil, "id")
		// Current model holds the FK: Comment.authorId
		allChildren, err := loadRelation(
			ctx, q, records,
			func(p *Comment) (string, bool) {
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
			func(p *Comment, children []*User) {
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

	return nil
}
