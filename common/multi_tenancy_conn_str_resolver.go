package common

import (
	"context"
	"github.com/goxiaoy/go-saas/data"
)

type TenantStoreCreator func() TenantStore

type MultiTenancyConnStrResolver struct {
	//use creator to prevent circular dependency
	tsc TenantStoreCreator
	*data.DefaultConnStrResolver
}

// NewMultiTenancyConnStrResolver from tenant
func NewMultiTenancyConnStrResolver(tsc TenantStoreCreator, opt *data.ConnStrOption) *MultiTenancyConnStrResolver {
	return &MultiTenancyConnStrResolver{
		tsc:                    tsc,
		DefaultConnStrResolver: data.NewDefaultConnStrResolver(opt),
	}
}

func (m *MultiTenancyConnStrResolver) Resolve(ctx context.Context, key string) string {
	tenantInfo := FromCurrentTenant(ctx)
	id := tenantInfo.GetId()
	if len(id) == 0 {
		//use default
		return m.DefaultConnStrResolver.Resolve(ctx, key)
	}
	ts := m.tsc()
	tenant, _ := ts.GetByNameOrId(ctx, id)
	if tenant.Conn == nil {
		//not found
		//use default
		return m.DefaultConnStrResolver.Resolve(ctx, key)
	}
	if key == "" {
		//get default
		ret := (*tenant).Conn.Default()
		if ret == "" {
			return m.Opt.Conn.Default()
		}
		return ret
	}
	//get key
	ret := tenant.Conn.GetOrDefault(key)
	if ret != "" {
		return ret
	}
	ret = m.Opt.Conn.GetOrDefault(key)
	if ret != "" {
		return ret
	}
	//still not found. fallback
	ret = (*tenant).Conn.Default()
	if ret == "" {
		return m.Opt.Conn.Default()
	}
	return ret
}
