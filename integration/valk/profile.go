package valk

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"time"
)

// Profile represents the database model
type Profile struct {
	Id        string    `db:"id" json:"id"`
	Bio       *string   `db:"bio" json:"bio,omitempty"`
	UserId    string    `db:"userId" json:"userId"`
	CreatedAt time.Time `db:"createdAt" json:"createdAt"`
	User      *User     `json:"user,omitempty"`
}

// ProfileCreate is used for hooks only — the Create API uses FieldAssignment
type ProfileCreate struct {
	Id        *string    `json:"id"`
	Bio       *string    `json:"bio"`
	UserId    string     `json:"userId"`
	CreatedAt *time.Time `json:"createdAt"`
}

// colMask returns a bit mask of columns that are set
func (s *ProfileCreate) colMask() uint64 {
	var mask uint64
	mask |= 1 << 0
	if s.Bio != nil {
		mask |= 1 << 1
	}
	mask |= 1 << 2
	mask |= 1 << 3
	return mask
}

type ProfileSelect struct {
	Id        bool            `json:"id"`
	Bio       bool            `json:"bio"`
	UserId    bool            `json:"userId"`
	CreatedAt bool            `json:"createdAt"`
	User      UserSelectQuery `json:"user,omitempty"`
}

type ProfileOmit struct {
	Id        bool `json:"id"`
	Bio       bool `json:"bio"`
	UserId    bool `json:"userId"`
	CreatedAt bool `json:"createdAt"`
}

type ProfileSelectQuery interface {
	GetRelationParams() (*ProfileSelect, *ProfileOmit, QueryParams[Profile])
}

func (s *ProfileSelect) GetRelationParams() (*ProfileSelect, *ProfileOmit, QueryParams[Profile]) {
	return s, nil, QueryParams[Profile]{}
}

type ProfileQueryBuilder struct {
	selects *ProfileSelect
	omits   *ProfileOmit
	where   []PredicateOf[Profile]
	take    *int
	skip    *int
	orderBy []OrderBy[Profile]
}

func (b *ProfileQueryBuilder) Where(preds ...PredicateOf[Profile]) *ProfileQueryBuilder {
	b.where = append(b.where, preds...)
	return b
}

func (b *ProfileQueryBuilder) Take(limit int) *ProfileQueryBuilder {
	b.take = &limit
	return b
}

func (b *ProfileQueryBuilder) Skip(offset int) *ProfileQueryBuilder {
	b.skip = &offset
	return b
}

func (b *ProfileQueryBuilder) OrderBy(orders ...OrderBy[Profile]) *ProfileQueryBuilder {
	b.orderBy = append(b.orderBy, orders...)
	return b
}

func (b *ProfileQueryBuilder) Select(s ProfileSelect) *ProfileQueryBuilder {
	b.selects = &s
	return b
}

func (b *ProfileQueryBuilder) Omit(o ProfileOmit) *ProfileQueryBuilder {
	b.omits = &o
	return b
}

func (b *ProfileQueryBuilder) GetRelationParams() (*ProfileSelect, *ProfileOmit, QueryParams[Profile]) {
	if b == nil {
		return nil, nil, QueryParams[Profile]{}
	}
	return b.selects, b.omits, QueryParams[Profile]{
		Where:   b.where,
		Take:    b.take,
		Skip:    b.skip,
		OrderBy: b.orderBy,
	}
}

type ProfileCreateQuery = func(ctx context.Context, args *ProfileCreate) (*Profile, error)
type ProfileCreateManyQuery = func(ctx context.Context, args []*ProfileCreate) (int64, error)
type ProfileCreateManyAndReturnQuery = func(ctx context.Context, args []*ProfileCreate) ([]*Profile, error)
type ProfileFindUniqueQuery = func(ctx context.Context, where UniquePredicate[Profile], additional []PredicateOf[Profile], selects *ProfileSelect, omits *ProfileOmit) (*Profile, error)
type ProfileFindFirstQuery = func(ctx context.Context, params QueryParams[Profile], selects *ProfileSelect, omits *ProfileOmit) (*Profile, error)
type ProfileFindManyQuery = func(ctx context.Context, params QueryParams[Profile], selects *ProfileSelect, omits *ProfileOmit) ([]*Profile, error)

type ProfileExtension struct {
	Create              func(ctx context.Context, input *ProfileCreate, next ProfileCreateQuery) (*Profile, error)
	CreateMany          func(ctx context.Context, inputs []*ProfileCreate, next ProfileCreateManyQuery) (int64, error)
	CreateManyAndReturn func(ctx context.Context, inputs []*ProfileCreate, next ProfileCreateManyAndReturnQuery) ([]*Profile, error)
	FindUnique          func(ctx context.Context, where UniquePredicate[Profile], additional []PredicateOf[Profile], selects *ProfileSelect, omits *ProfileOmit, next ProfileFindUniqueQuery) (*Profile, error)
	FindFirst           func(ctx context.Context, params QueryParams[Profile], selects *ProfileSelect, omits *ProfileOmit, next ProfileFindFirstQuery) (*Profile, error)
	FindMany            func(ctx context.Context, params QueryParams[Profile], selects *ProfileSelect, omits *ProfileOmit, next ProfileFindManyQuery) ([]*Profile, error)
}

type ProfileDelegate struct {
	client     *Queries
	extensions []ProfileExtension
}

func (d *ProfileDelegate) Use(exts ...ProfileExtension) {
	d.extensions = append(d.extensions, exts...)
}

func (m *Profile) ScanFields(cols []string) []any {
	targets := make([]any, len(cols))
	for i, col := range cols {
		switch col {
		case "id":
			targets[i] = &m.Id
		case "bio":
			targets[i] = &m.Bio
		case "userId":
			targets[i] = &m.UserId
		case "createdAt":
			targets[i] = &m.CreatedAt
		}
	}
	return targets
}

var profileDefaultCols = []string{
	"id",
	"bio",
	"userId",
	"createdAt",
}

var profilePKCols = []string{
	"id",
}

func selectProfileCols(selects *ProfileSelect, omits *ProfileOmit, forceCols ...string) []string {
	if selects == nil && omits == nil && len(forceCols) == 0 {
		return profileDefaultCols
	}

	anySelected := selects != nil && (selects.Id || selects.Bio || selects.UserId || selects.CreatedAt || selects.User != nil)

	specs := []colSpec{
		{"id", selects != nil && selects.Id, omits != nil && omits.Id, selects != nil && selects.hasAnyRelation()},
		{"bio", selects != nil && selects.Bio, omits != nil && omits.Bio, false},
		{"userId", selects != nil && selects.UserId, omits != nil && omits.UserId, selects != nil && selects.User != nil},
		{"createdAt", selects != nil && selects.CreatedAt, omits != nil && omits.CreatedAt, false},
	}

	cols := computeCols(specs, selects != nil, anySelected)

	for _, f := range forceCols {
		if !slices.Contains(cols, f) {
			cols = append(cols, f)
		}
	}

	return cols
}

func (s *ProfileSelect) hasAnyRelation() bool {
	if s == nil {
		return false
	}
	return s.User != nil
}

type ProfileCreateBuilder struct {
	*CreateBuilder[Profile, ProfileSelect, ProfileOmit]
}

func (b *ProfileCreateBuilder) OnConflict(target UniqueConstraintTarget) *ProfileConflictBuilder[ProfileCreateBuilder] {
	return &ProfileConflictBuilder[ProfileCreateBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (b *ProfileCreateBuilder) SetId(v string) *ProfileCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "id", Val: v})
	return b
}
func (b *ProfileCreateBuilder) SetBio(v string) *ProfileCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "bio", Val: v})
	return b
}
func (b *ProfileCreateBuilder) SetUserId(v string) *ProfileCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "userId", Val: v})
	return b
}
func (b *ProfileCreateBuilder) SetCreatedAt(v time.Time) *ProfileCreateBuilder {
	b.assignments = append(b.assignments, FieldAssignment{Col: "createdAt", Val: v})
	return b
}

func (d *ProfileDelegate) Create(assignments ...FieldAssignment) *ProfileCreateBuilder {
	return &ProfileCreateBuilder{
		CreateBuilder: &CreateBuilder[Profile, ProfileSelect, ProfileOmit]{
			assignments: assignments,
			execFunc:    d.executeCreate,
		},
	}
}

const (
	providedProfileId        uint64 = 1 << 0
	providedProfileBio       uint64 = 1 << 1
	providedProfileUserId    uint64 = 1 << 2
	providedProfileCreatedAt uint64 = 1 << 3
)

func assignmentsToProfileCreate(assignments []FieldAssignment) (ProfileCreate, error) {
	var input ProfileCreate
	var errs ValidationError
	var provided uint64

	for _, a := range assignments {
		switch a.Col {
		case "id":
			provided |= providedProfileId
			if v, ok := a.Val.(string); ok {
				input.Id = &v
				ValidateString(&errs, "id", v, false, 0, false, false)
			} else {
				errs.Add("id", a.Val, "type", "field id must be of type string")
			}
		case "bio":
			provided |= providedProfileBio
			if v, ok := a.Val.(string); ok {
				input.Bio = &v
				ValidateString(&errs, "bio", v, false, 0, false, false)
			} else {
				errs.Add("bio", a.Val, "type", "field bio must be of type string")
			}
		case "userId":
			provided |= providedProfileUserId
			if v, ok := a.Val.(string); ok {
				input.UserId = v
				ValidateString(&errs, "userId", v, true, 0, false, false)
			} else {
				errs.Add("userId", a.Val, "type", "field userId must be of type string")
			}
		case "createdAt":
			provided |= providedProfileCreatedAt
			if v, ok := a.Val.(time.Time); ok {
				input.CreatedAt = &v
			} else {
				errs.Add("createdAt", a.Val, "type", "field createdAt must be of type time.Time")
			}
		}
	}
	if provided&providedProfileUserId == 0 {
		errs.Add("userId", "", "required", "field UserId is required")
	}

	if errs.HasErrors() {
		return input, errs
	}
	return input, nil
}

func (s *ProfileCreate) ToColsVals() (cols []string, vals []any) {
	cols = make([]string, 0, 4)
	vals = make([]any, 0, 4)
	cols = append(cols, "id")
	if s.Id != nil {
		vals = append(vals, *s.Id)
	} else {
		vals = append(vals, generateCUID())
	}
	if s.Bio != nil {
		cols = append(cols, "bio")
		vals = append(vals, *s.Bio)
	}
	cols = append(cols, "userId")
	vals = append(vals, s.UserId)
	cols = append(cols, "createdAt")
	if s.CreatedAt != nil {
		vals = append(vals, *s.CreatedAt)
	} else {
		vals = append(vals, time.Now())
	}
	return
}

func partitionProfileInputs(dialect Dialect, inputs []*ProfileCreate) [][]*ProfileCreate {
	if !dialect.SupportsBulkInsert {
		result := make([][]*ProfileCreate, len(inputs))
		for i, input := range inputs {
			result[i] = []*ProfileCreate{input}
		}
		return result
	}

	if !dialect.SupportsDefaultKeyword {
		groups := make(map[uint64][]*ProfileCreate)
		var masks []uint64
		for _, input := range inputs {
			mask := input.colMask()
			if _, exists := groups[mask]; !exists {
				masks = append(masks, mask)
			}
			groups[mask] = append(groups[mask], input)
		}
		result := make([][]*ProfileCreate, len(masks))
		for i, mask := range masks {
			result[i] = groups[mask]
		}
		return result
	}

	return [][]*ProfileCreate{inputs}
}

func (d *ProfileDelegate) executeCreate(ctx context.Context, assignments []FieldAssignment, selects *ProfileSelect, omits *ProfileOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*Profile, error) {
	input, err := assignmentsToProfileCreate(assignments)
	if err != nil {
		return nil, err
	}

	cols, vals := input.ToColsVals()
	returningCols := selectProfileCols(selects, omits)
	pkCols := profilePKCols

	if len(d.extensions) == 0 {
		hasRelations := selects.hasAnyRelation()
		if hasRelations {
			var res *Profile
			err = d.client.transaction(ctx, func(txQ *Queries) error {
				var err error
				res, err = txQ.Profile.runCreate(ctx, cols, vals, returningCols, pkCols, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.Profile.loadRelations(ctx, []*Profile{res}, selects)
			})
			return res, err
		}
		return d.runCreate(ctx, cols, vals, returningCols, pkCols, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args *ProfileCreate) (*Profile, error) {
		cols, vals := args.ToColsVals()
		returningCols := selectProfileCols(selects, omits)
		pkCols := profilePKCols

		hasRelations := selects.hasAnyRelation()
		var res *Profile
		var err error
		if hasRelations {
			err = d.client.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = txQ.Profile.runCreate(c, cols, vals, returningCols, pkCols, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.Profile.loadRelations(c, []*Profile{res}, selects)
			})
		} else {
			res, err = d.runCreate(c, cols, vals, returningCols, pkCols, conflictTarget, conflictAction)
		}
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.Create != nil {
			next, hook := curr, ext.Create
			curr = func(c context.Context, input *ProfileCreate) (*Profile, error) {
				return hook(c, input, next)
			}
		}
	}

	return curr(ctx, &input)
}

type ProfileCreateManyBuilder struct {
	*CreateManyBuilder[Profile]
}

func (b *ProfileCreateManyBuilder) OnConflict(target UniqueConstraintTarget) *ProfileConflictBuilder[ProfileCreateManyBuilder] {
	return &ProfileConflictBuilder[ProfileCreateManyBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

type ProfileCreateManyAndReturnBuilder struct {
	*CreateManyAndReturnBuilder[Profile, ProfileSelect, ProfileOmit]
}

func (b *ProfileCreateManyAndReturnBuilder) OnConflict(target UniqueConstraintTarget) *ProfileConflictBuilder[ProfileCreateManyAndReturnBuilder] {
	return &ProfileConflictBuilder[ProfileCreateManyAndReturnBuilder]{
		builder:        b,
		conflictTarget: target,
		setAction: func(action ConflictAction, target UniqueConstraintTarget) {
			b.conflictAction = &action
			b.conflictTarget = target
		},
	}
}

func (d *ProfileDelegate) CreateMany(builders ...*ProfileCreateBuilder) *ProfileCreateManyBuilder {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &ProfileCreateManyBuilder{
		CreateManyBuilder: &CreateManyBuilder[Profile]{
			records:  records,
			execFunc: d.executeCreateMany,
		},
	}
}

func (d *ProfileDelegate) CreateManyAndReturn(builders ...*ProfileCreateBuilder) *ProfileCreateManyAndReturnBuilder {
	records := make([]RecordInput, len(builders))
	for i, b := range builders {
		records[i] = RecordInput{Assignments: b.assignments}
	}
	return &ProfileCreateManyAndReturnBuilder{
		CreateManyAndReturnBuilder: &CreateManyAndReturnBuilder[Profile, ProfileSelect, ProfileOmit]{
			records:  records,
			execFunc: d.executeCreateManyAndReturn,
		},
	}
}

func (d *ProfileDelegate) executeCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs := make([]*ProfileCreate, len(records))
	for i, rec := range records {
		input, err := assignmentsToProfileCreate(rec.Assignments)
		if err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &input
	}

	if len(d.extensions) == 0 {
		return d.runCreateMany(ctx, inputs, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args []*ProfileCreate) (int64, error) {
		return d.runCreateMany(c, args, conflictTarget, conflictAction)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.CreateMany != nil {
			next, hook := curr, ext.CreateMany
			curr = func(c context.Context, inputs []*ProfileCreate) (int64, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
}

func (d *ProfileDelegate) executeCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *ProfileSelect, omits *ProfileOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*Profile, error) {
	inputs := make([]*ProfileCreate, len(records))
	for i, rec := range records {
		input, err := assignmentsToProfileCreate(rec.Assignments)
		if err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		inputs[i] = &input
	}

	if len(d.extensions) == 0 {
		hasRelations := selects != nil && selects.hasAnyRelation()
		if hasRelations {
			var res []*Profile
			err := d.client.transaction(ctx, func(txQ *Queries) error {
				var err error
				res, err = txQ.Profile.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.Profile.loadRelations(ctx, res, selects)
			})
			return res, err
		}
		return d.runCreateManyAndReturn(ctx, inputs, selects, omits, conflictTarget, conflictAction)
	}

	curr := func(c context.Context, args []*ProfileCreate) ([]*Profile, error) {
		hasRelations := selects != nil && selects.hasAnyRelation()
		if hasRelations {
			var res []*Profile
			err := d.client.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = txQ.Profile.runCreateManyAndReturn(c, args, selects, omits, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.Profile.loadRelations(c, res, selects)
			})
			return res, err
		}
		return d.runCreateManyAndReturn(c, args, selects, omits, conflictTarget, conflictAction)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.CreateManyAndReturn != nil {
			next, hook := curr, ext.CreateManyAndReturn
			curr = func(c context.Context, inputs []*ProfileCreate) ([]*Profile, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
}

func (d *ProfileDelegate) runCreate(
	ctx context.Context,
	cols []string,
	vals []any,
	returningCols []string,
	pkCols []string,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) (*Profile, error) {
	query, clauseArgs := buildSingleInsertSQL(d.client, "Profile", cols, returningCols, pkCols, conflictTarget, conflictAction, len(vals))
	if len(clauseArgs) > 0 {
		vals = append(vals, clauseArgs...)
	}

	var res Profile
	if d.client.dialect.SupportsReturning {
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

func (d *ProfileDelegate) runCreateFallback(ctx context.Context, query string, vals []any, cols []string, returningCols []string, pkCols []string) (*Profile, error) {
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
	selectSb.Grow(64 + len(returningCols)*15 + len("Profile") + len(pkCols)*15)
	selectSb.WriteString("SELECT ")
	for i, col := range returningCols {
		if i > 0 {
			selectSb.WriteString(", ")
		}
		d.client.dialect.WriteQuotedIdent(&selectSb, col)
	}
	selectSb.WriteString(" FROM ")
	d.client.dialect.WriteQuotedIdent(&selectSb, "Profile")
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

	var res Profile
	if rows.Next() {
		if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
			return nil, err
		}
		return &res, nil
	}
	return nil, rows.Err()
}

func (d *ProfileDelegate) buildBulkInsertSQL(q *Queries, batch []*ProfileCreate, paramStartIdx int) (cols []string, vals []any, queryStr string) {
	var colMask uint64
	for _, input := range batch {
		colMask |= input.colMask()
	}

	cols = make([]string, 0, 4)
	for i, c := range profileDefaultCols {
		if colMask&(1<<i) != 0 {
			cols = append(cols, c)
		}
	}

	vals = make([]any, 0, len(batch)*len(cols))
	var sb strings.Builder
	sb.Grow(128 + len(batch)*len(cols)*10)
	sb.WriteString("INSERT INTO ")
	q.dialect.WriteQuotedIdent(&sb, "Profile")
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
			case "bio":
				if input.Bio != nil {
					vals = append(vals, *input.Bio)
				} else {
					writeDefault = true
				}
			case "userId":
				vals = append(vals, input.UserId)
			case "createdAt":
				if input.CreatedAt != nil {
					vals = append(vals, *input.CreatedAt)
				} else {
					vals = append(vals, time.Now())
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

func (d *ProfileDelegate) runCreateMany(ctx context.Context, inputs []*ProfileCreate, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	if len(inputs) == 0 {
		return 0, nil
	}

	batches := partitionProfileInputs(d.client.dialect, inputs)

	var count int64
	for _, batch := range batches {
		cols, vals, queryStr := d.buildBulkInsertSQL(d.client, batch, 1)

		var conflictCols []string
		if conflictTarget != nil {
			conflictCols = conflictTarget.UniqueColumns()
		}
		var nonConflictCols []string
		pkCols := profilePKCols
		if conflictAction != nil && conflictAction.Type == ConflictActionUpdateNewValues {
			nonConflictCols = computeNonConflictCols(cols, conflictCols, pkCols)
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

func (d *ProfileDelegate) runCreateManyAndReturn(
	ctx context.Context,
	inputs []*ProfileCreate,
	selects *ProfileSelect,
	omits *ProfileOmit,
	conflictTarget UniqueConstraintTarget,
	conflictAction *ConflictAction,
) ([]*Profile, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	batches := partitionProfileInputs(d.client.dialect, inputs)
	returningCols := selectProfileCols(selects, omits)
	hasRelations := selects != nil && selects.hasAnyRelation()

	recordsOut := make([]*Profile, 0, len(inputs))

	runBatch := func(txQ *Queries, batch []*ProfileCreate) error {
		cols, vals, queryStr := d.buildBulkInsertSQL(txQ, batch, 1)

		var conflictCols []string
		if conflictTarget != nil {
			conflictCols = conflictTarget.UniqueColumns()
		}
		var nonConflictCols []string
		pkCols := profilePKCols
		if conflictAction != nil && conflictAction.Type == ConflictActionUpdateNewValues {
			nonConflictCols = computeNonConflictCols(cols, conflictCols, pkCols)
		}
		clause, clauseArgs := txQ.dialect.BuildConflictClause(conflictCols, conflictAction, nonConflictCols, len(vals)+1)
		queryStr += clause
		vals = append(vals, clauseArgs...)

		if txQ.dialect.SupportsReturning && len(returningCols) > 0 {
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
				var res Profile
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
		selectSb.Grow(64 + len(returningCols)*15 + len("Profile") + len(batch)*15)
		selectSb.WriteString("SELECT ")
		for i, col := range returningCols {
			if i > 0 {
				selectSb.WriteString(", ")
			}
			txQ.dialect.WriteQuotedIdent(&selectSb, col)
		}
		selectSb.WriteString(" FROM ")
		txQ.dialect.WriteQuotedIdent(&selectSb, "Profile")
		selectSb.WriteString(" WHERE ")
		txQ.dialect.WriteQuotedIdent(&selectSb, pkCols[0])
		selectSb.WriteString(" >= ")
		txQ.dialect.WritePlaceholder(&selectSb, 1)
		selectSb.WriteString(" AND ")
		txQ.dialect.WriteQuotedIdent(&selectSb, pkCols[0])
		selectSb.WriteString(" < ")
		txQ.dialect.WritePlaceholder(&selectSb, 2)

		rows, err := txQ.query(ctx, selectSb.String(), lastID, lastID+int64(len(batch)))
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var res Profile
			if err := rows.Scan(res.ScanFields(returningCols)...); err != nil {
				return err
			}
			recordsOut = append(recordsOut, &res)
		}
		return rows.Err()
	}

	// Always wrap in transaction if we have multiple batches OR if we need to load relations
	if len(batches) > 1 || hasRelations || !d.client.dialect.SupportsReturning {
		err := d.client.transaction(ctx, func(txQ *Queries) error {
			for _, batch := range batches {
				if err := runBatch(txQ, batch); err != nil {
					return err
				}
			}
			if hasRelations {
				return txQ.Profile.loadRelations(ctx, recordsOut, selects)
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

type ProfileConflictBuilder[B any] struct {
	builder        *B
	setAction      func(ConflictAction, UniqueConstraintTarget)
	conflictTarget UniqueConstraintTarget
}

func (cb *ProfileConflictBuilder[B]) Ignore() *B {
	cb.setAction(ConflictAction{Type: ConflictActionIgnore}, cb.conflictTarget)
	return cb.builder
}

func (cb *ProfileConflictBuilder[B]) UpdateNewValues() *B {
	cb.setAction(ConflictAction{Type: ConflictActionUpdateNewValues}, cb.conflictTarget)
	return cb.builder
}

func (cb *ProfileConflictBuilder[B]) Update(fn func(u *ProfileUpsert)) *B {
	var up ConflictUpdate
	u := newProfileUpsert(&up)
	fn(u)
	cb.setAction(ConflictAction{
		Type:        ConflictActionUpdateCustom,
		Assignments: up.assignments,
		Args:        up.args,
	}, cb.conflictTarget)
	return cb.builder
}

type ProfileUpsert struct {
	Id        fieldUpsert[string]
	Bio       fieldUpsert[*string]
	UserId    fieldUpsert[string]
	CreatedAt fieldUpsert[time.Time]
}

func newProfileUpsert(up *ConflictUpdate) *ProfileUpsert {
	return &ProfileUpsert{
		Id:        fieldUpsert[string]{column: "id", update: up},
		Bio:       fieldUpsert[*string]{column: "bio", update: up},
		UserId:    fieldUpsert[string]{column: "userId", update: up},
		CreatedAt: fieldUpsert[time.Time]{column: "createdAt", update: up},
	}
}
func (d *ProfileDelegate) FindUnique(where UniquePredicate[Profile], additional ...PredicateOf[Profile]) *FindUniqueBuilder[Profile, ProfileSelect, ProfileOmit] {
	return &FindUniqueBuilder[Profile, ProfileSelect, ProfileOmit]{
		where:      where,
		additional: additional,
		execFunc:   d.executeFindUnique,
	}
}

func (d *ProfileDelegate) FindFirst(preds ...PredicateOf[Profile]) *FindFirstBuilder[Profile, ProfileSelect, ProfileOmit] {
	return &FindFirstBuilder[Profile, ProfileSelect, ProfileOmit]{
		where:    preds,
		execFunc: d.executeFindFirst,
	}
}

func (d *ProfileDelegate) FindMany(preds ...PredicateOf[Profile]) *FindManyBuilder[Profile, ProfileSelect, ProfileOmit] {
	return &FindManyBuilder[Profile, ProfileSelect, ProfileOmit]{
		where:    preds,
		execFunc: d.executeFindMany,
	}
}

func (d *ProfileDelegate) executeFindUnique(ctx context.Context, where UniquePredicate[Profile], additional []PredicateOf[Profile], selects *ProfileSelect, omits *ProfileOmit) (*Profile, error) {
	if len(d.extensions) == 0 {
		return d.runFindUnique(ctx, where, additional, selects, omits)
	}

	curr := func(c context.Context, w UniquePredicate[Profile], add []PredicateOf[Profile], sel *ProfileSelect, o *ProfileOmit) (*Profile, error) {
		return d.runFindUnique(c, w, add, sel, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindUnique != nil {
			next, hook := curr, ext.FindUnique
			curr = func(c context.Context, w UniquePredicate[Profile], add []PredicateOf[Profile], sel *ProfileSelect, o *ProfileOmit) (*Profile, error) {
				return hook(c, w, add, sel, o, next)
			}
		}
	}

	return curr(ctx, where, additional, selects, omits)
}

func (d *ProfileDelegate) executeFindFirst(
	ctx context.Context,
	params QueryParams[Profile],
	selects *ProfileSelect,
	omits *ProfileOmit,
) (*Profile, error) {
	if len(d.extensions) == 0 {
		return d.runFindFirst(ctx, params, selects, omits)
	}

	curr := func(c context.Context, p QueryParams[Profile], sel *ProfileSelect, o *ProfileOmit) (*Profile, error) {
		return d.runFindFirst(c, p, sel, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindFirst != nil {
			next, hook := curr, ext.FindFirst
			curr = func(c context.Context, p QueryParams[Profile], sel *ProfileSelect, o *ProfileOmit) (*Profile, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}

func (d *ProfileDelegate) executeFindMany(
	ctx context.Context,
	params QueryParams[Profile],
	selects *ProfileSelect,
	omits *ProfileOmit,
) ([]*Profile, error) {
	if len(d.extensions) == 0 {
		return d.runFindMany(ctx, params, selects, omits)
	}

	curr := func(c context.Context, p QueryParams[Profile], sel *ProfileSelect, o *ProfileOmit) ([]*Profile, error) {
		return d.runFindMany(c, p, sel, o)
	}

	for _, ext := range slices.Backward(d.extensions) {
		if ext.FindMany != nil {
			next, hook := curr, ext.FindMany
			curr = func(c context.Context, p QueryParams[Profile], sel *ProfileSelect, o *ProfileOmit) ([]*Profile, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}

func (d *ProfileDelegate) runFindUnique(ctx context.Context, where UniquePredicate[Profile], additional []PredicateOf[Profile], selects *ProfileSelect, omits *ProfileOmit) (*Profile, error) {
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
	allPreds := append([]PredicateOf[Profile]{where}, additional...)
	whereClause, vals := CompilePredicates(d.client.dialect, allPreds)
	if whereClause != "" {
		whereClause = " WHERE " + whereClause
	}
	returningCols := selectProfileCols(selects, omits)

	var res *Profile
	var err error
	if selects.hasAnyRelation() {
		err = d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Profile.queryOne(ctx, whereClause, vals, returningCols, nil)
			if err != nil {
				return err
			}
			if res == nil {
				return nil
			}
			return txQ.Profile.loadRelations(ctx, []*Profile{res}, selects)
		})
	} else {
		res, err = d.queryOne(ctx, whereClause, vals, returningCols, nil)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *ProfileDelegate) runFindFirst(
	ctx context.Context,
	params QueryParams[Profile],
	selects *ProfileSelect,
	omits *ProfileOmit,
) (*Profile, error) {
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
	returningCols := selectProfileCols(selects, omits)

	var res *Profile
	var err error
	if selects.hasAnyRelation() {
		err = d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = txQ.Profile.queryOne(ctx, whereClause, vals, returningCols, params.Skip)
			if err != nil {
				return err
			}
			if res == nil {
				return nil
			}
			return txQ.Profile.loadRelations(ctx, []*Profile{res}, selects)
		})
	} else {
		res, err = d.queryOne(ctx, whereClause, vals, returningCols, params.Skip)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *ProfileDelegate) runFindMany(
	ctx context.Context,
	params QueryParams[Profile],
	selects *ProfileSelect,
	omits *ProfileOmit,
) ([]*Profile, error) {
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
	returningCols := selectProfileCols(selects, omits)

	var results []*Profile
	var err error
	if selects.hasAnyRelation() {
		err = d.client.transaction(ctx, func(txQ *Queries) error {
			var err error
			results, err = txQ.Profile.queryMany(ctx, whereClause, vals, returningCols, params.Take, params.Skip)
			if err != nil {
				return err
			}
			if len(results) == 0 {
				return nil
			}
			return txQ.Profile.loadRelations(ctx, results, selects)
		})
	} else {
		results, err = d.queryMany(ctx, whereClause, vals, returningCols, params.Take, params.Skip)
	}
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (d *ProfileDelegate) queryOne(ctx context.Context, whereClause string, whereVals []any, returningCols []string, skip *int) (*Profile, error) {
	limitOne := 1
	query := buildSelectSQL(d.client, "Profile", returningCols, whereClause, &limitOne, skip)
	stmt, err := d.client.prepare(ctx, query)
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRowContext(ctx, whereVals...)
	var res Profile
	if err := row.Scan(res.ScanFields(returningCols)...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (d *ProfileDelegate) queryMany(ctx context.Context, whereClause string, whereVals []any, returningCols []string, take *int, skip *int) ([]*Profile, error) {
	query := buildSelectSQL(d.client, "Profile", returningCols, whereClause, take, skip)
	stmt, err := d.client.prepare(ctx, query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.QueryContext(ctx, whereVals...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*Profile, 0)
	for rows.Next() {
		var res Profile
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
func (d *ProfileDelegate) loadRelations(ctx context.Context, records []*Profile, selects *ProfileSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.User != nil {
		relationSelects, relationOmits, relationParams := selects.User.GetRelationParams()
		returningCols := selectUserCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: Profile.userId
		allChildren, err := loadRelation(
			ctx, d.client, records,
			directKey(func(p *Profile) string { return p.UserId }),
			"User",
			"id",
			returningCols,
			scanInto(returningCols, (*User).ScanFields),
			directKey(func(c *User) string { return c.Id }),
			setOne(func(p *Profile, c *User) { p.User = c }),
			relationParams,
		)
		if err != nil {
			return fmt.Errorf("loading user: %w", err)
		}
		if err := d.client.User.loadRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
