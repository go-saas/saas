package seed

import "context"

type SeedContributor interface {
	Seed(ctx context.Context, sCtx *Context)
}
