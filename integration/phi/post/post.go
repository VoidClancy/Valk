package post

import (
	"context"
	"integration/phi"
)

type Select = phi.PostSelect
type Omit = phi.PostOmit
type QueryBuilder = phi.PostQueryBuilder
type CreateBuilder = phi.PostCreateBuilder
type Upsert = phi.PostUpsert
type ConflictBuilder[B any] = phi.PostConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...phi.FieldAssignmentOf[phi.Post]) phi.RecordInput {
	raw := make([]phi.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = phi.FieldAssignment(a)
	}
	return phi.RecordInput{Assignments: raw}
}

func And(preds ...phi.PredicateOf[phi.Post]) phi.PredicateOf[phi.Post] {
	return phi.And(preds...)
}

func Or(preds ...phi.PredicateOf[phi.Post]) phi.PredicateOf[phi.Post] {
	return phi.Or(preds...)
}

func Not(pred phi.PredicateOf[phi.Post]) phi.PredicateOf[phi.Post] {
	return phi.Not(pred)
}

var Id = phi.StringUniqueField[phi.Post]{Column: "id"}

var Title = phi.StringField[phi.Post]{Column: "title"}

var Content = phi.OptionalStringField[phi.Post]{StringField: phi.StringField[phi.Post]{Column: "content"}}

var Published = phi.Field[phi.Post, bool]{Column: "published"}

var AuthorId = phi.StringField[phi.Post]{Column: "authorId"}

var Tags = phi.ArrayField[phi.Post, string]{Column: "tags"}

type CreateInput = phi.PostCreate
type Create = phi.PostCreate

type CreateArgs = phi.PostCreateArgs
type CreateManyArgs = phi.PostCreateManyArgs
type CreateManyAndReturnArgs = phi.PostCreateManyAndReturnArgs

type CreateQuery = phi.PostCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*phi.Post, error)

type CreateManyQuery = phi.PostCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = phi.PostCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*phi.Post, error)

type FindUniqueArgs = phi.PostFindUniqueArgs
type FindUniqueQuery = phi.PostFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*phi.Post, error)

type FindFirstArgs = phi.PostFindFirstArgs
type FindFirstQuery = phi.PostFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*phi.Post, error)

type FindManyArgs = phi.PostFindManyArgs
type FindManyQuery = phi.PostFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*phi.Post, error)

type CountArgs = phi.PostCountArgs
type CountQuery = phi.PostCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = phi.PostUpdate

type Extension = phi.PostExtension

// ConflictUpdate creates a custom ConflictAction for Post upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *phi.PostUpsert)) *phi.ConflictAction {
	return phi.PostConflictUpdate(fn)
}
