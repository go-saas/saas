package schema

import (
	"context"
	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/go-saas/saas/examples/ent/tenant/ent/hook"
	"github.com/go-saas/saas/examples/ent/tenant/ent/intercept"

	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/go-saas/saas"
	"github.com/go-saas/saas/data"
)

type HasTenant struct {
	mixin.Schema
}

func (HasTenant) Fields() []ent.Field {
	return []ent.Field{
		field.String("tenant_id").Optional().GoType(&sql.NullString{}),
	}
}

func (h HasTenant) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		intercept.TraverseFunc(func(ctx context.Context, q intercept.Query) error {
			e := data.FromMultiTenancyDataFilter(ctx)
			if !e {
				// Skip tenant filter
				return nil
			}
			ct, _ := saas.FromCurrentTenant(ctx)
			h.P(ct, q)
			return nil
		}),
	}
}

func (h HasTenant) P(t saas.TenantInfo, w interface{ WhereP(...func(*sql.Selector)) }) {
	if len(t.GetId()) == 0 {
		w.WhereP(
			sql.FieldIsNull(h.Fields()[0].Descriptor().Name))
		return
	}
	w.WhereP(
		sql.FieldEQ(h.Fields()[0].Descriptor().Name, t.GetId()))
}

func (h HasTenant) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				type hasTenant interface {
					SetOp(ent.Op)
					SetTenantID(ss *sql.NullString)
					WhereP(...func(*sql.Selector))
				}
				return ent.MutateFunc(func(ctx context.Context, mutation ent.Mutation) (ent.Value, error) {
					if hf, ok := mutation.(hasTenant); ok {
						ct, _ := saas.FromCurrentTenant(ctx)
						at := data.FromAutoSetTenantId(ctx)
						if ok && at {
							if ct.GetId() != "" {
								//normalize tenant side only
								hf.SetTenantID(&sql.NullString{
									String: ct.GetId(),
									Valid:  true,
								})
							}
						}
					}
					return next.Mutate(ctx, mutation)
				})
			},
			ent.OpCreate|ent.OpUpdate|ent.OpUpdateOne,
		),
	}
}
