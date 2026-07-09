package valk

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"
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
	Id       bool        `json:"id"`
	Textify  bool        `json:"textify"`
	Dummy3   bool        `json:"dummy3"`
	Dummy1   bool        `json:"dummy1"`
	Dummy2   bool        `json:"dummy2"`
	PostId   bool        `json:"postId"`
	AuthorId bool        `json:"authorId"`
	Meta     bool        `json:"meta"`
	Post     *PostSelect `json:"post,omitempty"`
	Author   *UserSelect `json:"author,omitempty"`
}

// CommentOmit specifies which fields to exclude
type CommentOmit struct {
	Id       bool      `json:"id"`
	Textify  bool      `json:"textify"`
	Dummy3   bool      `json:"dummy3"`
	Dummy1   bool      `json:"dummy1"`
	Dummy2   bool      `json:"dummy2"`
	PostId   bool      `json:"postId"`
	AuthorId bool      `json:"authorId"`
	Meta     bool      `json:"meta"`
	Post     *PostOmit `json:"post,omitempty"`
	Author   *UserOmit `json:"author,omitempty"`
}

type CommentDelegate struct {
	client       *Queries
	beforeCreate func(context.Context, *CommentCreate) error
	afterCreate  func(context.Context, *Comment) error
}

func (d *CommentDelegate) BeforeCreate(hook func(context.Context, *CommentCreate) error) {
	d.beforeCreate = hook
}

func (d *CommentDelegate) AfterCreate(hook func(context.Context, *Comment) error) {
	d.afterCreate = hook
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

func (d *CommentDelegate) Create(assignments ...FieldAssignment) *CreateBuilder[Comment, CommentSelect, CommentOmit] {
	return &CreateBuilder[Comment, CommentSelect, CommentOmit]{
		client:      d.client,
		assignments: assignments,
		execFunc:    d.client.executeCommentCreate,
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
				if strings.Contains(v, "\x00") {
					errs.Add("id", v, "safety", "string cannot contain null bytes")
				}
				if !utf8.ValidString(v) {
					errs.Add("id", v, "safety", "string must be valid UTF-8")
				}
			}
		case "textify":
		case "dummy3":
			if v, ok := a.Val.(string); ok {
				if v == "" {
					errs.Add("dummy3", v, "required", "field dummy3 is required")
				}
				if strings.Contains(v, "\x00") {
					errs.Add("dummy3", v, "safety", "string cannot contain null bytes")
				}
				if !utf8.ValidString(v) {
					errs.Add("dummy3", v, "safety", "string must be valid UTF-8")
				}
			}
		case "dummy1":
		case "dummy2":
			if v, ok := a.Val.(string); ok {
				if v == "" {
					errs.Add("dummy2", v, "required", "field dummy2 is required")
				}
				if strings.Contains(v, "\x00") {
					errs.Add("dummy2", v, "safety", "string cannot contain null bytes")
				}
				if !utf8.ValidString(v) {
					errs.Add("dummy2", v, "safety", "string must be valid UTF-8")
				}
			}
		case "postId":
			if v, ok := a.Val.(string); ok {
				if v == "" {
					errs.Add("postId", v, "required", "field postId is required")
				}
				if strings.Contains(v, "\x00") {
					errs.Add("postId", v, "safety", "string cannot contain null bytes")
				}
				if !utf8.ValidString(v) {
					errs.Add("postId", v, "safety", "string must be valid UTF-8")
				}
			}
		case "authorId":
			if v, ok := a.Val.(string); ok {
				if v == "" {
					errs.Add("authorId", v, "required", "field authorId is required")
				}
				if strings.Contains(v, "\x00") {
					errs.Add("authorId", v, "safety", "string cannot contain null bytes")
				}
				if !utf8.ValidString(v) {
					errs.Add("authorId", v, "safety", "string must be valid UTF-8")
				}
			}
		case "meta":
		}
	}
	if !provided["dummy3"] {
		errs.Add("dummy3", "", "required", "field Dummy3 is required")
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

func (q *Queries) executeCommentCreate(ctx context.Context, assignments []FieldAssignment, selects *CommentSelect, omits *CommentOmit) (*Comment, error) {
	input := assignmentsToCommentCreate(assignments)

	if q.Comment.beforeCreate != nil {
		if err := q.Comment.beforeCreate(ctx, &input); err != nil {
			return nil, err
		}
	}

	if err := validateCommentCreate(assignments); err != nil {
		return nil, err
	}

	var cols []string
	var vals []any
	if input.Id != nil {
		cols = append(cols, "id")
		vals = append(vals, *input.Id)
	} else {
		cols = append(cols, "id")
		vals = append(vals, generateCUID())
	}
	cols = append(cols, "textify")
	vals = append(vals, input.Textify)
	cols = append(cols, "dummy3")
	vals = append(vals, input.Dummy3)
	cols = append(cols, "dummy1")
	vals = append(vals, input.Dummy1)
	cols = append(cols, "dummy2")
	vals = append(vals, input.Dummy2)
	cols = append(cols, "postId")
	vals = append(vals, input.PostId)
	cols = append(cols, "authorId")
	vals = append(vals, input.AuthorId)
	if input.Meta != nil {
		cols = append(cols, "meta")
		vals = append(vals, *input.Meta)
	}

	returningCols := q.selectCommentCols(selects, omits)

	scanFunc := func(res *Comment, cols []string) []any {
		return res.ScanFields(cols)
	}

	idCol := "id"

	hasRelations := selects.hasAnyRelation()

	var res *Comment
	var err error
	if hasRelations {
		err = q.transaction(ctx, func(txQ *Queries) error {
			var err error
			res, err = executeInsert(ctx, txQ, "Comment", cols, vals, returningCols, idCol, scanFunc)
			if err != nil {
				return err
			}
			return txQ.loadCommentRelations(ctx, []*Comment{res}, selects)
		})
	} else {
		res, err = executeInsert(ctx, q, "Comment", cols, vals, returningCols, idCol, scanFunc)
	}
	if err != nil {
		return nil, err
	}

	if q.Comment.afterCreate != nil {
		if err := q.Comment.afterCreate(ctx, res); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func commentRecordsToRowMaps(records []RecordInput) []map[string]any {
	rowMaps := make([]map[string]any, len(records))
	for i, rec := range records {
		m := make(map[string]any, len(rec.Assignments))
		for _, a := range rec.Assignments {
			m[a.Col] = a.Val
		}
		if _, ok := m["id"]; !ok {
			m["id"] = generateCUID()
		}
		rowMaps[i] = m
	}
	return rowMaps
}

func (d *CommentDelegate) CreateMany(records ...RecordInput) *CreateManyBuilder[Comment] {
	return &CreateManyBuilder[Comment]{
		client:   d.client,
		records:  records,
		execFunc: d.client.executeCommentCreateMany,
	}
}

func (d *CommentDelegate) CreateManyAndReturn(records ...RecordInput) *CreateManyAndReturnBuilder[Comment, CommentSelect, CommentOmit] {
	return &CreateManyAndReturnBuilder[Comment, CommentSelect, CommentOmit]{
		client:   d.client,
		records:  records,
		execFunc: d.client.executeCommentCreateManyAndReturn,
	}
}

func (q *Queries) executeCommentCreateMany(ctx context.Context, records []RecordInput) (int64, error) {
	if len(records) == 0 {
		return 0, nil
	}
	for i, rec := range records {
		if err := validateCommentCreate(rec.Assignments); err != nil {
			return 0, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
	}

	if q.dialect.SupportsBulkInsert() {
		rowMaps := commentRecordsToRowMaps(records)
		query, vals := buildBulkInsertSQL(q.dialect, "Comment", rowMaps, CommentColOrder, nil)
		res, err := q.exec(ctx, query, vals...)
		if err != nil {
			return 0, err
		}
		return res.RowsAffected()
	}

	var count int64
	err := q.transaction(ctx, func(txQ *Queries) error {
		for _, rec := range records {
			_, err := txQ.executeCommentCreate(ctx, rec.Assignments, nil, nil)
			if err != nil {
				return err
			}
			count++
		}
		return nil
	})
	return count, err
}

func (q *Queries) executeCommentCreateManyAndReturn(ctx context.Context, records []RecordInput, selects *CommentSelect, omits *CommentOmit) ([]*Comment, error) {
	if len(records) == 0 {
		return nil, nil
	}
	for i, rec := range records {
		if err := validateCommentCreate(rec.Assignments); err != nil {
			return nil, fmt.Errorf("validation failed at index %d: %w", i, err)
		}
	}

	hasRelations := selects.hasAnyRelation()
	returningCols := q.selectCommentCols(selects, omits)

	if q.dialect.SupportsBulkInsert() {
		rowMaps := commentRecordsToRowMaps(records)
		query, vals := buildBulkInsertSQL(q.dialect, "Comment", rowMaps, CommentColOrder, returningCols)
		recordsOut := make([]*Comment, 0)
		err := q.transaction(ctx, func(txQ *Queries) error {
			rows, err := txQ.query(ctx, query, vals...)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var record Comment
				if err := rows.Scan(record.ScanFields(returningCols)...); err != nil {
					return err
				}
				recordsOut = append(recordsOut, &record)
			}
			if err := rows.Err(); err != nil {
				return err
			}
			if hasRelations {
				return txQ.loadCommentRelations(ctx, recordsOut, selects)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return recordsOut, nil
	}

	recordsOut := make([]*Comment, 0)
	err := q.transaction(ctx, func(txQ *Queries) error {
		for _, rec := range records {
			res, err := txQ.executeCommentCreate(ctx, rec.Assignments, nil, nil)
			if err != nil {
				return err
			}
			recordsOut = append(recordsOut, res)
		}

		if hasRelations {
			return txQ.loadCommentRelations(ctx, recordsOut, selects)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return recordsOut, nil
}
func (d *CommentDelegate) FindUnique(where UniquePredicate) *FindUniqueBuilder[Comment, CommentSelect, CommentOmit] {
	return &FindUniqueBuilder[Comment, CommentSelect, CommentOmit]{
		client:   d.client,
		where:    where,
		execFunc: d.client.executeCommentFindUnique,
	}
}

func (d *CommentDelegate) FindFirst(preds ...Predicate) *FindFirstBuilder[Comment, CommentSelect, CommentOmit] {
	return &FindFirstBuilder[Comment, CommentSelect, CommentOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeCommentFindFirst,
	}
}

func (d *CommentDelegate) FindMany(preds ...Predicate) *FindManyBuilder[Comment, CommentSelect, CommentOmit] {
	return &FindManyBuilder[Comment, CommentSelect, CommentOmit]{
		client:   d.client,
		where:    preds,
		execFunc: d.client.executeCommentFindMany,
	}
}

func (q *Queries) executeCommentFindUnique(ctx context.Context, where UniquePredicate, selects *CommentSelect, omits *CommentOmit) (*Comment, error) {
	if where == nil {
		return nil, fmt.Errorf("at least one unique field must be set for FindUnique")
	}
	if err := where.Validate(); err != nil {
		return nil, err
	}
	whereClause, vals := CompilePredicates(q.dialect, []Predicate{where})
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
	)
}

func (q *Queries) executeCommentFindFirst(ctx context.Context, where []Predicate, selects *CommentSelect, omits *CommentOmit) (*Comment, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}
	whereClause, vals := CompilePredicates(q.dialect, where)
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
	)
}

func (q *Queries) executeCommentFindMany(ctx context.Context, where []Predicate, selects *CommentSelect, omits *CommentOmit) ([]*Comment, error) {
	for _, p := range where {
		if p != nil {
			if err := p.Validate(); err != nil {
				return nil, err
			}
		}
	}
	whereClause, vals := CompilePredicates(q.dialect, where)
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
	)
}
func (q *Queries) loadCommentRelations(ctx context.Context, records []*Comment, selects *CommentSelect) error {
	if selects == nil || len(records) == 0 {
		return nil
	}
	if selects.Post != nil {
		returningCols := q.selectPostCols(selects.Post, nil, "id")
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
		)
		if err != nil {
			return fmt.Errorf("loading post: %w", err)
		}
		if err := q.loadPostRelations(ctx, allChildren, selects.Post); err != nil {
			return err
		}
	}
	if selects.Author != nil {
		returningCols := q.selectUserCols(selects.Author, nil, "id")
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
		)
		if err != nil {
			return fmt.Errorf("loading author: %w", err)
		}
		if err := q.loadUserRelations(ctx, allChildren, selects.Author); err != nil {
			return err
		}
	}

	return nil
}
