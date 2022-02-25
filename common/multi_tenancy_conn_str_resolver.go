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

var _ data.ConnStrResolver = (*MultiTenancyConnStrResolver)(nil)

// NewMultiTenancyConnStrResolver from tenant
func NewMultiTenancyConnStrResolver(tsc TenantStoreCreator, opt *data.ConnStrOption) *MultiTenancyConnStrResolver {
	return &MultiTenancyConnStrResolver{
		tsc:                    tsc,
		DefaultConnStrResolver: data.NewDefaultConnStrResolver(opt),
	}
}

func (m *MultiTenancyConnStrResolver) Resolve(ctx context.Context, key string) (string, error) {
	tenantInfo := FromCurrentTenant(ctx)
	id := tenantInfo.GetId()
	if len(id) == 0 {
		//use default
		return m.DefaultConnStrResolver.Resolve(ctx, key)
	}
	ts := m.tsc()
	tenant, err := ts.GetByNameOrId(ctx, id)
	if err != nil {
		return "", err
	}
	if tenant.Conn == nil {
		//not found
		//use default
		return m.DefaultConnStrResolver.Resolve(ctx, key)
	}
	if key == "" {
		//get default
		ret := (*tenant).Conn.Default()
		if ret == "" {
			return m.Opt.Conn.Default(), nil
		}
		return ret, nil
	}
	//get key
	ret := tenant.Conn.GetOrDefault(key)
	if ret != "" {
		return ret, nil
	}
	ret = m.Opt.Conn.GetOrDefault(key)
	if ret != "" {
		return ret, nil
	}
	//still not found. fallback
	ret = (*tenant).Conn.Default()
	if ret == "" {
		return m.Opt.Conn.Default(), nil
	}
	return ret, nil
}
