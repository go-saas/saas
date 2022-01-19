package seed

import (
	"context"
	"github.com/goxiaoy/go-saas/common"
)

type Seeder interface {
	Seed(ctx context.Context) error
}

var _ Seeder = (*DefaultSeeder)(nil)

type DefaultSeeder struct {
	extra map[string]interface{}
	opt   *Option
}

func NewDefaultSeeder(opt *Option, extra map[string]interface{}) *DefaultSeeder {
	return &DefaultSeeder{
		opt:   opt,
		extra: extra,
	}
}

func (d *DefaultSeeder) Seed(ctx context.Context) error {
	for _, tenant := range d.opt.TenantIds {
		// change to next tenant
		newCtx := common.NewCurrentTenant(ctx, tenant, "")

		seedFn := func(ctx context.Context) error {
			sCtx := NewSeedContext(tenant, d.extra)
			//create seeder
			for _, contributor := range d.opt.Contributors {
				if err := contributor.Seed(ctx, sCtx); err != nil {
					return err
				}
			}
			return nil
		}
		if d.opt.uowMgr != nil {
			//run into uow
			err := d.opt.uowMgr.WithNew(newCtx, func(ctx context.Context) error {
				return seedFn(ctx)
			})
			if err != nil {
				return err
			}
		} else {
			return seedFn(newCtx)
		}
	}
	return nil
}
