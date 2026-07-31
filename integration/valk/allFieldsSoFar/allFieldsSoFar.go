package allFieldsSoFar

import (
	"context"
	"encoding/json"
	"integration/valk"
	"time"
)

type Select = valk.AllFieldsSoFarSelect
type Omit = valk.AllFieldsSoFarOmit
type QueryBuilder = valk.AllFieldsSoFarQueryBuilder
type CreateBuilder = valk.AllFieldsSoFarCreateBuilder
type Upsert = valk.AllFieldsSoFarUpsert
type ConflictBuilder[B any] = valk.AllFieldsSoFarConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...valk.FieldAssignmentOf[valk.AllFieldsSoFar]) valk.RecordInput {
	raw := make([]valk.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = valk.FieldAssignment{Col: a.Col, Val: a.Val}
	}
	return valk.RecordInput{Assignments: raw}
}

func And(preds ...valk.PredicateOf[valk.AllFieldsSoFar]) valk.PredicateOf[valk.AllFieldsSoFar] {
	return valk.And(preds...)
}

func Or(preds ...valk.PredicateOf[valk.AllFieldsSoFar]) valk.PredicateOf[valk.AllFieldsSoFar] {
	return valk.Or(preds...)
}

func Not(pred valk.PredicateOf[valk.AllFieldsSoFar]) valk.PredicateOf[valk.AllFieldsSoFar] {
	return valk.Not(pred)
}

var Id = valk.UniqueField[valk.AllFieldsSoFar, int32]{Column: "id"}

var StringReq = valk.StringField[valk.AllFieldsSoFar]{Column: "stringReq"}

var StringOpt = valk.OptionalStringField[valk.AllFieldsSoFar]{StringField: valk.StringField[valk.AllFieldsSoFar]{Column: "stringOpt"}}

var StringDefault = valk.StringField[valk.AllFieldsSoFar]{Column: "stringDefault"}

var StringVarchar = valk.StringField[valk.AllFieldsSoFar]{Column: "stringVarchar"}

var StringChar = valk.StringField[valk.AllFieldsSoFar]{Column: "stringChar"}

var BitVal = valk.StringField[valk.AllFieldsSoFar]{Column: "bitVal"}

var VarBitVal = valk.StringField[valk.AllFieldsSoFar]{Column: "varBitVal"}

var InetVal = valk.StringField[valk.AllFieldsSoFar]{Column: "inetVal"}

var XmlVal = valk.StringField[valk.AllFieldsSoFar]{Column: "xmlVal"}

var CuidDefault = valk.StringField[valk.AllFieldsSoFar]{Column: "cuidDefault"}

var Cuid1Default = valk.StringField[valk.AllFieldsSoFar]{Column: "cuid1Default"}

var Cuid2Default = valk.StringField[valk.AllFieldsSoFar]{Column: "cuid2Default"}

var UuidDefault = valk.StringField[valk.AllFieldsSoFar]{Column: "uuidDefault"}

var Uuid4Default = valk.StringField[valk.AllFieldsSoFar]{Column: "uuid4Default"}

var Uuid7Default = valk.StringField[valk.AllFieldsSoFar]{Column: "uuid7Default"}

var UlidDefault = valk.StringField[valk.AllFieldsSoFar]{Column: "ulidDefault"}

var NanoidDefault = valk.StringField[valk.AllFieldsSoFar]{Column: "nanoidDefault"}

var UuidDb = valk.StringField[valk.AllFieldsSoFar]{Column: "uuidDb"}

var IntReq = valk.Field[valk.AllFieldsSoFar, int32]{Column: "intReq"}

var IntOpt = valk.OptionalField[valk.AllFieldsSoFar, int32]{Field: valk.Field[valk.AllFieldsSoFar, int32]{Column: "intOpt"}}

var IntDefault = valk.Field[valk.AllFieldsSoFar, int32]{Column: "intDefault"}

var IntegerVal = valk.Field[valk.AllFieldsSoFar, int32]{Column: "integerVal"}

var SmallInt = valk.Field[valk.AllFieldsSoFar, int32]{Column: "smallInt"}

var TinyInt = valk.Field[valk.AllFieldsSoFar, int32]{Column: "tinyInt"}

var OidVal = valk.Field[valk.AllFieldsSoFar, int32]{Column: "oidVal"}

var BigIntReq = valk.Field[valk.AllFieldsSoFar, int64]{Column: "bigIntReq"}

var BigIntOpt = valk.OptionalField[valk.AllFieldsSoFar, int64]{Field: valk.Field[valk.AllFieldsSoFar, int64]{Column: "bigIntOpt"}}

var FloatReq = valk.Field[valk.AllFieldsSoFar, float64]{Column: "floatReq"}

var FloatOpt = valk.OptionalField[valk.AllFieldsSoFar, float64]{Field: valk.Field[valk.AllFieldsSoFar, float64]{Column: "floatOpt"}}

var RealVal = valk.Field[valk.AllFieldsSoFar, float64]{Column: "realVal"}

var DecimalReq = valk.Field[valk.AllFieldsSoFar, string]{Column: "decimalReq"}

var DecimalOpt = valk.OptionalField[valk.AllFieldsSoFar, string]{Field: valk.Field[valk.AllFieldsSoFar, string]{Column: "decimalOpt"}}

var DecimalPrecise = valk.Field[valk.AllFieldsSoFar, string]{Column: "decimalPrecise"}

var MoneyVal = valk.Field[valk.AllFieldsSoFar, string]{Column: "moneyVal"}

var BoolReq = valk.Field[valk.AllFieldsSoFar, bool]{Column: "boolReq"}

var BoolOpt = valk.OptionalField[valk.AllFieldsSoFar, bool]{Field: valk.Field[valk.AllFieldsSoFar, bool]{Column: "boolOpt"}}

var BoolDefault = valk.Field[valk.AllFieldsSoFar, bool]{Column: "boolDefault"}

var DateTimeReq = valk.Field[valk.AllFieldsSoFar, time.Time]{Column: "dateTimeReq"}

var DateTimeOpt = valk.OptionalField[valk.AllFieldsSoFar, time.Time]{Field: valk.Field[valk.AllFieldsSoFar, time.Time]{Column: "dateTimeOpt"}}

var DateTimeDefault = valk.Field[valk.AllFieldsSoFar, time.Time]{Column: "dateTimeDefault"}

var UpdatedAt = valk.Field[valk.AllFieldsSoFar, time.Time]{Column: "updatedAt"}

var DateTimeTz = valk.Field[valk.AllFieldsSoFar, time.Time]{Column: "dateTimeTz"}

var TimestampVal = valk.Field[valk.AllFieldsSoFar, time.Time]{Column: "timestampVal"}

var TimeVal = valk.Field[valk.AllFieldsSoFar, time.Time]{Column: "timeVal"}

var TimetzVal = valk.Field[valk.AllFieldsSoFar, time.Time]{Column: "timetzVal"}

var JsonReq = valk.Field[valk.AllFieldsSoFar, json.RawMessage]{Column: "jsonReq"}

var JsonOpt = valk.OptionalField[valk.AllFieldsSoFar, json.RawMessage]{Field: valk.Field[valk.AllFieldsSoFar, json.RawMessage]{Column: "jsonOpt"}}

var JsonVal = valk.Field[valk.AllFieldsSoFar, json.RawMessage]{Column: "jsonVal"}

var BytesReq = valk.Field[valk.AllFieldsSoFar, []byte]{Column: "bytesReq"}

var BytesOpt = valk.OptionalField[valk.AllFieldsSoFar, []byte]{Field: valk.Field[valk.AllFieldsSoFar, []byte]{Column: "bytesOpt"}}

var HstoreField = valk.OptionalField[valk.AllFieldsSoFar, map[string]*string]{Field: valk.Field[valk.AllFieldsSoFar, map[string]*string]{Column: "hstoreField"}}

var LtreeField = valk.OptionalField[valk.AllFieldsSoFar, string]{Field: valk.Field[valk.AllFieldsSoFar, string]{Column: "ltreeField"}}

var CitextField = valk.OptionalField[valk.AllFieldsSoFar, string]{Field: valk.Field[valk.AllFieldsSoFar, string]{Column: "citextField"}}

type CreateInput = valk.AllFieldsSoFarCreate
type Create = valk.AllFieldsSoFarCreate

type CreateArgs = valk.AllFieldsSoFarCreateArgs
type CreateManyArgs = valk.AllFieldsSoFarCreateManyArgs
type CreateManyAndReturnArgs = valk.AllFieldsSoFarCreateManyAndReturnArgs

type CreateQuery = valk.AllFieldsSoFarCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*valk.AllFieldsSoFar, error)

type CreateManyQuery = valk.AllFieldsSoFarCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = valk.AllFieldsSoFarCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*valk.AllFieldsSoFar, error)

type FindUniqueArgs = valk.AllFieldsSoFarFindUniqueArgs
type FindUniqueQuery = valk.AllFieldsSoFarFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*valk.AllFieldsSoFar, error)

type FindFirstArgs = valk.AllFieldsSoFarFindFirstArgs
type FindFirstQuery = valk.AllFieldsSoFarFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*valk.AllFieldsSoFar, error)

type FindManyArgs = valk.AllFieldsSoFarFindManyArgs
type FindManyQuery = valk.AllFieldsSoFarFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*valk.AllFieldsSoFar, error)

type CountArgs = valk.AllFieldsSoFarCountArgs
type CountQuery = valk.AllFieldsSoFarCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = valk.AllFieldsSoFarUpdate

type Extension = valk.AllFieldsSoFarExtension

// ConflictUpdate creates a custom ConflictAction for AllFieldsSoFar upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *valk.AllFieldsSoFarUpsert)) *valk.ConflictAction {
	return valk.AllFieldsSoFarConflictUpdate(fn)
}
