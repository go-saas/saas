package common

import (
	"context"
	"github.com/goxiaoy/go-saas/data"
)

type (
	//ClientProvider resolve by dsn string
	ClientProvider[TClient interface{}] interface {
		Get(ctx context.Context, dsn string) (TClient, error)
	}

	ClientProviderFunc[TClient interface{}] func(ctx context.Context, dsn string) (TClient, error)

	DbProvider[TClient interface{}] interface {
		// Get instance by key
		Get(ctx context.Context, key string) TClient
	}

	DefaultDbProvider[TClient interface{}] struct {
		cs data.ConnStrResolver
		cp ClientProvider[TClient]
	}
)

func (c ClientProviderFunc[TClient]) Get(ctx context.Context, dsn string) (TClient, error) {
	return c(ctx, dsn)
}

func NewDbProvider[TClient interface{}](cs data.ConnStrResolver, cp ClientProvider[TClient]) (d *DefaultDbProvider[TClient]) {
	d = &DefaultDbProvider[TClient]{
		cs: cs,
		cp: cp,
	}
	return
}

func (d *DefaultDbProvider[TClient]) Get(ctx context.Context, key string) TClient {
	//resolve connection string
	s, err := d.cs.Resolve(ctx, key)
	if err != nil {
		panic(err)
	}
	c, err := d.cp.Get(ctx, s)
	if err != nil {
		panic(err)
	}
	return c
}
