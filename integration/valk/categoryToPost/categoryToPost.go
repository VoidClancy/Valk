package categoryToPost

import (
	"integration/valk"
)

type Select = valk.CategoryToPostSelect
type Omit = valk.CategoryToPostOmit
type QueryBuilder = valk.CategoryToPostQueryBuilder
type CreateBuilder = valk.CategoryToPostCreateBuilder
type Upsert = valk.CategoryToPostUpsert
type ConflictBuilder[B any] = valk.CategoryToPostConflictBuilder[B]

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

var PostIdCategoryId = valk.CompositeUniqueConstraint[valk.CategoryToPost]{
	Name: "PostIdCategoryId",
	Columns: []string{
		"postId",
		"categoryId",
	},
}

// Helper for compound primary key: PostIdCategoryId
func PostIdCategoryIdUnique(postId string, categoryId int32) valk.UniquePredicate[valk.CategoryToPost] {
	return valk.UniquePredicate[valk.CategoryToPost]{
		Data: valk.And(
			valk.Predicate[valk.CategoryToPost]{
				Data: valk.PredicateData{
					Column:   "postId",
					Operator: "=",
					Value:    postId,
				},
			},
			valk.Predicate[valk.CategoryToPost]{
				Data: valk.PredicateData{
					Column:   "categoryId",
					Operator: "=",
					Value:    categoryId,
				},
			},
		).ToPredicateData(),
	}
}
