package data

import (
	"context"
)

type (
	//soft delete status
	multiTenancyDataFilterCtx struct{}
	autoSetTenantIdCtx        struct{}
)

func NewEnableMultiTenancyDataFilter(ctx context.Context) context.Context {
	return context.WithValue(ctx, multiTenancyDataFilterCtx{}, true)
}
func NewDisableMultiTenancyDataFilter(ctx context.Context) context.Context {
	return context.WithValue(ctx, multiTenancyDataFilterCtx{}, false)
}

//FromMultiTenancyDataFilter resolve where apply multi tenancy data filter, default true
func FromMultiTenancyDataFilter(ctx context.Context) bool {
	v := ctx.Value(multiTenancyDataFilterCtx{})
	if v == nil {
		return true
	}
	return v.(bool)
}
func NewEnableAutoSetTenantId(ctx context.Context) context.Context {
	return context.WithValue(ctx, autoSetTenantIdCtx{}, true)
}
func NewDisableAutoSetTenantId(ctx context.Context) context.Context {
	return context.WithValue(ctx, autoSetTenantIdCtx{}, false)
}
func FromAutoSetTenantId(ctx context.Context) bool {
	v := ctx.Value(autoSetTenantIdCtx{})
	if v == nil {
		return true
	}
	return v.(bool)
}
