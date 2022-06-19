package data

import (
	"context"
)

type (
	//soft delete status
	multiTenancyDataFilterCtx struct{}
	autoSetTenantIdCtx        struct{}
)

func NewMultiTenancyDataFilter(ctx context.Context, enable ...bool) context.Context {
	v := true
	if len(enable) > 0 {
		v = enable[0]
	}
	return context.WithValue(ctx, multiTenancyDataFilterCtx{}, v)
}

//FromMultiTenancyDataFilter resolve where apply multi tenancy data filter, default true
func FromMultiTenancyDataFilter(ctx context.Context) bool {
	v := ctx.Value(multiTenancyDataFilterCtx{})
	if v == nil {
		return true
	}
	return v.(bool)
}

func NewAutoSetTenantId(ctx context.Context, enable ...bool) context.Context {
	v := true
	if len(enable) > 0 {
		v = enable[0]
	}
	return context.WithValue(ctx, autoSetTenantIdCtx{}, v)
}

func FromAutoSetTenantId(ctx context.Context) bool {
	v := ctx.Value(autoSetTenantIdCtx{})
	if v == nil {
		return true
	}
	return v.(bool)
}
