package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Profile struct {
	ent.Schema
}

func (Profile) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().StorageKey("id"),
		field.String("bio").Optional().Nillable(),
		field.String("user_id").Unique().StorageKey("userId"),
		field.Time("created_at").Default(time.Now).StorageKey("createdAt"),
	}
}

func (Profile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("profile").
			Unique().
			Required().
			Field("user_id"),
	}
}

func (Profile) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "Profile"},
	}
}
