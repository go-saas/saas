package main

import (
	"context"
	"github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/go-saas/seed"
)

type MigrationSeeder struct {
	dbProvider gorm.DbProvider
}

func NewMigrationSeeder(dbProvider gorm.DbProvider) *MigrationSeeder {
	return &MigrationSeeder{dbProvider: dbProvider}
}

func (m *MigrationSeeder) Seed(ctx context.Context, sCtx *seed.Context) error {
	db := m.dbProvider.Get(ctx, "")
	err := db.AutoMigrate(&Post{})
	if err != nil {
		return err
	}
	return nil
}
