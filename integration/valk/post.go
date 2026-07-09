package valk

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"
)

// Post represents the database model
type Post struct {
	Id         string            `db:"id" json:"id"`
	Title      string            `db:"title" json:"title"`
	Content    *string           `db:"content" json:"content,omitempty"`
	Published  bool              `db:"published" json:"published"`
	AuthorId   string            `db:"authorId" json:"authorId"`
	Author     *User             `json:"author,omitempty"`
	Comments   []*Comment        `json:"comments,omitempty"`
	Categories []*CategoryToPost `json:"categories,omitempty"`
}

// PostCreate is used for hooks only — the Create API uses FieldAssignment
type PostCreate struct {
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
	client       *Queries
	beforeCreate func(context.Context, *PostCreate) error
	afterCreate  func(context.Context, *Post) error
}

func (d *PostDelegate) BeforeCreate(hook func(context.Context, *PostCreate) error) {
	d.beforeCreate = hook
}

func (d *PostDelegate) AfterCreate(hook func(context.Context, *Post) error) {
	d.afterCreate = hook
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
		{"id", selects != nil && selects.Id, omits != nil && omits.Id, selects != nil && selects.hasAnyRelation()},
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

func (d *PostDelegate) Create(assignments ...FieldAssignment) *CreateBuilder[Post, PostSelect, PostOmit] {
	return &CreateBuilder[Post, PostSelect, PostOmit]{
		client:      d.client,
		assignments: assignments,
		execFunc:    d.client.executePostCreate,
	}
}

func validatePostCreate(assignments []FieldAssignment) error {
	errs := &ValidationError{}

	provided := make(map[string]bool)
	for _, a := range assignments {
		provided[a.Col] = true
		switch a.Col {
		case "id":
			if v, ok := a.Val.(string); ok {
				if strings.Contains(v, "\x00") {
					errs.Add("id", v, "safety", "string cannot contain null bytes")
				}
				if !utf8.ValidString(v) {
					errs.Add("id", v, "safety", "string must be valid UTF-8")
				}
			}
		case "title":
			if v, ok := a.Val.(string); ok {
				if v == "" {
					errs.Add("title", v, "required", "field title is required")
				}
				if strings.Contains(v, "\x00") {
					errs.Add("title", v, "safety", "string cannot contain null bytes")
				}
				if !utf8.ValidString(v) {
					errs.Add("title", v, "safety", "string must be valid UTF-8")
				}
			}
		case "content":
		case "published":
		case "authorId":
			if v, ok := a.Val.(string); ok {
				if v == "" {
					errs.Add("authorId", v, "required", "field authorId is required")
				}
				if strings.Contains(v, "\x00") {
					errs.Add("authorId", v, "safety", "string cannot contain null bytes")
				}
				if !utf8.ValidString(v) {
					errs.Add("authorId", v, "safety", "string must be valid UTF-8")
				}
			}
		}
	}
	if !provided["title"] {
		errs.Add("title", "", "required", "field Title is required")
	}
	if !provided["authorId"] {
		errs.Add("authorId", "", "required", "field AuthorId is required")
	}

	if errs.HasErrors() {
		return *errs
	}
	return nil
}

func assignmentsToPostCreate(assignments []FieldAssignment) PostCreate {
	var input PostCreate
	for _, a := range assignments {
		switch a.Col {
		case "id":
			if v, ok := a.Val.(string); ok {
				input.Id = &v
			}
		case "title":
			if v, ok := a.Val.(string); ok {
				input.Title = v
			}
		case "content":
			if v, ok := a.Val.(string); ok {
				input.Content = &v
			}
		case "published":
			if v, ok := a.Val.(bool); ok {
				input.Published = &v
			}
		case "authorId":
			if v, ok := a.Val.(string); ok {
				input.AuthorId = v
			}
		}
	}
	return input
}

func (s *PostCreate) ToRowMap() map[string]any {
	m := make(map[string]any, 5)
	if s.Id != nil {
		m["id"] = *s.Id
	} else {
		m["id"] = generateCUID()
	}
	m["title"] = s.Title
	if s.Content != nil {
		m["content"] = *s.Content
	}
	if s.Published != nil {
		m["published"] = *s.Published
	}
	m["authorId"] = s.AuthorId
	return m
}

func (q *Queries) executePostCreate(ctx context.Context, assignments []FieldAssignment, selects *PostSelect, omits *PostOmit) (*Post, error) {
	input := assignmentsToPostCreate(assignments)

	if q.Post.beforeCreate != nil {
		if err := q.Post.beforeCreate(ctx, &input); err != nil {
			return nil, err
		}
	}

	if err := validatePostCreate(assignments); err != nil {
		return nil, err
	}

	var cols []string
	var vals []any
	if input.Id != nil {
		cols = append(cols, "id")
		vals = append(vals, *input.Id)
	} else {
		cols = append(cols, "id")
		vals = append(vals, generateCUID())
	}
	cols = append(cols, "title")
	vals = append(vals, input.Title)
	if input.Content != nil {
		cols = append(cols, "content")
		vals = append(vals, *input.Content)
	}
	if input.Published != nil {
		cols = append(cols, "published")
		vals = append(vals, *input.Published)
	}
	cols = append(cols, "authorId")
	vals = append(vals, input.AuthorId)

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

	if q.Post.afterCreate != nil {
		if err := q.Post.afterCreate(ctx, res); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (d *PostDelegate) CreateMany(records ...RecordInput) *CreateManyBuilder[Post] {
	return &CreateManyBuilder[Post]{
		client:   d.client,
		records:  records,
		execFunc: d.client.executePostCreateMany,
	}
}

func (d *PostDelegate) CreateManyAndReturn(records ...RecordInput) *CreateManyAndReturnBuilder[Post, PostSelect, PostOmit] {
	return &CreateManyAndReturnBuilder[Post, PostSelect, PostOmit]{
		client:   d.client,
		records:  records,
		execFunc: d.client.executePostCreateManyAndReturn,
	}
}

func (q *Queries) executePostCreateMany(ctx context.Context, records []RecordInput) (int64, error) {
	rowMaps := make([]map[string]any, len(records))
	for i, rec := range records {
		if err := validatePostCreate(rec.Assignments); err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		input := assignmentsToPostCreate(rec.Assignments)
		if q.Post.beforeCreate != nil {
			if err := q.Post.beforeCreate(ctx, &input); err != nil {
				return 0, err
			}
		}
		rowMaps[i] = input.ToRowMap()
	}
	return executeCreateMany(ctx, q, rowMaps, "Post", PostColOrder)
}

func (q *Queries) executePostCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *PostSelect, omits *PostOmit) ([]*Post, error) {
	rowMaps := make([]map[string]any, len(records))
	idCol := "id"
	for i, rec := range records {
		if err := validatePostCreate(rec.Assignments); err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		input := assignmentsToPostCreate(rec.Assignments)
		if q.Post.beforeCreate != nil {
			if err := q.Post.beforeCreate(ctx, &input); err != nil {
				return nil, err
			}
		}
		rowMaps[i] = input.ToRowMap()
	}
	results, err := executeCreateManyAndReturn(ctx, q, rowMaps, "Post", PostColOrder, selects, omits,
		q.selectPostCols,
		q.loadPostRelations,
		(*Post).ScanFields,
		(*PostSelect).hasAnyRelation,
		idCol,
	)
	if err != nil {
		return nil, err
	}
	if q.Post.afterCreate != nil {
		for _, r := range results {
			if err := q.Post.afterCreate(ctx, r); err != nil {
				return nil, err
			}
		}
	}
	return results, nil
}
func (d *PostDelegate) FindUnique(where UniquePredicate) *FindUniqueBuilder[Post, PostSelect, PostOmit] {
	return &FindUniqueBuilder[Post, PostSelect, PostOmit]{
		client:   d.client,
		where:    where,
		execFunc: d.client.executePostFindUnique,
	}
}

func (d *PostDelegate) FindFirst(preds ...Predicate) *FindFirstBuilder[Post, PostSelect, PostOmit] {
	return &FindFirstBuilder[Post, PostSelect, PostOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executePostFindFirst,
	}
}

func (d *PostDelegate) FindMany(preds ...Predicate) *FindManyBuilder[Post, PostSelect, PostOmit] {
	return &FindManyBuilder[Post, PostSelect, PostOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executePostFindMany,
	}
}

func (q *Queries) executePostFindUnique(ctx context.Context, where UniquePredicate, selects *PostSelect, omits *PostOmit) (*Post, error) {
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
	returningCols := q.selectPostCols(selects, omits)
	return executeSingleWithRelations(ctx, q, "Post", whereClause, vals, returningCols,
		func(res *Post, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*Post) error {
			return txQ.loadPostRelations(ctx, results, selects)
		},
	)
}

func (q *Queries) executePostFindFirst(ctx context.Context, where []Predicate, selects *PostSelect, omits *PostOmit) (*Post, error) {
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
	returningCols := q.selectPostCols(selects, omits)
	return executeSingleWithRelations(ctx, q, "Post", whereClause, vals, returningCols,
		func(res *Post, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*Post) error {
			return txQ.loadPostRelations(ctx, results, selects)
		},
	)
}

func (q *Queries) executePostFindMany(ctx context.Context, where []Predicate, selects *PostSelect, omits *PostOmit) ([]*Post, error) {
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
	returningCols := q.selectPostCols(selects, omits)
	return executeManyWithRelations(ctx, q, "Post", whereClause, vals, returningCols,
		func(res *Post, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*Post) error {
			return txQ.loadPostRelations(ctx, results, selects)
		},
	)
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
