package common

import "context"

// TenantConfigProvider resolve tenant config from current context
type TenantConfigProvider interface {
	// Get tenant config
	Get(ctx context.Context) (TenantConfig, context.Context, error)
}

type DefaultTenantConfigProvider struct {
	tr TenantResolver
	ts TenantStore
}

func NewDefaultTenantConfigProvider(tr TenantResolver, ts TenantStore) TenantConfigProvider {
	return &DefaultTenantConfigProvider{
		tr: tr,
		ts: ts,
	}
}

func (d *DefaultTenantConfigProvider) Get(ctx context.Context) (TenantConfig, context.Context, error) {
	rr, ctx, err := d.tr.Resolve(ctx)
	if err != nil {
		return TenantConfig{}, ctx, err
	}
	if rr.TenantIdOrName != "" {
		//tenant side

		//read from cache
		if cfg, ok := FromTenantConfigContext(ctx, rr.TenantIdOrName); ok {
			return *cfg, ctx, nil
		}
		//get config from tenant store
		cfg, err := d.ts.GetByNameOrId(ctx, rr.TenantIdOrName)
		if err != nil {
			return TenantConfig{}, ctx, err
		}
		return *cfg, NewTenantConfigContext(ctx, cfg.ID, cfg), nil
	}
	// host side
	return TenantConfig{}, ctx, nil

}
