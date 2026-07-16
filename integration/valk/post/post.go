package post

import (
	"integration/valk"
)

type Select = valk.PostSelect
type Omit = valk.PostOmit
type QueryBuilder = valk.PostQueryBuilder
type CreateBuilder = valk.PostCreateBuilder
type Upsert = valk.PostUpsert
type ConflictBuilder[B any] = valk.PostConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...valk.FieldAssignment) valk.RecordInput {
	return valk.RecordInput{Assignments: assignments}
}

func And(preds ...valk.PredicateOf[valk.Post]) valk.PredicateOf[valk.Post] {
	return valk.And(preds...)
}

func Or(preds ...valk.PredicateOf[valk.Post]) valk.PredicateOf[valk.Post] {
	return valk.Or(preds...)
}

func Not(pred valk.PredicateOf[valk.Post]) valk.PredicateOf[valk.Post] {
	return valk.Not(pred)
}

var Id = valk.StringUniqueField[valk.Post]{Column: "id"}

var Title = valk.StringField[valk.Post]{Column: "title"}

var Content = valk.StringField[valk.Post]{Column: "content"}

var Published = valk.Field[valk.Post, bool]{Column: "published"}

var AuthorId = valk.StringField[valk.Post]{Column: "authorId"}
