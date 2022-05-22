package schema

import (
	"context"
	"database/sql"
	"entgo.io/ent"
	"entgo.io/ent/entql"
	"github.com/goxiaoy/go-saas/examples/ent/shared/ent/privacy"

	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
)

type HasTenant struct {
	mixin.Schema
}

func (HasTenant) Fields() []ent.Field {
	return []ent.Field{
		field.String("tenant_id").Optional().GoType(&sql.NullString{}),
	}
}

func (HasTenant) Policy() ent.Policy {
	return privacy.Policy{
		Query: privacy.QueryPolicy{
			FilterTenantRule(),
		},
	}
}

func FilterTenantRule() privacy.QueryMutationRule {
	type hasTenant interface {
		Where(p entql.P)
		WhereTenantID(p entql.StringP)
	}
	return privacy.FilterFunc(func(ctx context.Context, f privacy.Filter) error {
		ct, _ := common.FromCurrentTenant(ctx)
		e := data.FromMultiTenancyDataFilter(ctx)
		hf, ok := f.(hasTenant)
		if e && ok {
			//apply data filter
			if ct.GetId() == "" {
				//host side
				hf.Where(entql.FieldNil("tenant_id"))
			} else {
				hf.WhereTenantID(entql.StringEQ(ct.GetId()))
			}
		}

		return privacy.Skip
	})
}
