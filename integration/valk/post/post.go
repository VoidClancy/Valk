package post

import (
	"context"
	"integration/valk"
)

type Select = valk.PostSelect
type Omit = valk.PostOmit
type QueryBuilder = valk.PostQueryBuilder
type CreateBuilder = valk.PostCreateBuilder
type Upsert = valk.PostUpsert
type ConflictBuilder[B any] = valk.PostConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...valk.FieldAssignmentOf[valk.Post]) valk.RecordInput {
	raw := make([]valk.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = valk.FieldAssignment{Col: a.Col, Val: a.Val}
	}
	return valk.RecordInput{Assignments: raw}
}

func And(preds ...valk.PredicateOf[valk.Post]) valk.PredicateOf[valk.Post] {
	return valk.And(preds...)
}

func Or(preds ...valk.PredicateOf[valk.Post]) valk.PredicateOf[valk.Post] {
	return valk.Or(preds...)
}

func Not(pred valk.PredicateOf[valk.Post]) valk.PredicateOf[valk.Post] {
	return valk.Not(pred)
}

var Id = valk.StringUniqueField[valk.Post]{Column: "id"}

var Title = valk.StringField[valk.Post]{Column: "title"}

var Content = valk.OptionalStringField[valk.Post]{StringField: valk.StringField[valk.Post]{Column: "content"}}

var Published = valk.Field[valk.Post, bool]{Column: "published"}

var AuthorId = valk.StringField[valk.Post]{Column: "authorId"}

var Tags = valk.ArrayField[valk.Post, string]{Column: "tags"}

type CreateInput = valk.PostCreate
type Create = valk.PostCreate

type CreateArgs = valk.PostCreateArgs
type CreateManyArgs = valk.PostCreateManyArgs
type CreateManyAndReturnArgs = valk.PostCreateManyAndReturnArgs

type CreateQuery = valk.PostCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*valk.Post, error)

type CreateManyQuery = valk.PostCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = valk.PostCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*valk.Post, error)

type FindUniqueArgs = valk.PostFindUniqueArgs
type FindUniqueQuery = valk.PostFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*valk.Post, error)

type FindFirstArgs = valk.PostFindFirstArgs
type FindFirstQuery = valk.PostFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*valk.Post, error)

type FindManyArgs = valk.PostFindManyArgs
type FindManyQuery = valk.PostFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*valk.Post, error)

type CountArgs = valk.PostCountArgs
type CountQuery = valk.PostCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = valk.PostUpdate

type Extension = valk.PostExtension

// ConflictUpdate creates a custom ConflictAction for Post upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *valk.PostUpsert)) *valk.ConflictAction {
	return valk.PostConflictUpdate(fn)
}
