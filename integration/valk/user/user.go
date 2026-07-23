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

var EmailPhone = valk.CompositeUniqueConstraint[valk.User]{
	Name: "emailPhone",
	Columns: []string{
		"email",
		"phoneNum",
	},
}

// Helper for compound unique constraint: emailPhone
func EmailPhoneUnique(email string, phoneNum string) valk.UniquePredicate[valk.User] {
	return valk.UniquePredicate[valk.User]{
		Data: valk.And(
			valk.Predicate[valk.User]{
				Data: valk.PredicateData{
					Column:   "email",
					Operator: "=",
					Value:    email,
				},
			},
			valk.Predicate[valk.User]{
				Data: valk.PredicateData{
					Column:   "phoneNum",
					Operator: "=",
					Value:    phoneNum,
				},
			},
		).ToPredicateData(),
	}
}

type CreateInput = valk.UserCreate
type CreateQuery = valk.UserCreateQuery
type CreateHook = func(context.Context, *CreateInput, CreateQuery) (*valk.User, error)

type CreateManyQuery = valk.UserCreateManyQuery
type CreateManyHook = func(context.Context, []*CreateInput, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = valk.UserCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, []*CreateInput, CreateManyAndReturnQuery) ([]*valk.User, error)

type FindUniqueQuery = valk.UserFindUniqueQuery
type FindUniqueHook = func(context.Context, valk.UniquePredicate[valk.User], []valk.PredicateOf[valk.User], *valk.UserSelect, *valk.UserOmit, FindUniqueQuery) (*valk.User, error)

type FindFirstQuery = valk.UserFindFirstQuery
type FindFirstHook = func(context.Context, valk.QueryParams[valk.User], *valk.UserSelect, *valk.UserOmit, FindFirstQuery) (*valk.User, error)

type FindManyQuery = valk.UserFindManyQuery
type FindManyHook = func(context.Context, valk.QueryParams[valk.User], *valk.UserSelect, *valk.UserOmit, FindManyQuery) ([]*valk.User, error)

type Extension = valk.UserExtension
