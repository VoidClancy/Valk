package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Category struct {
	ent.Schema
}

func (Category) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().StorageKey("id"),
		field.String("name").Unique(),
	}
}

func (Category) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("posts", CategoryToPost.Type),
	}
}

func (Category) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "Category"},
	}
}
