package user

import (
	"integration/valk"
)

type Select = valk.UserSelect
type Omit = valk.UserOmit
type QueryBuilder = valk.UserQueryBuilder
type CreateBuilder = valk.UserCreateBuilder

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...valk.FieldAssignment) valk.RecordInput {
	return valk.RecordInput{Assignments: assignments}
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

var ReferredById = valk.StringField[valk.User]{Column: "referredById"}

// Helper for compound unique constraint: emailPhone
func EmailPhoneUnique(email string, phoneNum string) valk.UniquePredicate[valk.User] {
	return valk.UniquePredicate[valk.User]{
		Data: valk.And[valk.User](
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
