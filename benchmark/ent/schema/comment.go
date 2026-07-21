package schema

import (
	"encoding/json"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Comment struct {
	ent.Schema
}

func (Comment) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().StorageKey("id"),
		field.Int("textify"),
		field.String("dummy3"),
		field.Int("dummy1"),
		field.String("dummy2"),
		field.String("post_id").StorageKey("postId"),
		field.String("author_id").StorageKey("authorId"),
		field.JSON("meta", json.RawMessage{}).Optional(),
	}
}

func (Comment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("post", Post.Type).
			Ref("comments").
			Unique().
			Required().
			Field("post_id"),
		edge.From("author", User.Type).
			Ref("comments").
			Unique().
			Required().
			Field("author_id"),
	}
}

func (Comment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "Comment"},
	}
}
