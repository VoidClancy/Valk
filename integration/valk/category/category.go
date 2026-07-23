package category

import (
	"context"
	"integration/valk"
)

type Select = valk.CategorySelect
type Omit = valk.CategoryOmit
type QueryBuilder = valk.CategoryQueryBuilder
type CreateBuilder = valk.CategoryCreateBuilder
type Upsert = valk.CategoryUpsert
type ConflictBuilder[B any] = valk.CategoryConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...valk.FieldAssignmentOf[valk.Category]) valk.RecordInput {
	raw := make([]valk.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = valk.FieldAssignment{Col: a.Col, Val: a.Val}
	}
	return valk.RecordInput{Assignments: raw}
}

func And(preds ...valk.PredicateOf[valk.Category]) valk.PredicateOf[valk.Category] {
	return valk.And(preds...)
}

func Or(preds ...valk.PredicateOf[valk.Category]) valk.PredicateOf[valk.Category] {
	return valk.Or(preds...)
}

func Not(pred valk.PredicateOf[valk.Category]) valk.PredicateOf[valk.Category] {
	return valk.Not(pred)
}

var Id = valk.UniqueField[valk.Category, int32]{Column: "id"}

var Name = valk.StringUniqueField[valk.Category]{Column: "name"}

type CreateInput = valk.CategoryCreate
type CreateQuery = valk.CategoryCreateQuery
type CreateHook = func(context.Context, *CreateInput, CreateQuery) (*valk.Category, error)

type CreateManyQuery = valk.CategoryCreateManyQuery
type CreateManyHook = func(context.Context, []*CreateInput, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = valk.CategoryCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, []*CreateInput, CreateManyAndReturnQuery) ([]*valk.Category, error)

type FindUniqueQuery = valk.CategoryFindUniqueQuery
type FindUniqueHook = func(context.Context, valk.UniquePredicate[valk.Category], []valk.PredicateOf[valk.Category], *valk.CategorySelect, *valk.CategoryOmit, FindUniqueQuery) (*valk.Category, error)

type FindFirstQuery = valk.CategoryFindFirstQuery
type FindFirstHook = func(context.Context, valk.QueryParams[valk.Category], *valk.CategorySelect, *valk.CategoryOmit, FindFirstQuery) (*valk.Category, error)

type FindManyQuery = valk.CategoryFindManyQuery
type FindManyHook = func(context.Context, valk.QueryParams[valk.Category], *valk.CategorySelect, *valk.CategoryOmit, FindManyQuery) ([]*valk.Category, error)

type Extension = valk.CategoryExtension
