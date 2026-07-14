package categoryToPost

import (
	"integration/valk"
)

type Select = valk.CategoryToPostSelect
type Omit = valk.CategoryToPostOmit
type QueryBuilder = valk.CategoryToPostQueryBuilder
type CreateBuilder = valk.CategoryToPostCreateBuilder

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...valk.FieldAssignment) valk.RecordInput {
	return valk.RecordInput{Assignments: assignments}
}

func And(preds ...valk.PredicateOf[valk.CategoryToPost]) valk.PredicateOf[valk.CategoryToPost] {
	return valk.And(preds...)
}

func Or(preds ...valk.PredicateOf[valk.CategoryToPost]) valk.PredicateOf[valk.CategoryToPost] {
	return valk.Or(preds...)
}

func Not(pred valk.PredicateOf[valk.CategoryToPost]) valk.PredicateOf[valk.CategoryToPost] {
	return valk.Not(pred)
}

var PostId = valk.StringField[valk.CategoryToPost]{Column: "postId"}

var CategoryId = valk.Field[valk.CategoryToPost, int32]{Column: "categoryId"}
