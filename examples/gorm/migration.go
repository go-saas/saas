package main

import (
	"context"
	"github.com/go-saas/saas/gorm"
	"github.com/go-saas/saas/seed"
)

type MigrationSeeder struct {
	dbProvider gorm.DbProvider
}

func NewMigrationSeeder(dbProvider gorm.DbProvider) *MigrationSeeder {
	return &MigrationSeeder{dbProvider: dbProvider}
}

func (m *MigrationSeeder) Seed(ctx context.Context, sCtx *seed.Context) error {
	db := m.dbProvider.Get(ctx, "")
	if sCtx.TenantId == "" {
		//host add tenant database
		err := db.AutoMigrate(&Tenant{}, &TenantConn{})
		if err != nil {
			return err
		}
	}
	err := db.AutoMigrate(&Post{})
	if err != nil {
		return err
	}
	return nil
}
