package common

import (
	"context"
	"github.com/goxiaoy/go-saas/data"
)

type MultiTenancyConnStrResolver struct {
	currentTenant CurrentTenant
	ts TenantStore
	*data.DefaultConnStrResolver
}

func NewMultiTenancyConnStrResolver(currentTenant CurrentTenant,ts TenantStore,opt data.ConnStrOption) *MultiTenancyConnStrResolver {
	return &MultiTenancyConnStrResolver{
		currentTenant:          currentTenant,
		ts:                     ts,
		DefaultConnStrResolver: &data.DefaultConnStrResolver{Opt: opt},
	}
}


//direct return value from option value
func (m MultiTenancyConnStrResolver) Resolve(ctx context.Context, key string) string {
	id := m.currentTenant.Id(ctx)
	if !m.currentTenant.IsAvailable(ctx){
		//use default
		return m.DefaultConnStrResolver.Resolve(ctx,key)
	}
	tenant,_ := m.ts.GetByNameOrId(ctx,id)
	if tenant.Conn ==nil{
		//not found
		//use default
		return m.DefaultConnStrResolver.Resolve(ctx,key)
	}
	if key==""{
		//get default
		ret := (*tenant).Conn.Default()
		if ret==""{
			return m.Opt.Conn.Default()
		}
		return  ret
	}
	//get key
	ret := tenant.Conn.GetOrDefault(key)
	if ret!=""{
		return ret
	}
	ret = m.Opt.Conn.GetOrDefault(key)
	if ret!=""{
		return ret
	}
	//still not found. fallback
	ret = (*tenant).Conn.Default()
	if ret==""{
		return m.Opt.Conn.Default()
	}
	return  ret
}

