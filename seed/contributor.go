package seed

import (
	"context"
)

type Contributor interface {
	Seed(ctx context.Context, sCtx *Context) error
}

type chainContributor struct {
	seeds []Contributor
}

var _ Contributor = (*chainContributor)(nil)

func (c *chainContributor) Seed(ctx context.Context, sCtx *Context) error {
	for _, seed := range c.seeds {
		if err := seed.Seed(ctx, sCtx); err != nil {
			return err
		}
	}
	return nil
}

func Chain(seeds ...Contributor) Contributor {
	return &chainContributor{seeds: seeds}
}
