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

var Content = valk.StringField[valk.Post]{Column: "content"}

var Published = valk.Field[valk.Post, bool]{Column: "published"}

var AuthorId = valk.StringField[valk.Post]{Column: "authorId"}

type CreateInput = valk.PostCreate
type CreateQuery = valk.PostCreateQuery
type CreateHook = func(context.Context, *CreateInput, CreateQuery) (*valk.Post, error)

type CreateManyQuery = valk.PostCreateManyQuery
type CreateManyHook = func(context.Context, []*CreateInput, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = valk.PostCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, []*CreateInput, CreateManyAndReturnQuery) ([]*valk.Post, error)

type FindUniqueQuery = valk.PostFindUniqueQuery
type FindUniqueHook = func(context.Context, valk.UniquePredicate[valk.Post], []valk.PredicateOf[valk.Post], *valk.PostSelect, *valk.PostOmit, FindUniqueQuery) (*valk.Post, error)

type FindFirstQuery = valk.PostFindFirstQuery
type FindFirstHook = func(context.Context, valk.QueryParams[valk.Post], *valk.PostSelect, *valk.PostOmit, FindFirstQuery) (*valk.Post, error)

type FindManyQuery = valk.PostFindManyQuery
type FindManyHook = func(context.Context, valk.QueryParams[valk.Post], *valk.PostSelect, *valk.PostOmit, FindManyQuery) ([]*valk.Post, error)

type Extension = valk.PostExtension
