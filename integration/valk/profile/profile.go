package profile

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

type Select = valk.ProfileSelect
type Omit = valk.ProfileOmit

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

var Id = valk.StringUniqueField{Column: "id"}

var Bio = valk.StringField{Column: "bio"}

var UserId = valk.StringUniqueField{Column: "userId"}

var CreatedAt = valk.Field[time.Time]{Column: "createdAt"}
