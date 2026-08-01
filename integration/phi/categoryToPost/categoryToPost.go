package categoryToPost

import (
	"context"
	"integration/phi"
)

type Select = phi.CategoryToPostSelect
type Omit = phi.CategoryToPostOmit
type QueryBuilder = phi.CategoryToPostQueryBuilder
type CreateBuilder = phi.CategoryToPostCreateBuilder
type Upsert = phi.CategoryToPostUpsert
type ConflictBuilder[B any] = phi.CategoryToPostConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...phi.FieldAssignmentOf[phi.CategoryToPost]) phi.RecordInput {
	raw := make([]phi.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = phi.FieldAssignment(a)
	}
	return phi.RecordInput{Assignments: raw}
}

func And(preds ...phi.PredicateOf[phi.CategoryToPost]) phi.PredicateOf[phi.CategoryToPost] {
	return phi.And(preds...)
}

func Or(preds ...phi.PredicateOf[phi.CategoryToPost]) phi.PredicateOf[phi.CategoryToPost] {
	return phi.Or(preds...)
}

func Not(pred phi.PredicateOf[phi.CategoryToPost]) phi.PredicateOf[phi.CategoryToPost] {
	return phi.Not(pred)
}

var PostId = phi.StringField[phi.CategoryToPost]{Column: "postId"}

var CategoryId = phi.Field[phi.CategoryToPost, int32]{Column: "categoryId"}

type postId_CategoryId struct {
	phi.CompositeUniqueConstraint[phi.CategoryToPost]
	Column string
}

var PostId_CategoryId = postId_CategoryId{
	CompositeUniqueConstraint: phi.CompositeUniqueConstraint[phi.CategoryToPost]{
		Name: "PostId_CategoryId",
		Columns: []string{
			"postId",
			"categoryId",
		},
	},
	Column: "PostId_CategoryId",
}

func (f postId_CategoryId) EQ(postId string, categoryId int32) phi.UniquePredicate[phi.CategoryToPost] {
	return phi.UniquePredicate[phi.CategoryToPost]{
		Data: phi.PredicateData{
			Column:   "PostId_CategoryId",
			Operator: "AND",
			Value: map[string]any{
				"postId":     postId,
				"categoryId": categoryId,
			},
			IsLogical: true,
			Children: []phi.PredicateData{
				{Column: "postId", Operator: "=", Value: postId},
				{Column: "categoryId", Operator: "=", Value: categoryId},
			},
		},
	}
}

type CreateInput = phi.CategoryToPostCreate
type Create = phi.CategoryToPostCreate

type CreateArgs = phi.CategoryToPostCreateArgs
type CreateManyArgs = phi.CategoryToPostCreateManyArgs
type CreateManyAndReturnArgs = phi.CategoryToPostCreateManyAndReturnArgs

type CreateQuery = phi.CategoryToPostCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*phi.CategoryToPost, error)

type CreateManyQuery = phi.CategoryToPostCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = phi.CategoryToPostCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*phi.CategoryToPost, error)

type FindUniqueArgs = phi.CategoryToPostFindUniqueArgs
type FindUniqueQuery = phi.CategoryToPostFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*phi.CategoryToPost, error)

type FindFirstArgs = phi.CategoryToPostFindFirstArgs
type FindFirstQuery = phi.CategoryToPostFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*phi.CategoryToPost, error)

type FindManyArgs = phi.CategoryToPostFindManyArgs
type FindManyQuery = phi.CategoryToPostFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*phi.CategoryToPost, error)

type CountArgs = phi.CategoryToPostCountArgs
type CountQuery = phi.CategoryToPostCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = phi.CategoryToPostUpdate

type Extension = phi.CategoryToPostExtension

// ConflictUpdate creates a custom ConflictAction for CategoryToPost upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *phi.CategoryToPostUpsert)) *phi.ConflictAction {
	return phi.CategoryToPostConflictUpdate(fn)
}
