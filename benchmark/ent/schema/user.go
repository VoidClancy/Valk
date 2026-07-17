package schema

import (
	"entgo.io/ent"
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
		field.String("role").Default("student"),
		field.Int32("login_count").Default(0).StorageKey("loginCount"),
	}
}

func (User) Edges() []ent.Edge {
	return nil
}
