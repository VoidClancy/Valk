package valk

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
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

// colMask returns a bit mask of columns that are set
func (s *PostCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	mask |= 1 << 1
	if s.Content != nil {
		mask |= 1 << 2
	}
	if s.Published != nil {
		mask |= 1 << 3
	}
	mask |= 1 << 4
	return mask
}

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
type PostDeleteManyQuery = func(ctx context.Context, preds []PredicateOf[Post]) (int64, error)
type PostDeleteQuery = func(ctx context.Context, where UniquePredicate[Post], selects *PostSelect, omits *PostOmit) (*Post, error)
type PostCountQuery = func(ctx context.Context, params QueryParams[Post]) (int64, error)

type PostExtension struct {
	Create              func(ctx context.Context, input *PostCreate, next PostCreateQuery) (*Post, error)
	CreateMany          func(ctx context.Context, inputs []*PostCreate, next PostCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, inputs []*PostCreate, next PostCreateManyAndReturnQuery) ([]*Post, error)
	FindUnique          func(ctx context.Context, where UniquePredicate[Post], additional []PredicateOf[Post], selects *PostSelect, omits *PostOmit, next PostFindUniqueQuery) (*Post, error)
	FindFirst           func(ctx context.Context, params QueryParams[Post], selects *PostSelect, omits *PostOmit, next PostFindFirstQuery) (*Post, error)
	FindMany            func(ctx context.Context, params QueryParams[Post], selects *PostSelect, omits *PostOmit, next PostFindManyQuery) ([]*Post, error)
	DeleteMany          func(ctx context.Context, preds []PredicateOf[Post], next PostDeleteManyQuery) (int64, error)
	Delete              func(ctx context.Context, where UniquePredicate[Post], selects *PostSelect, omits *PostOmit, next PostDeleteQuery) (*Post, error)
	Count               func(ctx context.Context, params QueryParams[Post], next PostCountQuery) (int64, error)
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

var postPKCols = []string{
	"id",
}

func selectPostCols(selects *PostSelect, omits *PostOmit, forceCols ...string) []string {
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
			assignments: assignments,
			execFunc:    d.executeCreate,
		},
	}
}

const (
	providedPostId        uint64 = 1 << 0
	providedPostTitle     uint64 = 1 << 1
	providedPostContent   uint64 = 1 << 2
	providedPostPublished uint64 = 1 << 3
	providedPostAuthorId  uint64 = 1 << 4
)

func assignmentsToPostCreate(assignments []FieldAssignment) (PostCreate, error) {
	var input PostCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "id":
			provided |= providedPostId
			if v, ok := a.Val.(string); ok {
				input.Id = &v
				ValidateString(&errs, "id", v, false, 0, false, false)
			} else {
				errs.Add("id", a.Val, "type", "field id must be of type string")
			}
		case "title":
			provided |= providedPostTitle
			if v, ok := a.Val.(string); ok {
				input.Title = v
				ValidateString(&errs, "title", v, true, 0, false, false)
			} else {
				errs.Add("title", a.Val, "type", "field title must be of type string")
			}
		case "content":
			provided |= providedPostContent
			if v, ok := a.Val.(string); ok {
				input.Content = &v
				ValidateString(&errs, "content", v, false, 0, false, false)
			} else {
				errs.Add("content", a.Val, "type", "field content must be of type string")
			}
		case "published":
			provided |= providedPostPublished
			if v, ok := a.Val.(bool); ok {
				input.Published = &v
			} else {
				errs.Add("published", a.Val, "type", "field published must be of type bool")
			}
		case "authorId":
			provided |= providedPostAuthorId
			if v, ok := a.Val.(string); ok {
				input.AuthorId = v
				ValidateString(&errs, "authorId", v, true, 0, false, false)
			} else {
				errs.Add("authorId", a.Val, "type", "field authorId must be of type string")
			}
		}
	}
	if provided&providedPostTitle == 0 {
		errs.Add("title", "", "required", "field Title is required")
	}
	if provided&providedPostAuthorId == 0 {
		errs.Add("authorId", "", "required", "field AuthorId is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *PostCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 5)
	vals = make([]any, 0, 5)
	cols = append(cols, "id")
	if s.Id != nil {
		vals = append(vals, *s.Id)
	} else {
		vals = append(vals, generateCUID())
	}
	cols = append(cols, "title")
	vals = append(vals, s.Title)
	if s.Content != nil {
		cols = append(cols, "content")
		vals = append(vals, *s.Content)
	}
	if s.Published != nil {
		cols = append(cols, "published")
		vals = append(vals, *s.Published)
	}
	cols = append(cols, "authorId")
	vals = append(vals, s.AuthorId)
	return
}

func partitionPostInputs(dialect Dialect, inputs []*PostCreate) [][]*PostCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*PostCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*PostCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*PostCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*PostCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*PostCreate{inputs}
}

func (d *PostDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *PostSelect, omits *PostOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*Post, error) {
	input, err := assignmentsToPostCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectPostCols(selects, omits)

	if len(d.extensions) == 0 {
		hasRelations := selects.hasAnyRelation()
		if hasRelations {
			var res *Post
			err = d.client.transaction(ctx, func(txQ *Queries) error {
				var err error
				res, err = txQ.Post.runCreate(ctx, cols, vals, returningCols, postPKCols, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.Post.loadRelations(ctx, []*Post{res}, selects)
			})
			return res, err
		}
		return d.runCreate(ctx, cols, vals, returningCols, postPKCols, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args *PostCreate) (*Post, error) {
		cols, vals := args.ToColsVals()
		returningCols := selectPostCols(selects, omits)

		hasRelations := selects.hasAnyRelation()
		var res *Post
		var err error
		if hasRelations {
			err = d.client.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = txQ.Post.runCreate(c, cols, vals, returningCols, postPKCols, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.Post.loadRelations(c, []*Post{res}, selects)
			})
		} else {
			res, err = d.runCreate(c, cols, vals, returningCols, postPKCols, conflictTarget, conflictAction)
		}
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	for _, ext := range slices.Backward(d.extensions) {
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
			records:  records,
			execFunc: d.executeCreateMany,
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
			records:  records,
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func (d *PostDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs := make([]*PostCreate, len(records))
	for i, rec := range records {
		input, err := assignmentsToPostCreate(rec.Assignments)
		if err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &input
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args []*PostCreate) (int64, error) {
		return d.runCreateMany(c, args, conflictTarget, conflictAction)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.CreateMany != nil {
			next, hook := curr, ext.CreateMany
			curr = func(c context.Context, inputs []*PostCreate) (int64, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
}

func (d *PostDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *PostSelect, omits *PostOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*Post, error) {
	inputs := make([]*PostCreate, len(records))
	for i, rec := range records {
		input, err := assignmentsToPostCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &input
	}

	if len(d.extensions) == 0 {
		hasRelations := selects != nil && selects.hasAnyRelation()
		if hasRelations {
			var res []*Post
			err := d.client.transaction(ctx, func(txQ *Queries) error {
				var err error
				res, err = txQ.Post.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.Post.loadRelations(ctx, res, selects)
			})
			return res, err
		}
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args []*PostCreate) ([]*Post, error) {
		hasRelations := selects != nil && selects.hasAnyRelation()
		if hasRelations {
			var res []*Post
			err := d.client.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = txQ.Post.runCreateManyAndReturn(c, args, selects, omits, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.Post.loadRelations(c, res, selects)
			})
			return res, err
		}
		return d.runCreateManyAndReturn(c, args, selects, omits, conflictTarget, conflictAction)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.CreateManyAndReturn != nil {
			next, hook := curr, ext.CreateManyAndReturn
			curr = func(c context.Context, inputs []*PostCreate) ([]*Post, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
}

func (d *PostDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	pkCols []string,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*Post, error) {
	query, clauseArgs := buildSingleInsertSQL(d.client, "Post", cols, returningCols, pkCols, conflictTarget, conflictAction, len(vals))
	if len(clauseArgs) > 0 {
		vals = append(vals, clauseArgs...)
	}

	var res Post
	if d.client.dialect.SupportsInsertReturning {
		rows, err := d.client.query(ctx, query, vals...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		if rows.Next() {
			if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
				return nil, err
			}
			return &res, nil
		}
		return nil, rows.Err()
	}

	return d.runCreateFallback(ctx, query, vals, cols, returningCols, pkCols)
}

func (d *PostDelegate) runCreateFallback(ctx context.Context, query string, vals []any, cols []string, returningCols []string, pkCols []string) (*Post, error) {
	result, err := d.client.exec(ctx, query, vals...)
	if err != nil {
		return nil, err
	}

	var pkVals []any
	for _, pkCol := range pkCols {
		var val any
		for i, c := range cols {
			if c == pkCol {
				val = vals[i]
				break
			}
		}
		if val == nil && len(pkCols) == 1 {
			lastID, err := result.LastInsertId()
			if err != nil {
				return nil, err
			}
			val = lastID
		}
		pkVals = append(pkVals, val)
	}

	var selectSb strings.Builder
	selectSb.Grow(64 + len(returningCols)*15 + len("Post") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "Post")
	selectSb.WriteString(" WHERE ")
	for i, pkCol := range pkCols {
		if i > 0 {
			selectSb.WriteString(" AND ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, pkCol)
		selectSb.WriteString(" = ")
		d.client.dialect.WritePlaceholder(&selectSb, i+1)
	}

	rows, err := d.client.query(ctx, selectSb.String(), pkVals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res Post
	if rows.Next() {
		if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
			return nil, err
		}
		return &res, nil
	}
	return nil, rows.Err()
}

func (d *PostDelegate) buildBulkInsertSQL(q *Queries, batch []*PostCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 5)
	for i, c := range postDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "Post")
	sb.WriteString(" (")
	for i, col := range cols {
		if i > 0 {
			sb.WriteString(", ")
		}
		q.dialect.WriteQuotedIdent(&sb, col)
	}
	sb.WriteString(") VALUES ")

	paramIdx := paramStartIdx
	for ri, input := range batch {
		if ri > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("(")
		for j, col := range cols {
			if j > 0 {
				sb.WriteString(", ")
			}
			writeDefault := false
			switch col {
			case "id":
				if input.Id != nil {
					vals = append(vals, *input.Id)
				} else {
					vals = append(vals, generateCUID())
				}
			case "title":
				vals = append(vals, input.Title)
			case "content":
				if input.Content != nil {
					vals = append(vals, *input.Content)
				} else {
					writeDefault = true
				}
			case "published":
				if input.Published != nil {
					vals = append(vals, *input.Published)
				} else {
					writeDefault = true
				}
			case "authorId":
				vals = append(vals, input.AuthorId)
			}
			if writeDefault {
				sb.WriteString("DEFAULT")
			} else {
				q.dialect.WritePlaceholder(&sb, paramIdx)
				paramIdx++
			}
		}
		sb.WriteString(")")
	}
	queryStr = sb.String()
	return cols, vals, queryStr
}

func (d *PostDelegate) runCreateMany(ctx context.Context, inputs []*PostCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionPostInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)

		var conflictCols []string
		if conflictTarget != nil {
			conflictCols = conflictTarget.UniqueColumns()
		}
		var nonConflictCols []string
		if conflictAction != nil && conflictAction.Type == ConflictActionUpdateNewValues {
			nonConflictCols = computeNonConflictCols(cols, conflictCols, postPKCols)
		}
		clause, clauseArgs := d.client.dialect.BuildConflictClause(conflictCols, conflictAction, nonConflictCols, len(vals)+1)
		queryStr += clause
		vals = append(vals, clauseArgs...)

		result, err := d.client.exec(ctx, queryStr, vals...)
		if err != nil {
			return 0, err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		count += affected
	}
	return count, nil
}

func (d *PostDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*PostCreate,
	selects *PostSelect,
	omits *PostOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*Post, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionPostInputs(d.client.dialect, inputs)
	returningCols := selectPostCols(selects, omits)
	hasRelations := selects != nil && selects.hasAnyRelation()

	recordsOut := make([]*Post, 0, len(inputs))

	runBatch := func(txQ *Queries, batch []*PostCreate) error {
		cols, vals, queryStr := d.buildBulkInsertSQL(txQ, batch, 1)

		var conflictCols []string
		if conflictTarget != nil {
			conflictCols = conflictTarget.UniqueColumns()
		}
		var nonConflictCols []string
		if conflictAction != nil && conflictAction.Type == ConflictActionUpdateNewValues {
			nonConflictCols = computeNonConflictCols(cols, conflictCols, postPKCols)
		}
		clause, clauseArgs := txQ.dialect.BuildConflictClause(conflictCols, conflictAction, nonConflictCols, len(vals)+1)
		queryStr += clause
		vals = append(vals, clauseArgs...)

		if txQ.dialect.SupportsInsertReturning && len(returningCols) > 0 {
			var retSb strings.Builder
			retSb.Grow(12 + len(returningCols)*15)
			retSb.WriteString(" RETURNING ")
			for i, col := range returningCols {
				if i > 0 {
					retSb.WriteString(", ")
				}
				txQ.dialect.WriteQuotedIdent(&retSb, col)
			}
			queryStr += retSb.String()
			rows, err := txQ.query(ctx, queryStr, vals...)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var res Post
				if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
					return err
				}
				recordsOut = append(recordsOut, &res)
			}
			return rows.Err()
		}

		// Fallback for dialects without RETURNING (MySQL)
		result, err := txQ.exec(ctx, queryStr, vals...)
		if err != nil {
			return err
		}

		// We need to fetch the inserted records for this batch
		// Note: MySQL bulk inserts only return the ID of the FIRST inserted row
		lastID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// Query back the rows by IDs (assuming autoincrement ID and single PK)
		// If composite PK, it's more complex, but this is a standard fallback
		var selectSb strings.Builder
		selectSb.Grow(64 + len(returningCols)*15 + len("Post") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			txQ.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		txQ.dialect.WriteQuotedIdent(&selectSb, "Post")
		selectSb.WriteString(" WHERE ")
		txQ.dialect.WriteQuotedIdent(&selectSb, postPKCols[0])
		selectSb.WriteString(" >= ")
		txQ.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		txQ.dialect.WriteQuotedIdent(&selectSb, postPKCols[0])
		selectSb.WriteString(" < ")
		txQ.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := txQ.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var res Post
			if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
				return err
			}
			recordsOut = append(recordsOut, &res)
		}
		return rows.Err()
	}

	// Always wrap in transaction if we have multiple batches OR if we need to load relations
	if len(batches) > 1 || hasRelations || !d.client.dialect.SupportsInsertReturning {
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			for _, batch := range batches {
				if err := runBatch(txQ, batch); err != nil {
					return err
				}
			}
			if hasRelations {
				return txQ.Post.loadRelations(ctx, recordsOut, selects)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		if err := runBatch(d.client, batches[0]); err != nil {
			return nil, err
		}
	}

	return recordsOut, nil
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
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *PostDelegate) FindFirst(preds ...PredicateOf[Post]) *FindFirstBuilder[Post, PostSelect, PostOmit] {
	return &FindFirstBuilder[Post, PostSelect, PostOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *PostDelegate) FindMany(preds ...PredicateOf[Post]) *FindManyBuilder[Post, PostSelect, PostOmit] {
	return &FindManyBuilder[Post, PostSelect, PostOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *PostDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[Post], additional []PredicateOf[Post], selects *PostSelect, omits *PostOmit) (*Post, error) {
	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, where, additional, selects, omits)
	}

	curr := func(c context.Context, w UniquePredicate[Post], add []PredicateOf[Post], sel *PostSelect, o *PostOmit) (*Post, error) {
		return d.runFindUnique(c, w, add, sel, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindUnique != nil {
			next, hook := curr, ext.FindUnique
			curr = func(c context.Context, w UniquePredicate[Post], add []PredicateOf[Post], sel *PostSelect, o *PostOmit) (*Post, error) {
				return hook(c, w, add, sel, o, next)
			}
		}
	}

	return curr(ctx, where, additional, selects, omits)
}

func (d *PostDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[Post],
	selects *PostSelect,
	omits *PostOmit,
) (*Post, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	curr := func(c context.Context, p QueryParams[Post], sel *PostSelect, o *PostOmit) (*Post, error) {
		return d.runFindFirst(c, p, sel, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindFirst != nil {
			next, hook := curr, ext.FindFirst
			curr = func(c context.Context, p QueryParams[Post], sel *PostSelect, o *PostOmit) (*Post, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}

func (d *PostDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[Post],
	selects *PostSelect,
	omits *PostOmit,
) ([]*Post, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	curr := func(c context.Context, p QueryParams[Post], sel *PostSelect, o *PostOmit) ([]*Post, error) {
		return d.runFindMany(c, p, sel, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindMany != nil {
			next, hook := curr, ext.FindMany
			curr = func(c context.Context, p QueryParams[Post], sel *PostSelect, o *PostOmit) ([]*Post, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}

func (d *PostDelegate) runFindUnique(ctx context.Context, where UniquePredicate[Post], additional []PredicateOf[Post], selects *PostSelect, omits *PostOmit) (*Post, error) {
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
	whereClause, vals := CompilePredicates(d.client.dialect, allPreds)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	returningCols := selectPostCols(selects, omits)

	var res *Post
	var err error
	if selects.hasAnyRelation() {
		err = d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Post.queryOne(ctx, whereClause, vals, returningCols, nil)
			if err != nil {
				return err
			}
			if res == nil {
				return nil
			}
			return txQ.Post.loadRelations(ctx, []*Post{res}, selects)
		})
	} else {
		res, err = d.queryOne(ctx, whereClause, vals, returningCols, nil)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *PostDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[Post],
	selects *PostSelect,
	omits *PostOmit,
) (*Post, error) {
	for _, pr := range params.Where {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return nil, err
			}
		}
	}
	whereClause, vals := CompilePredicates(d.client.dialect, params.Where)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	returningCols := selectPostCols(selects, omits)

	var res *Post
	var err error
	if selects.hasAnyRelation() {
		err = d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Post.queryOne(ctx, whereClause, vals, returningCols, params.Skip)
			if err != nil {
				return err
			}
			if res == nil {
				return nil
			}
			return txQ.Post.loadRelations(ctx, []*Post{res}, selects)
		})
	} else {
		res, err = d.queryOne(ctx, whereClause, vals, returningCols, params.Skip)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *PostDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[Post],
	selects *PostSelect,
	omits *PostOmit,
) ([]*Post, error) {
	for _, pr := range params.Where {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return nil, err
			}
		}
	}
	whereClause, vals := CompilePredicates(d.client.dialect, params.Where)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	returningCols := selectPostCols(selects, omits)

	var results []*Post
	var err error
	if selects.hasAnyRelation() {
		err = d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			results, err = txQ.Post.queryMany(ctx, whereClause, vals, returningCols, params.Take, params.Skip)
			if err != nil {
				return err
			}
			if len(results) == 0 {
				return nil
			}
			return txQ.Post.loadRelations(ctx, results, selects)
		})
	} else {
		results, err = d.queryMany(ctx, whereClause, vals, returningCols, params.Take, params.Skip)
	}
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (d *PostDelegate) queryOne(ctx context.Context, whereClause string, whereVals []any, returningCols []string, skip *int) (*Post, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "Post", returningCols, whereClause, &limitOne, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		return nil, nil
	}

	var res Post
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (d *PostDelegate) queryMany(ctx context.Context, whereClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*Post, error) {
	query := buildSelectSQL(d.client, "Post", returningCols, whereClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*Post, 0)
	for rows.Next() {
		var res Post
		if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
			return nil, err
		}
		results = append(results, &res)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
func (d *PostDelegate) DeleteMany(preds ...PredicateOf[Post]) *DeleteManyBuilder[Post] {
	return &DeleteManyBuilder[Post]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *PostDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[Post]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	curr := func(c context.Context, p []PredicateOf[Post]) (int64, error) {
		return d.runDeleteMany(c, p)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.DeleteMany != nil {
			next, hook := curr, ext.DeleteMany
			curr = func(c context.Context, p []PredicateOf[Post]) (int64, error) {
				return hook(c, p, next)
			}
		}
	}

	return curr(ctx, preds)
}

func (d *PostDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[Post]) (int64, error) {
	for _, pr := range preds {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return 0, err
			}
		}
	}

	whereClause, vals := CompilePredicates(d.client.dialect, preds)

	var sb strings.Builder
	sb.WriteString("DELETE FROM ")
	d.client.dialect.WriteQuotedIdent(&sb, "Post")
	if whereClause != "" {
		sb.WriteString(" WHERE ")
		sb.WriteString(whereClause)
	}

	result, err := d.client.exec(ctx, sb.String(), vals...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (d *PostDelegate) Delete(where UniquePredicate[Post]) *DeleteBuilder[Post, PostSelect, PostOmit] {
	return &DeleteBuilder[Post, PostSelect, PostOmit]{
		where:    where,
		execFunc: d.executeDelete,
	}
}

func (d *PostDelegate) executeDelete(ctx context.Context, where UniquePredicate[Post], selects *PostSelect, omits *PostOmit) (*Post, error) {
	if len(d.extensions) == 0 {
		return d.runDelete(ctx, where, selects, omits)
	}

	curr := func(c context.Context, w UniquePredicate[Post], s *PostSelect, o *PostOmit) (*Post, error) {
		return d.runDelete(c, w, s, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Delete != nil {
			next, hook := curr, ext.Delete
			curr = func(c context.Context, w UniquePredicate[Post], s *PostSelect, o *PostOmit) (*Post, error) {
				return hook(c, w, s, o, next)
			}
		}
	}

	return curr(ctx, where, selects, omits)
}

func (d *PostDelegate) runDelete(ctx context.Context, where UniquePredicate[Post], selects *PostSelect, omits *PostOmit) (*Post, error) {
	if err := where.Validate(); err != nil {
		return nil, err
	}

	returningCols := selectPostCols(selects, omits, postPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *Post
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Post.executeFindUnique(ctx, where, nil, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return sql.ErrNoRows
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "Post")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[Post]
			pkPreds = append(pkPreds, Predicate[Post]{
				Data: PredicateData{
					Column:   "id",
					Operator: "=",
					Value:    res.Id,
				},
			})

			whereClause, vals := CompilePredicates(txQ.dialect, pkPreds)
			deleteSb.WriteString(whereClause)

			_, err = txQ.exec(ctx, deleteSb.String(), vals...)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	// Dialect supports RETURNING, and no relations need loading: run direct DELETE ... RETURNING
	var sb strings.Builder
	sb.WriteString("DELETE FROM ")
	d.client.dialect.WriteQuotedIdent(&sb, "Post")

	whereClause, vals := CompilePredicates(d.client.dialect, []PredicateOf[Post]{where})
	if whereClause != "" {
		sb.WriteString(" WHERE ")
		sb.WriteString(whereClause)
	}

	sb.WriteString(" RETURNING ")
	for i, col := range returningCols {
		if i > 0 {
			sb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&sb, col)
	}

	rows, err := d.client.query(ctx, sb.String(), vals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	}

	var row Post
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, err
	}
	return &row, nil
}
func (d *PostDelegate) Count(preds ...PredicateOf[Post]) *CountBuilder[Post] {
	return &CountBuilder[Post]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *PostDelegate) executeCount(ctx context.Context, params QueryParams[Post]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	curr := func(c context.Context, p QueryParams[Post]) (int64, error) {
		return d.runCount(c, p)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Count != nil {
			next, hook := curr, ext.Count
			curr = func(c context.Context, p QueryParams[Post]) (int64, error) {
				return hook(c, p, next)
			}
		}
	}

	return curr(ctx, params)
}

func (d *PostDelegate) runCount(ctx context.Context, params QueryParams[Post]) (int64, error) {
	for _, pr := range params.Where {
		if pr != nil {
			if err := pr.Validate(); err != nil {
				return 0, err
			}
		}
	}

	whereClause, vals := CompilePredicates(d.client.dialect, params.Where)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}

	var query string
	if params.Take != nil || params.Skip != nil {
		var subQuery strings.Builder
		subQuery.WriteString("SELECT 1 FROM ")
		d.client.dialect.WriteQuotedIdent(&subQuery, "Post")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "Post")
		if whereClause != "" {
			sb.WriteString(whereClause)
		}
		query = sb.String()
	}

	rows, err := d.client.query(ctx, query, vals...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var count int64
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	return count, nil
}
func (d *PostDelegate) loadRelations(ctx context.Context, records []*Post, selects *PostSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Author != nil {
		relationSelects, relationOmits, relationParams := selects.Author.GetRelationParams()
		returningCols := selectUserCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: Post.authorId
		allChildren, err := loadRelation(
			ctx, d.client, records,
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
		if err := d.client.User.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Comments != nil {
		relationSelects, relationOmits, relationParams := selects.Comments.GetRelationParams()
		returningCols := selectCommentCols(relationSelects, relationOmits, "postId")
		// Inverse holds the FK: Comment.postId
		allChildren, err := loadRelation(
			ctx, d.client, records,
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
		if err := d.client.Comment.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Categories != nil {
		relationSelects, relationOmits, relationParams := selects.Categories.GetRelationParams()
		returningCols := selectCategoryToPostCols(relationSelects, relationOmits, "postId")
		// Inverse holds the FK: CategoryToPost.postId
		allChildren, err := loadRelation(
			ctx, d.client, records,
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
		if err := d.client.CategoryToPost.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
