package categoryToPost

import (
	"context"
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

func Record(assignments ...valk.FieldAssignmentOf[valk.CategoryToPost]) valk.RecordInput {
	raw := make([]valk.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = valk.FieldAssignment{Col: a.Col, Val: a.Val}
	}
	return valk.RecordInput{Assignments: raw}
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

type CreateInput = valk.CategoryToPostCreate
type CreateArgs = valk.CategoryToPostCreateArgs
type CreateManyArgs = valk.CategoryToPostCreateManyArgs
type CreateManyAndReturnArgs = valk.CategoryToPostCreateManyAndReturnArgs

type CreateQuery = valk.CategoryToPostCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*valk.CategoryToPost, error)

type CreateManyQuery = valk.CategoryToPostCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = valk.CategoryToPostCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*valk.CategoryToPost, error)

type FindUniqueArgs = valk.CategoryToPostFindUniqueArgs
type FindUniqueQuery = valk.CategoryToPostFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*valk.CategoryToPost, error)

type FindFirstArgs = valk.CategoryToPostFindFirstArgs
type FindFirstQuery = valk.CategoryToPostFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*valk.CategoryToPost, error)

type FindManyArgs = valk.CategoryToPostFindManyArgs
type FindManyQuery = valk.CategoryToPostFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*valk.CategoryToPost, error)

type CountArgs = valk.CategoryToPostCountArgs
type CountQuery = valk.CategoryToPostCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

type Extension = valk.CategoryToPostExtension

// ConflictUpdate creates a custom ConflictAction for CategoryToPost upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *valk.CategoryToPostUpsert)) *valk.ConflictAction {
	return valk.CategoryToPostConflictUpdate(fn)
}
