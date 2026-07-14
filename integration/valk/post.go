package valk

import (
	"context"
	"fmt"
	"slices"
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
	Id         bool                      `json:"id"`
	Title      bool                      `json:"title"`
	Content    bool                      `json:"content"`
	Published  bool                      `json:"published"`
	AuthorId   bool                      `json:"authorId"`
	Author     UserSelectQuery           `json:"author,omitempty"`
	Comments   CommentSelectQuery        `json:"comments,omitempty"`
	Categories CategoryToPostSelectQuery `json:"categories,omitempty"`
}

// PostOmit specifies which fields to exclude
type PostOmit struct {
	Id        bool `json:"id"`
	Title     bool `json:"title"`
	Content   bool `json:"content"`
	Published bool `json:"published"`
	AuthorId  bool `json:"authorId"`
}

type PostSelectQuery interface {
	GetRelationParams() (*PostSelect, *PostOmit, QueryParams[Post])
}

func (s *PostSelect) GetRelationParams() (*PostSelect, *PostOmit, QueryParams[Post]) {
	return s, nil, QueryParams[Post]{}
}

// PostQueryBuilder builds a query for the relation Post
type PostQueryBuilder struct {
	selects *PostSelect
	omits   *PostOmit
	where   []PredicateOf[Post]
	take    *int
	skip    *int
	orderBy []OrderBy
}

func (b *PostQueryBuilder) Where(preds ...PredicateOf[Post]) *PostQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *PostQueryBuilder) Take(limit int) *PostQueryBuilder {
	b.take = &limit
	return b
}

func (b *PostQueryBuilder) Skip(offset int) *PostQueryBuilder {
	b.skip = &offset
	return b
}

func (b *PostQueryBuilder) OrderBy(orders ...OrderBy) *PostQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *PostQueryBuilder) Select(s PostSelect) *PostQueryBuilder {
	b.selects = &s
	return b
}

func (b *PostQueryBuilder) Omit(o PostOmit) *PostQueryBuilder {
	b.omits = &o
	return b
}

func (b *PostQueryBuilder) GetRelationParams() (*PostSelect, *PostOmit, QueryParams[Post]) {
	if b == nil {
		return nil, nil, QueryParams[Post]{}
	}
	return b.selects, b.omits, QueryParams[Post]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
	}
}

type PostDelegate struct {
	client          *Queries
	beforeCreate    func(context.Context, *PostCreate) error
	afterCreate     func(context.Context, []*Post) error
	afterCreateMany func(context.Context, []PostCreate, int64) error
}

func (d *PostDelegate) BeforeCreate(hook func(context.Context, *PostCreate) error) {
	d.beforeCreate = hook
}

func (d *PostDelegate) AfterCreate(hook func(context.Context, []*Post) error) {
	d.afterCreate = hook
}

func (d *PostDelegate) AfterCreateMany(hook func(context.Context, []PostCreate, int64) error) {
	d.afterCreateMany = hook
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

type PostCreateBuilder struct {
	*CreateBuilder[Post, PostSelect, PostOmit]
}

func (b *PostCreateBuilder) SetId(v string) *PostCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}
func (b *PostCreateBuilder) SetTitle(v string) *PostCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "title", Val: v})
	return b
}
func (b *PostCreateBuilder) SetContent(v string) *PostCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "content", Val: v})
	return b
}
func (b *PostCreateBuilder) SetPublished(v bool) *PostCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "published", Val: v})
	return b
}
func (b *PostCreateBuilder) SetAuthorId(v string) *PostCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "authorId", Val: v})
	return b
}

func (d *PostDelegate) Create(assignments ...FieldAssignment) *PostCreateBuilder {
	return &PostCreateBuilder{
		CreateBuilder: &CreateBuilder[Post, PostSelect, PostOmit]{
			client:      d.client,
			assignments: assignments,
			execFunc:    d.client.executePostCreate,
		},
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
				ValidateString(errs, "id", v, false, 0, false, false)
			} else {
				errs.Add("id", a.Val, "type", "field id must be of type string")
			}
		case "title":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "title", v, true, 0, false, false)
			} else {
				errs.Add("title", a.Val, "type", "field title must be of type string")
			}
		case "content":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "content", v, false, 0, false, false)
			} else {
				errs.Add("content", a.Val, "type", "field content must be of type string")
			}
		case "published":
			if _, ok := a.Val.(bool); !ok {
				errs.Add("published", a.Val, "type", "field published must be of type bool")
			}
		case "authorId":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "authorId", v, true, 0, false, false)
			} else {
				errs.Add("authorId", a.Val, "type", "field authorId must be of type string")
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

	rowMap := input.ToRowMap()
	var cols []string
	var vals []any
	for _, col := range PostColOrder {
		if val, ok := rowMap[col]; ok {
			cols = append(cols, col)
			vals = append(vals, val)
		}
	}

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
		if err := q.Post.afterCreate(ctx, []*Post{res}); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (d *PostDelegate) CreateMany(builders ...*PostCreateBuilder) *CreateManyBuilder[Post] {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &CreateManyBuilder[Post]{
		client:   d.client,
		records:  records,
		execFunc: d.client.executePostCreateMany,
	}
}

func (d *PostDelegate) CreateManyAndReturn(builders ...*PostCreateBuilder) *CreateManyAndReturnBuilder[Post, PostSelect, PostOmit] {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &CreateManyAndReturnBuilder[Post, PostSelect, PostOmit]{
		client:   d.client,
		records:  records,
		execFunc: d.client.executePostCreateManyAndReturn,
	}
}

func (q *Queries) executePostCreateMany(ctx context.Context, records []RecordInput) (int64, error) {
	rowMaps := make([]map[string]any, len(records))
	inputs := make([]PostCreate, len(records))
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
		inputs[i] = input
	}
	count, err := executeCreateMany(ctx, q, rowMaps, "Post", PostColOrder)
	if err != nil {
		return 0, err
	}
	if q.Post.afterCreateMany != nil {
		if err := q.Post.afterCreateMany(ctx, inputs, count); err != nil {
			return 0, err
		}
	}
	return count, nil
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
		if err := q.Post.afterCreate(ctx, results); err != nil {
			return nil, err
		}
	}
	return results, nil
}
func (d *PostDelegate) FindUnique(where UniquePredicate[Post], additional ...PredicateOf[Post]) *FindUniqueBuilder[Post, PostSelect, PostOmit] {
	return &FindUniqueBuilder[Post, PostSelect, PostOmit]{
		client:     d.client,
		where:      where,
		additional: additional,
		execFunc:   d.client.executePostFindUnique,
	}
}

func (d *PostDelegate) FindFirst(preds ...PredicateOf[Post]) *FindFirstBuilder[Post, PostSelect, PostOmit] {
	return &FindFirstBuilder[Post, PostSelect, PostOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executePostFindFirst,
	}
}

func (d *PostDelegate) FindMany(preds ...PredicateOf[Post]) *FindManyBuilder[Post, PostSelect, PostOmit] {
	return &FindManyBuilder[Post, PostSelect, PostOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executePostFindMany,
	}
}

func (q *Queries) executePostFindUnique(ctx context.Context, where UniquePredicate[Post], additional []PredicateOf[Post], selects *PostSelect, omits *PostOmit) (*Post, error) {
	if err := where.Validate(); err != nil {
		return nil, err
	}
	for _, p := range additional {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}
	allPreds := append([]PredicateOf[Post]{where}, additional...)
	whereClause, vals := CompilePredicates(q.dialect, allPreds)
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
		nil,
	)
}

func (q *Queries) executePostFindFirst(
	ctx context.Context,
	params QueryParams[Post],
	selects *PostSelect,
	omits *PostOmit,
) (*Post, error) {
	for _, p := range params.Where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}
	whereClause, vals := CompilePredicates(q.dialect, params.Where)
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
		params.Skip,
	)
}

func (q *Queries) executePostFindMany(
	ctx context.Context,
	params QueryParams[Post],
	selects *PostSelect,
	omits *PostOmit,
) ([]*Post, error) {
	for _, p := range params.Where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}
	whereClause, vals := CompilePredicates(q.dialect, params.Where)
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
		params.Take,
		params.Skip,
	)
}
func (q *Queries) loadPostRelations(ctx context.Context, records []*Post, selects *PostSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Author != nil {
		relationSelects, relationOmits, relationParams := selects.Author.GetRelationParams()
		returningCols := q.selectUserCols(relationSelects, relationOmits, "id")
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
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading author: %w", err)
		}
		if err := q.loadUserRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Comments != nil {
		relationSelects, relationOmits, relationParams := selects.Comments.GetRelationParams()
		returningCols := q.selectCommentCols(relationSelects, relationOmits, "postId")
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
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading comments: %w", err)
		}
		if err := q.loadCommentRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Categories != nil {
		relationSelects, relationOmits, relationParams := selects.Categories.GetRelationParams()
		returningCols := q.selectCategoryToPostCols(relationSelects, relationOmits, "postId")
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
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading categories: %w", err)
		}
		if err := q.loadCategoryToPostRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
