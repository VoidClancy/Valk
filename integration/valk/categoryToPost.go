package valk

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
)

// CategoryToPost represents the database model
type CategoryToPost struct {
	PostId     string    `db:"postId" json:"postId"`
	CategoryId int32     `db:"categoryId" json:"categoryId"`
	Post       *Post     `json:"post,omitempty"`
	Category   *Category `json:"category,omitempty"`
}

// CategoryToPostCreate is used for hooks only — the Create API uses FieldAssignment
type CategoryToPostCreate struct {
	PostId     string `json:"postId"`
	CategoryId int32  `json:"categoryId"`
}

// colMask returns a bit mask of columns that are set
func (s *CategoryToPostCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	mask |= 1 << 1
	return mask
}

type CategoryToPostSelect struct {
	PostId     bool            `json:"postId"`
	CategoryId bool            `json:"categoryId"`
	Post       *PostSelect     `json:"post,omitempty"`
	Category   *CategorySelect `json:"category,omitempty"`
}

type CategoryToPostOmit struct {
	PostId     bool `json:"postId"`
	CategoryId bool `json:"categoryId"`
}

type CategoryToPostSelectQuery interface {
	GetRelationParams() (*CategoryToPostSelect, *CategoryToPostOmit, QueryParams[CategoryToPost])
}

func (s *CategoryToPostSelect) GetRelationParams() (*CategoryToPostSelect, *CategoryToPostOmit, QueryParams[CategoryToPost]) {
	return s, nil, QueryParams[CategoryToPost]{}
}

type CategoryToPostQueryBuilder struct {
	selects *CategoryToPostSelect
	omits   *CategoryToPostOmit
	where   []PredicateOf[CategoryToPost]
	take    *int
	skip    *int
	orderBy []OrderBy[CategoryToPost]
}

func (b *CategoryToPostQueryBuilder) Where(preds ...PredicateOf[CategoryToPost]) *CategoryToPostQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *CategoryToPostQueryBuilder) Take(limit int) *CategoryToPostQueryBuilder {
	b.take = &limit
	return b
}

func (b *CategoryToPostQueryBuilder) Skip(offset int) *CategoryToPostQueryBuilder {
	b.skip = &offset
	return b
}

func (b *CategoryToPostQueryBuilder) OrderBy(orders ...OrderBy[CategoryToPost]) *CategoryToPostQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *CategoryToPostQueryBuilder) Select(s CategoryToPostSelect) *CategoryToPostQueryBuilder {
	b.selects = &s
	return b
}

func (b *CategoryToPostQueryBuilder) Omit(o CategoryToPostOmit) *CategoryToPostQueryBuilder {
	b.omits = &o
	return b
}

func (b *CategoryToPostQueryBuilder) GetRelationParams() (*CategoryToPostSelect, *CategoryToPostOmit, QueryParams[CategoryToPost]) {
	if b == nil {
		return nil, nil, QueryParams[CategoryToPost]{}
	}
	return b.selects, b.omits, QueryParams[CategoryToPost]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
	}
}

type CategoryToPostCreateQuery = func(ctx context.Context, args *CategoryToPostCreate) (*CategoryToPost, error)
type CategoryToPostCreateManyQuery = func(ctx context.Context, args []*CategoryToPostCreate) (int64, error)
type CategoryToPostCreateManyAndReturnQuery = func(ctx context.Context, args []*CategoryToPostCreate) ([]*CategoryToPost, error)
type CategoryToPostFindUniqueQuery = func(ctx context.Context, where UniquePredicate[CategoryToPost], additional []PredicateOf[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error)
type CategoryToPostFindFirstQuery = func(ctx context.Context, params QueryParams[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error)
type CategoryToPostFindManyQuery = func(ctx context.Context, params QueryParams[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit) ([]*CategoryToPost, error)
type CategoryToPostDeleteManyQuery = func(ctx context.Context, preds []PredicateOf[CategoryToPost]) (int64, error)
type CategoryToPostDeleteQuery = func(ctx context.Context, where UniquePredicate[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error)
type CategoryToPostCountQuery = func(ctx context.Context, params QueryParams[CategoryToPost]) (int64, error)

type CategoryToPostExtension struct {
	Create              func(ctx context.Context, input *CategoryToPostCreate, next CategoryToPostCreateQuery) (*CategoryToPost, error)
	CreateMany          func(ctx context.Context, inputs []*CategoryToPostCreate, next CategoryToPostCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, inputs []*CategoryToPostCreate, next CategoryToPostCreateManyAndReturnQuery) ([]*CategoryToPost, error)
	FindUnique          func(ctx context.Context, where UniquePredicate[CategoryToPost], additional []PredicateOf[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit, next CategoryToPostFindUniqueQuery) (*CategoryToPost, error)
	FindFirst           func(ctx context.Context, params QueryParams[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit, next CategoryToPostFindFirstQuery) (*CategoryToPost, error)
	FindMany            func(ctx context.Context, params QueryParams[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit, next CategoryToPostFindManyQuery) ([]*CategoryToPost, error)
	DeleteMany          func(ctx context.Context, preds []PredicateOf[CategoryToPost], next CategoryToPostDeleteManyQuery) (int64, error)
	Delete              func(ctx context.Context, where UniquePredicate[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit, next CategoryToPostDeleteQuery) (*CategoryToPost, error)
	Count               func(ctx context.Context, params QueryParams[CategoryToPost], next CategoryToPostCountQuery) (int64, error)
}

type CategoryToPostDelegate struct {
	client     *Queries
	extensions []CategoryToPostExtension
}

func (d *CategoryToPostDelegate) Use(exts ...CategoryToPostExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *CategoryToPost) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "postId":
			targets[i] = &m.PostId
		case "categoryId":
			targets[i] = &m.CategoryId
		}
	}
	return targets
}

var categoryToPostDefaultCols = []string{
	"postId",
	"categoryId",
}

var categoryToPostPKCols = []string{
	"postId",
	"categoryId",
}

func selectCategoryToPostCols(selects *CategoryToPostSelect, omits *CategoryToPostOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return categoryToPostDefaultCols
	}

	anySelected := selects != nil && (selects.PostId || selects.CategoryId || selects.Post != nil || selects.Category != nil)

	specs := []colSpec{
		{"postId", selects != nil && selects.PostId, omits != nil && omits.PostId, selects != nil && selects.Post != nil},
		{"categoryId", selects != nil && selects.CategoryId, omits != nil && omits.CategoryId, selects != nil && selects.Category != nil},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *CategoryToPostSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Post != nil || s.Category != nil
}

type CategoryToPostCreateBuilder struct {
	*CreateBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]
}

func (b *CategoryToPostCreateBuilder) OnConflict(target UniqueConstraintTarget) *CategoryToPostConflictBuilder[CategoryToPostCreateBuilder] {
	return &CategoryToPostConflictBuilder[CategoryToPostCreateBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *CategoryToPostCreateBuilder) SetPostId(v string) *CategoryToPostCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "postId", Val: v})
	return b
}
func (b *CategoryToPostCreateBuilder) SetCategoryId(v int32) *CategoryToPostCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "categoryId", Val: v})
	return b
}

func (d *CategoryToPostDelegate) Create(assignments ...FieldAssignment) *CategoryToPostCreateBuilder {
	return &CategoryToPostCreateBuilder{
		CreateBuilder: &CreateBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
			assignments: assignments,
			execFunc:    d.executeCreate,
		},
	}
}

const (
	providedCategoryToPostPostId     uint64 = 1 << 0
	providedCategoryToPostCategoryId uint64 = 1 << 1
)

func assignmentsToCategoryToPostCreate(assignments []FieldAssignment) (CategoryToPostCreate, error) {
	var input CategoryToPostCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "postId":
			provided |= providedCategoryToPostPostId
			if v, ok := a.Val.(string); ok {
				input.PostId = v
				ValidateString(&errs, "postId", v, true, 0, false, false)
			} else {
				errs.Add("postId", a.Val, "type", "field postId must be of type string")
			}
		case "categoryId":
			provided |= providedCategoryToPostCategoryId
			if v, ok := a.Val.(int32); ok {
				input.CategoryId = v
				ValidateInt32(&errs, "categoryId", v, "")
			} else {
				errs.Add("categoryId", a.Val, "type", "field categoryId must be of type int32")
			}
		}
	}
	if provided&providedCategoryToPostPostId == 0 {
		errs.Add("postId", "", "required", "field PostId is required")
	}
	if provided&providedCategoryToPostCategoryId == 0 {
		errs.Add("categoryId", nil, "required", "field CategoryId is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *CategoryToPostCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 2)
	vals = make([]any, 0, 2)
	cols = append(cols, "postId")
	vals = append(vals, s.PostId)
	cols = append(cols, "categoryId")
	vals = append(vals, s.CategoryId)
	return
}

func partitionCategoryToPostInputs(dialect Dialect, inputs []*CategoryToPostCreate) [][]*CategoryToPostCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*CategoryToPostCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*CategoryToPostCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*CategoryToPostCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*CategoryToPostCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*CategoryToPostCreate{inputs}
}

func (d *CategoryToPostDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *CategoryToPostSelect, omits *CategoryToPostOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*CategoryToPost, error) {
	input, err := assignmentsToCategoryToPostCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectCategoryToPostCols(selects, omits)

	if len(d.extensions) == 0 {
		hasRelations := selects.hasAnyRelation()
		if hasRelations {
			var res *CategoryToPost
			err = d.client.transaction(ctx, func(txQ *Queries) error {
				var err error
				res, err = txQ.CategoryToPost.runCreate(ctx, cols, vals, returningCols, categoryToPostPKCols, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.CategoryToPost.loadRelations(ctx, []*CategoryToPost{res}, selects)
			})
			return res, err
		}
		return d.runCreate(ctx, cols, vals, returningCols, categoryToPostPKCols, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args *CategoryToPostCreate) (*CategoryToPost, error) {
		cols, vals := args.ToColsVals()
		returningCols := selectCategoryToPostCols(selects, omits)

		hasRelations := selects.hasAnyRelation()
		var res *CategoryToPost
		var err error
		if hasRelations {
			err = d.client.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = txQ.CategoryToPost.runCreate(c, cols, vals, returningCols, categoryToPostPKCols, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.CategoryToPost.loadRelations(c, []*CategoryToPost{res}, selects)
			})
		} else {
			res, err = d.runCreate(c, cols, vals, returningCols, categoryToPostPKCols, conflictTarget, conflictAction)
		}
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Create != nil {
			next, hook := curr, ext.Create
			curr = func(c context.Context, input *CategoryToPostCreate) (*CategoryToPost, error) {
				return hook(c, input, next)
			}
		}
	}

	return curr(ctx, &input)
}

type CategoryToPostCreateManyBuilder struct {
	*CreateManyBuilder[CategoryToPost]
}

func (b *CategoryToPostCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *CategoryToPostConflictBuilder[CategoryToPostCreateManyBuilder] {
	return &CategoryToPostConflictBuilder[CategoryToPostCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type CategoryToPostCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]
}

func (b *CategoryToPostCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *CategoryToPostConflictBuilder[CategoryToPostCreateManyAndReturnBuilder] {
	return &CategoryToPostConflictBuilder[CategoryToPostCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (d *CategoryToPostDelegate) CreateMany(builders ...*CategoryToPostCreateBuilder) *CategoryToPostCreateManyBuilder {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &CategoryToPostCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[CategoryToPost]{
			records:  records,
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *CategoryToPostDelegate) CreateManyAndReturn(builders ...*CategoryToPostCreateBuilder) *CategoryToPostCreateManyAndReturnBuilder {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &CategoryToPostCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
			records:  records,
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func (d *CategoryToPostDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs := make([]*CategoryToPostCreate, len(records))
	for i, rec := range records {
		input, err := assignmentsToCategoryToPostCreate(rec.Assignments)
		if err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &input
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args []*CategoryToPostCreate) (int64, error) {
		return d.runCreateMany(c, args, conflictTarget, conflictAction)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.CreateMany != nil {
			next, hook := curr, ext.CreateMany
			curr = func(c context.Context, inputs []*CategoryToPostCreate) (int64, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
}

func (d *CategoryToPostDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *CategoryToPostSelect, omits *CategoryToPostOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*CategoryToPost, error) {
	inputs := make([]*CategoryToPostCreate, len(records))
	for i, rec := range records {
		input, err := assignmentsToCategoryToPostCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &input
	}

	if len(d.extensions) == 0 {
		hasRelations := selects != nil && selects.hasAnyRelation()
		if hasRelations {
			var res []*CategoryToPost
			err := d.client.transaction(ctx, func(txQ *Queries) error {
				var err error
				res, err = txQ.CategoryToPost.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.CategoryToPost.loadRelations(ctx, res, selects)
			})
			return res, err
		}
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args []*CategoryToPostCreate) ([]*CategoryToPost, error) {
		hasRelations := selects != nil && selects.hasAnyRelation()
		if hasRelations {
			var res []*CategoryToPost
			err := d.client.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = txQ.CategoryToPost.runCreateManyAndReturn(c, args, selects, omits, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.CategoryToPost.loadRelations(c, res, selects)
			})
			return res, err
		}
		return d.runCreateManyAndReturn(c, args, selects, omits, conflictTarget, conflictAction)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.CreateManyAndReturn != nil {
			next, hook := curr, ext.CreateManyAndReturn
			curr = func(c context.Context, inputs []*CategoryToPostCreate) ([]*CategoryToPost, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
}

func (d *CategoryToPostDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	pkCols []string,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*CategoryToPost, error) {
	query, clauseArgs := buildSingleInsertSQL(d.client, "CategoryToPost", cols, returningCols, pkCols, conflictTarget, conflictAction, len(vals))
	if len(clauseArgs) > 0 {
		vals = append(vals, clauseArgs...)
	}

	var res CategoryToPost
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

func (d *CategoryToPostDelegate) runCreateFallback(ctx context.Context, query string, vals []any, cols []string, returningCols []string, pkCols []string) (*CategoryToPost, error) {
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
	selectSb.Grow(64 + len(returningCols)*15 + len("CategoryToPost") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "CategoryToPost")
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

	var res CategoryToPost
	if rows.Next() {
		if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
			return nil, err
		}
		return &res, nil
	}
	return nil, rows.Err()
}

func (d *CategoryToPostDelegate) buildBulkInsertSQL(q *Queries, batch []*CategoryToPostCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 2)
	for i, c := range categoryToPostDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "CategoryToPost")
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
			case "postId":
				vals = append(vals, input.PostId)
			case "categoryId":
				vals = append(vals, input.CategoryId)
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

func (d *CategoryToPostDelegate) runCreateMany(ctx context.Context, inputs []*CategoryToPostCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionCategoryToPostInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)

		var conflictCols []string
		if conflictTarget != nil {
			conflictCols = conflictTarget.UniqueColumns()
		}
		var nonConflictCols []string
		if conflictAction != nil && conflictAction.Type == ConflictActionUpdateNewValues {
			nonConflictCols = computeNonConflictCols(cols, conflictCols, categoryToPostPKCols)
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

func (d *CategoryToPostDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*CategoryToPostCreate,
	selects *CategoryToPostSelect,
	omits *CategoryToPostOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*CategoryToPost, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionCategoryToPostInputs(d.client.dialect, inputs)
	returningCols := selectCategoryToPostCols(selects, omits)
	hasRelations := selects != nil && selects.hasAnyRelation()

	recordsOut := make([]*CategoryToPost, 0, len(inputs))

	runBatch := func(txQ *Queries, batch []*CategoryToPostCreate) error {
		cols, vals, queryStr := d.buildBulkInsertSQL(txQ, batch, 1)

		var conflictCols []string
		if conflictTarget != nil {
			conflictCols = conflictTarget.UniqueColumns()
		}
		var nonConflictCols []string
		if conflictAction != nil && conflictAction.Type == ConflictActionUpdateNewValues {
			nonConflictCols = computeNonConflictCols(cols, conflictCols, categoryToPostPKCols)
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
				var res CategoryToPost
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
		selectSb.Grow(64 + len(returningCols)*15 + len("CategoryToPost") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			txQ.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		txQ.dialect.WriteQuotedIdent(&selectSb, "CategoryToPost")
		selectSb.WriteString(" WHERE ")
		txQ.dialect.WriteQuotedIdent(&selectSb, categoryToPostPKCols[0])
		selectSb.WriteString(" >= ")
		txQ.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		txQ.dialect.WriteQuotedIdent(&selectSb, categoryToPostPKCols[0])
		selectSb.WriteString(" < ")
		txQ.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := txQ.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var res CategoryToPost
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
				return txQ.CategoryToPost.loadRelations(ctx, recordsOut, selects)
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

type CategoryToPostConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *CategoryToPostConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *CategoryToPostConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *CategoryToPostConflictBuilder[B]) Update(fn func(u *CategoryToPostUpsert)) *B {
	var up ConflictUpdate
	u := newCategoryToPostUpsert(&up)
	fn(u)
	cb.setAction(ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}, cb.conflictTarget)
	return cb.builder
}

type CategoryToPostUpsert struct {
	PostId     fieldUpsert[string]
	CategoryId numericFieldUpsert[int32]
}

func newCategoryToPostUpsert(up *ConflictUpdate) *CategoryToPostUpsert {
	return &CategoryToPostUpsert{
		PostId: fieldUpsert[string]{column: "postId", update: up},
		CategoryId: numericFieldUpsert[int32]{
			fieldUpsert: fieldUpsert[int32]{column: "categoryId", update: up},
			tableName:   "CategoryToPost",
		},
	}
}
func (d *CategoryToPostDelegate) FindUnique(where UniquePredicate[CategoryToPost], additional ...PredicateOf[CategoryToPost]) *FindUniqueBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit] {
	return &FindUniqueBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *CategoryToPostDelegate) FindFirst(preds ...PredicateOf[CategoryToPost]) *FindFirstBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit] {
	return &FindFirstBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *CategoryToPostDelegate) FindMany(preds ...PredicateOf[CategoryToPost]) *FindManyBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit] {
	return &FindManyBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *CategoryToPostDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[CategoryToPost], additional []PredicateOf[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error) {
	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, where, additional, selects, omits)
	}

	curr := func(c context.Context, w UniquePredicate[CategoryToPost], add []PredicateOf[CategoryToPost], sel *CategoryToPostSelect, o *CategoryToPostOmit) (*CategoryToPost, error) {
		return d.runFindUnique(c, w, add, sel, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindUnique != nil {
			next, hook := curr, ext.FindUnique
			curr = func(c context.Context, w UniquePredicate[CategoryToPost], add []PredicateOf[CategoryToPost], sel *CategoryToPostSelect, o *CategoryToPostOmit) (*CategoryToPost, error) {
				return hook(c, w, add, sel, o, next)
			}
		}
	}

	return curr(ctx, where, additional, selects, omits)
}

func (d *CategoryToPostDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[CategoryToPost],
	selects *CategoryToPostSelect,
	omits *CategoryToPostOmit,
) (*CategoryToPost, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	curr := func(c context.Context, p QueryParams[CategoryToPost], sel *CategoryToPostSelect, o *CategoryToPostOmit) (*CategoryToPost, error) {
		return d.runFindFirst(c, p, sel, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindFirst != nil {
			next, hook := curr, ext.FindFirst
			curr = func(c context.Context, p QueryParams[CategoryToPost], sel *CategoryToPostSelect, o *CategoryToPostOmit) (*CategoryToPost, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}

func (d *CategoryToPostDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[CategoryToPost],
	selects *CategoryToPostSelect,
	omits *CategoryToPostOmit,
) ([]*CategoryToPost, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	curr := func(c context.Context, p QueryParams[CategoryToPost], sel *CategoryToPostSelect, o *CategoryToPostOmit) ([]*CategoryToPost, error) {
		return d.runFindMany(c, p, sel, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindMany != nil {
			next, hook := curr, ext.FindMany
			curr = func(c context.Context, p QueryParams[CategoryToPost], sel *CategoryToPostSelect, o *CategoryToPostOmit) ([]*CategoryToPost, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}

func (d *CategoryToPostDelegate) runFindUnique(ctx context.Context, where UniquePredicate[CategoryToPost], additional []PredicateOf[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error) {
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
	allPreds := append([]PredicateOf[CategoryToPost]{where}, additional...)
	whereClause, vals := CompilePredicates(d.client.dialect, allPreds)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	returningCols := selectCategoryToPostCols(selects, omits)

	var res *CategoryToPost
	var err error
	if selects.hasAnyRelation() {
		err = d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.CategoryToPost.queryOne(ctx, whereClause, vals, returningCols, nil)
			if err != nil {
				return err
			}
			if res == nil {
				return nil
			}
			return txQ.CategoryToPost.loadRelations(ctx, []*CategoryToPost{res}, selects)
		})
	} else {
		res, err = d.queryOne(ctx, whereClause, vals, returningCols, nil)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *CategoryToPostDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[CategoryToPost],
	selects *CategoryToPostSelect,
	omits *CategoryToPostOmit,
) (*CategoryToPost, error) {
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
	returningCols := selectCategoryToPostCols(selects, omits)

	var res *CategoryToPost
	var err error
	if selects.hasAnyRelation() {
		err = d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.CategoryToPost.queryOne(ctx, whereClause, vals, returningCols, params.Skip)
			if err != nil {
				return err
			}
			if res == nil {
				return nil
			}
			return txQ.CategoryToPost.loadRelations(ctx, []*CategoryToPost{res}, selects)
		})
	} else {
		res, err = d.queryOne(ctx, whereClause, vals, returningCols, params.Skip)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *CategoryToPostDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[CategoryToPost],
	selects *CategoryToPostSelect,
	omits *CategoryToPostOmit,
) ([]*CategoryToPost, error) {
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
	returningCols := selectCategoryToPostCols(selects, omits)

	var results []*CategoryToPost
	var err error
	if selects.hasAnyRelation() {
		err = d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			results, err = txQ.CategoryToPost.queryMany(ctx, whereClause, vals, returningCols, params.Take, params.Skip)
			if err != nil {
				return err
			}
			if len(results) == 0 {
				return nil
			}
			return txQ.CategoryToPost.loadRelations(ctx, results, selects)
		})
	} else {
		results, err = d.queryMany(ctx, whereClause, vals, returningCols, params.Take, params.Skip)
	}
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (d *CategoryToPostDelegate) queryOne(ctx context.Context, whereClause string, whereVals []any, returningCols []string, skip *int) (*CategoryToPost, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "CategoryToPost", returningCols, whereClause, &limitOne, skip)
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

	var res CategoryToPost
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (d *CategoryToPostDelegate) queryMany(ctx context.Context, whereClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*CategoryToPost, error) {
	query := buildSelectSQL(d.client, "CategoryToPost", returningCols, whereClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*CategoryToPost, 0)
	for rows.Next() {
		var res CategoryToPost
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
func (d *CategoryToPostDelegate) DeleteMany(preds ...PredicateOf[CategoryToPost]) *DeleteManyBuilder[CategoryToPost] {
	return &DeleteManyBuilder[CategoryToPost]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *CategoryToPostDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[CategoryToPost]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	curr := func(c context.Context, p []PredicateOf[CategoryToPost]) (int64, error) {
		return d.runDeleteMany(c, p)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.DeleteMany != nil {
			next, hook := curr, ext.DeleteMany
			curr = func(c context.Context, p []PredicateOf[CategoryToPost]) (int64, error) {
				return hook(c, p, next)
			}
		}
	}

	return curr(ctx, preds)
}

func (d *CategoryToPostDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[CategoryToPost]) (int64, error) {
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
	d.client.dialect.WriteQuotedIdent(&sb, "CategoryToPost")
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

func (d *CategoryToPostDelegate) Delete(where UniquePredicate[CategoryToPost]) *DeleteBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit] {
	return &DeleteBuilder[CategoryToPost, CategoryToPostSelect, CategoryToPostOmit]{
		where:    where,
		execFunc: d.executeDelete,
	}
}

func (d *CategoryToPostDelegate) executeDelete(ctx context.Context, where UniquePredicate[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error) {
	if len(d.extensions) == 0 {
		return d.runDelete(ctx, where, selects, omits)
	}

	curr := func(c context.Context, w UniquePredicate[CategoryToPost], s *CategoryToPostSelect, o *CategoryToPostOmit) (*CategoryToPost, error) {
		return d.runDelete(c, w, s, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Delete != nil {
			next, hook := curr, ext.Delete
			curr = func(c context.Context, w UniquePredicate[CategoryToPost], s *CategoryToPostSelect, o *CategoryToPostOmit) (*CategoryToPost, error) {
				return hook(c, w, s, o, next)
			}
		}
	}

	return curr(ctx, where, selects, omits)
}

func (d *CategoryToPostDelegate) runDelete(ctx context.Context, where UniquePredicate[CategoryToPost], selects *CategoryToPostSelect, omits *CategoryToPostOmit) (*CategoryToPost, error) {
	if err := where.Validate(); err != nil {
		return nil, err
	}

	returningCols := selectCategoryToPostCols(selects, omits, categoryToPostPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *CategoryToPost
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.CategoryToPost.executeFindUnique(ctx, where, nil, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return sql.ErrNoRows
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "CategoryToPost")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[CategoryToPost]
			pkPreds = append(pkPreds, Predicate[CategoryToPost]{
				Data: PredicateData{
					Column:   "postId",
					Operator: "=",
					Value:    res.PostId,
				},
			})
			pkPreds = append(pkPreds, Predicate[CategoryToPost]{
				Data: PredicateData{
					Column:   "categoryId",
					Operator: "=",
					Value:    res.CategoryId,
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
	d.client.dialect.WriteQuotedIdent(&sb, "CategoryToPost")

	whereClause, vals := CompilePredicates(d.client.dialect, []PredicateOf[CategoryToPost]{where})
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

	var row CategoryToPost
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, err
	}
	return &row, nil
}
func (d *CategoryToPostDelegate) Count(preds ...PredicateOf[CategoryToPost]) *CountBuilder[CategoryToPost] {
	return &CountBuilder[CategoryToPost]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *CategoryToPostDelegate) executeCount(ctx context.Context, params QueryParams[CategoryToPost]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	curr := func(c context.Context, p QueryParams[CategoryToPost]) (int64, error) {
		return d.runCount(c, p)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Count != nil {
			next, hook := curr, ext.Count
			curr = func(c context.Context, p QueryParams[CategoryToPost]) (int64, error) {
				return hook(c, p, next)
			}
		}
	}

	return curr(ctx, params)
}

func (d *CategoryToPostDelegate) runCount(ctx context.Context, params QueryParams[CategoryToPost]) (int64, error) {
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
		d.client.dialect.WriteQuotedIdent(&subQuery, "CategoryToPost")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "CategoryToPost")
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
func (d *CategoryToPostDelegate) loadRelations(ctx context.Context, records []*CategoryToPost, selects *CategoryToPostSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Post != nil {
		relationSelects, relationOmits, relationParams := selects.Post.GetRelationParams()
		returningCols := selectPostCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: CategoryToPost.postId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *CategoryToPost) string { return p.PostId }),
			"Post",
			"id",
			returningCols,
			scanInto(returningCols, (*Post).ScanFields),
			directKey(func(c *Post) string { return c.Id }),
			setOne(func(p *CategoryToPost, c *Post) { p.Post = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading post: %w", err)
		}
		if err := d.client.Post.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}
	if selects.Category != nil {
		relationSelects, relationOmits, relationParams := selects.Category.GetRelationParams()
		returningCols := selectCategoryCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: CategoryToPost.categoryId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *CategoryToPost) int32 { return p.CategoryId }),
			"Category",
			"id",
			returningCols,
			scanInto(returningCols, (*Category).ScanFields),
			directKey(func(c *Category) int32 { return c.Id }),
			setOne(func(p *CategoryToPost, c *Category) { p.Category = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading category: %w", err)
		}
		if err := d.client.Category.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
