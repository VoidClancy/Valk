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
type CreateArgs = valk.CategoryCreateArgs
type CreateManyArgs = valk.CategoryCreateManyArgs
type CreateManyAndReturnArgs = valk.CategoryCreateManyAndReturnArgs

type CreateQuery = valk.CategoryCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*valk.Category, error)

type CreateManyQuery = valk.CategoryCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = valk.CategoryCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*valk.Category, error)

type FindUniqueArgs = valk.CategoryFindUniqueArgs
type FindUniqueQuery = valk.CategoryFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*valk.Category, error)

type FindFirstArgs = valk.CategoryFindFirstArgs
type FindFirstQuery = valk.CategoryFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*valk.Category, error)

type FindManyArgs = valk.CategoryFindManyArgs
type FindManyQuery = valk.CategoryFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*valk.Category, error)

type CountArgs = valk.CategoryCountArgs
type CountQuery = valk.CategoryCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

type Extension = valk.CategoryExtension

// ConflictUpdate creates a custom ConflictAction for Category upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *valk.CategoryUpsert)) *valk.ConflictAction {
	return valk.CategoryConflictUpdate(fn)
}
