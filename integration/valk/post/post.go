package post

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

type Select = valk.PostSelect
type Omit = valk.PostOmit
type Create = valk.PostCreate

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

var Title = valk.StringField{Column: "title"}

var Content = valk.StringField{Column: "content"}

var Published = valk.Field[bool]{Column: "published"}

var AuthorId = valk.StringField{Column: "authorId"}
