package user

import (
	"fmt"
	"integration/valk"
)

type UniquePredicate struct {
	valk.StandardPredicate
}

func (UniquePredicate) IsUnique() {}

func (p UniquePredicate) Validate() error {
	if p.StandardPredicate.Data.Column == "" && len(p.StandardPredicate.Data.Children) == 0 {
		return fmt.Errorf("at least one unique field must be set for FindUnique")
	}
	return p.StandardPredicate.Validate()
}

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

func And(preds ...valk.Predicate) valk.Predicate {
	return valk.And(preds...)
}

func Or(preds ...valk.Predicate) valk.Predicate {
	return valk.Or(preds...)
}

func Not(pred valk.Predicate) valk.Predicate {
	return valk.Not(pred)
}

var Id = valk.StringUniqueField{Column: "id"}

var Email = valk.StringUniqueField{Column: "email"}

var PhoneNum = valk.StringUniqueField{Column: "phoneNum"}

var Password = valk.StringField{Column: "password"}

var Role = valk.Field[valk.UserRoleType]{Column: "role"}

var RoleOptional = valk.Field[valk.UserRoleType]{Column: "roleOptional"}

var ReferredById = valk.StringField{Column: "referredById"}

// Helper for compound unique constraint: emailPhone
func EmailPhoneUnique(email string, phoneNum string) UniquePredicate {
	return UniquePredicate{
		StandardPredicate: valk.StandardPredicate{
			Data: valk.And(
				valk.StandardPredicate{
					Data: valk.PredicateData{
						Column:   "email",
						Operator: "=",
						Value:    email,
					},
				},
				valk.StandardPredicate{
					Data: valk.PredicateData{
						Column:   "phoneNum",
						Operator: "=",
						Value:    phoneNum,
					},
				},
			).ToPredicateData(),
		},
	}
}
