package categoryToPost

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

type Select = valk.CategoryToPostSelect
type Omit = valk.CategoryToPostOmit
type Create = valk.CategoryToPostCreate

func And(preds ...valk.Predicate) valk.Predicate {
	return valk.And(preds...)
}

func Or(preds ...valk.Predicate) valk.Predicate {
	return valk.Or(preds...)
}

func Not(pred valk.Predicate) valk.Predicate {
	return valk.Not(pred)
}

var PostId = valk.StringField{Column: "postId"}

var CategoryId = valk.Field[int32]{Column: "categoryId"}
