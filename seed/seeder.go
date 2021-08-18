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
		sCtx := NewSeedContext(tenant, d.extra)
		//create seeder
		for _, contributor := range d.opt.Contributors {
			if err := contributor.Seed(newCtx, sCtx); err != nil {
				return err
			}
		}
	}
	return nil
}
