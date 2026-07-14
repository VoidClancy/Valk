package profile

import (
	"integration/valk"
	"time"
)

type Select = valk.ProfileSelect
type Omit = valk.ProfileOmit
type QueryBuilder = valk.ProfileQueryBuilder
type CreateBuilder = valk.ProfileCreateBuilder

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
