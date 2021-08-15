package data

import "context"

type Seeder interface {
	Seed(ctx context.Context, sCtx SeedContext)
}

type DefaultSeeder struct {
	opt SeedOption
}

func (d DefaultSeeder) Seed(ctx context.Context, sCtx SeedContext) {
	//create seeder
	for _, contributor := range d.opt.Contributors {
		contributor.Seed(ctx, sCtx)
	}
}
