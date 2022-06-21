package common

import (
	"context"
	"github.com/goxiaoy/go-saas/data"
)

type TenantStoreCreator func() TenantStore

type MultiTenancyConnStrResolver struct {
	//use creator to prevent circular dependency
	tsc  TenantStoreCreator
	conn data.ConnStrings
}

var _ data.ConnStrResolver = (*MultiTenancyConnStrResolver)(nil)

// NewMultiTenancyConnStrResolver from tenant
func NewMultiTenancyConnStrResolver(tsc TenantStoreCreator, conn data.ConnStrings) *MultiTenancyConnStrResolver {
	return &MultiTenancyConnStrResolver{
		tsc:  tsc,
		conn: conn,
	}
}

func (m *MultiTenancyConnStrResolver) Resolve(ctx context.Context, key string) (string, error) {
	tenantInfo, _ := FromCurrentTenant(ctx)
	id := tenantInfo.GetId()
	if len(id) == 0 {
		//use default
		return m.conn.Resolve(ctx, key)
	}

	var tenantConfig *TenantConfig
	//read from cache
	if tenant, ok := FromTenantConfigContext(ctx, id); ok {
		tenantConfig = tenant
	} else {
		ts := m.tsc()
		tenant, err := ts.GetByNameOrId(ctx, id)
		if err != nil {
			return "", err
		}
		tenantConfig = tenant
	}

	if tenantConfig.Conn == nil {
		//not found
		//use default
		return m.conn.Resolve(ctx, key)
	}

	//get key
	ret, err := tenantConfig.Conn.Resolve(ctx, key)
	if err != nil {
		return "", err
	}
	if ret != "" {
		return ret, nil
	}
	//still not found
	return m.conn.Resolve(ctx, key)
}
