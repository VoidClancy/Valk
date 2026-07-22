package valk

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
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

// colMask returns a bit mask of columns that are set
func (s *CommentCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	mask |= 1 << 1
	mask |= 1 << 2
	mask |= 1 << 3
	mask |= 1 << 4
	mask |= 1 << 5
	mask |= 1 << 6
	if s.Meta != nil {
		mask |= 1 << 7
	}
	return mask
}

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
type CommentDeleteManyQuery = func(ctx context.Context, preds []PredicateOf[Comment]) (int64, error)
type CommentDeleteQuery = func(ctx context.Context, where UniquePredicate[Comment], selects *CommentSelect, omits *CommentOmit) (*Comment, error)
type CommentCountQuery = func(ctx context.Context, params QueryParams[Comment]) (int64, error)

type CommentExtension struct {
	Create              func(ctx context.Context, input *CommentCreate, next CommentCreateQuery) (*Comment, error)
	CreateMany          func(ctx context.Context, inputs []*CommentCreate, next CommentCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, inputs []*CommentCreate, next CommentCreateManyAndReturnQuery) ([]*Comment, error)
	FindUnique          func(ctx context.Context, where UniquePredicate[Comment], additional []PredicateOf[Comment], selects *CommentSelect, omits *CommentOmit, next CommentFindUniqueQuery) (*Comment, error)
	FindFirst           func(ctx context.Context, params QueryParams[Comment], selects *CommentSelect, omits *CommentOmit, next CommentFindFirstQuery) (*Comment, error)
	FindMany            func(ctx context.Context, params QueryParams[Comment], selects *CommentSelect, omits *CommentOmit, next CommentFindManyQuery) ([]*Comment, error)
	DeleteMany          func(ctx context.Context, preds []PredicateOf[Comment], next CommentDeleteManyQuery) (int64, error)
	Delete              func(ctx context.Context, where UniquePredicate[Comment], selects *CommentSelect, omits *CommentOmit, next CommentDeleteQuery) (*Comment, error)
	Count               func(ctx context.Context, params QueryParams[Comment], next CommentCountQuery) (int64, error)
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

var commentPKCols = []string{
	"id",
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

func partitionCommentInputs(dialect Dialect, inputs []*CommentCreate) [][]*CommentCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*CommentCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*CommentCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*CommentCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*CommentCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*CommentCreate{inputs}
}

func (d *CommentDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *CommentSelect, omits *CommentOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*Comment, error) {
	input, err := assignmentsToCommentCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectCommentCols(selects, omits)

	if len(d.extensions) == 0 {
		hasRelations := selects.hasAnyRelation()
		if hasRelations {
			var res *Comment
			err = d.client.transaction(ctx, func(txQ *Queries) error {
				var err error
				res, err = txQ.Comment.runCreate(ctx, cols, vals, returningCols, commentPKCols, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.Comment.loadRelations(ctx, []*Comment{res}, selects)
			})
			return res, err
		}
		return d.runCreate(ctx, cols, vals, returningCols, commentPKCols, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args *CommentCreate) (*Comment, error) {
		cols, vals := args.ToColsVals()
		returningCols := selectCommentCols(selects, omits)

		hasRelations := selects.hasAnyRelation()
		var res *Comment
		var err error
		if hasRelations {
			err = d.client.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = txQ.Comment.runCreate(c, cols, vals, returningCols, commentPKCols, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.Comment.loadRelations(c, []*Comment{res}, selects)
			})
		} else {
			res, err = d.runCreate(c, cols, vals, returningCols, commentPKCols, conflictTarget, conflictAction)
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

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args []*CommentCreate) (int64, error) {
		return d.runCreateMany(c, args, conflictTarget, conflictAction)
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

	if len(d.extensions) == 0 {
		hasRelations := selects != nil && selects.hasAnyRelation()
		if hasRelations {
			var res []*Comment
			err := d.client.transaction(ctx, func(txQ *Queries) error {
				var err error
				res, err = txQ.Comment.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.Comment.loadRelations(ctx, res, selects)
			})
			return res, err
		}
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args []*CommentCreate) ([]*Comment, error) {
		hasRelations := selects != nil && selects.hasAnyRelation()
		if hasRelations {
			var res []*Comment
			err := d.client.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = txQ.Comment.runCreateManyAndReturn(c, args, selects, omits, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.Comment.loadRelations(c, res, selects)
			})
			return res, err
		}
		return d.runCreateManyAndReturn(c, args, selects, omits, conflictTarget, conflictAction)
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

func (d *CommentDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	pkCols []string,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*Comment, error) {
	query, clauseArgs := buildSingleInsertSQL(d.client, "Comment", cols, returningCols, pkCols, conflictTarget, conflictAction, len(vals))
	if len(clauseArgs) > 0 {
		vals = append(vals, clauseArgs...)
	}

	var res Comment
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

func (d *CommentDelegate) runCreateFallback(ctx context.Context, query string, vals []any, cols []string, returningCols []string, pkCols []string) (*Comment, error) {
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
	selectSb.Grow(64 + len(returningCols)*15 + len("Comment") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "Comment")
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

	var res Comment
	if rows.Next() {
		if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
			return nil, err
		}
		return &res, nil
	}
	return nil, rows.Err()
}

func (d *CommentDelegate) buildBulkInsertSQL(q *Queries, batch []*CommentCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 8)
	for i, c := range commentDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "Comment")
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
			case "textify":
				vals = append(vals, input.Textify)
			case "dummy3":
				vals = append(vals, input.Dummy3)
			case "dummy1":
				vals = append(vals, input.Dummy1)
			case "dummy2":
				vals = append(vals, input.Dummy2)
			case "postId":
				vals = append(vals, input.PostId)
			case "authorId":
				vals = append(vals, input.AuthorId)
			case "meta":
				if input.Meta != nil {
					vals = append(vals, *input.Meta)
				} else {
					writeDefault = true
				}
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

func (d *CommentDelegate) runCreateMany(ctx context.Context, inputs []*CommentCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionCommentInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)

		var conflictCols []string
		if conflictTarget != nil {
			conflictCols = conflictTarget.UniqueColumns()
		}
		var nonConflictCols []string
		if conflictAction != nil && conflictAction.Type == ConflictActionUpdateNewValues {
			nonConflictCols = computeNonConflictCols(cols, conflictCols, commentPKCols)
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

func (d *CommentDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*CommentCreate,
	selects *CommentSelect,
	omits *CommentOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*Comment, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionCommentInputs(d.client.dialect, inputs)
	returningCols := selectCommentCols(selects, omits)
	hasRelations := selects != nil && selects.hasAnyRelation()

	recordsOut := make([]*Comment, 0, len(inputs))

	runBatch := func(txQ *Queries, batch []*CommentCreate) error {
		cols, vals, queryStr := d.buildBulkInsertSQL(txQ, batch, 1)

		var conflictCols []string
		if conflictTarget != nil {
			conflictCols = conflictTarget.UniqueColumns()
		}
		var nonConflictCols []string
		if conflictAction != nil && conflictAction.Type == ConflictActionUpdateNewValues {
			nonConflictCols = computeNonConflictCols(cols, conflictCols, commentPKCols)
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
				var res Comment
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
		selectSb.Grow(64 + len(returningCols)*15 + len("Comment") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			txQ.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		txQ.dialect.WriteQuotedIdent(&selectSb, "Comment")
		selectSb.WriteString(" WHERE ")
		txQ.dialect.WriteQuotedIdent(&selectSb, commentPKCols[0])
		selectSb.WriteString(" >= ")
		txQ.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		txQ.dialect.WriteQuotedIdent(&selectSb, commentPKCols[0])
		selectSb.WriteString(" < ")
		txQ.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := txQ.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var res Comment
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
				return txQ.Comment.loadRelations(ctx, recordsOut, selects)
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
	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, where, additional, selects, omits)
	}

	curr := func(c context.Context, w UniquePredicate[Comment], add []PredicateOf[Comment], sel *CommentSelect, o *CommentOmit) (*Comment, error) {
		return d.runFindUnique(c, w, add, sel, o)
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
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	curr := func(c context.Context, p QueryParams[Comment], sel *CommentSelect, o *CommentOmit) (*Comment, error) {
		return d.runFindFirst(c, p, sel, o)
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
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	curr := func(c context.Context, p QueryParams[Comment], sel *CommentSelect, o *CommentOmit) ([]*Comment, error) {
		return d.runFindMany(c, p, sel, o)
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

func (d *CommentDelegate) runFindUnique(ctx context.Context, where UniquePredicate[Comment], additional []PredicateOf[Comment], selects *CommentSelect, omits *CommentOmit) (*Comment, error) {
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
	whereClause, vals := CompilePredicates(d.client.dialect, allPreds)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	returningCols := selectCommentCols(selects, omits)

	var res *Comment
	var err error
	if selects.hasAnyRelation() {
		err = d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Comment.queryOne(ctx, whereClause, vals, returningCols, nil)
			if err != nil {
				return err
			}
			if res == nil {
				return nil
			}
			return txQ.Comment.loadRelations(ctx, []*Comment{res}, selects)
		})
	} else {
		res, err = d.queryOne(ctx, whereClause, vals, returningCols, nil)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *CommentDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[Comment],
	selects *CommentSelect,
	omits *CommentOmit,
) (*Comment, error) {
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
	returningCols := selectCommentCols(selects, omits)

	var res *Comment
	var err error
	if selects.hasAnyRelation() {
		err = d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Comment.queryOne(ctx, whereClause, vals, returningCols, params.Skip)
			if err != nil {
				return err
			}
			if res == nil {
				return nil
			}
			return txQ.Comment.loadRelations(ctx, []*Comment{res}, selects)
		})
	} else {
		res, err = d.queryOne(ctx, whereClause, vals, returningCols, params.Skip)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *CommentDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[Comment],
	selects *CommentSelect,
	omits *CommentOmit,
) ([]*Comment, error) {
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
	returningCols := selectCommentCols(selects, omits)

	var results []*Comment
	var err error
	if selects.hasAnyRelation() {
		err = d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			results, err = txQ.Comment.queryMany(ctx, whereClause, vals, returningCols, params.Take, params.Skip)
			if err != nil {
				return err
			}
			if len(results) == 0 {
				return nil
			}
			return txQ.Comment.loadRelations(ctx, results, selects)
		})
	} else {
		results, err = d.queryMany(ctx, whereClause, vals, returningCols, params.Take, params.Skip)
	}
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (d *CommentDelegate) queryOne(ctx context.Context, whereClause string, whereVals []any, returningCols []string, skip *int) (*Comment, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "Comment", returningCols, whereClause, &limitOne, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return nil, nil
	}

	var res Comment
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		return nil, err
	}
	return &res, nil
}

func (d *CommentDelegate) queryMany(ctx context.Context, whereClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*Comment, error) {
	query := buildSelectSQL(d.client, "Comment", returningCols, whereClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*Comment, 0)
	for rows.Next() {
		var res Comment
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
func (d *CommentDelegate) DeleteMany(preds ...PredicateOf[Comment]) *DeleteManyBuilder[Comment] {
	return &DeleteManyBuilder[Comment]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *CommentDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[Comment]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	curr := func(c context.Context, p []PredicateOf[Comment]) (int64, error) {
		return d.runDeleteMany(c, p)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.DeleteMany != nil {
			next, hook := curr, ext.DeleteMany
			curr = func(c context.Context, p []PredicateOf[Comment]) (int64, error) {
				return hook(c, p, next)
			}
		}
	}

	return curr(ctx, preds)
}

func (d *CommentDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[Comment]) (int64, error) {
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
	d.client.dialect.WriteQuotedIdent(&sb, "Comment")
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

func (d *CommentDelegate) Delete(where UniquePredicate[Comment]) *DeleteBuilder[Comment, CommentSelect, CommentOmit] {
	return &DeleteBuilder[Comment, CommentSelect, CommentOmit]{
		where:    where,
		execFunc: d.executeDelete,
	}
}

func (d *CommentDelegate) executeDelete(ctx context.Context, where UniquePredicate[Comment], selects *CommentSelect, omits *CommentOmit) (*Comment, error) {
	if len(d.extensions) == 0 {
		return d.runDelete(ctx, where, selects, omits)
	}

	curr := func(c context.Context, w UniquePredicate[Comment], s *CommentSelect, o *CommentOmit) (*Comment, error) {
		return d.runDelete(c, w, s, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Delete != nil {
			next, hook := curr, ext.Delete
			curr = func(c context.Context, w UniquePredicate[Comment], s *CommentSelect, o *CommentOmit) (*Comment, error) {
				return hook(c, w, s, o, next)
			}
		}
	}

	return curr(ctx, where, selects, omits)
}

func (d *CommentDelegate) runDelete(ctx context.Context, where UniquePredicate[Comment], selects *CommentSelect, omits *CommentOmit) (*Comment, error) {
	if err := where.Validate(); err != nil {
		return nil, err
	}

	returningCols := selectCommentCols(selects, omits, commentPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *Comment
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Comment.executeFindUnique(ctx, where, nil, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return sql.ErrNoRows
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "Comment")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[Comment]
			pkPreds = append(pkPreds, Predicate[Comment]{
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
	d.client.dialect.WriteQuotedIdent(&sb, "Comment")

	whereClause, vals := CompilePredicates(d.client.dialect, []PredicateOf[Comment]{where})
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

	var row Comment
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, err
	}
	return &row, nil
}
func (d *CommentDelegate) Count(preds ...PredicateOf[Comment]) *CountBuilder[Comment] {
	return &CountBuilder[Comment]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *CommentDelegate) executeCount(ctx context.Context, params QueryParams[Comment]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	curr := func(c context.Context, p QueryParams[Comment]) (int64, error) {
		return d.runCount(c, p)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Count != nil {
			next, hook := curr, ext.Count
			curr = func(c context.Context, p QueryParams[Comment]) (int64, error) {
				return hook(c, p, next)
			}
		}
	}

	return curr(ctx, params)
}

func (d *CommentDelegate) runCount(ctx context.Context, params QueryParams[Comment]) (int64, error) {
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
		d.client.dialect.WriteQuotedIdent(&subQuery, "Comment")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "Comment")
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
