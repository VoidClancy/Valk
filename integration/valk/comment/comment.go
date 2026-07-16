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

func Record(assignments ...valk.FieldAssignment) valk.RecordInput {
	return valk.RecordInput{Assignments: assignments}
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
type CreateQuery = valk.CommentCreateQuery
type CreateHook = func(context.Context, *CreateInput, CreateQuery) (*valk.Comment, error)

type CreateManyQuery = valk.CommentCreateManyQuery
type CreateManyHook = func(context.Context, []*CreateInput, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = valk.CommentCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, []*CreateInput, CreateManyAndReturnQuery) ([]*valk.Comment, error)

type FindUniqueQuery = valk.CommentFindUniqueQuery
type FindUniqueHook = func(context.Context, valk.UniquePredicate[valk.Comment], []valk.PredicateOf[valk.Comment], *valk.CommentSelect, *valk.CommentOmit, FindUniqueQuery) (*valk.Comment, error)

type FindFirstQuery = valk.CommentFindFirstQuery
type FindFirstHook = func(context.Context, valk.QueryParams[valk.Comment], *valk.CommentSelect, *valk.CommentOmit, FindFirstQuery) (*valk.Comment, error)

type FindManyQuery = valk.CommentFindManyQuery
type FindManyHook = func(context.Context, valk.QueryParams[valk.Comment], *valk.CommentSelect, *valk.CommentOmit, FindManyQuery) ([]*valk.Comment, error)

type Extension = valk.CommentExtension
