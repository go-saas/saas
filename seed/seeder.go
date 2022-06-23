package seed

import (
	"context"
	"github.com/goxiaoy/go-saas"
)

type Seeder interface {
	Seed(ctx context.Context, option *Option) error
}

var _ Seeder = (*DefaultSeeder)(nil)

type DefaultSeeder struct {
	contrib []Contrib
}

func NewDefaultSeeder(contrib ...Contrib) *DefaultSeeder {
	return &DefaultSeeder{
		contrib: contrib,
	}
}

func (d *DefaultSeeder) Seed(ctx context.Context, option *Option) error {
	for _, tenant := range option.TenantIds {
		// change to next tenant
		ctx = saas.NewCurrentTenant(ctx, tenant, "")

		seedFn := func(ctx context.Context) error {
			sCtx := NewSeedContext(tenant, option.Extra)
			//create seeder
			for _, contributor := range d.contrib {
				if err := contributor.Seed(ctx, sCtx); err != nil {
					return err
				}
			}
			return nil
		}
		if err := seedFn(ctx); err != nil {
			return err
		}
	}
	return nil
}
