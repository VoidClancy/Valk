package defaultsTest

import (
	"context"
	"integration/valk"
	"time"
)

type Select = valk.DefaultsTestSelect
type Omit = valk.DefaultsTestOmit
type QueryBuilder = valk.DefaultsTestQueryBuilder
type CreateBuilder = valk.DefaultsTestCreateBuilder
type Upsert = valk.DefaultsTestUpsert
type ConflictBuilder[B any] = valk.DefaultsTestConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...valk.FieldAssignment) valk.RecordInput {
	return valk.RecordInput{Assignments: assignments}
}

func And(preds ...valk.PredicateOf[valk.DefaultsTest]) valk.PredicateOf[valk.DefaultsTest] {
	return valk.And(preds...)
}

func Or(preds ...valk.PredicateOf[valk.DefaultsTest]) valk.PredicateOf[valk.DefaultsTest] {
	return valk.Or(preds...)
}

func Not(pred valk.PredicateOf[valk.DefaultsTest]) valk.PredicateOf[valk.DefaultsTest] {
	return valk.Not(pred)
}

var Uuid4 = valk.StringUniqueField[valk.DefaultsTest]{Column: "uuid4"}

var Uuid7 = valk.StringField[valk.DefaultsTest]{Column: "uuid7"}

var UuidNoArgs = valk.StringField[valk.DefaultsTest]{Column: "uuidNoArgs"}

var Cuid1 = valk.StringField[valk.DefaultsTest]{Column: "cuid1"}

var Cuid2 = valk.StringField[valk.DefaultsTest]{Column: "cuid2"}

var CuidNoArgs = valk.StringField[valk.DefaultsTest]{Column: "cuidNoArgs"}

var Ulid = valk.StringField[valk.DefaultsTest]{Column: "ulid"}

var Nanoid = valk.StringField[valk.DefaultsTest]{Column: "nanoid"}

var Now = valk.Field[valk.DefaultsTest, time.Time]{Column: "now"}

type CreateInput = valk.DefaultsTestCreate
type CreateQuery = valk.DefaultsTestCreateQuery
type CreateHook = func(context.Context, *CreateInput, CreateQuery) (*valk.DefaultsTest, error)

type CreateManyQuery = valk.DefaultsTestCreateManyQuery
type CreateManyHook = func(context.Context, []*CreateInput, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = valk.DefaultsTestCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, []*CreateInput, CreateManyAndReturnQuery) ([]*valk.DefaultsTest, error)

type FindUniqueQuery = valk.DefaultsTestFindUniqueQuery
type FindUniqueHook = func(context.Context, valk.UniquePredicate[valk.DefaultsTest], []valk.PredicateOf[valk.DefaultsTest], *valk.DefaultsTestSelect, *valk.DefaultsTestOmit, FindUniqueQuery) (*valk.DefaultsTest, error)

type FindFirstQuery = valk.DefaultsTestFindFirstQuery
type FindFirstHook = func(context.Context, valk.QueryParams[valk.DefaultsTest], *valk.DefaultsTestSelect, *valk.DefaultsTestOmit, FindFirstQuery) (*valk.DefaultsTest, error)

type FindManyQuery = valk.DefaultsTestFindManyQuery
type FindManyHook = func(context.Context, valk.QueryParams[valk.DefaultsTest], *valk.DefaultsTestSelect, *valk.DefaultsTestOmit, FindManyQuery) ([]*valk.DefaultsTest, error)

type Extension = valk.DefaultsTestExtension
