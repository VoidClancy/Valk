package valk

import (
	"context"
	"fmt"
	"slices"
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

// ProfileSelect specifies which fields to include
type ProfileSelect struct {
	Id        bool            `json:"id"`
	Bio       bool            `json:"bio"`
	UserId    bool            `json:"userId"`
	CreatedAt bool            `json:"createdAt"`
	User      UserSelectQuery `json:"user,omitempty"`
}

// ProfileOmit specifies which fields to exclude
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

// ProfileQueryBuilder builds a query for the relation Profile
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

func (q *Queries) selectProfileCols(selects *ProfileSelect, omits *ProfileOmit, forceCols ...string) []string {
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

var ProfileColOrder = []string{
	"id",
	"bio",
	"userId",
	"createdAt",
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
			client:      d.client,
			assignments: assignments,
			execFunc:    d.client.executeProfileCreate,
		},
	}
}

func validateProfileCreate(assignments []FieldAssignment) error {
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
		case "bio":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "bio", v, false, 0, false, false)
			} else {
				errs.Add("bio", a.Val, "type", "field bio must be of type string")
			}
		case "userId":
			if v, ok := a.Val.(string); ok {
				ValidateString(errs, "userId", v, true, 0, false, false)
			} else {
				errs.Add("userId", a.Val, "type", "field userId must be of type string")
			}
		case "createdAt":
			if _, ok := a.Val.(time.Time); !ok {
				errs.Add("createdAt", a.Val, "type", "field createdAt must be of type time.Time")
			}
		}
	}
	if !provided["userId"] {
		errs.Add("userId", "", "required", "field UserId is required")
	}

	if errs.HasErrors() {
		return *errs
	}
	return nil
}

func assignmentsToProfileCreate(assignments []FieldAssignment) ProfileCreate {
	var input ProfileCreate
	for _, a := range assignments {
		switch a.Col {
		case "id":
			if v, ok := a.Val.(string); ok {
				input.Id = &v
			}
		case "bio":
			if v, ok := a.Val.(string); ok {
				input.Bio = &v
			}
		case "userId":
			if v, ok := a.Val.(string); ok {
				input.UserId = v
			}
		case "createdAt":
			if v, ok := a.Val.(time.Time); ok {
				input.CreatedAt = &v
			}
		}
	}
	return input
}

func (s *ProfileCreate) ToRowMap() map[string]any {
	m := make(map[string]any, 4)
	if s.Id != nil {
		m["id"] = *s.Id
	} else {
		m["id"] = generateCUID()
	}
	if s.Bio != nil {
		m["bio"] = *s.Bio
	}
	m["userId"] = s.UserId
	if s.CreatedAt != nil {
		m["createdAt"] = *s.CreatedAt
	} else {
		m["createdAt"] = time.Now()
	}
	return m
}

func (q *Queries) executeProfileCreate(ctx context.Context, assignments []FieldAssignment, selects *ProfileSelect, omits *ProfileOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (*Profile, error) {
	input := assignmentsToProfileCreate(assignments)

	curr := func(c context.Context, args *ProfileCreate) (*Profile, error) {
		if err := validateProfileCreate(assignments); err != nil {
			return nil, err
		}

		rowMap := args.ToRowMap()
		cols, vals := mapToColsVals(rowMap, ProfileColOrder)

		returningCols := q.selectProfileCols(selects, omits)

		scanFunc := func(res *Profile, cols []string) []any {
			return res.ScanFields(cols)
		}

		pkCols := []string{
			"id",
		}

		hasRelations := selects.hasAnyRelation()

		var res *Profile
		var err error
		if hasRelations {
			err = q.transaction(c, func(txQ *Queries) error {
				var err error
				res, err = executeInsert(c, txQ, "Profile", cols, vals, returningCols, pkCols, scanFunc, conflictTarget, conflictAction)
				if err != nil {
					return err
				}
				return txQ.loadProfileRelations(c, []*Profile{res}, selects)
			})
		} else {
			res, err = executeInsert(c, q, "Profile", cols, vals, returningCols, pkCols, scanFunc, conflictTarget, conflictAction)
		}
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	for _, ext := range slices.Backward(q.Profile.extensions) {
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
			client:   d.client,
			records:  records,
			execFunc: d.client.executeProfileCreateMany,
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
			client:   d.client,
			records:  records,
			execFunc: d.client.executeProfileCreateManyAndReturn,
		},
	}
}

func (q *Queries) executeProfileCreateMany(ctx context.Context, records []RecordInput, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) (int64, error) {
	inputs := make([]*ProfileCreate, len(records))
	for i, rec := range records {
		if err := validateProfileCreate(rec.Assignments); err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		input := assignmentsToProfileCreate(rec.Assignments)
		inputs[i] = &input
	}

	curr := func(c context.Context, args []*ProfileCreate) (int64, error) {
		rowMaps := make([]map[string]any, len(args))
		for i, input := range args {
			rowMaps[i] = input.ToRowMap()
		}

		pkCols := []string{
			"id",
		}

		return executeCreateMany(c, q, rowMaps, "Profile", ProfileColOrder, pkCols, conflictTarget, conflictAction)
	}

	for _, ext := range slices.Backward(q.Profile.extensions) {
		if ext.CreateMany != nil {
			next, hook := curr, ext.CreateMany
			curr = func(c context.Context, inputs []*ProfileCreate) (int64, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
}

func (q *Queries) executeProfileCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *ProfileSelect, omits *ProfileOmit, conflictTarget UniqueConstraintTarget, conflictAction *ConflictAction) ([]*Profile, error) {
	inputs := make([]*ProfileCreate, len(records))
	for i, rec := range records {
		if err := validateProfileCreate(rec.Assignments); err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
		input := assignmentsToProfileCreate(rec.Assignments)
		inputs[i] = &input
	}

	curr := func(c context.Context, args []*ProfileCreate) ([]*Profile, error) {
		rowMaps := make([]map[string]any, len(args))
		for i, input := range args {
			rowMaps[i] = input.ToRowMap()
		}

		pkCols := []string{
			"id",
		}

		return executeCreateManyAndReturn(c, q, rowMaps, "Profile", ProfileColOrder, selects, omits,
			q.selectProfileCols,
			q.loadProfileRelations,
			(*Profile).ScanFields,
			(*ProfileSelect).hasAnyRelation,
			pkCols,
			conflictTarget,
			conflictAction,
		)
	}

	for _, ext := range slices.Backward(q.Profile.extensions) {
		if ext.CreateManyAndReturn != nil {
			next, hook := curr, ext.CreateManyAndReturn
			curr = func(c context.Context, inputs []*ProfileCreate) ([]*Profile, error) {
				return hook(c, inputs, next)
			}
		}
	}

	return curr(ctx, inputs)
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
		client:     d.client,
		where:      where,
		additional: additional,
		execFunc:   d.client.executeProfileFindUnique,
	}
}

func (d *ProfileDelegate) FindFirst(preds ...PredicateOf[Profile]) *FindFirstBuilder[Profile, ProfileSelect, ProfileOmit] {
	return &FindFirstBuilder[Profile, ProfileSelect, ProfileOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeProfileFindFirst,
	}
}

func (d *ProfileDelegate) FindMany(preds ...PredicateOf[Profile]) *FindManyBuilder[Profile, ProfileSelect, ProfileOmit] {
	return &FindManyBuilder[Profile, ProfileSelect, ProfileOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeProfileFindMany,
	}
}

func (q *Queries) executeProfileFindUnique(ctx context.Context, where UniquePredicate[Profile], additional []PredicateOf[Profile], selects *ProfileSelect, omits *ProfileOmit) (*Profile, error) {
	curr := func(c context.Context, w UniquePredicate[Profile], add []PredicateOf[Profile], sel *ProfileSelect, o *ProfileOmit) (*Profile, error) {
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
		allPreds := append([]PredicateOf[Profile]{w}, add...)
		whereClause, vals := CompilePredicates(q.dialect, allPreds)
		if whereClause != "" {
			whereClause = " WHERE " + whereClause
		}
		returningCols := q.selectProfileCols(sel, o)
		return executeSingleWithRelations(c, q, "Profile", whereClause, vals, returningCols,
			func(res *Profile, cols []string) []any { return res.ScanFields(cols) },
			sel.hasAnyRelation(),
			func(ctx context.Context, txQ *Queries, results []*Profile) error {
				return txQ.loadProfileRelations(ctx, results, sel)
			},
			nil,
		)
	}

	for _, ext := range slices.Backward(q.Profile.extensions) {
		if ext.FindUnique != nil {
			next, hook := curr, ext.FindUnique
			curr = func(c context.Context, w UniquePredicate[Profile], add []PredicateOf[Profile], sel *ProfileSelect, o *ProfileOmit) (*Profile, error) {
				return hook(c, w, add, sel, o, next)
			}
		}
	}

	return curr(ctx, where, additional, selects, omits)
}

func (q *Queries) executeProfileFindFirst(
	ctx context.Context,
	params QueryParams[Profile],
	selects *ProfileSelect,
	omits *ProfileOmit,
) (*Profile, error) {
	curr := func(c context.Context, p QueryParams[Profile], sel *ProfileSelect, o *ProfileOmit) (*Profile, error) {
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
		returningCols := q.selectProfileCols(sel, o)
		return executeSingleWithRelations(c, q, "Profile", whereClause, vals, returningCols,
			func(res *Profile, cols []string) []any { return res.ScanFields(cols) },
			sel.hasAnyRelation(),
			func(ctx context.Context, txQ *Queries, results []*Profile) error {
				return txQ.loadProfileRelations(ctx, results, sel)
			},
			p.Skip,
		)
	}

	for _, ext := range slices.Backward(q.Profile.extensions) {
		if ext.FindFirst != nil {
			next, hook := curr, ext.FindFirst
			curr = func(c context.Context, p QueryParams[Profile], sel *ProfileSelect, o *ProfileOmit) (*Profile, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}

func (q *Queries) executeProfileFindMany(
	ctx context.Context,
	params QueryParams[Profile],
	selects *ProfileSelect,
	omits *ProfileOmit,
) ([]*Profile, error) {
	curr := func(c context.Context, p QueryParams[Profile], sel *ProfileSelect, o *ProfileOmit) ([]*Profile, error) {
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
		returningCols := q.selectProfileCols(sel, o)
		return executeManyWithRelations(c, q, "Profile", whereClause, vals, returningCols,
			func(res *Profile, cols []string) []any { return res.ScanFields(cols) },
			sel.hasAnyRelation(),
			func(ctx context.Context, txQ *Queries, results []*Profile) error {
				return txQ.loadProfileRelations(ctx, results, sel)
			},
			p.Take,
			p.Skip,
		)
	}

	for _, ext := range slices.Backward(q.Profile.extensions) {
		if ext.FindMany != nil {
			next, hook := curr, ext.FindMany
			curr = func(c context.Context, p QueryParams[Profile], sel *ProfileSelect, o *ProfileOmit) ([]*Profile, error) {
				return hook(c, p, sel, o, next)
			}
		}
	}

	return curr(ctx, params, selects, omits)
}
func (q *Queries) loadProfileRelations(ctx context.Context, records []*Profile, selects *ProfileSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.User != nil {
		relationSelects, relationOmits, relationParams := selects.User.GetRelationParams()
		returningCols := q.selectUserCols(relationSelects, relationOmits, "id")
		// Current model holds the FK: Profile.userId
		allChildren, err := loadRelation(
			ctx, q, records,
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
		if err := q.loadUserRelations(ctx, allChildren, relationSelects); err != nil {
			return err
		}
	}

	return nil
}
