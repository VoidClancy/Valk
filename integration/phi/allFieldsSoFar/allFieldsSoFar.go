package allFieldsSoFar

import (
	"context"
	"encoding/json"
	"integration/phi"
	"time"
)

type Select = phi.AllFieldsSoFarSelect
type Omit = phi.AllFieldsSoFarOmit
type QueryBuilder = phi.AllFieldsSoFarQueryBuilder
type CreateBuilder = phi.AllFieldsSoFarCreateBuilder
type Upsert = phi.AllFieldsSoFarUpsert
type ConflictBuilder[B any] = phi.AllFieldsSoFarConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...phi.FieldAssignmentOf[phi.AllFieldsSoFar]) phi.RecordInput {
	raw := make([]phi.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = phi.FieldAssignment(a)
	}
	return phi.RecordInput{Assignments: raw}
}

func And(preds ...phi.PredicateOf[phi.AllFieldsSoFar]) phi.PredicateOf[phi.AllFieldsSoFar] {
	return phi.And(preds...)
}

func Or(preds ...phi.PredicateOf[phi.AllFieldsSoFar]) phi.PredicateOf[phi.AllFieldsSoFar] {
	return phi.Or(preds...)
}

func Not(pred phi.PredicateOf[phi.AllFieldsSoFar]) phi.PredicateOf[phi.AllFieldsSoFar] {
	return phi.Not(pred)
}

var Id = phi.UniqueField[phi.AllFieldsSoFar, int32]{Column: "id"}

var StringReq = phi.StringField[phi.AllFieldsSoFar]{Column: "stringReq"}

var StringOpt = phi.OptionalStringField[phi.AllFieldsSoFar]{StringField: phi.StringField[phi.AllFieldsSoFar]{Column: "stringOpt"}}

var StringDefault = phi.StringField[phi.AllFieldsSoFar]{Column: "stringDefault"}

var StringVarchar = phi.StringField[phi.AllFieldsSoFar]{Column: "stringVarchar"}

var StringChar = phi.StringField[phi.AllFieldsSoFar]{Column: "stringChar"}

var BitVal = phi.StringField[phi.AllFieldsSoFar]{Column: "bitVal"}

var VarBitVal = phi.StringField[phi.AllFieldsSoFar]{Column: "varBitVal"}

var InetVal = phi.StringField[phi.AllFieldsSoFar]{Column: "inetVal"}

var XmlVal = phi.StringField[phi.AllFieldsSoFar]{Column: "xmlVal"}

var CuidDefault = phi.StringField[phi.AllFieldsSoFar]{Column: "cuidDefault"}

var Cuid1Default = phi.StringField[phi.AllFieldsSoFar]{Column: "cuid1Default"}

var Cuid2Default = phi.StringField[phi.AllFieldsSoFar]{Column: "cuid2Default"}

var UuidDefault = phi.StringField[phi.AllFieldsSoFar]{Column: "uuidDefault"}

var Uuid4Default = phi.StringField[phi.AllFieldsSoFar]{Column: "uuid4Default"}

var Uuid7Default = phi.StringField[phi.AllFieldsSoFar]{Column: "uuid7Default"}

var UlidDefault = phi.StringField[phi.AllFieldsSoFar]{Column: "ulidDefault"}

var NanoidDefault = phi.StringField[phi.AllFieldsSoFar]{Column: "nanoidDefault"}

var UuidDb = phi.StringField[phi.AllFieldsSoFar]{Column: "uuidDb"}

var IntReq = phi.Field[phi.AllFieldsSoFar, int32]{Column: "intReq"}

var IntOpt = phi.OptionalField[phi.AllFieldsSoFar, int32]{Field: phi.Field[phi.AllFieldsSoFar, int32]{Column: "intOpt"}}

var IntDefault = phi.Field[phi.AllFieldsSoFar, int32]{Column: "intDefault"}

var IntegerVal = phi.Field[phi.AllFieldsSoFar, int32]{Column: "integerVal"}

var SmallInt = phi.Field[phi.AllFieldsSoFar, int32]{Column: "smallInt"}

var TinyInt = phi.Field[phi.AllFieldsSoFar, int32]{Column: "tinyInt"}

var OidVal = phi.Field[phi.AllFieldsSoFar, int32]{Column: "oidVal"}

var BigIntReq = phi.Field[phi.AllFieldsSoFar, int64]{Column: "bigIntReq"}

var BigIntOpt = phi.OptionalField[phi.AllFieldsSoFar, int64]{Field: phi.Field[phi.AllFieldsSoFar, int64]{Column: "bigIntOpt"}}

var FloatReq = phi.Field[phi.AllFieldsSoFar, float64]{Column: "floatReq"}

var FloatOpt = phi.OptionalField[phi.AllFieldsSoFar, float64]{Field: phi.Field[phi.AllFieldsSoFar, float64]{Column: "floatOpt"}}

var RealVal = phi.Field[phi.AllFieldsSoFar, float64]{Column: "realVal"}

var DecimalReq = phi.Field[phi.AllFieldsSoFar, string]{Column: "decimalReq"}

var DecimalOpt = phi.OptionalField[phi.AllFieldsSoFar, string]{Field: phi.Field[phi.AllFieldsSoFar, string]{Column: "decimalOpt"}}

var DecimalPrecise = phi.Field[phi.AllFieldsSoFar, string]{Column: "decimalPrecise"}

var MoneyVal = phi.Field[phi.AllFieldsSoFar, string]{Column: "moneyVal"}

var BoolReq = phi.Field[phi.AllFieldsSoFar, bool]{Column: "boolReq"}

var BoolOpt = phi.OptionalField[phi.AllFieldsSoFar, bool]{Field: phi.Field[phi.AllFieldsSoFar, bool]{Column: "boolOpt"}}

var BoolDefault = phi.Field[phi.AllFieldsSoFar, bool]{Column: "boolDefault"}

var DateTimeReq = phi.Field[phi.AllFieldsSoFar, time.Time]{Column: "dateTimeReq"}

var DateTimeOpt = phi.OptionalField[phi.AllFieldsSoFar, time.Time]{Field: phi.Field[phi.AllFieldsSoFar, time.Time]{Column: "dateTimeOpt"}}

var DateTimeDefault = phi.Field[phi.AllFieldsSoFar, time.Time]{Column: "dateTimeDefault"}

var UpdatedAt = phi.Field[phi.AllFieldsSoFar, time.Time]{Column: "updatedAt"}

var DateTimeTz = phi.Field[phi.AllFieldsSoFar, time.Time]{Column: "dateTimeTz"}

var TimestampVal = phi.Field[phi.AllFieldsSoFar, time.Time]{Column: "timestampVal"}

var TimeVal = phi.Field[phi.AllFieldsSoFar, time.Time]{Column: "timeVal"}

var TimetzVal = phi.Field[phi.AllFieldsSoFar, time.Time]{Column: "timetzVal"}

var JsonReq = phi.Field[phi.AllFieldsSoFar, json.RawMessage]{Column: "jsonReq"}

var JsonOpt = phi.OptionalField[phi.AllFieldsSoFar, json.RawMessage]{Field: phi.Field[phi.AllFieldsSoFar, json.RawMessage]{Column: "jsonOpt"}}

var JsonVal = phi.Field[phi.AllFieldsSoFar, json.RawMessage]{Column: "jsonVal"}

var BytesReq = phi.Field[phi.AllFieldsSoFar, []byte]{Column: "bytesReq"}

var BytesOpt = phi.OptionalField[phi.AllFieldsSoFar, []byte]{Field: phi.Field[phi.AllFieldsSoFar, []byte]{Column: "bytesOpt"}}

var HstoreField = phi.OptionalField[phi.AllFieldsSoFar, map[string]*string]{Field: phi.Field[phi.AllFieldsSoFar, map[string]*string]{Column: "hstoreField"}}

var LtreeField = phi.OptionalField[phi.AllFieldsSoFar, string]{Field: phi.Field[phi.AllFieldsSoFar, string]{Column: "ltreeField"}}

var CitextField = phi.OptionalField[phi.AllFieldsSoFar, string]{Field: phi.Field[phi.AllFieldsSoFar, string]{Column: "citextField"}}

type CreateInput = phi.AllFieldsSoFarCreate
type Create = phi.AllFieldsSoFarCreate

type CreateArgs = phi.AllFieldsSoFarCreateArgs
type CreateManyArgs = phi.AllFieldsSoFarCreateManyArgs
type CreateManyAndReturnArgs = phi.AllFieldsSoFarCreateManyAndReturnArgs

type CreateQuery = phi.AllFieldsSoFarCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*phi.AllFieldsSoFar, error)

type CreateManyQuery = phi.AllFieldsSoFarCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = phi.AllFieldsSoFarCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*phi.AllFieldsSoFar, error)

type FindUniqueArgs = phi.AllFieldsSoFarFindUniqueArgs
type FindUniqueQuery = phi.AllFieldsSoFarFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*phi.AllFieldsSoFar, error)

type FindFirstArgs = phi.AllFieldsSoFarFindFirstArgs
type FindFirstQuery = phi.AllFieldsSoFarFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*phi.AllFieldsSoFar, error)

type FindManyArgs = phi.AllFieldsSoFarFindManyArgs
type FindManyQuery = phi.AllFieldsSoFarFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*phi.AllFieldsSoFar, error)

type CountArgs = phi.AllFieldsSoFarCountArgs
type CountQuery = phi.AllFieldsSoFarCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = phi.AllFieldsSoFarUpdate

type Extension = phi.AllFieldsSoFarExtension

// ConflictUpdate creates a custom ConflictAction for AllFieldsSoFar upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *phi.AllFieldsSoFarUpsert)) *phi.ConflictAction {
	return phi.AllFieldsSoFarConflictUpdate(fn)
}
