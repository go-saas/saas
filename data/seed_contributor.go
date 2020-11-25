package data

import "context"

type SeedContributor interface {
	Seed(ctx context.Context,sCtx SeedContext)
}
