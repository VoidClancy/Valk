package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type CategoryToPost struct {
	ent.Schema
}

func (CategoryToPost) Fields() []ent.Field {
	return []ent.Field{
		field.String("post_id").StorageKey("postId"),
		field.Int("category_id").StorageKey("categoryId"),
	}
}

func (CategoryToPost) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("post", Post.Type).
			Ref("categories").
			Unique().
			Required().
			Field("post_id"),
		edge.From("category", Category.Type).
			Ref("posts").
			Unique().
			Required().
			Field("category_id"),
	}
}

func (CategoryToPost) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "CategoryToPost"},
	}
}
