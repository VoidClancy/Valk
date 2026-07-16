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

type CommentDelegate struct {
	client          *Queries
	beforeCreate    func(context.Context, *CommentCreate) error
	afterCreate     func(context.Context, []*Comment) error
	afterCreateMany func(context.Context, []CommentCreate, int64) error
}

func (d *CommentDelegate) BeforeCreate(hook func(context.Context, *CommentCreate) error) {
	d.beforeCreate = hook
}

func (d *CommentDelegate) AfterCreate(hook func(context.Context, []*Comment) error) {
	d.afterCreate = hook
}

func (d *CommentDelegate) AfterCreateMany(hook func(context.Context, []CommentCreate, int64) error) {
	d.afterCreateMany = hook
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

func (q *Queries) selectCommentCols(selects *CommentSelect, omits *CommentOmit, forceCols ...string) []string {
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

var CommentColOrder = []string{
	"id",
	"textify",
	"dummy3",
	"dummy1",
	"dummy2",
	"postId",
	"authorId",
	"meta",
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
			client:      d.client,
			assignments: assignments,
			execFunc:    d.client.executeCommentCreate,
		},
	}
}

func validateCommentCreate(assignments []FieldAssignment) error {
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
		case "textify":
			if v, ok := a.Val.(int32); ok {
				ValidateInt32(errs, "textify", v, "")
			} else {
				errs.Add("textify", a.Val, "type", "field textify must be of type int32")
			}
		case "dummy3":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "dummy3", v, true, 0, false, false)
			} else {
				errs.Add("dummy3", a.Val, "type", "field dummy3 must be of type string")
			}
		case "dummy1":
			if v, ok := a.Val.(int32); ok {
				ValidateInt32(errs, "dummy1", v, "")
			} else {
				errs.Add("dummy1", a.Val, "type", "field dummy1 must be of type int32")
			}
		case "dummy2":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "dummy2", v, true, 0, false, false)
			} else {
				errs.Add("dummy2", a.Val, "type", "field dummy2 must be of type string")
			}
		case "postId":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "postId", v, true, 0, false, false)
			} else {
				errs.Add("postId", a.Val, "type", "field postId must be of type string")
			}
		case "authorId":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "authorId", v, true, 0, false, false)
			} else {
				errs.Add("authorId", a.Val, "type", "field authorId must be of type string")
			}
		case "meta":
			if _, ok := a.Val.(json.RawMessage); !ok {
				errs.Add("meta", a.Val, "type", "field meta must be of type json.RawMessage")
			}
		}
	}
	if !provided["textify"] {
		errs.Add("textify", nil, "required", "field Textify is required")
	}
	if !provided["dummy3"] {
		errs.Add("dummy3", "", "required", "field Dummy3 is required")
	}
	if !provided["dummy1"] {
		errs.Add("dummy1", nil, "required", "field Dummy1 is required")
	}
	if !provided["dummy2"] {
		errs.Add("dummy2", "", "required", "field Dummy2 is required")
	}
	if !provided["postId"] {
		errs.Add("postId", "", "required", "field PostId is required")
	}
	if !provided["authorId"] {
		errs.Add("authorId", "", "required", "field AuthorId is required")
	}

	if errs.HasErrors() {
		return *errs
	}
	return nil
}

func assignmentsToCommentCreate(assignments []FieldAssignment) CommentCreate {
	var input CommentCreate
	for _, a := range assignments {
		switch a.Col {
		case "id":
			if v, ok := a.Val.(string); ok {
				input.Id = &v
			}
		case "textify":
			if v, ok := a.Val.(int32); ok {
				input.Textify = v
			}
		case "dummy3":
			if v, ok := a.Val.(string); ok {
				input.Dummy3 = v
			}
		case "dummy1":
			if v, ok := a.Val.(int32); ok {
				input.Dummy1 = v
			}
		case "dummy2":
			if v, ok := a.Val.(string); ok {
				input.Dummy2 = v
			}
		case "postId":
			if v, ok := a.Val.(string); ok {
				input.PostId = v
			}
		case "authorId":
			if v, ok := a.Val.(string); ok {
				input.AuthorId = v
			}
		case "meta":
			if v, ok := a.Val.(json.RawMessage); ok {
				input.Meta = &v
			}
		}
	}
	return input
}

func (s *CommentCreate) ToRowMap() map[string]any {
	m := make(map[string]any, 8)
	if s.Id != nil {
		m["id"] = *s.Id
	} else {
		m["id"] = generateCUID()
	}
	m["textify"] = s.Textify
	m["dummy3"] = s.Dummy3
	m["dummy1"] = s.Dummy1
	m["dummy2"] = s.Dummy2
	m["postId"] = s.PostId
	m["authorId"] = s.AuthorId
	if s.Meta != nil {
		m["meta"] = *s.Meta
	}
	return m
}

func (q *Queries) executeCommentCreate(ctx context.Context, assignments []FieldAssignment, selects *CommentSelect, omits *CommentOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*Comment, error) {
	input := assignmentsToCommentCreate(assignments)

	if q.Comment.beforeCreate != nil {
		if err := q.Comment.beforeCreate(ctx, &input); err != nil {
			return nil, err
		}
	}

	if err := validateCommentCreate(assignments); err != nil {
		return nil, err
	}

	rowMap := input.ToRowMap()
	cols, vals := mapToColsVals(rowMap, CommentColOrder)

	returningCols := q.selectCommentCols(selects, omits)

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
		err = q.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = executeInsert(ctx, txQ, "Comment", cols, vals, returningCols, pkCols, scanFunc, conflictTarget, conflictAction)
			if err != nil {
				return err
			}
			return txQ.loadCommentRelations(ctx, []*Comment{res}, selects)
		})
	} else {
		res, err = executeInsert(ctx, q, "Comment", cols, vals, returningCols, pkCols, scanFunc, conflictTarget, conflictAction)
	}
	if err != nil {
		return nil, err
	}

	if q.Comment.afterCreate != nil {
		if err := q.Comment.afterCreate(ctx, []*Comment{res}); err != nil {
			return nil, err
		}
	}

	return res, nil
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
			client:   d.client,
			records:  records,
			execFunc: d.client.executeCommentCreateMany,
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
			client:   d.client,
			records:  records,
			execFunc: d.client.executeCommentCreateManyAndReturn,
		},
	}
}

func (q *Queries) executeCommentCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	rowMaps := make([]map[string]any, len(records))
	inputs := make([]CommentCreate, len(records))
	for i, rec := range records {
		if err := validateCommentCreate(rec.Assignments); err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		input := assignmentsToCommentCreate(rec.Assignments)
		if q.Comment.beforeCreate != nil {
			if err := q.Comment.beforeCreate(ctx, &input); err != nil {
				return 0, err
			}
		}
		rowMaps[i] = input.ToRowMap()
		inputs[i] = input
	}
	pkCols := []string{
		"id",
	}
	count, err := executeCreateMany(ctx, q, rowMaps, "Comment", CommentColOrder, pkCols, conflictTarget, conflictAction)
	if err != nil {
		return 0, err
	}
	if q.Comment.afterCreateMany != nil {
		if err := q.Comment.afterCreateMany(ctx, inputs, count); err != nil {
			return 0, err
		}
	}
	return count, nil
}

func (q *Queries) executeCommentCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *CommentSelect, omits *CommentOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*Comment, error) {
	rowMaps := make([]map[string]any, len(records))
	pkCols := []string{
		"id",
	}
	for i, rec := range records {
		if err := validateCommentCreate(rec.Assignments); err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		input := assignmentsToCommentCreate(rec.Assignments)
		if q.Comment.beforeCreate != nil {
			if err := q.Comment.beforeCreate(ctx, &input); err != nil {
				return nil, err
			}
		}
		rowMaps[i] = input.ToRowMap()
	}
	results, err := executeCreateManyAndReturn(ctx, q, rowMaps, "Comment", CommentColOrder, selects, omits,
		q.selectCommentCols,
		q.loadCommentRelations,
		(*Comment).ScanFields,
		(*CommentSelect).hasAnyRelation,
		pkCols,
		conflictTarget,
		conflictAction,
	)
	if err != nil {
		return nil, err
	}
	if q.Comment.afterCreate != nil {
		if err := q.Comment.afterCreate(ctx, results); err != nil {
			return nil, err
		}
	}
	return results, nil
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
		client:     d.client,
		where:      where,
		additional: additional,
		execFunc:   d.client.executeCommentFindUnique,
	}
}

func (d *CommentDelegate) FindFirst(preds ...PredicateOf[Comment]) *FindFirstBuilder[Comment, CommentSelect, CommentOmit] {
	return &FindFirstBuilder[Comment, CommentSelect, CommentOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeCommentFindFirst,
	}
}

func (d *CommentDelegate) FindMany(preds ...PredicateOf[Comment]) *FindManyBuilder[Comment, CommentSelect, CommentOmit] {
	return &FindManyBuilder[Comment, CommentSelect, CommentOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeCommentFindMany,
	}
}

func (q *Queries) executeCommentFindUnique(ctx context.Context, where UniquePredicate[Comment], additional []PredicateOf[Comment], selects *CommentSelect, omits *CommentOmit) (*Comment, error) {
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
	allPreds := append([]PredicateOf[Comment]{where}, additional...)
	whereClause, vals := CompilePredicates(q.dialect, allPreds)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	returningCols := q.selectCommentCols(selects, omits)
	return executeSingleWithRelations(ctx, q, "Comment", whereClause, vals, returningCols,
		func(res *Comment, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*Comment) error {
			return txQ.loadCommentRelations(ctx, results, selects)
		},
		nil,
	)
}

func (q *Queries) executeCommentFindFirst(
	ctx context.Context,
	params QueryParams[Comment],
	selects *CommentSelect,
	omits *CommentOmit,
) (*Comment, error) {
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
	returningCols := q.selectCommentCols(selects, omits)
	return executeSingleWithRelations(ctx, q, "Comment", whereClause, vals, returningCols,
		func(res *Comment, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*Comment) error {
			return txQ.loadCommentRelations(ctx, results, selects)
		},
		params.Skip,
	)
}

func (q *Queries) executeCommentFindMany(
	ctx context.Context,
	params QueryParams[Comment],
	selects *CommentSelect,
	omits *CommentOmit,
) ([]*Comment, error) {
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
	returningCols := q.selectCommentCols(selects, omits)
	return executeManyWithRelations(ctx, q, "Comment", whereClause, vals, returningCols,
		func(res *Comment, cols []string) []any { return res.ScanFields(cols) },
		selects.hasAnyRelation(),
		func(ctx context.Context, txQ *Queries, results []*Comment) error {
			return txQ.loadCommentRelations(ctx, results, selects)
		},
		params.Take,
		params.Skip,
	)
}
func (q *Queries) loadCommentRelations(ctx context.Context, records []*Comment, selects *CommentSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Post != nil {
		relationSelects, relationOmits, relationParams := selects.Post.GetRelationParams()
		returningCols := q.selectPostCols(relationSelects, relationOmits, "id")
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
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading post: %w", err)
		}
		if err := q.loadPostRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Author != nil {
		relationSelects, relationOmits, relationParams := selects.Author.GetRelationParams()
		returningCols := q.selectUserCols(relationSelects, relationOmits, "id")
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
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading author: %w", err)
		}
		if err := q.loadUserRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
