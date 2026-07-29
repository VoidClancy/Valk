package user

import (
	"context"
	"integration/valk"
)

type Select = valk.UserSelect
type Omit = valk.UserOmit
type QueryBuilder = valk.UserQueryBuilder
type CreateBuilder = valk.UserCreateBuilder
type Upsert = valk.UserUpsert
type ConflictBuilder[B any] = valk.UserConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...valk.FieldAssignmentOf[valk.User]) valk.RecordInput {
	raw := make([]valk.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = valk.FieldAssignment{Col: a.Col, Val: a.Val}
	}
	return valk.RecordInput{Assignments: raw}
}

func And(preds ...valk.PredicateOf[valk.User]) valk.PredicateOf[valk.User] {
	return valk.And(preds...)
}

func Or(preds ...valk.PredicateOf[valk.User]) valk.PredicateOf[valk.User] {
	return valk.Or(preds...)
}

func Not(pred valk.PredicateOf[valk.User]) valk.PredicateOf[valk.User] {
	return valk.Not(pred)
}

var Id = valk.StringUniqueField[valk.User]{Column: "id"}

var Email = valk.StringUniqueField[valk.User]{Column: "email"}

var PhoneNum = valk.StringUniqueField[valk.User]{Column: "phoneNum"}

var Password = valk.StringField[valk.User]{Column: "password"}

var Role = valk.Field[valk.User, valk.UserRoleType]{Column: "role"}

var RoleOptional = valk.Field[valk.User, valk.UserRoleType]{Column: "roleOptional"}

var LoginCount = valk.Field[valk.User, int32]{Column: "loginCount"}

var ReferredById = valk.StringField[valk.User]{Column: "referredById"}

type emailPhone struct {
	valk.CompositeUniqueConstraint[valk.User]
}

var EmailPhone = emailPhone{
	CompositeUniqueConstraint: valk.CompositeUniqueConstraint[valk.User]{
		Name: "emailPhone",
		Columns: []string{
			"email",
			"phoneNum",
		},
	},
}

func (f emailPhone) EQ(email string, phoneNum string) valk.UniquePredicate[valk.User] {
	return valk.UniquePredicate[valk.User]{
		Data: valk.PredicateData{
			Column:   "emailPhone",
			Operator: "AND",
			Value: map[string]any{
				"email":    email,
				"phoneNum": phoneNum,
			},
			IsLogical: true,
			Children: []valk.PredicateData{
				{Column: "email", Operator: "=", Value: email},
				{Column: "phoneNum", Operator: "=", Value: phoneNum},
			},
		},
	}
}

type CreateInput = valk.UserCreate
type Create = valk.UserCreate

type CreateArgs = valk.UserCreateArgs
type CreateManyArgs = valk.UserCreateManyArgs
type CreateManyAndReturnArgs = valk.UserCreateManyAndReturnArgs

type CreateQuery = valk.UserCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*valk.User, error)

type CreateManyQuery = valk.UserCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = valk.UserCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*valk.User, error)

type FindUniqueArgs = valk.UserFindUniqueArgs
type FindUniqueQuery = valk.UserFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*valk.User, error)

type FindFirstArgs = valk.UserFindFirstArgs
type FindFirstQuery = valk.UserFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*valk.User, error)

type FindManyArgs = valk.UserFindManyArgs
type FindManyQuery = valk.UserFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*valk.User, error)

type CountArgs = valk.UserCountArgs
type CountQuery = valk.UserCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = valk.UserUpdate

type Extension = valk.UserExtension

// ConflictUpdate creates a custom ConflictAction for User upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *valk.UserUpsert)) *valk.ConflictAction {
	return valk.UserConflictUpdate(fn)
}
