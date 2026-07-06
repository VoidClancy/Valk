package valkyrie

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

var _ = time.Time{}
var _ = fmt.Sprintf
var _ = strings.Join
var _ = context.Background
var _ = sql.LevelDefault
var _ = slices.Contains[[]string, string]
var _ = utf8.ValidString

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

	anySelected := selects != nil && (selects.Id || selects.Title || selects.Content || selects.Published || selects.AuthorId || selects.Author != nil || selects.Comments != nil || selects.Categories != nil)

	specs := []colSpec{
		{"id", selects != nil && selects.Id, omits != nil && omits.Id, false},
		{"title", selects != nil && selects.Title, omits != nil && omits.Title, false},
		{"content", selects != nil && selects.Content, omits != nil && omits.Content, false},
		{"published", selects != nil && selects.Published, omits != nil && omits.Published, false},
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

func (input PostCreateInput) Validate() error {
	errs := &ValidationError{}
	if input.Id != nil {
		val := *input.Id
		if strings.Contains(val, "\x00") {
			errs.Add("id", val, "safety", "string cannot contain null bytes")
		}
		if !utf8.ValidString(val) {
			errs.Add("id", val, "safety", "string must be valid UTF-8")
		}
	}
	if input.Title == "" {
		errs.Add("title", input.Title, "required", "field Title is required")
	}
	if strings.Contains(input.Title, "\x00") {
		errs.Add("title", input.Title, "safety", "string cannot contain null bytes")
	}
	if !utf8.ValidString(input.Title) {
		errs.Add("title", input.Title, "safety", "string must be valid UTF-8")
	}
	if input.AuthorId == "" {
		errs.Add("authorId", input.AuthorId, "required", "field AuthorId is required")
	}
	if strings.Contains(input.AuthorId, "\x00") {
		errs.Add("authorId", input.AuthorId, "safety", "string cannot contain null bytes")
	}
	if !utf8.ValidString(input.AuthorId) {
		errs.Add("authorId", input.AuthorId, "safety", "string must be valid UTF-8")
	}

	if errs.HasErrors() {
		return *errs
	}
	return nil
}

var PostColOrder = []string{
	"id",
	"title",
	"content",
	"published",
	"authorId",
}

func (s *PostSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Author != nil || s.Comments != nil || s.Categories != nil
}

func (d *PostDelegate) Create(input PostCreateInput) *CreateBuilder[Post, PostCreateInput, PostSelect, PostOmit] {
	return &CreateBuilder[Post, PostCreateInput, PostSelect, PostOmit]{
		client:   d.client,
		input:    input,
		execFunc: d.client.executePostCreate,
	}
}

func (q *Queries) executePostCreate(ctx context.Context, input PostCreateInput, selects *PostSelect, omits *PostOmit) (*Post, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	m := q.PostInputToMap(input)
	cols, vals := mapToColsVals(m, PostColOrder)

	returningCols := q.selectPostCols(selects, omits)

	scanFunc := func(res *Post, cols []string) []any {
		return res.ScanFields(cols)
	}

	idCol := "id"

	hasRelations := selects.hasAnyRelation()

	var res *Post
	var err error
	if hasRelations {
		err = q.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = executeInsert(ctx, txQ, "Post", cols, vals, returningCols, idCol, scanFunc)
			if err != nil {
				return err
			}
			return txQ.loadPostRelations(ctx, []*Post{res}, selects)
		})
	} else {
		res, err = executeInsert(ctx, q, "Post", cols, vals, returningCols, idCol, scanFunc)
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (q *Queries) PostInputToMap(input PostCreateInput) map[string]any {
	m := make(map[string]any)
	if input.Id != nil {
		m["id"] = *input.Id
	} else {
		m["id"] = generateCUID()
	}
	m["title"] = input.Title
	if input.Content != nil {
		m["content"] = *input.Content
	}
	if input.Published != nil {
		m["published"] = *input.Published
	}
	m["authorId"] = input.AuthorId
	return m
}

func (d *PostDelegate) CreateMany(inputs []PostCreateInput) *CreateManyBuilder[Post, PostCreateInput] {
	return &CreateManyBuilder[Post, PostCreateInput]{
		client:   d.client,
		inputs:   inputs,
		execFunc: d.client.executePostCreateMany,
	}
}

func (d *PostDelegate) CreateManyAndReturn(inputs []PostCreateInput) *CreateManyAndReturnBuilder[Post, PostCreateInput, PostSelect, PostOmit] {
	return &CreateManyAndReturnBuilder[Post, PostCreateInput, PostSelect, PostOmit]{
		client:   d.client,
		inputs:   inputs,
		execFunc: d.client.executePostCreateManyAndReturn,
	}
}

func (q *Queries) executePostCreateMany(ctx context.Context, inputs []PostCreateInput) (int64, error) {
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
			rowMaps[i] = q.PostInputToMap(input)
		}
		query, vals := buildBulkInsertSQL(q.dialect, "Post", rowMaps, PostColOrder, nil)
		res, err := q.exec(ctx, query, vals...)
		if err != nil {
			return 0, err
		}
		return res.RowsAffected()
	}

	var count int64
	err := q.transaction(ctx, func(txQ *Queries) error {
		for _, input := range inputs {
			_, err := txQ.executePostCreate(ctx, input, nil, nil)
			if err != nil {
				return err
			}
			count++
		}
		return nil
	})
	return count, err
}

func (q *Queries) executePostCreateManyAndReturn(ctx context.Context, inputs []PostCreateInput, selects *PostSelect, omits *PostOmit) ([]*Post, error) {
	if len(inputs) == 0 {
		return nil, nil
	}
	for i, input := range inputs {
		if err := input.Validate(); err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
	}

	hasRelations := selects.hasAnyRelation()
	returningCols := q.selectPostCols(selects, omits)

	if q.dialect.SupportsBulkInsert() {
		rowMaps := make([]map[string]any, len(inputs))
		for i, input := range inputs {
			rowMaps[i] = q.PostInputToMap(input)
		}
		query, vals := buildBulkInsertSQL(q.dialect, "Post", rowMaps, PostColOrder, returningCols)
		var records []*Post
		err := q.transaction(ctx, func(txQ *Queries) error {
			rows, err := txQ.query(ctx, query, vals...)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var record Post
				if err := rows.Scan(record.ScanFields(returningCols)...); err != nil {
					return err
				}
				records = append(records, &record)
			}
			if err := rows.Err(); err != nil {
				return err
			}
			if hasRelations {
				return txQ.loadPostRelations(ctx, records, selects)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return records, nil
	}

	// Fallback to loop inside transaction
	var records []*Post
	err := q.transaction(ctx, func(txQ *Queries) error {
		for _, input := range inputs {
			res, err := txQ.executePostCreate(ctx, input, nil, nil)
			if err != nil {
				return err
			}
			records = append(records, res)
		}

		if hasRelations {
			return txQ.loadPostRelations(ctx, records, selects)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return records, nil
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
			directKey(func(p *Post) string { return p.AuthorId }),
			"User",
			"id",
			returningCols,
			scanInto(returningCols, (*User).ScanFields),
			directKey(func(c *User) string { return c.Id }),
			setOne(func(p *Post, c *User) { p.Author = c }),
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
			directKey(func(p *Post) string { return p.Id }),
			"Comment",
			"postId",
			returningCols,
			scanInto(returningCols, (*Comment).ScanFields),
			directKey(func(c *Comment) string { return c.PostId }),
			appendMany(func(p *Post) *[]*Comment { return &p.Comments }),
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
			directKey(func(p *Post) string { return p.Id }),
			"CategoryToPost",
			"postId",
			returningCols,
			scanInto(returningCols, (*CategoryToPost).ScanFields),
			directKey(func(c *CategoryToPost) string { return c.PostId }),
			appendMany(func(p *Post) *[]*CategoryToPost { return &p.Categories }),
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
