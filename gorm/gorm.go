package gorm

import (
	"context"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"gorm.io/gorm"
)

// MultiTenancy entity
type MultiTenancy struct {
	TenantId HasTenant `gorm:"index"`
}

type DbProvider common.DbProvider[*gorm.DB]
type ClientProvider common.ClientProvider[*gorm.DB]
type ClientProviderFunc common.ClientProviderFunc[*gorm.DB]

func (c ClientProviderFunc) Get(ctx context.Context, dsn string) (*gorm.DB, error) {
	return c(ctx, dsn)
}

func NewDbProvider(cs data.ConnStrResolver, cp ClientProvider) DbProvider {
	return common.NewDbProvider[*gorm.DB](cs, cp)
}
