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

func (input CommentCreateInput) Validate() error {
	if input.Dummy3 == "" {
		return fmt.Errorf("field Dummy3 is required")
	}
	if input.Dummy2 == "" {
		return fmt.Errorf("field Dummy2 is required")
	}
	if input.PostId == "" {
		return fmt.Errorf("field PostId is required")
	}
	if input.AuthorId == "" {
		return fmt.Errorf("field AuthorId is required")
	}
	return nil
}

var CommentColOrder = []string{
	"id",
	"textify",
	"dummy3",
	"dummy1",
	"dummy2",
	"postId",
	"authorId",
}

func (s *CommentSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Post != nil || s.Author != nil
}

func (d *CommentDelegate) Create(input CommentCreateInput) *CreateBuilder[Comment, CommentCreateInput, CommentSelect, CommentOmit] {
	return &CreateBuilder[Comment, CommentCreateInput, CommentSelect, CommentOmit]{
		client:   d.client,
		input:    input,
		execFunc: d.client.executeCommentCreate,
	}
}

func (q *Queries) executeCommentCreate(ctx context.Context, input CommentCreateInput, selects *CommentSelect, omits *CommentOmit) (*Comment, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	m := q.CommentInputToMap(input)
	cols, vals := mapToColsVals(m, CommentColOrder)

	returningCols := q.selectCommentCols(selects, omits)

	scanFunc := func(res *Comment, cols []string) []any {
		return res.ScanFields(cols)
	}

	idCol := "id"

	hasRelations := selects.hasAnyRelation()

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

func (q *Queries) CommentInputToMap(input CommentCreateInput) map[string]any {
	m := make(map[string]any)
	if input.Id != nil {
		m["id"] = *input.Id
	} else {
		m["id"] = generateCUID()
	}
	m["textify"] = input.Textify
	m["dummy3"] = input.Dummy3
	m["dummy1"] = input.Dummy1
	m["dummy2"] = input.Dummy2
	m["postId"] = input.PostId
	m["authorId"] = input.AuthorId
	return m
}

func (d *CommentDelegate) CreateMany(inputs []CommentCreateInput) *CreateManyBuilder[Comment, CommentCreateInput] {
	return &CreateManyBuilder[Comment, CommentCreateInput]{
		client:   d.client,
		inputs:   inputs,
		execFunc: d.client.executeCommentCreateMany,
	}
}

func (d *CommentDelegate) CreateManyAndReturn(inputs []CommentCreateInput) *CreateManyAndReturnBuilder[Comment, CommentCreateInput, CommentSelect, CommentOmit] {
	return &CreateManyAndReturnBuilder[Comment, CommentCreateInput, CommentSelect, CommentOmit]{
		client:   d.client,
		inputs:   inputs,
		execFunc: d.client.executeCommentCreateManyAndReturn,
	}
}

func (q *Queries) executeCommentCreateMany(ctx context.Context, inputs []CommentCreateInput) (int64, error) {
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
			rowMaps[i] = q.CommentInputToMap(input)
		}
		query, vals := buildBulkInsertSQL(q.dialect, "Comment", rowMaps, CommentColOrder, nil)
		res, err := q.exec(ctx, query, vals...)
		if err != nil {
			return 0, err
		}
		return res.RowsAffected()
	}

	var count int64
	err := q.transaction(ctx, func(txQ *Queries) error {
		for _, input := range inputs {
			_, err := txQ.executeCommentCreate(ctx, input, nil, nil)
			if err != nil {
				return err
			}
			count++
		}
		return nil
	})
	return count, err
}

func (q *Queries) executeCommentCreateManyAndReturn(ctx context.Context, inputs []CommentCreateInput, selects *CommentSelect, omits *CommentOmit) ([]*Comment, error) {
	if len(inputs) == 0 {
		return nil, nil
	}
	for i, input := range inputs {
		if err := input.Validate(); err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
	}

	hasRelations := selects.hasAnyRelation()
	returningCols := q.selectCommentCols(selects, omits)

	if q.dialect.SupportsBulkInsert() {
		rowMaps := make([]map[string]any, len(inputs))
		for i, input := range inputs {
			rowMaps[i] = q.CommentInputToMap(input)
		}
		query, vals := buildBulkInsertSQL(q.dialect, "Comment", rowMaps, CommentColOrder, returningCols)
		var records []*Comment
		err := q.transaction(ctx, func(txQ *Queries) error {
			rows, err := txQ.query(ctx, query, vals...)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var record Comment
				if err := rows.Scan(record.ScanFields(returningCols)...); err != nil {
					return err
				}
				records = append(records, &record)
			}
			if err := rows.Err(); err != nil {
				return err
			}
			if hasRelations {
				return txQ.loadCommentRelations(ctx, records, selects)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return records, nil
	}

	// Fallback to loop inside transaction
	var records []*Comment
	err := q.transaction(ctx, func(txQ *Queries) error {
		for _, input := range inputs {
			res, err := txQ.executeCommentCreate(ctx, input, nil, nil)
			if err != nil {
				return err
			}
			records = append(records, res)
		}

		if hasRelations {
			return txQ.loadCommentRelations(ctx, records, selects)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return records, nil
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
