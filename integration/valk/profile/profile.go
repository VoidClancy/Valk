package profile

import (
	"context"
	"integration/valk"
	"time"
)

type Select = valk.ProfileSelect
type Omit = valk.ProfileOmit
type QueryBuilder = valk.ProfileQueryBuilder
type CreateBuilder = valk.ProfileCreateBuilder
type Upsert = valk.ProfileUpsert
type ConflictBuilder[B any] = valk.ProfileConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...valk.FieldAssignment) valk.RecordInput {
	return valk.RecordInput{Assignments: assignments}
}

func And(preds ...valk.PredicateOf[valk.Profile]) valk.PredicateOf[valk.Profile] {
	return valk.And(preds...)
}

func Or(preds ...valk.PredicateOf[valk.Profile]) valk.PredicateOf[valk.Profile] {
	return valk.Or(preds...)
}

func Not(pred valk.PredicateOf[valk.Profile]) valk.PredicateOf[valk.Profile] {
	return valk.Not(pred)
}

var Id = valk.StringUniqueField[valk.Profile]{Column: "id"}

var Bio = valk.StringField[valk.Profile]{Column: "bio"}

var UserId = valk.StringUniqueField[valk.Profile]{Column: "userId"}

var CreatedAt = valk.Field[valk.Profile, time.Time]{Column: "createdAt"}

type CreateInput = valk.ProfileCreate
type CreateQuery = valk.ProfileCreateQuery
type CreateHook = func(context.Context, *CreateInput, CreateQuery) (*valk.Profile, error)

type CreateManyQuery = valk.ProfileCreateManyQuery
type CreateManyHook = func(context.Context, []*CreateInput, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = valk.ProfileCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, []*CreateInput, CreateManyAndReturnQuery) ([]*valk.Profile, error)

type FindUniqueQuery = valk.ProfileFindUniqueQuery
type FindUniqueHook = func(context.Context, valk.UniquePredicate[valk.Profile], []valk.PredicateOf[valk.Profile], *valk.ProfileSelect, *valk.ProfileOmit, FindUniqueQuery) (*valk.Profile, error)

type FindFirstQuery = valk.ProfileFindFirstQuery
type FindFirstHook = func(context.Context, valk.QueryParams[valk.Profile], *valk.ProfileSelect, *valk.ProfileOmit, FindFirstQuery) (*valk.Profile, error)

type FindManyQuery = valk.ProfileFindManyQuery
type FindManyHook = func(context.Context, valk.QueryParams[valk.Profile], *valk.ProfileSelect, *valk.ProfileOmit, FindManyQuery) ([]*valk.Profile, error)

type Extension = valk.ProfileExtension
