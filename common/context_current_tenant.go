package common

import "context"

type ContextCurrentTenant struct {

}

type tenantKey string

const TenantIdKey tenantKey = "tenantId"

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
	newInfo := NewBasicTenantInfo(id,name)
	newCtx :=context.WithValue(ctx,TenantIdKey,*newInfo)
	//TODO ??? performance
	return newCtx,func() context.Context {
		return context.WithValue(newCtx,TenantIdKey,current)
	}
}

func getCurrent(ctx context.Context) BasicTenantInfo {
	value,ok:= ctx.Value(TenantIdKey).(BasicTenantInfo)
	if ok{
		return value
	}
	return BasicTenantInfo{}
}

