package gorm

import (
	"context"
	"errors"
	"fmt"
	"github.com/goxiaoy/go-saas/data"
	"github.com/goxiaoy/uow"
	gorm2 "github.com/goxiaoy/uow/gorm"
	"gorm.io/gorm"
)

const DbKind = "gorm"

type DbProvider interface {
	// Get gorm db instance by key
	Get(ctx context.Context, key string) *gorm.DB
}

type DefaultDbProvider struct {
	cs     data.ConnStrResolver
	c      *Config
	opener DbOpener
}

func NewDefaultDbProvider(cs data.ConnStrResolver, c *Config, opener DbOpener) (d *DefaultDbProvider) {
	d = &DefaultDbProvider{
		cs:     cs,
		c:      c,
		opener: opener,
	}
	return
}

func (d *DefaultDbProvider) Get(ctx context.Context, key string) *gorm.DB {
	//resolve connection string
	s := d.cs.Resolve(ctx, key)
	// try resolve unit of work
	u, ok := uow.FromCurrentUow(ctx)
	if ok {
		// get transaction db form current unit of work
		tx, err := u.GetTxDb(ctx, DbKind, s)
		if err != nil {
			panic(err)
		}
		g, ok := tx.(*gorm2.TransactionDb)
		if !ok {
			panic(errors.New(fmt.Sprintf("%s is not a *gorm.DB instance", s)))
		}
		return g.DB
	}
	g, err := d.opener.Open(d.c, s)
	if err != nil {
		panic(err)
	}
	return g.WithContext(ctx)
}
