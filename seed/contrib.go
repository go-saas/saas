package seed

import (
	"context"
)

type Contrib interface {
	Seed(ctx context.Context, sCtx *Context) error
}

type chainContrib struct {
	seeds []Contrib
}

var _ Contrib = (*chainContrib)(nil)

func (c *chainContrib) Seed(ctx context.Context, sCtx *Context) error {
	for _, seed := range c.seeds {
		if err := seed.Seed(ctx, sCtx); err != nil {
			return err
		}
	}
	return nil
}

func Chain(seeds ...Contrib) Contrib {
	return &chainContrib{seeds: seeds}
}
