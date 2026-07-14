package comment

import (
	"encoding/json"
	"integration/valk"
)

type Select = valk.CommentSelect
type Omit = valk.CommentOmit
type QueryBuilder = valk.CommentQueryBuilder
type CreateBuilder = valk.CommentCreateBuilder

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...valk.FieldAssignment) valk.RecordInput {
	return valk.RecordInput{Assignments: assignments}
}

func And(preds ...valk.PredicateOf[valk.Comment]) valk.PredicateOf[valk.Comment] {
	return valk.And(preds...)
}

func Or(preds ...valk.PredicateOf[valk.Comment]) valk.PredicateOf[valk.Comment] {
	return valk.Or(preds...)
}

func Not(pred valk.PredicateOf[valk.Comment]) valk.PredicateOf[valk.Comment] {
	return valk.Not(pred)
}

var Id = valk.StringUniqueField[valk.Comment]{Column: "id"}

var Textify = valk.Field[valk.Comment, int32]{Column: "textify"}

var Dummy3 = valk.StringField[valk.Comment]{Column: "dummy3"}

var Dummy1 = valk.Field[valk.Comment, int32]{Column: "dummy1"}

var Dummy2 = valk.StringField[valk.Comment]{Column: "dummy2"}

var PostId = valk.StringField[valk.Comment]{Column: "postId"}

var AuthorId = valk.StringField[valk.Comment]{Column: "authorId"}

var Meta = valk.Field[valk.Comment, json.RawMessage]{Column: "meta"}
