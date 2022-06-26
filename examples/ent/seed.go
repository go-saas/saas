package main

import (
	"context"
	"fmt"
	"github.com/go-saas/saas/data"
	"github.com/go-saas/saas/examples/ent/shared/ent"
	"github.com/go-saas/saas/seed"
)

type Seed struct {
	shared SharedDbProvider
	tenant TenantDbProvider
}

func NewSeed(shared SharedDbProvider, tenant TenantDbProvider) *Seed {
	return &Seed{shared: shared, tenant: tenant}
}

func (s *Seed) Seed(ctx context.Context, sCtx *seed.Context) error {

	if sCtx.TenantId == "" {
		//seed host
		c := s.shared.Get(ctx, "")

		c3, err := c.TenantConn.Create().SetKey(data.Default).SetValue("./tenant3.db?_fk=1").Save(ctx)
		if err != nil {
			return err
		}

		tenants := make([]*ent.TenantCreate, 3)
		tenants[0] = c.Tenant.Create().SetID(1).SetName("Test1")
		tenants[1] = c.Tenant.Create().SetID(2).SetName("Test2")
		tenants[2] = c.Tenant.Create().SetID(3).SetName("Test3").AddConn(c3)

		err = c.Tenant.CreateBulk(tenants...).OnConflict().UpdateNewValues().Exec(ctx)
		if err != nil {
			return err
		}

		if err := c.Post.Create().SetID(1).SetTitle("Host Side").SetDescription("Hello Host").OnConflict().UpdateNewValues().Exec(ctx); err != nil {
			return err
		}

	} else if sCtx.TenantId == "1" {
		c := s.tenant.Get(ctx, "")
		for i := 1; i < 2; i++ {
			if err := c.Post.Create().SetID(10 + i).SetTitle(fmt.Sprintf("Tenant %s Post %v", sCtx.TenantId, i)).
				SetDescription(fmt.Sprintf("Tenant %s ", sCtx.TenantId)).OnConflict().UpdateNewValues().Exec(ctx); err != nil {
				return err
			}
		}

	} else if sCtx.TenantId == "2" {
		c := s.tenant.Get(ctx, "")
		for i := 1; i < 3; i++ {
			if err := c.Post.Create().SetID(20 + i).SetTitle(fmt.Sprintf("Tenant %s Post %v", sCtx.TenantId, i)).
				SetDescription(fmt.Sprintf("Tenant %s ", sCtx.TenantId)).OnConflict().UpdateNewValues().Exec(ctx); err != nil {
				return err
			}
		}
	} else if sCtx.TenantId == "3" {
		c := s.tenant.Get(ctx, "")
		for i := 1; i < 4; i++ {
			if err := c.Post.Create().SetID(30 + i).SetTitle(fmt.Sprintf("Tenant %s Post %v", sCtx.TenantId, i)).
				SetDescription(fmt.Sprintf("Tenant %s ", sCtx.TenantId)).OnConflict().UpdateNewValues().Exec(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}
