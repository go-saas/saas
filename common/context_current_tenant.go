package common

import (
	"context"
)

type ContextCurrentTenant struct {

}


func (c ContextCurrentTenant) IsAvailable(ctx context.Context) bool {
	current := getCurrent(ctx)
	return current.Id != ""
}

func (c ContextCurrentTenant) Id(ctx context.Context) string {
	current := getCurrent(ctx)
	return  current.Id
}

func (c ContextCurrentTenant) Change(ctx context.Context, id string, name string) (context.Context,CancelFunc) {
	current := getCurrent(ctx)
	newCtx := NewCurrentTenant(ctx,id,name)
	//TODO ??? performance
	return newCtx,func() context.Context {
		return NewCurrentTenant(ctx,current.Id,current.Name)
	}
}

func getCurrent(ctx context.Context) BasicTenantInfo {
	return FromCurrentTenant(ctx)
}

