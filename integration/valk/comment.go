package valk

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
)

// Comment represents the database model
type Comment struct {
	Id       string           `db:"id" json:"id"`
	Textify  int32            `db:"textify" json:"textify"`
	Dummy3   string           `db:"dummy3" json:"dummy3"`
	Dummy1   int32            `db:"dummy1" json:"dummy1"`
	Dummy2   string           `db:"dummy2" json:"dummy2"`
	PostId   string           `db:"postId" json:"postId"`
	AuthorId string           `db:"authorId" json:"authorId"`
	Meta     *json.RawMessage `db:"meta" json:"meta,omitempty"`
	Post     *Post            `json:"post,omitempty"`
	Author   *User            `json:"author,omitempty"`
}

// CommentCreate is used for hooks only — the Create API uses FieldAssignment
type CommentCreate struct {
	Id       *string          `json:"id"`
	Textify  int32            `json:"textify"`
	Dummy3   string           `json:"dummy3"`
	Dummy1   int32            `json:"dummy1"`
	Dummy2   string           `json:"dummy2"`
	PostId   string           `json:"postId"`
	AuthorId string           `json:"authorId"`
	Meta     *json.RawMessage `json:"meta"`
}

// CommentSelect specifies which fields to include
type CommentSelect struct {
	Id       bool            `json:"id"`
	Textify  bool            `json:"textify"`
	Dummy3   bool            `json:"dummy3"`
	Dummy1   bool            `json:"dummy1"`
	Dummy2   bool            `json:"dummy2"`
	PostId   bool            `json:"postId"`
	AuthorId bool            `json:"authorId"`
	Meta     bool            `json:"meta"`
	Post     PostSelectQuery `json:"post,omitempty"`
	Author   UserSelectQuery `json:"author,omitempty"`
}

// CommentOmit specifies which fields to exclude
type CommentOmit struct {
	Id       bool `json:"id"`
	Textify  bool `json:"textify"`
	Dummy3   bool `json:"dummy3"`
	Dummy1   bool `json:"dummy1"`
	Dummy2   bool `json:"dummy2"`
	PostId   bool `json:"postId"`
	AuthorId bool `json:"authorId"`
	Meta     bool `json:"meta"`
}

type CommentSelectQuery interface {
	GetRelationParams() (*CommentSelect, *CommentOmit, QueryParams[Comment])
}

func (s *CommentSelect) GetRelationParams() (*CommentSelect, *CommentOmit, QueryParams[Comment]) {
	return s, nil, QueryParams[Comment]{}
}

// CommentQueryBuilder builds a query for the relation Comment
type CommentQueryBuilder struct {
	selects *CommentSelect
	omits   *CommentOmit
	where   []PredicateOf[Comment]
	take    *int
	skip    *int
	orderBy []OrderBy[Comment]
}

func (b *CommentQueryBuilder) Where(preds ...PredicateOf[Comment]) *CommentQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *CommentQueryBuilder) Take(limit int) *CommentQueryBuilder {
	b.take = &limit
	return b
}

func (b *CommentQueryBuilder) Skip(offset int) *CommentQueryBuilder {
	b.skip = &offset
	return b
}

func (b *CommentQueryBuilder) OrderBy(orders ...OrderBy[Comment]) *CommentQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *CommentQueryBuilder) Select(s CommentSelect) *CommentQueryBuilder {
	b.selects = &s
	return b
}

func (b *CommentQueryBuilder) Omit(o CommentOmit) *CommentQueryBuilder {
	b.omits = &o
	return b
}

func (b *CommentQueryBuilder) GetRelationParams() (*CommentSelect, *CommentOmit, QueryParams[Comment]) {
	if b == nil {
		return nil, nil, QueryParams[Comment]{}
	}
	return b.selects, b.omits, QueryParams[Comment]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
	}
}

type CommentCreateQuery = func(ctx context.Context, args *CommentCreate) (*Comment, error)
type CommentCreateManyQuery = func(ctx context.Context, args []*CommentCreate) (int64, error)
type CommentCreateManyAndReturnQuery = func(ctx context.Context, args []*CommentCreate) ([]*Comment, error)
type CommentFindUniqueQuery = func(ctx context.Context, where UniquePredicate[Comment], additional []PredicateOf[Comment], selects *CommentSelect, omits *CommentOmit) (*Comment, error)
type CommentFindFirstQuery = func(ctx context.Context, params QueryParams[Comment], selects *CommentSelect, omits *CommentOmit) (*Comment, error)
type CommentFindManyQuery = func(ctx context.Context, params QueryParams[Comment], selects *CommentSelect, omits *CommentOmit) ([]*Comment, error)

type CommentExtension struct {
	Create              func(ctx context.Context, input *CommentCreate, next CommentCreateQuery) (*Comment, error)
	CreateMany          func(ctx context.Context, inputs []*CommentCreate, next CommentCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, inputs []*CommentCreate, next CommentCreateManyAndReturnQuery) ([]*Comment, error)
	FindUnique          func(ctx context.Context, where UniquePredicate[Comment], additional []PredicateOf[Comment], selects *CommentSelect, omits *CommentOmit, next CommentFindUniqueQuery) (*Comment, error)
	FindFirst           func(ctx context.Context, params QueryParams[Comment], selects *CommentSelect, omits *CommentOmit, next CommentFindFirstQuery) (*Comment, error)
	FindMany            func(ctx context.Context, params QueryParams[Comment], selects *CommentSelect, omits *CommentOmit, next CommentFindManyQuery) ([]*Comment, error)
}

type CommentDelegate struct {
	client     *Queries
	extensions []CommentExtension
}

func (d *CommentDelegate) Use(exts ...CommentExtension) {
	d.extensions = append(d.extensions, exts...)
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
		case "meta":
			targets[i] = &m.Meta
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
	"meta",
}

func selectCommentCols(selects *CommentSelect, omits *CommentOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return commentDefaultCols
	}

	anySelected := selects != nil && (selects.Id || selects.Textify || selects.Dummy3 || selects.Dummy1 || selects.Dummy2 || selects.PostId || selects.AuthorId || selects.Meta || selects.Post != nil || selects.Author != nil)

	specs := []colSpec{
		{"id", selects != nil && selects.Id, omits != nil && omits.Id, selects != nil && selects.hasAnyRelation()},
		{"textify", selects != nil && selects.Textify, omits != nil && omits.Textify, false},
		{"dummy3", selects != nil && selects.Dummy3, omits != nil && omits.Dummy3, false},
		{"dummy1", selects != nil && selects.Dummy1, omits != nil && omits.Dummy1, false},
		{"dummy2", selects != nil && selects.Dummy2, omits != nil && omits.Dummy2, false},
		{"postId", selects != nil && selects.PostId, omits != nil && omits.PostId, selects != nil && selects.Post != nil},
		{"authorId", selects != nil && selects.AuthorId, omits != nil && omits.AuthorId, selects != nil && selects.Author != nil},
		{"meta", selects != nil && selects.Meta, omits != nil && omits.Meta, false},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *CommentSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Post != nil || s.Author != nil
}

type CommentCreateBuilder struct {
	*CreateBuilder[Comment, CommentSelect, CommentOmit]
}

func (b *CommentCreateBuilder) OnConflict(target UniqueConstraintTarget) *CommentConflictBuilder[CommentCreateBuilder] {
	return &CommentConflictBuilder[CommentCreateBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *CommentCreateBuilder) SetId(v string) *CommentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}
func (b *CommentCreateBuilder) SetTextify(v int32) *CommentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "textify", Val: v})
	return b
}
func (b *CommentCreateBuilder) SetDummy3(v string) *CommentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "dummy3", Val: v})
	return b
}
func (b *CommentCreateBuilder) SetDummy1(v int32) *CommentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "dummy1", Val: v})
	return b
}
func (b *CommentCreateBuilder) SetDummy2(v string) *CommentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "dummy2", Val: v})
	return b
}
func (b *CommentCreateBuilder) SetPostId(v string) *CommentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "postId", Val: v})
	return b
}
func (b *CommentCreateBuilder) SetAuthorId(v string) *CommentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "authorId", Val: v})
	return b
}
func (b *CommentCreateBuilder) SetMeta(v json.RawMessage) *CommentCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "meta", Val: v})
	return b
}

func (d *CommentDelegate) Create(assignments ...FieldAssignment) *CommentCreateBuilder {
	return &CommentCreateBuilder{
		CreateBuilder: &CreateBuilder[Comment, CommentSelect, CommentOmit]{
			assignments: assignments,
			execFunc:    d.executeCreate,
		},
	}
}

const (
	providedCommentId       uint64 = 1 << 0
	providedCommentTextify  uint64 = 1 << 1
	providedCommentDummy3   uint64 = 1 << 2
	providedCommentDummy1   uint64 = 1 << 3
	providedCommentDummy2   uint64 = 1 << 4
	providedCommentPostId   uint64 = 1 << 5
	providedCommentAuthorId uint64 = 1 << 6
	providedCommentMeta     uint64 = 1 << 7
)

func assignmentsToCommentCreate(assignments []FieldAssignment) (CommentCreate, error) {
	var input CommentCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "id":
			provided |= providedCommentId
			if v, ok := a.Val.(string); ok {
				input.Id = &v
				ValidateString(&errs, "id", v, false, 0, false, false)
			} else {
				errs.Add("id", a.Val, "type", "field id must be of type string")
			}
		case "textify":
			provided |= providedCommentTextify
			if v, ok := a.Val.(int32); ok {
				input.Textify = v
				ValidateInt32(&errs, "textify", v, "")
			} else {
				errs.Add("textify", a.Val, "type", "field textify must be of type int32")
			}
		case "dummy3":
			provided |= providedCommentDummy3
			if v, ok := a.Val.(string); ok {
				input.Dummy3 = v
				ValidateString(&errs, "dummy3", v, true, 0, false, false)
			} else {
				errs.Add("dummy3", a.Val, "type", "field dummy3 must be of type string")
			}
		case "dummy1":
			provided |= providedCommentDummy1
			if v, ok := a.Val.(int32); ok {
				input.Dummy1 = v
				ValidateInt32(&errs, "dummy1", v, "")
			} else {
				errs.Add("dummy1", a.Val, "type", "field dummy1 must be of type int32")
			}
		case "dummy2":
			provided |= providedCommentDummy2
			if v, ok := a.Val.(string); ok {
				input.Dummy2 = v
				ValidateString(&errs, "dummy2", v, true, 0, false, false)
			} else {
				errs.Add("dummy2", a.Val, "type", "field dummy2 must be of type string")
			}
		case "postId":
			provided |= providedCommentPostId
			if v, ok := a.Val.(string); ok {
				input.PostId = v
				ValidateString(&errs, "postId", v, true, 0, false, false)
			} else {
				errs.Add("postId", a.Val, "type", "field postId must be of type string")
			}
		case "authorId":
			provided |= providedCommentAuthorId
			if v, ok := a.Val.(string); ok {
				input.AuthorId = v
				ValidateString(&errs, "authorId", v, true, 0, false, false)
			} else {
				errs.Add("authorId", a.Val, "type", "field authorId must be of type string")
			}
		case "meta":
			provided |= providedCommentMeta
			if v, ok := a.Val.(json.RawMessage); ok {
				input.Meta = &v
			} else {
				errs.Add("meta", a.Val, "type", "field meta must be of type json.RawMessage")
			}
		}
	}
	if provided&providedCommentTextify == 0 {
		errs.Add("textify", nil, "required", "field Textify is required")
	}
	if provided&providedCommentDummy3 == 0 {
		errs.Add("dummy3", "", "required", "field Dummy3 is required")
	}
	if provided&providedCommentDummy1 == 0 {
		errs.Add("dummy1", nil, "required", "field Dummy1 is required")
	}
	if provided&providedCommentDummy2 == 0 {
		errs.Add("dummy2", "", "required", "field Dummy2 is required")
	}
	if provided&providedCommentPostId == 0 {
		errs.Add("postId", "", "required", "field PostId is required")
	}
	if provided&providedCommentAuthorId == 0 {
		errs.Add("authorId", "", "required", "field AuthorId is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *CommentCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 8)
	vals = make([]any, 0, 8)
	cols = append(cols, "id")
	if s.Id != nil {
		vals = append(vals, *s.Id)
	} else {
		vals = append(vals, generateCUID())
	}
	cols = append(cols, "textify")
	vals = append(vals, s.Textify)
	cols = append(cols, "dummy3")
	vals = append(vals, s.Dummy3)
	cols = append(cols, "dummy1")
	vals = append(vals, s.Dummy1)
	cols = append(cols, "dummy2")
	vals = append(vals, s.Dummy2)
	cols = append(cols, "postId")
	vals = append(vals, s.PostId)
	cols = append(cols, "authorId")
	vals = append(vals, s.AuthorId)
	if s.Meta != nil {
		cols = append(cols, "meta")
		vals = append(vals, *s.Meta)
	}
	return
}

func (s *CommentCreate) ToRowMap() map[string]any {
	cols, vals := s.ToColsVals()
	m := make(map[string]any, len(cols))
	for i, c := range cols {
		m[c] = vals[i]
	}
	return m
}

func (d *CommentDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *CommentSelect, omits *CommentOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*Comment, error) {
	input, err := assignmentsToCommentCreate(assignments)
	if err != nil {
		return nil, err
	}

	curr := func(c context.Context, args *CommentCreate) (*Comment, error) {
		cols, vals := args.ToColsVals()

		returningCols := selectCommentCols(selects, omits)

		scanFunc := func(res *Comment, cols []string) []any {
			return res.ScanFields(cols)
		}

		pkCols := []string{
			"id",
		}

		hasRelations := selects.hasAnyRelation()

		var res *Comment
		var err error
		if hasRelations {
			err = d.client.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = executeInsert(c, txQ, "Comment", cols, vals, returningCols, pkCols, scanFunc, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.Comment.loadRelations(c, []*Comment{res}, selects)
			})
		} else {
			res, err = executeInsert(c, d.client, "Comment", cols, vals, returningCols, pkCols, scanFunc, conflictTarget, conflictAction)
		}
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Create != nil {
			next, hook := curr, ext.Create
			curr = func(c context.Context, input *CommentCreate) (*Comment, error) {
				return hook(c, input, next)
			}
		}
	}

	return curr(ctx, &input)
}

type CommentCreateManyBuilder struct {
	*CreateManyBuilder[Comment]
}

func (b *CommentCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *CommentConflictBuilder[CommentCreateManyBuilder] {
	return &CommentConflictBuilder[CommentCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type CommentCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[Comment, CommentSelect, CommentOmit]
}

func (b *CommentCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *CommentConflictBuilder[CommentCreateManyAndReturnBuilder] {
	return &CommentConflictBuilder[CommentCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (d *CommentDelegate) CreateMany(builders ...*CommentCreateBuilder) *CommentCreateManyBuilder {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &CommentCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[Comment]{
			records:  records,
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *CommentDelegate) CreateManyAndReturn(builders ...*CommentCreateBuilder) *CommentCreateManyAndReturnBuilder {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &CommentCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[Comment, CommentSelect, CommentOmit]{
			records:  records,
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func (d *CommentDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs := make([]*CommentCreate, len(records))
	for i, rec := range records {
		input, err := assignmentsToCommentCreate(rec.Assignments)
		if err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &input
	}

	curr := func(c context.Context, args []*CommentCreate) (int64, error) {
		rowMaps := make([]map[string]any, len(args))
		for i, input := range args {
			rowMaps[i] = input.ToRowMap()
		}

		pkCols := []string{
			"id",
		}

		return executeCreateMany(c, d.client, rowMaps, "Comment", commentDefaultCols, pkCols, conflictTarget, conflictAction)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.CreateMany != nil {
			next, hook := curr, ext.CreateMany
			curr = func(c context.Context, inputs []*CommentCreate) (int64, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
}

func (d *CommentDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *CommentSelect, omits *CommentOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*Comment, error) {
	inputs := make([]*CommentCreate, len(records))
	for i, rec := range records {
		input, err := assignmentsToCommentCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &input
	}

	curr := func(c context.Context, args []*CommentCreate) ([]*Comment, error) {
		rowMaps := make([]map[string]any, len(args))
		for i, input := range args {
			rowMaps[i] = input.ToRowMap()
		}

		pkCols := []string{
			"id",
		}

		return executeCreateManyAndReturn(c, d.client, rowMaps, "Comment", commentDefaultCols, selects, omits,
			selectCommentCols,
			func(ctx context.Context, txQ *Queries, results []*Comment, sel *CommentSelect) error {
				return txQ.Comment.loadRelations(ctx, results, sel)
			},
			(*Comment).ScanFields,
			(*CommentSelect).hasAnyRelation,
			pkCols,
			conflictTarget,
			conflictAction,
		)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.CreateManyAndReturn != nil {
			next, hook := curr, ext.CreateManyAndReturn
			curr = func(c context.Context, inputs []*CommentCreate) ([]*Comment, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
}

type CommentConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *CommentConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *CommentConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *CommentConflictBuilder[B]) Update(fn func(u *CommentUpsert)) *B {
	var up ConflictUpdate
	u := newCommentUpsert(&up)
	fn(u)
	cb.setAction(ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}, cb.conflictTarget)
	return cb.builder
}

type CommentUpsert struct {
	Id       fieldUpsert[string]
	Textify  numericFieldUpsert[int32]
	Dummy3   fieldUpsert[string]
	Dummy1   numericFieldUpsert[int32]
	Dummy2   fieldUpsert[string]
	PostId   fieldUpsert[string]
	AuthorId fieldUpsert[string]
	Meta     fieldUpsert[*json.RawMessage]
}

func newCommentUpsert(up *ConflictUpdate) *CommentUpsert {
	return &CommentUpsert{
		Id: fieldUpsert[string]{column: "id", update: up},
		Textify: numericFieldUpsert[int32]{
			fieldUpsert: fieldUpsert[int32]{column: "textify", update: up},
			tableName:   "Comment",
		},
		Dummy3: fieldUpsert[string]{column: "dummy3", update: up},
		Dummy1: numericFieldUpsert[int32]{
			fieldUpsert: fieldUpsert[int32]{column: "dummy1", update: up},
			tableName:   "Comment",
		},
		Dummy2:   fieldUpsert[string]{column: "dummy2", update: up},
		PostId:   fieldUpsert[string]{column: "postId", update: up},
		AuthorId: fieldUpsert[string]{column: "authorId", update: up},
		Meta:     fieldUpsert[*json.RawMessage]{column: "meta", update: up},
	}
}
func (d *CommentDelegate) FindUnique(where UniquePredicate[Comment], additional ...PredicateOf[Comment]) *FindUniqueBuilder[Comment, CommentSelect, CommentOmit] {
	return &FindUniqueBuilder[Comment, CommentSelect, CommentOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *CommentDelegate) FindFirst(preds ...PredicateOf[Comment]) *FindFirstBuilder[Comment, CommentSelect, CommentOmit] {
	return &FindFirstBuilder[Comment, CommentSelect, CommentOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *CommentDelegate) FindMany(preds ...PredicateOf[Comment]) *FindManyBuilder[Comment, CommentSelect, CommentOmit] {
	return &FindManyBuilder[Comment, CommentSelect, CommentOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *CommentDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[Comment], additional []PredicateOf[Comment], selects *CommentSelect, omits *CommentOmit) (*Comment, error) {
	curr := func(c context.Context, w UniquePredicate[Comment], add []PredicateOf[Comment], sel *CommentSelect, o *CommentOmit) (*Comment, error) {
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
		allPreds := append([]PredicateOf[Comment]{w}, add...)
		whereClause, vals := CompilePredicates(d.client.dialect, allPreds)
		if whereClause != "" {
			whereClause = " WHERE " + whereClause
		}
		returningCols := selectCommentCols(sel, o)
		return executeSingleWithRelations(c, d.client, "Comment", whereClause, vals, returningCols,
			func(res *Comment, cols []string) []any { return res.ScanFields(cols) },
			sel.hasAnyRelation(),
			func(ctx context.Context, txQ *Queries, results []*Comment) error {
				return txQ.Comment.loadRelations(ctx, results, sel)
			},
			nil,
		)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindUnique != nil {
			next, hook := curr, ext.FindUnique
			curr = func(c context.Context, w UniquePredicate[Comment], add []PredicateOf[Comment], sel *CommentSelect, o *CommentOmit) (*Comment, error) {
				return hook(c, w, add, sel, o, next)
			}
		}
	}

	return curr(ctx, where, additional, selects, omits)
}

func (d *CommentDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[Comment],
	selects *CommentSelect,
	omits *CommentOmit,
) (*Comment, error) {
	curr := func(c context.Context, p QueryParams[Comment], sel *CommentSelect, o *CommentOmit) (*Comment, error) {
		for _, pr := range p.Where {
			if pr != nil {
				if err := pr.Validate(); err != nil {
					return nil, err
				}
			}
		}
		whereClause, vals := CompilePredicates(d.client.dialect, p.Where)
		if whereClause != "" {
			whereClause = " WHERE " + whereClause
		}
		returningCols := selectCommentCols(sel, o)
		return executeSingleWithRelations(c, d.client, "Comment", whereClause, vals, returningCols,
			func(res *Comment, cols []string) []any { return res.ScanFields(cols) },
			sel.hasAnyRelation(),
			func(ctx context.Context, txQ *Queries, results []*Comment) error {
				return txQ.Comment.loadRelations(ctx, results, sel)
			},
			p.Skip,
		)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindFirst != nil {
			next, hook := curr, ext.FindFirst
			curr = func(c context.Context, p QueryParams[Comment], sel *CommentSelect, o *CommentOmit) (*Comment, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}

func (d *CommentDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[Comment],
	selects *CommentSelect,
	omits *CommentOmit,
) ([]*Comment, error) {
	curr := func(c context.Context, p QueryParams[Comment], sel *CommentSelect, o *CommentOmit) ([]*Comment, error) {
		for _, pr := range p.Where {
			if pr != nil {
				if err := pr.Validate(); err != nil {
					return nil, err
				}
			}
		}
		whereClause, vals := CompilePredicates(d.client.dialect, p.Where)
		if whereClause != "" {
			whereClause = " WHERE " + whereClause
		}
		returningCols := selectCommentCols(sel, o)
		return executeManyWithRelations(c, d.client, "Comment", whereClause, vals, returningCols,
			func(res *Comment, cols []string) []any { return res.ScanFields(cols) },
			sel.hasAnyRelation(),
			func(ctx context.Context, txQ *Queries, results []*Comment) error {
				return txQ.Comment.loadRelations(ctx, results, sel)
			},
			p.Take,
			p.Skip,
		)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindMany != nil {
			next, hook := curr, ext.FindMany
			curr = func(c context.Context, p QueryParams[Comment], sel *CommentSelect, o *CommentOmit) ([]*Comment, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}
func (d *CommentDelegate) loadRelations(ctx context.Context, records []*Comment, selects *CommentSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Post != nil {
		relationSelects, relationOmits, relationParams := selects.Post.GetRelationParams()
		returningCols := selectPostCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: Comment.postId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *Comment) string { return p.PostId }),
			"Post",
			"id",
			returningCols,
			scanInto(returningCols, (*Post).ScanFields),
			directKey(func(c *Post) string { return c.Id }),
			setOne(func(p *Comment, c *Post) { p.Post = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading post: %w", err)
		}
		if err := d.client.Post.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Author != nil {
		relationSelects, relationOmits, relationParams := selects.Author.GetRelationParams()
		returningCols := selectUserCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: Comment.authorId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *Comment) string { return p.AuthorId }),
			"User",
			"id",
			returningCols,
			scanInto(returningCols, (*User).ScanFields),
			directKey(func(c *User) string { return c.Id }),
			setOne(func(p *Comment, c *User) { p.Author = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading author: %w", err)
		}
		if err := d.client.User.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
