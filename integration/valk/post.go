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
	orderBy []OrderBy[Post]
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

func (b *PostQueryBuilder) OrderBy(orders ...OrderBy[Post]) *PostQueryBuilder {
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

type PostCreateQuery = func(ctx context.Context, args *PostCreate) (*Post, error)
type PostCreateManyQuery = func(ctx context.Context, args []*PostCreate) (int64, error)
type PostCreateManyAndReturnQuery = func(ctx context.Context, args []*PostCreate) ([]*Post, error)
type PostFindUniqueQuery = func(ctx context.Context, where UniquePredicate[Post], additional []PredicateOf[Post], selects *PostSelect, omits *PostOmit) (*Post, error)
type PostFindFirstQuery = func(ctx context.Context, params QueryParams[Post], selects *PostSelect, omits *PostOmit) (*Post, error)
type PostFindManyQuery = func(ctx context.Context, params QueryParams[Post], selects *PostSelect, omits *PostOmit) ([]*Post, error)

type PostExtension struct {
	Create              func(ctx context.Context, input *PostCreate, next PostCreateQuery) (*Post, error)
	CreateMany          func(ctx context.Context, inputs []*PostCreate, next PostCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, inputs []*PostCreate, next PostCreateManyAndReturnQuery) ([]*Post, error)
	FindUnique          func(ctx context.Context, where UniquePredicate[Post], additional []PredicateOf[Post], selects *PostSelect, omits *PostOmit, next PostFindUniqueQuery) (*Post, error)
	FindFirst           func(ctx context.Context, params QueryParams[Post], selects *PostSelect, omits *PostOmit, next PostFindFirstQuery) (*Post, error)
	FindMany            func(ctx context.Context, params QueryParams[Post], selects *PostSelect, omits *PostOmit, next PostFindManyQuery) ([]*Post, error)
}

type PostDelegate struct {
	client     *Queries
	extensions []PostExtension
}

func (d *PostDelegate) Use(exts ...PostExtension) {
	d.extensions = append(d.extensions, exts...)
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

func (b *PostCreateBuilder) OnConflict(target UniqueConstraintTarget) *PostConflictBuilder[PostCreateBuilder] {
	return &PostConflictBuilder[PostCreateBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
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
	} else {
		m["published"] = rawDefault{}
	}
	m["authorId"] = s.AuthorId
	return m
}

func (q *Queries) executePostCreate(ctx context.Context, assignments []FieldAssignment, selects *PostSelect, omits *PostOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*Post, error) {
	input := assignmentsToPostCreate(assignments)

	curr := func(c context.Context, args *PostCreate) (*Post, error) {
		if err := validatePostCreate(assignments); err != nil {
			return nil, err
		}

		rowMap := args.ToRowMap()
		cols, vals := mapToColsVals(rowMap, PostColOrder)

		returningCols := q.selectPostCols(selects, omits)

		scanFunc := func(res *Post, cols []string) []any {
			return res.ScanFields(cols)
		}

		pkCols := []string{
			"id",
		}

		hasRelations := selects.hasAnyRelation()

		var res *Post
		var err error
		if hasRelations {
			err = q.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = executeInsert(c, txQ, "Post", cols, vals, returningCols, pkCols, scanFunc, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.loadPostRelations(c, []*Post{res}, selects)
			})
		} else {
			res, err = executeInsert(c, q, "Post", cols, vals, returningCols, pkCols, scanFunc, conflictTarget, conflictAction)
		}
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	for _, ext := range slices.Backward(q.Post.extensions) {
		if ext.Create != nil {
			next, hook := curr, ext.Create
			curr = func(c context.Context, input *PostCreate) (*Post, error) {
				return hook(c, input, next)
			}
		}
	}

	return curr(ctx, &input)
}

type PostCreateManyBuilder struct {
	*CreateManyBuilder[Post]
}

func (b *PostCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *PostConflictBuilder[PostCreateManyBuilder] {
	return &PostConflictBuilder[PostCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type PostCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[Post, PostSelect, PostOmit]
}

func (b *PostCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *PostConflictBuilder[PostCreateManyAndReturnBuilder] {
	return &PostConflictBuilder[PostCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (d *PostDelegate) CreateMany(builders ...*PostCreateBuilder) *PostCreateManyBuilder {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &PostCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[Post]{
			client:   d.client,
			records:  records,
			execFunc: d.client.executePostCreateMany,
		},
	}
}

func (d *PostDelegate) CreateManyAndReturn(builders ...*PostCreateBuilder) *PostCreateManyAndReturnBuilder {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &PostCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[Post, PostSelect, PostOmit]{
			client:   d.client,
			records:  records,
			execFunc: d.client.executePostCreateManyAndReturn,
		},
	}
}

func (q *Queries) executePostCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs := make([]*PostCreate, len(records))
	for i, rec := range records {
		if err := validatePostCreate(rec.Assignments); err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		input := assignmentsToPostCreate(rec.Assignments)
		inputs[i] = &input
	}

	curr := func(c context.Context, args []*PostCreate) (int64, error) {
		rowMaps := make([]map[string]any, len(args))
		for i, input := range args {
			rowMaps[i] = input.ToRowMap()
		}

		pkCols := []string{
			"id",
		}

		return executeCreateMany(c, q, rowMaps, "Post", PostColOrder, pkCols, conflictTarget, conflictAction)
	}

	for _, ext := range slices.Backward(q.Post.extensions) {
		if ext.CreateMany != nil {
			next, hook := curr, ext.CreateMany
			curr = func(c context.Context, inputs []*PostCreate) (int64, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
}

func (q *Queries) executePostCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *PostSelect, omits *PostOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*Post, error) {
	inputs := make([]*PostCreate, len(records))
	for i, rec := range records {
		if err := validatePostCreate(rec.Assignments); err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		input := assignmentsToPostCreate(rec.Assignments)
		inputs[i] = &input
	}

	curr := func(c context.Context, args []*PostCreate) ([]*Post, error) {
		rowMaps := make([]map[string]any, len(args))
		for i, input := range args {
			rowMaps[i] = input.ToRowMap()
		}

		pkCols := []string{
			"id",
		}

		return executeCreateManyAndReturn(c, q, rowMaps, "Post", PostColOrder, selects, omits,
			q.selectPostCols,
			q.loadPostRelations,
			(*Post).ScanFields,
			(*PostSelect).hasAnyRelation,
			pkCols,
			conflictTarget,
			conflictAction,
		)
	}

	for _, ext := range slices.Backward(q.Post.extensions) {
		if ext.CreateManyAndReturn != nil {
			next, hook := curr, ext.CreateManyAndReturn
			curr = func(c context.Context, inputs []*PostCreate) ([]*Post, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
}

type PostConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *PostConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *PostConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *PostConflictBuilder[B]) Update(fn func(u *PostUpsert)) *B {
	var up ConflictUpdate
	u := newPostUpsert(&up)
	fn(u)
	cb.setAction(ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}, cb.conflictTarget)
	return cb.builder
}

type PostUpsert struct {
	Id        fieldUpsert[string]
	Title     fieldUpsert[string]
	Content   fieldUpsert[*string]
	Published fieldUpsert[bool]
	AuthorId  fieldUpsert[string]
}

func newPostUpsert(up *ConflictUpdate) *PostUpsert {
	return &PostUpsert{
		Id:        fieldUpsert[string]{column: "id", update: up},
		Title:     fieldUpsert[string]{column: "title", update: up},
		Content:   fieldUpsert[*string]{column: "content", update: up},
		Published: fieldUpsert[bool]{column: "published", update: up},
		AuthorId:  fieldUpsert[string]{column: "authorId", update: up},
	}
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
	curr := func(c context.Context, w UniquePredicate[Post], add []PredicateOf[Post], sel *PostSelect, o *PostOmit) (*Post, error) {
		if err := w.Validate(); err != nil {
			return nil, err
		}
		for _, p := range add {
			if p != nil {
				if err := p.Validate(); err != nil {
					return nil, err
				}
			}
		}
		allPreds := append([]PredicateOf[Post]{w}, add...)
		whereClause, vals := CompilePredicates(q.dialect, allPreds)
		if whereClause != "" {
			whereClause = " WHERE " + whereClause
		}
		returningCols := q.selectPostCols(sel, o)
		return executeSingleWithRelations(c, q, "Post", whereClause, vals, returningCols,
			func(res *Post, cols []string) []any { return res.ScanFields(cols) },
			sel.hasAnyRelation(),
			func(ctx context.Context, txQ *Queries, results []*Post) error {
				return txQ.loadPostRelations(ctx, results, sel)
			},
			nil,
		)
	}

	for _, ext := range slices.Backward(q.Post.extensions) {
		if ext.FindUnique != nil {
			next, hook := curr, ext.FindUnique
			curr = func(c context.Context, w UniquePredicate[Post], add []PredicateOf[Post], sel *PostSelect, o *PostOmit) (*Post, error) {
				return hook(c, w, add, sel, o, next)
			}
		}
	}

	return curr(ctx, where, additional, selects, omits)
}

func (q *Queries) executePostFindFirst(
	ctx context.Context,
	params QueryParams[Post],
	selects *PostSelect,
	omits *PostOmit,
) (*Post, error) {
	curr := func(c context.Context, p QueryParams[Post], sel *PostSelect, o *PostOmit) (*Post, error) {
		for _, pr := range p.Where {
			if pr != nil {
				if err := pr.Validate(); err != nil {
					return nil, err
				}
			}
		}
		whereClause, vals := CompilePredicates(q.dialect, p.Where)
		if whereClause != "" {
			whereClause = " WHERE " + whereClause
		}
		returningCols := q.selectPostCols(sel, o)
		return executeSingleWithRelations(c, q, "Post", whereClause, vals, returningCols,
			func(res *Post, cols []string) []any { return res.ScanFields(cols) },
			sel.hasAnyRelation(),
			func(ctx context.Context, txQ *Queries, results []*Post) error {
				return txQ.loadPostRelations(ctx, results, sel)
			},
			p.Skip,
		)
	}

	for _, ext := range slices.Backward(q.Post.extensions) {
		if ext.FindFirst != nil {
			next, hook := curr, ext.FindFirst
			curr = func(c context.Context, p QueryParams[Post], sel *PostSelect, o *PostOmit) (*Post, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}

func (q *Queries) executePostFindMany(
	ctx context.Context,
	params QueryParams[Post],
	selects *PostSelect,
	omits *PostOmit,
) ([]*Post, error) {
	curr := func(c context.Context, p QueryParams[Post], sel *PostSelect, o *PostOmit) ([]*Post, error) {
		for _, pr := range p.Where {
			if pr != nil {
				if err := pr.Validate(); err != nil {
					return nil, err
				}
			}
		}
		whereClause, vals := CompilePredicates(q.dialect, p.Where)
		if whereClause != "" {
			whereClause = " WHERE " + whereClause
		}
		returningCols := q.selectPostCols(sel, o)
		return executeManyWithRelations(c, q, "Post", whereClause, vals, returningCols,
			func(res *Post, cols []string) []any { return res.ScanFields(cols) },
			sel.hasAnyRelation(),
			func(ctx context.Context, txQ *Queries, results []*Post) error {
				return txQ.loadPostRelations(ctx, results, sel)
			},
			p.Take,
			p.Skip,
		)
	}

	for _, ext := range slices.Backward(q.Post.extensions) {
		if ext.FindMany != nil {
			next, hook := curr, ext.FindMany
			curr = func(c context.Context, p QueryParams[Post], sel *PostSelect, o *PostOmit) ([]*Post, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
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
