package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Post struct {
	ent.Schema
}

func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().StorageKey("id"),
		field.String("title"),
		field.String("content").Optional().Nillable(),
		field.Bool("published").Default(false),
		field.String("author_id").StorageKey("authorId"),
	}
}

func (Post) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("author", User.Type).
			Ref("posts").
			Unique().
			Required().
			Field("author_id"),
		edge.To("comments", Comment.Type),
		edge.To("categories", CategoryToPost.Type),
	}
}

func (Post) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "Post"},
	}
}
