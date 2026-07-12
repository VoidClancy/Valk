package defaultsTest

import (
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

type Select = valk.DefaultsTestSelect
type Omit = valk.DefaultsTestOmit
type QueryBuilder = valk.DefaultsTestQueryBuilder

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

var Uuid4 = valk.StringUniqueField{Column: "uuid4"}

var Uuid7 = valk.StringField{Column: "uuid7"}

var UuidNoArgs = valk.StringField{Column: "uuidNoArgs"}

var Cuid1 = valk.StringField{Column: "cuid1"}

var Cuid2 = valk.StringField{Column: "cuid2"}

var CuidNoArgs = valk.StringField{Column: "cuidNoArgs"}

var Ulid = valk.StringField{Column: "ulid"}

var Nanoid = valk.StringField{Column: "nanoid"}

var Now = valk.Field[time.Time]{Column: "now"}
