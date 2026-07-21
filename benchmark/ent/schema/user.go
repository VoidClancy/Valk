package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().StorageKey("id"),
		field.String("email").Unique(),
		field.String("phone_num").Unique().StorageKey("phoneNum"),
		field.String("password").Optional().Nillable(),
		field.String("role").Default("STUDENT"),
		field.String("role_optional").Optional().Nillable().StorageKey("roleOptional"),
		field.Int32("login_count").Default(0).StorageKey("loginCount"),
		field.String("referred_by_id").Optional().Nillable().StorageKey("referredById"),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("referrals", User.Type).
			From("referred_by").
			Unique().
			Field("referred_by_id"),
		edge.To("profile", Profile.Type).Unique(),
		edge.To("posts", Post.Type),
		edge.To("comments", Comment.Type),
	}
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "User"},
	}
}
