package data

import "context"

type Seeder interface {
	Seed(ctx context.Context, sCtx SeedContext)
}
