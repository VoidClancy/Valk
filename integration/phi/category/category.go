package category

import (
	"context"
	"integration/phi"
)

type Select = phi.CategorySelect
type Omit = phi.CategoryOmit
type QueryBuilder = phi.CategoryQueryBuilder
type CreateBuilder = phi.CategoryCreateBuilder
type Upsert = phi.CategoryUpsert
type ConflictBuilder[B any] = phi.CategoryConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...phi.FieldAssignmentOf[phi.Category]) phi.RecordInput {
	raw := make([]phi.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = phi.FieldAssignment(a)
	}
	return phi.RecordInput{Assignments: raw}
}

func And(preds ...phi.PredicateOf[phi.Category]) phi.PredicateOf[phi.Category] {
	return phi.And(preds...)
}

func Or(preds ...phi.PredicateOf[phi.Category]) phi.PredicateOf[phi.Category] {
	return phi.Or(preds...)
}

func Not(pred phi.PredicateOf[phi.Category]) phi.PredicateOf[phi.Category] {
	return phi.Not(pred)
}

var Id = phi.UniqueField[phi.Category, int32]{Column: "id"}

var Name = phi.StringUniqueField[phi.Category]{Column: "name"}

type CreateInput = phi.CategoryCreate
type Create = phi.CategoryCreate

type CreateArgs = phi.CategoryCreateArgs
type CreateManyArgs = phi.CategoryCreateManyArgs
type CreateManyAndReturnArgs = phi.CategoryCreateManyAndReturnArgs

type CreateQuery = phi.CategoryCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*phi.Category, error)

type CreateManyQuery = phi.CategoryCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = phi.CategoryCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*phi.Category, error)

type FindUniqueArgs = phi.CategoryFindUniqueArgs
type FindUniqueQuery = phi.CategoryFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*phi.Category, error)

type FindFirstArgs = phi.CategoryFindFirstArgs
type FindFirstQuery = phi.CategoryFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*phi.Category, error)

type FindManyArgs = phi.CategoryFindManyArgs
type FindManyQuery = phi.CategoryFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*phi.Category, error)

type CountArgs = phi.CategoryCountArgs
type CountQuery = phi.CategoryCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = phi.CategoryUpdate

type Extension = phi.CategoryExtension

// ConflictUpdate creates a custom ConflictAction for Category upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *phi.CategoryUpsert)) *phi.ConflictAction {
	return phi.CategoryConflictUpdate(fn)
}
