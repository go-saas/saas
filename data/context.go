package data

import (
	"context"
)

type (
	//soft delete status
	multiTenancyDataFilterCtx struct{}
)

func NewEnableMultiTenancyDataFilter(ctx context.Context)context.Context{
	return context.WithValue(ctx, multiTenancyDataFilterCtx{}, true)
}
func NewDisableMultiTenancyDataFilter(ctx context.Context)context.Context  {
	return context.WithValue(ctx, multiTenancyDataFilterCtx{}, false)
}
func FromMultiTenancyDataFilter(ctx context.Context) bool {
	v := ctx.Value(multiTenancyDataFilterCtx{})
	if v==nil{
		return true
	}
	return v.(bool)
}
