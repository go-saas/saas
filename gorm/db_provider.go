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

type DbProvider interface {
	// Get gorm db instance by key
	Get(ctx context.Context, key string) *gorm.DB
}

type DefaultDbProvider struct {
	cs     data.ConnStrResolver
	um     uow.Manager
	c      *Config
	opener DbOpener
}

func NewDefaultDbProvider(cs data.ConnStrResolver, c *Config, um uow.Manager, opener DbOpener) (d *DefaultDbProvider) {
	d = &DefaultDbProvider{
		cs:     cs,
		c:      c,
		um:     um,
		opener: opener,
	}
	return
}

func (d *DefaultDbProvider) Get(ctx context.Context, key string) *gorm.DB {
	//resolve connection string
	s := d.cs.Resolve(ctx, key)
	fk := fmt.Sprintf("gorm_%s", s)
	u, ok := uow.FromCurrentUow(ctx)
	if ok {
		// get transaction db form current unit of work
		tx, err := u.GetTxDb(ctx, fk)
		if err != nil {
			panic(err)
		}
		g, ok := tx.(*gorm2.TransactionDb)
		if !ok {
			panic(errors.New(fmt.Sprintf("%s is not a *gorm.DB instance", fk)))
		}
		return g.DB
	}
	g, err := d.opener.Open(d.c, s)
	if err != nil {
		panic(err)
	}
	return g.WithContext(ctx)
}
