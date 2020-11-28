package gorm

import (
	"context"
	"github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/go-saas/management/gorm/entity"
	g "gorm.io/gorm"
)

const ConnKey = "Tenant"

func GetDb(ctx context.Context, provider gorm.DbProvider) *g.DB {
	return provider.Get(ctx, ConnKey)
}

func AutoMigrate(f func(*g.DB), db *g.DB) error {
	if f != nil {
		f(db)
	}
	return db.AutoMigrate(
		new(entity.Tenant),
		new(entity.TenantConn),
		new(entity.TenantFeature),
	)
}

func AutoMigrateMySQL(f func(*g.DB), db *g.DB) error {
	if f != nil {
		f(db)
	}
	db.Set("gorm:table_options", "ENGINE=InnoDB")
	return db.AutoMigrate(
		new(entity.Tenant),
		new(entity.TenantConn),
		new(entity.TenantFeature),
	)
}
