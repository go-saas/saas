package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	sent "github.com/go-saas/saas/ent"
)

// Post holds the schema definition for the Post entity.
type Post struct {
	ent.Schema
}

// Fields of the Post.
func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("title"),
		field.String("description").Optional(),
		field.String("dsn").Optional(),
	}
}

// Edges of the Post.
func (Post) Edges() []ent.Edge {
	return nil
}

func (Post) Mixin() []ent.Mixin {
	return []ent.Mixin{
		sent.HasTenant{},
	}
}
