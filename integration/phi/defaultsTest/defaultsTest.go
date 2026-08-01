package defaultsTest

import (
	"context"
	"integration/phi"
	"time"
)

type Select = phi.DefaultsTestSelect
type Omit = phi.DefaultsTestOmit
type QueryBuilder = phi.DefaultsTestQueryBuilder
type CreateBuilder = phi.DefaultsTestCreateBuilder
type Upsert = phi.DefaultsTestUpsert
type ConflictBuilder[B any] = phi.DefaultsTestConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...phi.FieldAssignmentOf[phi.DefaultsTest]) phi.RecordInput {
	raw := make([]phi.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = phi.FieldAssignment(a)
	}
	return phi.RecordInput{Assignments: raw}
}

func And(preds ...phi.PredicateOf[phi.DefaultsTest]) phi.PredicateOf[phi.DefaultsTest] {
	return phi.And(preds...)
}

func Or(preds ...phi.PredicateOf[phi.DefaultsTest]) phi.PredicateOf[phi.DefaultsTest] {
	return phi.Or(preds...)
}

func Not(pred phi.PredicateOf[phi.DefaultsTest]) phi.PredicateOf[phi.DefaultsTest] {
	return phi.Not(pred)
}

var Uuid4 = phi.StringUniqueField[phi.DefaultsTest]{Column: "uuid4"}

var Uuid7 = phi.StringField[phi.DefaultsTest]{Column: "uuid7"}

var UuidNoArgs = phi.StringField[phi.DefaultsTest]{Column: "uuidNoArgs"}

var Cuid1 = phi.StringField[phi.DefaultsTest]{Column: "cuid1"}

var Cuid2 = phi.StringField[phi.DefaultsTest]{Column: "cuid2"}

var CuidNoArgs = phi.StringField[phi.DefaultsTest]{Column: "cuidNoArgs"}

var Ulid = phi.StringField[phi.DefaultsTest]{Column: "ulid"}

var Nanoid = phi.StringField[phi.DefaultsTest]{Column: "nanoid"}

var Now = phi.Field[phi.DefaultsTest, time.Time]{Column: "now"}

type CreateInput = phi.DefaultsTestCreate
type Create = phi.DefaultsTestCreate

type CreateArgs = phi.DefaultsTestCreateArgs
type CreateManyArgs = phi.DefaultsTestCreateManyArgs
type CreateManyAndReturnArgs = phi.DefaultsTestCreateManyAndReturnArgs

type CreateQuery = phi.DefaultsTestCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*phi.DefaultsTest, error)

type CreateManyQuery = phi.DefaultsTestCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = phi.DefaultsTestCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*phi.DefaultsTest, error)

type FindUniqueArgs = phi.DefaultsTestFindUniqueArgs
type FindUniqueQuery = phi.DefaultsTestFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*phi.DefaultsTest, error)

type FindFirstArgs = phi.DefaultsTestFindFirstArgs
type FindFirstQuery = phi.DefaultsTestFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*phi.DefaultsTest, error)

type FindManyArgs = phi.DefaultsTestFindManyArgs
type FindManyQuery = phi.DefaultsTestFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*phi.DefaultsTest, error)

type CountArgs = phi.DefaultsTestCountArgs
type CountQuery = phi.DefaultsTestCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = phi.DefaultsTestUpdate

type Extension = phi.DefaultsTestExtension

// ConflictUpdate creates a custom ConflictAction for DefaultsTest upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *phi.DefaultsTestUpsert)) *phi.ConflictAction {
	return phi.DefaultsTestConflictUpdate(fn)
}
