package seed

import (
	"context"
	"github.com/goxiaoy/uow"
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

type UowContributor struct {
	uow uow.Manager
	up  Contributor
}

var _ Contributor = (*UowContributor)(nil)

func NewUowContributor(uow uow.Manager, up Contributor) *UowContributor {
	return &UowContributor{
		uow: uow,
		up:  up,
	}
}

func (u *UowContributor) Seed(ctx context.Context, sCtx *Context) error {
	return u.uow.WithNew(ctx, func(ctx context.Context) error {
		return u.up.Seed(ctx, sCtx)
	})
}
