package comment

import (
	"context"
	"encoding/json"
	"integration/valk"
)

type Select = valk.CommentSelect
type Omit = valk.CommentOmit
type QueryBuilder = valk.CommentQueryBuilder
type CreateBuilder = valk.CommentCreateBuilder
type Upsert = valk.CommentUpsert
type ConflictBuilder[B any] = valk.CommentConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...valk.FieldAssignmentOf[valk.Comment]) valk.RecordInput {
	raw := make([]valk.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = valk.FieldAssignment{Col: a.Col, Val: a.Val}
	}
	return valk.RecordInput{Assignments: raw}
}

func And(preds ...valk.PredicateOf[valk.Comment]) valk.PredicateOf[valk.Comment] {
	return valk.And(preds...)
}

func Or(preds ...valk.PredicateOf[valk.Comment]) valk.PredicateOf[valk.Comment] {
	return valk.Or(preds...)
}

func Not(pred valk.PredicateOf[valk.Comment]) valk.PredicateOf[valk.Comment] {
	return valk.Not(pred)
}

var Id = valk.StringUniqueField[valk.Comment]{Column: "id"}

var Textify = valk.Field[valk.Comment, int32]{Column: "textify"}

var Dummy3 = valk.StringField[valk.Comment]{Column: "dummy3"}

var Dummy1 = valk.Field[valk.Comment, int32]{Column: "dummy1"}

var Dummy2 = valk.StringField[valk.Comment]{Column: "dummy2"}

var PostId = valk.StringField[valk.Comment]{Column: "postId"}

var AuthorId = valk.StringField[valk.Comment]{Column: "authorId"}

var Meta = valk.Field[valk.Comment, json.RawMessage]{Column: "meta"}

type CreateInput = valk.CommentCreate
type CreateArgs = valk.CommentCreateArgs
type CreateManyArgs = valk.CommentCreateManyArgs
type CreateManyAndReturnArgs = valk.CommentCreateManyAndReturnArgs

type CreateQuery = valk.CommentCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*valk.Comment, error)

type CreateManyQuery = valk.CommentCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = valk.CommentCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*valk.Comment, error)

type FindUniqueArgs = valk.CommentFindUniqueArgs
type FindUniqueQuery = valk.CommentFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*valk.Comment, error)

type FindFirstArgs = valk.CommentFindFirstArgs
type FindFirstQuery = valk.CommentFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*valk.Comment, error)

type FindManyArgs = valk.CommentFindManyArgs
type FindManyQuery = valk.CommentFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*valk.Comment, error)

type CountArgs = valk.CommentCountArgs
type CountQuery = valk.CommentCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

type Extension = valk.CommentExtension

// ConflictUpdate creates a custom ConflictAction for Comment upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *valk.CommentUpsert)) *valk.ConflictAction {
	return valk.CommentConflictUpdate(fn)
}
