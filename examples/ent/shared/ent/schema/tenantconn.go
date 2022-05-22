package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// TenantConn holds the schema definition for the TenantConn entity.
type TenantConn struct {
	ent.Schema
}

// Fields of the TenantConn.
func (TenantConn) Fields() []ent.Field {
	return []ent.Field{
		field.String("key"),
		field.String("value"),
	}
}

// Edges of the TenantConn.
func (TenantConn) Edges() []ent.Edge {
	return nil
}
func (TenantConn) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
