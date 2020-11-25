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

func (d *DefaultTenantConfigProvider) Get(ctx context.Context,store bool) (TenantConfig,context.Context,error) {
	rr := d.tr.Resolve(ctx)
	rc := ctx
	if store{
		//store into context
		rc = NewTenantResolveRes(ctx,&rr)
	}
	if rr.TenantIdOrName!=""{
		//tenant side
		//get config from tenant store
		cfg,err := d.ts.getByNameOrId(rr.TenantIdOrName)
		if err!=nil{
			return TenantConfig{},rc,err
		}
		return *cfg,rc,nil
		//check error
	}
	return TenantConfig{},rc,nil

}



