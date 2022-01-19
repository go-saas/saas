package common

import (
	"context"
)

type (
	currentTenantCtx struct{}
	tenantResolveRes struct{}
)

func NewCurrentTenant(ctx context.Context, id, name string) context.Context {
	return NewCurrentTenantInfo(ctx, NewBasicTenantInfo(id, name))
}

func NewCurrentTenantInfo(ctx context.Context, info TenantInfo) context.Context {
	return context.WithValue(ctx, currentTenantCtx{}, info)
}

func FromCurrentTenant(ctx context.Context) TenantInfo {
	value, ok := ctx.Value(currentTenantCtx{}).(TenantInfo)
	if ok {
		return value
	}
	return NewBasicTenantInfo("", "")
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
