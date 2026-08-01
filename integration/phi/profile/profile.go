package profile

import (
	"context"
	"integration/phi"
	"time"
)

type Select = phi.ProfileSelect
type Omit = phi.ProfileOmit
type QueryBuilder = phi.ProfileQueryBuilder
type CreateBuilder = phi.ProfileCreateBuilder
type Upsert = phi.ProfileUpsert
type ConflictBuilder[B any] = phi.ProfileConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...phi.FieldAssignmentOf[phi.Profile]) phi.RecordInput {
	raw := make([]phi.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = phi.FieldAssignment(a)
	}
	return phi.RecordInput{Assignments: raw}
}

func And(preds ...phi.PredicateOf[phi.Profile]) phi.PredicateOf[phi.Profile] {
	return phi.And(preds...)
}

func Or(preds ...phi.PredicateOf[phi.Profile]) phi.PredicateOf[phi.Profile] {
	return phi.Or(preds...)
}

func Not(pred phi.PredicateOf[phi.Profile]) phi.PredicateOf[phi.Profile] {
	return phi.Not(pred)
}

var Id = phi.StringUniqueField[phi.Profile]{Column: "id"}

var Bio = phi.OptionalStringField[phi.Profile]{StringField: phi.StringField[phi.Profile]{Column: "bio"}}

var UserId = phi.StringUniqueField[phi.Profile]{Column: "userId"}

var CreatedAt = phi.Field[phi.Profile, time.Time]{Column: "createdAt"}

type CreateInput = phi.ProfileCreate
type Create = phi.ProfileCreate

type CreateArgs = phi.ProfileCreateArgs
type CreateManyArgs = phi.ProfileCreateManyArgs
type CreateManyAndReturnArgs = phi.ProfileCreateManyAndReturnArgs

type CreateQuery = phi.ProfileCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*phi.Profile, error)

type CreateManyQuery = phi.ProfileCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = phi.ProfileCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*phi.Profile, error)

type FindUniqueArgs = phi.ProfileFindUniqueArgs
type FindUniqueQuery = phi.ProfileFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*phi.Profile, error)

type FindFirstArgs = phi.ProfileFindFirstArgs
type FindFirstQuery = phi.ProfileFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*phi.Profile, error)

type FindManyArgs = phi.ProfileFindManyArgs
type FindManyQuery = phi.ProfileFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*phi.Profile, error)

type CountArgs = phi.ProfileCountArgs
type CountQuery = phi.ProfileCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = phi.ProfileUpdate

type Extension = phi.ProfileExtension

// ConflictUpdate creates a custom ConflictAction for Profile upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *phi.ProfileUpsert)) *phi.ConflictAction {
	return phi.ProfileConflictUpdate(fn)
}
