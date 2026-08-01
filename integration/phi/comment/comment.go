package comment

import (
	"context"
	"encoding/json"
	"integration/phi"
)

type Select = phi.CommentSelect
type Omit = phi.CommentOmit
type QueryBuilder = phi.CommentQueryBuilder
type CreateBuilder = phi.CommentCreateBuilder
type Upsert = phi.CommentUpsert
type ConflictBuilder[B any] = phi.CommentConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...phi.FieldAssignmentOf[phi.Comment]) phi.RecordInput {
	raw := make([]phi.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = phi.FieldAssignment(a)
	}
	return phi.RecordInput{Assignments: raw}
}

func And(preds ...phi.PredicateOf[phi.Comment]) phi.PredicateOf[phi.Comment] {
	return phi.And(preds...)
}

func Or(preds ...phi.PredicateOf[phi.Comment]) phi.PredicateOf[phi.Comment] {
	return phi.Or(preds...)
}

func Not(pred phi.PredicateOf[phi.Comment]) phi.PredicateOf[phi.Comment] {
	return phi.Not(pred)
}

var Id = phi.StringUniqueField[phi.Comment]{Column: "id"}

var Textify = phi.Field[phi.Comment, int32]{Column: "textify"}

var Dummy3 = phi.StringField[phi.Comment]{Column: "dummy3"}

var Dummy1 = phi.Field[phi.Comment, int32]{Column: "dummy1"}

var Dummy2 = phi.StringField[phi.Comment]{Column: "dummy2"}

var PostId = phi.StringField[phi.Comment]{Column: "postId"}

var AuthorId = phi.StringField[phi.Comment]{Column: "authorId"}

var Meta = phi.OptionalField[phi.Comment, json.RawMessage]{Field: phi.Field[phi.Comment, json.RawMessage]{Column: "meta"}}

type CreateInput = phi.CommentCreate
type Create = phi.CommentCreate

type CreateArgs = phi.CommentCreateArgs
type CreateManyArgs = phi.CommentCreateManyArgs
type CreateManyAndReturnArgs = phi.CommentCreateManyAndReturnArgs

type CreateQuery = phi.CommentCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*phi.Comment, error)

type CreateManyQuery = phi.CommentCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = phi.CommentCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*phi.Comment, error)

type FindUniqueArgs = phi.CommentFindUniqueArgs
type FindUniqueQuery = phi.CommentFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*phi.Comment, error)

type FindFirstArgs = phi.CommentFindFirstArgs
type FindFirstQuery = phi.CommentFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*phi.Comment, error)

type FindManyArgs = phi.CommentFindManyArgs
type FindManyQuery = phi.CommentFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*phi.Comment, error)

type CountArgs = phi.CommentCountArgs
type CountQuery = phi.CommentCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = phi.CommentUpdate

type Extension = phi.CommentExtension

// ConflictUpdate creates a custom ConflictAction for Comment upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *phi.CommentUpsert)) *phi.ConflictAction {
	return phi.CommentConflictUpdate(fn)
}
