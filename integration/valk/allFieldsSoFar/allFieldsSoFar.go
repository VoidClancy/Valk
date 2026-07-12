package allFieldsSoFar

import (
	"encoding/json"
	"fmt"
	"integration/valk"
	"time"
)

type UniquePredicate struct {
	valk.StandardPredicate
}

func (UniquePredicate) IsUnique() {}

func (p UniquePredicate) Validate() error {
	if p.StandardPredicate.Data.Column == "" && len(p.StandardPredicate.Data.Children) == 0 {
		return fmt.Errorf("at least one unique field must be set for FindUnique")
	}
	return p.StandardPredicate.Validate()
}

type Select = valk.AllFieldsSoFarSelect
type Omit = valk.AllFieldsSoFarOmit
type QueryBuilder = valk.AllFieldsSoFarQueryBuilder

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...valk.FieldAssignment) valk.RecordInput {
	return valk.RecordInput{Assignments: assignments}
}

func And(preds ...valk.Predicate) valk.Predicate {
	return valk.And(preds...)
}

func Or(preds ...valk.Predicate) valk.Predicate {
	return valk.Or(preds...)
}

func Not(pred valk.Predicate) valk.Predicate {
	return valk.Not(pred)
}

var Id = valk.UniqueField[int32]{Column: "id"}

var StringReq = valk.StringField{Column: "stringReq"}

var StringOpt = valk.StringField{Column: "stringOpt"}

var StringDefault = valk.StringField{Column: "stringDefault"}

var StringVarchar = valk.StringField{Column: "stringVarchar"}

var StringChar = valk.StringField{Column: "stringChar"}

var BitVal = valk.StringField{Column: "bitVal"}

var VarBitVal = valk.StringField{Column: "varBitVal"}

var InetVal = valk.StringField{Column: "inetVal"}

var XmlVal = valk.StringField{Column: "xmlVal"}

var CuidDefault = valk.StringField{Column: "cuidDefault"}

var Cuid1Default = valk.StringField{Column: "cuid1Default"}

var Cuid2Default = valk.StringField{Column: "cuid2Default"}

var UuidDefault = valk.StringField{Column: "uuidDefault"}

var Uuid4Default = valk.StringField{Column: "uuid4Default"}

var Uuid7Default = valk.StringField{Column: "uuid7Default"}

var UlidDefault = valk.StringField{Column: "ulidDefault"}

var NanoidDefault = valk.StringField{Column: "nanoidDefault"}

var UuidDb = valk.StringField{Column: "uuidDb"}

var IntReq = valk.Field[int32]{Column: "intReq"}

var IntOpt = valk.Field[int32]{Column: "intOpt"}

var IntDefault = valk.Field[int32]{Column: "intDefault"}

var IntegerVal = valk.Field[int32]{Column: "integerVal"}

var SmallInt = valk.Field[int32]{Column: "smallInt"}

var TinyInt = valk.Field[int32]{Column: "tinyInt"}

var OidVal = valk.Field[int32]{Column: "oidVal"}

var BigIntReq = valk.Field[int64]{Column: "bigIntReq"}

var BigIntOpt = valk.Field[int64]{Column: "bigIntOpt"}

var FloatReq = valk.Field[float64]{Column: "floatReq"}

var FloatOpt = valk.Field[float64]{Column: "floatOpt"}

var RealVal = valk.Field[float64]{Column: "realVal"}

var DecimalReq = valk.Field[string]{Column: "decimalReq"}

var DecimalOpt = valk.Field[string]{Column: "decimalOpt"}

var DecimalPrecise = valk.Field[string]{Column: "decimalPrecise"}

var MoneyVal = valk.Field[string]{Column: "moneyVal"}

var BoolReq = valk.Field[bool]{Column: "boolReq"}

var BoolOpt = valk.Field[bool]{Column: "boolOpt"}

var BoolDefault = valk.Field[bool]{Column: "boolDefault"}

var DateTimeReq = valk.Field[time.Time]{Column: "dateTimeReq"}

var DateTimeOpt = valk.Field[time.Time]{Column: "dateTimeOpt"}

var DateTimeDefault = valk.Field[time.Time]{Column: "dateTimeDefault"}

var UpdatedAt = valk.Field[time.Time]{Column: "updatedAt"}

var DateTimeTz = valk.Field[time.Time]{Column: "dateTimeTz"}

var TimestampVal = valk.Field[time.Time]{Column: "timestampVal"}

var TimeVal = valk.Field[time.Time]{Column: "timeVal"}

var TimetzVal = valk.Field[time.Time]{Column: "timetzVal"}

var JsonReq = valk.Field[json.RawMessage]{Column: "jsonReq"}

var JsonOpt = valk.Field[json.RawMessage]{Column: "jsonOpt"}

var JsonVal = valk.Field[json.RawMessage]{Column: "jsonVal"}

var BytesReq = valk.Field[[]byte]{Column: "bytesReq"}

var BytesOpt = valk.Field[[]byte]{Column: "bytesOpt"}

var HstoreField = valk.Field[map[string]*string]{Column: "hstoreField"}

var LtreeField = valk.Field[string]{Column: "ltreeField"}

var CitextField = valk.Field[string]{Column: "citextField"}
