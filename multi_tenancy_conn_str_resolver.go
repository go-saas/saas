package saas

import (
	"context"
	"github.com/go-saas/saas/data"
)

type MultiTenancyConnStrResolver struct {
	//use creator to prevent circular dependency
	ts       TenantStore
	fallback data.ConnStrResolver
}

var _ data.ConnStrResolver = (*MultiTenancyConnStrResolver)(nil)

// NewMultiTenancyConnStrResolver from tenant
func NewMultiTenancyConnStrResolver(ts TenantStore, fallback data.ConnStrResolver) *MultiTenancyConnStrResolver {
	return &MultiTenancyConnStrResolver{
		ts:       ts,
		fallback: fallback,
	}
}

func (m *MultiTenancyConnStrResolver) Resolve(ctx context.Context, key string) (string, error) {
	tenantInfo, _ := FromCurrentTenant(ctx)
	id := tenantInfo.GetId()
	if len(id) == 0 {
		//skip query tenant store
		return m.fallback.Resolve(ctx, key)
	}

	var tenantConfig *TenantConfig
	//read from cache
	if tenant, ok := FromTenantConfigContext(ctx, id); ok {
		tenantConfig = tenant
	} else {
		tenant, err := m.ts.GetByNameOrId(ctx, id)
		if err != nil {
			return "", err
		}
		tenantConfig = tenant
	}

	if tenantConfig.Conn == nil {
		//not found
		return m.fallback.Resolve(ctx, key)
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
	return m.fallback.Resolve(ctx, key)
}
