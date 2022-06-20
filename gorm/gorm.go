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

type DbWrap struct {
	*gorm.DB
}

// NewDbWrap wrap gorm.DB into closable
func NewDbWrap(db *gorm.DB) *DbWrap {
	return &DbWrap{db}
}

func (d *DbWrap) Close() error {
	return closeDb(d.DB)
}

func closeDb(d *gorm.DB) error {
	sqlDB, err := d.DB()
	if err != nil {
		return err
	}
	cErr := sqlDB.Close()
	if cErr != nil {
		//todo logging
		//logger.Errorf("Gorm db close error: %s", err.Error())
		return cErr
	}
	return nil
}
