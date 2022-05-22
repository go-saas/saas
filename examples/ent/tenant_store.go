package main

import (
	"context"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/examples/ent/shared/ent"
	"github.com/goxiaoy/go-saas/examples/ent/shared/ent/tenant"
	"strconv"
)

type TenantStore struct {
	shared SharedDbProvider
}

func (t *TenantStore) GetByNameOrId(ctx context.Context, nameOrId string) (*common.TenantConfig, error) {
	ctx = common.NewCurrentTenant(ctx, "", "")
	db := t.shared.Get(ctx, "")
	i, err := strconv.Atoi(nameOrId)
	var te *ent.Tenant
	if err == nil {
		te, err = db.Tenant.Query().Where(tenant.Or(tenant.ID(i), tenant.Name(nameOrId))).First(ctx)
	} else {
		te, err = db.Tenant.Query().Where(tenant.Name(nameOrId)).First(ctx)
	}
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, common.ErrTenantNotFound
		}
	}
	ret := common.NewTenantConfig(strconv.Itoa(te.ID), te.Name, te.Region)
	conns, err := te.QueryConn().All(ctx)
	if err != nil {
		return nil, err
	}
	for _, conn := range conns {
		ret.Conn[conn.Key] = conn.Value
	}
	return ret, nil

}
