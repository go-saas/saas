package gorm

import (
	"context"
	"github.com/go-saas/saas"

	"github.com/go-saas/saas/data"
	"gorm.io/gorm"
)

// MultiTenancy entity
type MultiTenancy struct {
	TenantId HasTenant `gorm:"index"`
}

type DbProvider saas.DbProvider[*gorm.DB]
type ClientProvider saas.ClientProvider[*gorm.DB]
type ClientProviderFunc saas.ClientProviderFunc[*gorm.DB]

func (c ClientProviderFunc) Get(ctx context.Context, dsn string) (*gorm.DB, error) {
	return c(ctx, dsn)
}

func NewDbProvider(cs data.ConnStrResolver, cp ClientProvider) DbProvider {
	return saas.NewDbProvider[*gorm.DB](cs, cp)
}

type DbWrap struct {
	*gorm.DB
}

// NewDbWrap wrap gorm.DB into io.Close
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
