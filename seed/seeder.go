package seed

import (
	"context"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/uow"
)

type Seeder interface {
	Seed(ctx context.Context) error
}

var _ Seeder = (*DefaultSeeder)(nil)

type DefaultSeeder struct {
	extra  map[string]interface{}
	opt    *Option
	uowMgr uow.Manager
}

func NewDefaultSeeder(opt *Option, uowMgr uow.Manager, extra map[string]interface{}) *DefaultSeeder {
	return &DefaultSeeder{
		opt:    opt,
		extra:  extra,
		uowMgr: uowMgr,
	}
}

func (d *DefaultSeeder) Seed(ctx context.Context) error {
	for _, tenant := range d.opt.TenantIds {
		// change to next tenant
		newCtx := common.NewCurrentTenant(ctx, tenant, "")
		err := d.uowMgr.WithNew(newCtx, func(ctx context.Context) error {
			sCtx := NewSeedContext(tenant, d.extra)
			//create seeder
			for _, contributor := range d.opt.Contributors {
				if err := contributor.Seed(ctx, sCtx); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}
