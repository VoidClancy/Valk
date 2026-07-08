package comment

import (
	"encoding/json"
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
	return nil
}

type Select = valk.CommentSelect
type Omit = valk.CommentOmit
type Create = valk.CommentCreate

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

var Textify = valk.Field[int32]{Column: "textify"}

var Dummy3 = valk.StringField{Column: "dummy3"}

var Dummy1 = valk.Field[int32]{Column: "dummy1"}

var Dummy2 = valk.StringField{Column: "dummy2"}

var PostId = valk.StringField{Column: "postId"}

var AuthorId = valk.StringField{Column: "authorId"}

var Meta = valk.Field[json.RawMessage]{Column: "meta"}
