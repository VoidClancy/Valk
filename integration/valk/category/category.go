package category

import (
	"fmt"
	"integration/valk"
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

type Select = valk.CategorySelect
type Omit = valk.CategoryOmit
type QueryBuilder = valk.CategoryQueryBuilder

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

var Name = valk.StringUniqueField{Column: "name"}
