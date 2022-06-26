package main

import (
	"context"
	"github.com/go-saas/saas/seed"
)

type MigrationSeeder struct {
	shared SharedDbProvider
	tenant TenantDbProvider
}

func NewMigrationSeeder(shared SharedDbProvider, tenant TenantDbProvider) *MigrationSeeder {
	return &MigrationSeeder{shared: shared, tenant: tenant}
}

func (m *MigrationSeeder) Seed(ctx context.Context, sCtx *seed.Context) error {
	if sCtx.TenantId == "" {
		c := m.shared.Get(ctx, "")
		if err := c.Schema.Create(ctx); err != nil {
			return err
		}
	} else {
		c := m.tenant.Get(ctx, "")
		if err := c.Schema.Create(ctx); err != nil {
			return err
		}
	}
	return nil
}
