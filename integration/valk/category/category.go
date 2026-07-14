package category

import (
	"integration/valk"
)

type Select = valk.CategorySelect
type Omit = valk.CategoryOmit
type QueryBuilder = valk.CategoryQueryBuilder
type CreateBuilder = valk.CategoryCreateBuilder

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...valk.FieldAssignment) valk.RecordInput {
	return valk.RecordInput{Assignments: assignments}
}

func And(preds ...valk.PredicateOf[valk.Category]) valk.PredicateOf[valk.Category] {
	return valk.And(preds...)
}

func Or(preds ...valk.PredicateOf[valk.Category]) valk.PredicateOf[valk.Category] {
	return valk.Or(preds...)
}

func Not(pred valk.PredicateOf[valk.Category]) valk.PredicateOf[valk.Category] {
	return valk.Not(pred)
}

var Id = valk.UniqueField[valk.Category, int32]{Column: "id"}

var Name = valk.StringUniqueField[valk.Category]{Column: "name"}
