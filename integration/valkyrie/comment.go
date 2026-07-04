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

	anySelected := selects != nil && (selects.Id || selects.Textify || selects.Dummy3 || selects.Dummy1 || selects.Dummy2 || selects.PostId || selects.AuthorId || selects.Post != nil || selects.Author != nil)

	specs := []colSpec{
		{"id", selects != nil && selects.Id, omits != nil && omits.Id, false},
		{"textify", selects != nil && selects.Textify, omits != nil && omits.Textify, false},
		{"dummy3", selects != nil && selects.Dummy3, omits != nil && omits.Dummy3, false},
		{"dummy1", selects != nil && selects.Dummy1, omits != nil && omits.Dummy1, false},
		{"dummy2", selects != nil && selects.Dummy2, omits != nil && omits.Dummy2, false},
		{"postId", selects != nil && selects.PostId, omits != nil && omits.PostId, selects != nil && selects.Post != nil},
		{"authorId", selects != nil && selects.AuthorId, omits != nil && omits.AuthorId, selects != nil && selects.Author != nil},
	}

	cols := computeCols(specs, selects != nil, anySelected)

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
	hasRelations := selects != nil && (selects.Post != nil || selects.Author != nil)

	var res *Comment
	var err error
	if hasRelations {
		err = q.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = executeInsert(ctx, txQ, "Comment", cols, vals, returningCols, idCol, scanFunc)
			if err != nil {
				return err
			}
			return txQ.loadCommentRelations(ctx, []*Comment{res}, selects)
		})
	} else {
		res, err = executeInsert(ctx, q, "Comment", cols, vals, returningCols, idCol, scanFunc)
	}
	if err != nil {
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
			directKey(func(p *Comment) string { return p.PostId }),
			"Post",
			"id",
			returningCols,
			scanInto(returningCols, (*Post).ScanFields),
			directKey(func(c *Post) string { return c.Id }),
			setOne(func(p *Comment, c *Post) { p.Post = c }),
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
			directKey(func(p *Comment) string { return p.AuthorId }),
			"User",
			"id",
			returningCols,
			scanInto(returningCols, (*User).ScanFields),
			directKey(func(c *User) string { return c.Id }),
			setOne(func(p *Comment, c *User) { p.Author = c }),
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
