package saas

import (
	"context"
)

type (
	currentTenantCtx struct{}
	tenantResolveRes struct{}
	tenantConfigKey  string
)

func NewCurrentTenant(ctx context.Context, id, name string) context.Context {
	return NewCurrentTenantInfo(ctx, NewBasicTenantInfo(id, name))
}

func NewCurrentTenantInfo(ctx context.Context, info TenantInfo) context.Context {
	return context.WithValue(ctx, currentTenantCtx{}, info)
}

func FromCurrentTenant(ctx context.Context) (TenantInfo, bool) {
	value, ok := ctx.Value(currentTenantCtx{}).(TenantInfo)
	if ok {
		return value, true
	}
	return NewBasicTenantInfo("", ""), false
}

func NewTenantResolveRes(ctx context.Context, t *TenantResolveResult) context.Context {
	return context.WithValue(ctx, tenantResolveRes{}, t)
}

func FromTenantResolveRes(ctx context.Context) *TenantResolveResult {
	v, ok := ctx.Value(tenantResolveRes{}).(*TenantResolveResult)
	if ok {
		return v
	}
	return nil
}

func NewTenantConfigContext(ctx context.Context, tenantId string, cfg *TenantConfig) context.Context {
	return context.WithValue(ctx, tenantConfigKey(tenantId), cfg)
}

func FromTenantConfigContext(ctx context.Context, tenantId string) (*TenantConfig, bool) {
	v, ok := ctx.Value(tenantConfigKey(tenantId)).(*TenantConfig)
	if ok {
		return v, ok && v != nil
	}
	return nil, false
}
