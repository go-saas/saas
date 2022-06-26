package seed

import (
	"context"
	"github.com/go-saas/saas"
)

type Seeder interface {
	Seed(ctx context.Context, option ...Option) error
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

func (d *DefaultSeeder) Seed(ctx context.Context, options ...Option) error {
	opt := NewOption()
	for _, option := range options {
		option(opt)
	}
	for _, tenant := range opt.TenantIds {
		// change to next tenant
		ctx = saas.NewCurrentTenant(ctx, tenant, "")

		seedFn := func(ctx context.Context) error {
			sCtx := NewSeedContext(tenant, opt.Extra)
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
