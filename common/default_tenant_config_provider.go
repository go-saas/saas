package common

import "context"

type DefaultTenantConfigProvider struct {
	tr TenantResolver
	ts TenantStore
}

func NewDefaultTenantConfigProvider(tr TenantResolver,ts TenantStore) TenantConfigProvider  {
	return &DefaultTenantConfigProvider{
		tr: tr,
		ts: ts,
	}
}

func (d *DefaultTenantConfigProvider) Get(ctx context.Context) (TenantConfig,error) {
	//TODO how to cache??
	rr := d.tr.Resolve(ctx)
	if rr.TenantIdOrName!=""{
		//tenant side
		//get config from tenant store
		cfg,err := d.ts.getByNameOrId(rr.TenantIdOrName)
		if err!=nil{
			return TenantConfig{},err
		}
		return cfg,nil
		//check error
	}
	return TenantConfig{},nil

}



