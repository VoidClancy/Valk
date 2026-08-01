package user

import (
	"context"
	"integration/phi"
	"time"
)

type Select = phi.UserSelect
type Omit = phi.UserOmit
type QueryBuilder = phi.UserQueryBuilder
type CreateBuilder = phi.UserCreateBuilder
type Upsert = phi.UserUpsert
type ConflictBuilder[B any] = phi.UserConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...phi.FieldAssignmentOf[phi.User]) phi.RecordInput {
	raw := make([]phi.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = phi.FieldAssignment(a)
	}
	return phi.RecordInput{Assignments: raw}
}

func And(preds ...phi.PredicateOf[phi.User]) phi.PredicateOf[phi.User] {
	return phi.And(preds...)
}

func Or(preds ...phi.PredicateOf[phi.User]) phi.PredicateOf[phi.User] {
	return phi.Or(preds...)
}

func Not(pred phi.PredicateOf[phi.User]) phi.PredicateOf[phi.User] {
	return phi.Not(pred)
}

var Id = phi.StringUniqueField[phi.User]{Column: "id"}

var Email = phi.StringUniqueField[phi.User]{Column: "email"}

var PhoneNum = phi.StringUniqueField[phi.User]{Column: "phoneNum"}

var Password = phi.OptionalStringField[phi.User]{StringField: phi.StringField[phi.User]{Column: "password"}}

var Role = phi.Field[phi.User, phi.UserRoleType]{Column: "role"}

var RoleOptional = phi.OptionalField[phi.User, phi.UserRoleType]{Field: phi.Field[phi.User, phi.UserRoleType]{Column: "roleOptional"}}

var LoginCount = phi.Field[phi.User, int32]{Column: "loginCount"}

var ReferredById = phi.OptionalStringField[phi.User]{StringField: phi.StringField[phi.User]{Column: "referredById"}}

var CreatedAt = phi.Field[phi.User, time.Time]{Column: "createdAt"}

var UpdatedAt = phi.Field[phi.User, time.Time]{Column: "updatedAt"}

var UpdatedAtNoDefualt = phi.Field[phi.User, time.Time]{Column: "updatedAtNoDefualt"}

var UpdatedAtnoDecorator = phi.Field[phi.User, time.Time]{Column: "updatedAtnoDecorator"}

var UpdatedAtOptional = phi.OptionalField[phi.User, time.Time]{Field: phi.Field[phi.User, time.Time]{Column: "updatedAtOptional"}}

type emailPhone struct {
	phi.CompositeUniqueConstraint[phi.User]
	Column string
}

var EmailPhone = emailPhone{
	CompositeUniqueConstraint: phi.CompositeUniqueConstraint[phi.User]{
		Name: "emailPhone",
		Columns: []string{
			"email",
			"phoneNum",
		},
	},
	Column: "emailPhone",
}

func (f emailPhone) EQ(email string, phoneNum string) phi.UniquePredicate[phi.User] {
	return phi.UniquePredicate[phi.User]{
		Data: phi.PredicateData{
			Column:   "emailPhone",
			Operator: "AND",
			Value: map[string]any{
				"email":    email,
				"phoneNum": phoneNum,
			},
			IsLogical: true,
			Children: []phi.PredicateData{
				{Column: "email", Operator: "=", Value: email},
				{Column: "phoneNum", Operator: "=", Value: phoneNum},
			},
		},
	}
}

type CreateInput = phi.UserCreate
type Create = phi.UserCreate

type CreateArgs = phi.UserCreateArgs
type CreateManyArgs = phi.UserCreateManyArgs
type CreateManyAndReturnArgs = phi.UserCreateManyAndReturnArgs

type CreateQuery = phi.UserCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*phi.User, error)

type CreateManyQuery = phi.UserCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = phi.UserCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*phi.User, error)

type FindUniqueArgs = phi.UserFindUniqueArgs
type FindUniqueQuery = phi.UserFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*phi.User, error)

type FindFirstArgs = phi.UserFindFirstArgs
type FindFirstQuery = phi.UserFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*phi.User, error)

type FindManyArgs = phi.UserFindManyArgs
type FindManyQuery = phi.UserFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*phi.User, error)

type CountArgs = phi.UserCountArgs
type CountQuery = phi.UserCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = phi.UserUpdate

type Extension = phi.UserExtension

// ConflictUpdate creates a custom ConflictAction for User upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *phi.UserUpsert)) *phi.ConflictAction {
	return phi.UserConflictUpdate(fn)
}
