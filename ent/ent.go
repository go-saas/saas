package ent

import (
	"context"
	entgo "entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/field"
	"github.com/go-saas/saas"
	"github.com/go-saas/saas/data"
	"reflect"
)

type HasTenant struct {
	entgo.Schema
}

func (HasTenant) Fields() []entgo.Field {
	return []entgo.Field{
		field.String("tenant_id").Optional().GoType(&sql.NullString{}),
	}
}

type WhereP interface{ WhereP(...func(*sql.Selector)) }

func (h HasTenant) P(t saas.TenantInfo, w interface{ WhereP(...func(*sql.Selector)) }) {
	if len(t.GetId()) == 0 {
		w.WhereP(
			sql.FieldIsNull(h.Fields()[0].Descriptor().Name))
		return
	}
	w.WhereP(
		sql.FieldEQ(h.Fields()[0].Descriptor().Name, t.GetId()))
}

func (h HasTenant) Interceptors() []entgo.Interceptor {
	return []entgo.Interceptor{
		entgo.TraverseFunc(func(ctx context.Context, q entgo.Query) error {
			e := data.FromMultiTenancyDataFilter(ctx)
			if !e {
				// Skip tenant filter
				return nil
			}

			ct, _ := saas.FromCurrentTenant(ctx)
			//TODO we can not call WhereP directly because q does not implement it. So we need to use reflection
			addFilter := func(sqls ...func(*sql.Selector)) {
				in := make([]reflect.Value, len(sqls))
				for i := range in {
					in[i] = reflect.ValueOf(sqls[i])
				}
				reflect.ValueOf(q).MethodByName("Where").Call(in)
			}
			if len(ct.GetId()) == 0 {
				addFilter(sql.FieldIsNull(h.Fields()[0].Descriptor().Name))
			} else {
				addFilter(sql.FieldEQ(h.Fields()[0].Descriptor().Name, ct.GetId()))
			}
			//h.P(ct, q)
			return nil
		}),
	}
}

func (h HasTenant) Hooks() []entgo.Hook {
	return []entgo.Hook{
		On(
			func(next entgo.Mutator) entgo.Mutator {
				type hasTenant interface {
					SetOp(entgo.Op)
					SetTenantID(ss *sql.NullString)
					WhereP(...func(*sql.Selector))
				}
				return entgo.MutateFunc(func(ctx context.Context, mutation entgo.Mutation) (entgo.Value, error) {
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
			entgo.OpCreate|entgo.OpUpdate|entgo.OpUpdateOne|entgo.OpDeleteOne|entgo.OpDelete,
		),
	}
}

func On(hk entgo.Hook, op entgo.Op) entgo.Hook {
	return If(hk, HasOp(op))
}

// Condition is a hook condition function.
type Condition func(context.Context, entgo.Mutation) bool

// If executes the given hook under condition.
//
//	hook.If(ComputeAverage, And(HasFields(...), HasAddedFields(...)))
func If(hk entgo.Hook, cond Condition) entgo.Hook {
	return func(next entgo.Mutator) entgo.Mutator {
		return entgo.MutateFunc(func(ctx context.Context, m entgo.Mutation) (entgo.Value, error) {
			if cond(ctx, m) {
				return hk(next).Mutate(ctx, m)
			}
			return next.Mutate(ctx, m)
		})
	}
}

// HasOp is a condition testing mutation operation.
func HasOp(op entgo.Op) Condition {
	return func(_ context.Context, m entgo.Mutation) bool {
		return m.Op().Is(op)
	}
}
