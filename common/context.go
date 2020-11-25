package common

import (
	"context"
)

type (
	//soft delete status
	softDeleteCtx struct{}
	currentTenantCtx struct{}
	tenantResolveRes struct {}
)

func NewSoftDelete(ctx context.Context)context.Context{
	return context.WithValue(ctx, softDeleteCtx{}, true)
}
func NewNoSoftDelete(ctx context.Context)context.Context  {
	return context.WithValue(ctx, softDeleteCtx{}, false)
}
func FromSoftDelete(ctx context.Context) bool {
	v := ctx.Value(softDeleteCtx{})
	return v != nil && v.(bool)
}

func NewCurrentTenant(ctx context.Context, id string, name string)context.Context{
	newInfo := NewBasicTenantInfo(id,name)
	return context.WithValue(ctx, currentTenantCtx{},*newInfo)
}

func FromCurrentTenant(ctx context.Context) BasicTenantInfo {
	value,ok:= ctx.Value(currentTenantCtx{}).(BasicTenantInfo)
	if ok{
		return value
	}
	return BasicTenantInfo{}
}

func NewTenantResolveRes(ctx context.Context,t *TenantResolveResult) context.Context  {
	return context.WithValue(ctx, tenantResolveRes{}, t)
}

func FromTenantResolveRes(ctx context.Context) *TenantResolveResult {
	v,ok:= ctx.Value(tenantResolveRes{}).(*TenantResolveResult)
	if ok{
		return v
	}
	return nil
}
