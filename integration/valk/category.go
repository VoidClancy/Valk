package valk

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
)

// Category represents the database model
type Category struct {
	Id    int32             `db:"id" json:"id"`
	Name  string            `db:"name" json:"name"`
	Posts []*CategoryToPost `json:"posts,omitempty"`
}

// CategoryCreate is used for hooks only — the Create API uses FieldAssignment
type CategoryCreate struct {
	Id   *int32 `json:"id"`
	Name string `json:"name"`
}

// colMask returns a bit mask of columns that are set
func (s *CategoryCreate) colMask() uint64 {
	var mask uint64
	if s.Id != nil {
		mask |= 1 << 0
	}
	mask |= 1 << 1
	return mask
}

type CategorySelect struct {
	Id    bool                      `json:"id"`
	Name  bool                      `json:"name"`
	Posts CategoryToPostSelectQuery `json:"posts,omitempty"`
}

type CategoryOmit struct {
	Id   bool `json:"id"`
	Name bool `json:"name"`
}

type CategorySelectQuery interface {
	GetRelationParams() (*CategorySelect, *CategoryOmit, QueryParams[Category])
}

func (s *CategorySelect) GetRelationParams() (*CategorySelect, *CategoryOmit, QueryParams[Category]) {
	return s, nil, QueryParams[Category]{}
}

type CategoryQueryBuilder struct {
	selects *CategorySelect
	omits   *CategoryOmit
	where   []PredicateOf[Category]
	take    *int
	skip    *int
	orderBy []OrderBy[Category]
}

func (b *CategoryQueryBuilder) Where(preds ...PredicateOf[Category]) *CategoryQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *CategoryQueryBuilder) Take(limit int) *CategoryQueryBuilder {
	b.take = &limit
	return b
}

func (b *CategoryQueryBuilder) Skip(offset int) *CategoryQueryBuilder {
	b.skip = &offset
	return b
}

func (b *CategoryQueryBuilder) OrderBy(orders ...OrderBy[Category]) *CategoryQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *CategoryQueryBuilder) Select(s CategorySelect) *CategoryQueryBuilder {
	b.selects = &s
	return b
}

func (b *CategoryQueryBuilder) Omit(o CategoryOmit) *CategoryQueryBuilder {
	b.omits = &o
	return b
}

func (b *CategoryQueryBuilder) GetRelationParams() (*CategorySelect, *CategoryOmit, QueryParams[Category]) {
	if b == nil {
		return nil, nil, QueryParams[Category]{}
	}
	return b.selects, b.omits, QueryParams[Category]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
	}
}

type CategoryCreateQuery = func(ctx context.Context, args *CategoryCreate) (*Category, error)
type CategoryCreateManyQuery = func(ctx context.Context, args []*CategoryCreate) (int64, error)
type CategoryCreateManyAndReturnQuery = func(ctx context.Context, args []*CategoryCreate) ([]*Category, error)
type CategoryFindUniqueQuery = func(ctx context.Context, where UniquePredicate[Category], additional []PredicateOf[Category], selects *CategorySelect, omits *CategoryOmit) (*Category, error)
type CategoryFindFirstQuery = func(ctx context.Context, params QueryParams[Category], selects *CategorySelect, omits *CategoryOmit) (*Category, error)
type CategoryFindManyQuery = func(ctx context.Context, params QueryParams[Category], selects *CategorySelect, omits *CategoryOmit) ([]*Category, error)
type CategoryDeleteManyQuery = func(ctx context.Context, preds []PredicateOf[Category]) (int64, error)
type CategoryDeleteQuery = func(ctx context.Context, where UniquePredicate[Category], selects *CategorySelect, omits *CategoryOmit) (*Category, error)
type CategoryCountQuery = func(ctx context.Context, params QueryParams[Category]) (int64, error)

type CategoryExtension struct {
	Create              func(ctx context.Context, input *CategoryCreate, next CategoryCreateQuery) (*Category, error)
	CreateMany          func(ctx context.Context, inputs []*CategoryCreate, next CategoryCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, inputs []*CategoryCreate, next CategoryCreateManyAndReturnQuery) ([]*Category, error)
	FindUnique          func(ctx context.Context, where UniquePredicate[Category], additional []PredicateOf[Category], selects *CategorySelect, omits *CategoryOmit, next CategoryFindUniqueQuery) (*Category, error)
	FindFirst           func(ctx context.Context, params QueryParams[Category], selects *CategorySelect, omits *CategoryOmit, next CategoryFindFirstQuery) (*Category, error)
	FindMany            func(ctx context.Context, params QueryParams[Category], selects *CategorySelect, omits *CategoryOmit, next CategoryFindManyQuery) ([]*Category, error)
	DeleteMany          func(ctx context.Context, preds []PredicateOf[Category], next CategoryDeleteManyQuery) (int64, error)
	Delete              func(ctx context.Context, where UniquePredicate[Category], selects *CategorySelect, omits *CategoryOmit, next CategoryDeleteQuery) (*Category, error)
	Count               func(ctx context.Context, params QueryParams[Category], next CategoryCountQuery) (int64, error)
}

type CategoryDelegate struct {
	client     *Queries
	extensions []CategoryExtension
}

func (d *CategoryDelegate) Use(exts ...CategoryExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *Category) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "id":
			targets[i] = &m.Id
		case "name":
			targets[i] = &m.Name
		}
	}
	return targets
}

var categoryDefaultCols = []string{
	"id",
	"name",
}

var categoryPKCols = []string{
	"id",
}

func selectCategoryCols(selects *CategorySelect, omits *CategoryOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return categoryDefaultCols
	}

	anySelected := selects != nil && (selects.Id || selects.Name || selects.Posts != nil)

	specs := []colSpec{
		{"id", selects != nil && selects.Id, omits != nil && omits.Id, selects != nil && selects.hasAnyRelation()},
		{"name", selects != nil && selects.Name, omits != nil && omits.Name, false},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *CategorySelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.Posts != nil
}

type CategoryCreateBuilder struct {
	*CreateBuilder[Category, CategorySelect, CategoryOmit]
}

func (b *CategoryCreateBuilder) OnConflict(target UniqueConstraintTarget) *CategoryConflictBuilder[CategoryCreateBuilder] {
	return &CategoryConflictBuilder[CategoryCreateBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *CategoryCreateBuilder) SetId(v int32) *CategoryCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}
func (b *CategoryCreateBuilder) SetName(v string) *CategoryCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "name", Val: v})
	return b
}

func (d *CategoryDelegate) Create(assignments ...FieldAssignment) *CategoryCreateBuilder {
	return &CategoryCreateBuilder{
		CreateBuilder: &CreateBuilder[Category, CategorySelect, CategoryOmit]{
			assignments: assignments,
			execFunc:    d.executeCreate,
		},
	}
}

const (
	providedCategoryId   uint64 = 1 << 0
	providedCategoryName uint64 = 1 << 1
)

func assignmentsToCategoryCreate(assignments []FieldAssignment) (CategoryCreate, error) {
	var input CategoryCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "id":
			provided |= providedCategoryId
			if v, ok := a.Val.(int32); ok {
				input.Id = &v
				ValidateInt32(&errs, "id", v, "")
			} else {
				errs.Add("id", a.Val, "type", "field id must be of type int32")
			}
		case "name":
			provided |= providedCategoryName
			if v, ok := a.Val.(string); ok {
				input.Name = v
				ValidateString(&errs, "name", v, true, 0, false, false)
			} else {
				errs.Add("name", a.Val, "type", "field name must be of type string")
			}
		}
	}
	if provided&providedCategoryName == 0 {
		errs.Add("name", "", "required", "field Name is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *CategoryCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 2)
	vals = make([]any, 0, 2)
	if s.Id != nil {
		cols = append(cols, "id")
		vals = append(vals, *s.Id)
	}
	cols = append(cols, "name")
	vals = append(vals, s.Name)
	return
}

func partitionCategoryInputs(dialect Dialect, inputs []*CategoryCreate) [][]*CategoryCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*CategoryCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*CategoryCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*CategoryCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*CategoryCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*CategoryCreate{inputs}
}

func (d *CategoryDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *CategorySelect, omits *CategoryOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*Category, error) {
	input, err := assignmentsToCategoryCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectCategoryCols(selects, omits)

	if len(d.extensions) == 0 {
		hasRelations := selects.hasAnyRelation()
		if hasRelations {
			var res *Category
			err = d.client.transaction(ctx, func(txQ *Queries) error {
				var err error
				res, err = txQ.Category.runCreate(ctx, cols, vals, returningCols, categoryPKCols, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.Category.loadRelations(ctx, []*Category{res}, selects)
			})
			return res, err
		}
		return d.runCreate(ctx, cols, vals, returningCols, categoryPKCols, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args *CategoryCreate) (*Category, error) {
		cols, vals := args.ToColsVals()
		returningCols := selectCategoryCols(selects, omits)

		hasRelations := selects.hasAnyRelation()
		var res *Category
		var err error
		if hasRelations {
			err = d.client.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = txQ.Category.runCreate(c, cols, vals, returningCols, categoryPKCols, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.Category.loadRelations(c, []*Category{res}, selects)
			})
		} else {
			res, err = d.runCreate(c, cols, vals, returningCols, categoryPKCols, conflictTarget, conflictAction)
		}
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Create != nil {
			next, hook := curr, ext.Create
			curr = func(c context.Context, input *CategoryCreate) (*Category, error) {
				return hook(c, input, next)
			}
		}
	}

	return curr(ctx, &input)
}

type CategoryCreateManyBuilder struct {
	*CreateManyBuilder[Category]
}

func (b *CategoryCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *CategoryConflictBuilder[CategoryCreateManyBuilder] {
	return &CategoryConflictBuilder[CategoryCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type CategoryCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[Category, CategorySelect, CategoryOmit]
}

func (b *CategoryCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *CategoryConflictBuilder[CategoryCreateManyAndReturnBuilder] {
	return &CategoryConflictBuilder[CategoryCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (d *CategoryDelegate) CreateMany(builders ...*CategoryCreateBuilder) *CategoryCreateManyBuilder {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &CategoryCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[Category]{
			records:  records,
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *CategoryDelegate) CreateManyAndReturn(builders ...*CategoryCreateBuilder) *CategoryCreateManyAndReturnBuilder {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &CategoryCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[Category, CategorySelect, CategoryOmit]{
			records:  records,
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func (d *CategoryDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs := make([]*CategoryCreate, len(records))
	for i, rec := range records {
		input, err := assignmentsToCategoryCreate(rec.Assignments)
		if err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &input
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args []*CategoryCreate) (int64, error) {
		return d.runCreateMany(c, args, conflictTarget, conflictAction)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.CreateMany != nil {
			next, hook := curr, ext.CreateMany
			curr = func(c context.Context, inputs []*CategoryCreate) (int64, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
}

func (d *CategoryDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *CategorySelect, omits *CategoryOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*Category, error) {
	inputs := make([]*CategoryCreate, len(records))
	for i, rec := range records {
		input, err := assignmentsToCategoryCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &input
	}

	if len(d.extensions) == 0 {
		hasRelations := selects != nil && selects.hasAnyRelation()
		if hasRelations {
			var res []*Category
			err := d.client.transaction(ctx, func(txQ *Queries) error {
				var err error
				res, err = txQ.Category.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.Category.loadRelations(ctx, res, selects)
			})
			return res, err
		}
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args []*CategoryCreate) ([]*Category, error) {
		hasRelations := selects != nil && selects.hasAnyRelation()
		if hasRelations {
			var res []*Category
			err := d.client.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = txQ.Category.runCreateManyAndReturn(c, args, selects, omits, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.Category.loadRelations(c, res, selects)
			})
			return res, err
		}
		return d.runCreateManyAndReturn(c, args, selects, omits, conflictTarget, conflictAction)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.CreateManyAndReturn != nil {
			next, hook := curr, ext.CreateManyAndReturn
			curr = func(c context.Context, inputs []*CategoryCreate) ([]*Category, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
}

func (d *CategoryDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	pkCols []string,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*Category, error) {
	query, clauseArgs := buildSingleInsertSQL(d.client, "Category", cols, returningCols, pkCols, conflictTarget, conflictAction, len(vals))
	if len(clauseArgs) > 0 {
		vals = append(vals, clauseArgs...)
	}

	var res Category
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

func (d *CategoryDelegate) runCreateFallback(ctx context.Context, query string, vals []any, cols []string, returningCols []string, pkCols []string) (*Category, error) {
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
	selectSb.Grow(64 + len(returningCols)*15 + len("Category") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "Category")
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

	var res Category
	if rows.Next() {
		if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
			return nil, err
		}
		return &res, nil
	}
	return nil, rows.Err()
}

func (d *CategoryDelegate) buildBulkInsertSQL(q *Queries, batch []*CategoryCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 2)
	for i, c := range categoryDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "Category")
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
					writeDefault = true
				}
			case "name":
				vals = append(vals, input.Name)
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

func (d *CategoryDelegate) runCreateMany(ctx context.Context, inputs []*CategoryCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionCategoryInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)

		var conflictCols []string
		if conflictTarget != nil {
			conflictCols = conflictTarget.UniqueColumns()
		}
		var nonConflictCols []string
		if conflictAction != nil && conflictAction.Type == ConflictActionUpdateNewValues {
			nonConflictCols = computeNonConflictCols(cols, conflictCols, categoryPKCols)
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

func (d *CategoryDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*CategoryCreate,
	selects *CategorySelect,
	omits *CategoryOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*Category, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionCategoryInputs(d.client.dialect, inputs)
	returningCols := selectCategoryCols(selects, omits)
	hasRelations := selects != nil && selects.hasAnyRelation()

	recordsOut := make([]*Category, 0, len(inputs))

	runBatch := func(txQ *Queries, batch []*CategoryCreate) error {
		cols, vals, queryStr := d.buildBulkInsertSQL(txQ, batch, 1)

		var conflictCols []string
		if conflictTarget != nil {
			conflictCols = conflictTarget.UniqueColumns()
		}
		var nonConflictCols []string
		if conflictAction != nil && conflictAction.Type == ConflictActionUpdateNewValues {
			nonConflictCols = computeNonConflictCols(cols, conflictCols, categoryPKCols)
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
				var res Category
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
		selectSb.Grow(64 + len(returningCols)*15 + len("Category") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			txQ.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		txQ.dialect.WriteQuotedIdent(&selectSb, "Category")
		selectSb.WriteString(" WHERE ")
		txQ.dialect.WriteQuotedIdent(&selectSb, categoryPKCols[0])
		selectSb.WriteString(" >= ")
		txQ.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		txQ.dialect.WriteQuotedIdent(&selectSb, categoryPKCols[0])
		selectSb.WriteString(" < ")
		txQ.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := txQ.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var res Category
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
				return txQ.Category.loadRelations(ctx, recordsOut, selects)
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

type CategoryConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *CategoryConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *CategoryConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *CategoryConflictBuilder[B]) Update(fn func(u *CategoryUpsert)) *B {
	var up ConflictUpdate
	u := newCategoryUpsert(&up)
	fn(u)
	cb.setAction(ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}, cb.conflictTarget)
	return cb.builder
}

type CategoryUpsert struct {
	Id   numericFieldUpsert[int32]
	Name fieldUpsert[string]
}

func newCategoryUpsert(up *ConflictUpdate) *CategoryUpsert {
	return &CategoryUpsert{
		Id: numericFieldUpsert[int32]{
			fieldUpsert: fieldUpsert[int32]{column: "id", update: up},
			tableName:   "Category",
		},
		Name: fieldUpsert[string]{column: "name", update: up},
	}
}
func (d *CategoryDelegate) FindUnique(where UniquePredicate[Category], additional ...PredicateOf[Category]) *FindUniqueBuilder[Category, CategorySelect, CategoryOmit] {
	return &FindUniqueBuilder[Category, CategorySelect, CategoryOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *CategoryDelegate) FindFirst(preds ...PredicateOf[Category]) *FindFirstBuilder[Category, CategorySelect, CategoryOmit] {
	return &FindFirstBuilder[Category, CategorySelect, CategoryOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *CategoryDelegate) FindMany(preds ...PredicateOf[Category]) *FindManyBuilder[Category, CategorySelect, CategoryOmit] {
	return &FindManyBuilder[Category, CategorySelect, CategoryOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *CategoryDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[Category], additional []PredicateOf[Category], selects *CategorySelect, omits *CategoryOmit) (*Category, error) {
	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, where, additional, selects, omits)
	}

	curr := func(c context.Context, w UniquePredicate[Category], add []PredicateOf[Category], sel *CategorySelect, o *CategoryOmit) (*Category, error) {
		return d.runFindUnique(c, w, add, sel, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindUnique != nil {
			next, hook := curr, ext.FindUnique
			curr = func(c context.Context, w UniquePredicate[Category], add []PredicateOf[Category], sel *CategorySelect, o *CategoryOmit) (*Category, error) {
				return hook(c, w, add, sel, o, next)
			}
		}
	}

	return curr(ctx, where, additional, selects, omits)
}

func (d *CategoryDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[Category],
	selects *CategorySelect,
	omits *CategoryOmit,
) (*Category, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	curr := func(c context.Context, p QueryParams[Category], sel *CategorySelect, o *CategoryOmit) (*Category, error) {
		return d.runFindFirst(c, p, sel, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindFirst != nil {
			next, hook := curr, ext.FindFirst
			curr = func(c context.Context, p QueryParams[Category], sel *CategorySelect, o *CategoryOmit) (*Category, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}

func (d *CategoryDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[Category],
	selects *CategorySelect,
	omits *CategoryOmit,
) ([]*Category, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	curr := func(c context.Context, p QueryParams[Category], sel *CategorySelect, o *CategoryOmit) ([]*Category, error) {
		return d.runFindMany(c, p, sel, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindMany != nil {
			next, hook := curr, ext.FindMany
			curr = func(c context.Context, p QueryParams[Category], sel *CategorySelect, o *CategoryOmit) ([]*Category, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}

func (d *CategoryDelegate) runFindUnique(ctx context.Context, where UniquePredicate[Category], additional []PredicateOf[Category], selects *CategorySelect, omits *CategoryOmit) (*Category, error) {
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
	allPreds := append([]PredicateOf[Category]{where}, additional...)
	whereClause, vals := CompilePredicates(d.client.dialect, allPreds)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	returningCols := selectCategoryCols(selects, omits)

	var res *Category
	var err error
	if selects.hasAnyRelation() {
		err = d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Category.queryOne(ctx, whereClause, vals, returningCols, nil)
			if err != nil {
				return err
			}
			if res == nil {
				return nil
			}
			return txQ.Category.loadRelations(ctx, []*Category{res}, selects)
		})
	} else {
		res, err = d.queryOne(ctx, whereClause, vals, returningCols, nil)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *CategoryDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[Category],
	selects *CategorySelect,
	omits *CategoryOmit,
) (*Category, error) {
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
	returningCols := selectCategoryCols(selects, omits)

	var res *Category
	var err error
	if selects.hasAnyRelation() {
		err = d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Category.queryOne(ctx, whereClause, vals, returningCols, params.Skip)
			if err != nil {
				return err
			}
			if res == nil {
				return nil
			}
			return txQ.Category.loadRelations(ctx, []*Category{res}, selects)
		})
	} else {
		res, err = d.queryOne(ctx, whereClause, vals, returningCols, params.Skip)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *CategoryDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[Category],
	selects *CategorySelect,
	omits *CategoryOmit,
) ([]*Category, error) {
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
	returningCols := selectCategoryCols(selects, omits)

	var results []*Category
	var err error
	if selects.hasAnyRelation() {
		err = d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			results, err = txQ.Category.queryMany(ctx, whereClause, vals, returningCols, params.Take, params.Skip)
			if err != nil {
				return err
			}
			if len(results) == 0 {
				return nil
			}
			return txQ.Category.loadRelations(ctx, results, selects)
		})
	} else {
		results, err = d.queryMany(ctx, whereClause, vals, returningCols, params.Take, params.Skip)
	}
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (d *CategoryDelegate) queryOne(ctx context.Context, whereClause string, whereVals []any, returningCols []string, skip *int) (*Category, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "Category", returningCols, whereClause, &limitOne, skip)
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

	var res Category
	if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
		return nil, err
	}
	return &res, nil
}

func (d *CategoryDelegate) queryMany(ctx context.Context, whereClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*Category, error) {
	query := buildSelectSQL(d.client, "Category", returningCols, whereClause, take, skip)
	rows, err := d.client.query(ctx, query, whereVals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*Category, 0)
	for rows.Next() {
		var res Category
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
func (d *CategoryDelegate) DeleteMany(preds ...PredicateOf[Category]) *DeleteManyBuilder[Category] {
	return &DeleteManyBuilder[Category]{
		where:    preds,
		execFunc: d.executeDeleteMany,
	}
}

func (d *CategoryDelegate) executeDeleteMany(ctx context.Context, preds []PredicateOf[Category]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runDeleteMany(ctx, preds)
	}

	curr := func(c context.Context, p []PredicateOf[Category]) (int64, error) {
		return d.runDeleteMany(c, p)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.DeleteMany != nil {
			next, hook := curr, ext.DeleteMany
			curr = func(c context.Context, p []PredicateOf[Category]) (int64, error) {
				return hook(c, p, next)
			}
		}
	}

	return curr(ctx, preds)
}

func (d *CategoryDelegate) runDeleteMany(ctx context.Context, preds []PredicateOf[Category]) (int64, error) {
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
	d.client.dialect.WriteQuotedIdent(&sb, "Category")
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

func (d *CategoryDelegate) Delete(where UniquePredicate[Category]) *DeleteBuilder[Category, CategorySelect, CategoryOmit] {
	return &DeleteBuilder[Category, CategorySelect, CategoryOmit]{
		where:    where,
		execFunc: d.executeDelete,
	}
}

func (d *CategoryDelegate) executeDelete(ctx context.Context, where UniquePredicate[Category], selects *CategorySelect, omits *CategoryOmit) (*Category, error) {
	if len(d.extensions) == 0 {
		return d.runDelete(ctx, where, selects, omits)
	}

	curr := func(c context.Context, w UniquePredicate[Category], s *CategorySelect, o *CategoryOmit) (*Category, error) {
		return d.runDelete(c, w, s, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Delete != nil {
			next, hook := curr, ext.Delete
			curr = func(c context.Context, w UniquePredicate[Category], s *CategorySelect, o *CategoryOmit) (*Category, error) {
				return hook(c, w, s, o, next)
			}
		}
	}

	return curr(ctx, where, selects, omits)
}

func (d *CategoryDelegate) runDelete(ctx context.Context, where UniquePredicate[Category], selects *CategorySelect, omits *CategoryOmit) (*Category, error) {
	if err := where.Validate(); err != nil {
		return nil, err
	}

	returningCols := selectCategoryCols(selects, omits, categoryPKCols...)

	hasRelations := selects != nil && selects.hasAnyRelation()
	useTx := !d.client.dialect.SupportsDeleteReturning || hasRelations

	if useTx {
		var res *Category
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Category.executeFindUnique(ctx, where, nil, selects, omits)
			if err != nil {
				return err
			}
			if res == nil {
				return sql.ErrNoRows
			}

			// Build DELETE statement by PK
			var deleteSb strings.Builder
			deleteSb.WriteString("DELETE FROM ")
			txQ.dialect.WriteQuotedIdent(&deleteSb, "Category")
			deleteSb.WriteString(" WHERE ")

			var pkPreds []PredicateOf[Category]
			pkPreds = append(pkPreds, Predicate[Category]{
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
	d.client.dialect.WriteQuotedIdent(&sb, "Category")

	whereClause, vals := CompilePredicates(d.client.dialect, []PredicateOf[Category]{where})
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

	var row Category
	if err := rows.Scan(row.ScanFields(returningCols)...); err != nil {
		return nil, err
	}
	return &row, nil
}
func (d *CategoryDelegate) Count(preds ...PredicateOf[Category]) *CountBuilder[Category] {
	return &CountBuilder[Category]{
		where:    preds,
		execFunc: d.executeCount,
	}
}

func (d *CategoryDelegate) executeCount(ctx context.Context, params QueryParams[Category]) (int64, error) {
	if len(d.extensions) == 0 {
		return d.runCount(ctx, params)
	}

	curr := func(c context.Context, p QueryParams[Category]) (int64, error) {
		return d.runCount(c, p)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Count != nil {
			next, hook := curr, ext.Count
			curr = func(c context.Context, p QueryParams[Category]) (int64, error) {
				return hook(c, p, next)
			}
		}
	}

	return curr(ctx, params)
}

func (d *CategoryDelegate) runCount(ctx context.Context, params QueryParams[Category]) (int64, error) {
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
		d.client.dialect.WriteQuotedIdent(&subQuery, "Category")
		if whereClause != "" {
			subQuery.WriteString(whereClause)
		}
		subQuery.WriteString(d.client.dialect.FormatLimitOffset(params.Take, params.Skip))
		query = "SELECT COUNT(*) FROM (" + subQuery.String() + ") as sub"
	} else {
		var sb strings.Builder
		sb.WriteString("SELECT COUNT(*) FROM ")
		d.client.dialect.WriteQuotedIdent(&sb, "Category")
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
func (d *CategoryDelegate) loadRelations(ctx context.Context, records []*Category, selects *CategorySelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Posts != nil {
		relationSelects, relationOmits, relationParams := selects.Posts.GetRelationParams()
		returningCols := selectCategoryToPostCols(relationSelects, relationOmits, "categoryId")
		// Inverse holds the FK: CategoryToPost.categoryId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *Category) int32 { return p.Id }),
			"CategoryToPost",
			"categoryId",
			returningCols,
			scanInto(returningCols, (*CategoryToPost).ScanFields),
			directKey(func(c *CategoryToPost) int32 { return c.CategoryId }),
			appendMany(func(p *Category) *[]*CategoryToPost { return &p.Posts }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading posts: %w", err)
		}
		if err := d.client.CategoryToPost.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
