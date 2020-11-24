package common

import (
	"context"
)

type DefaultTenantResolver struct {
	//options
	o TenantResolveOption
	
}

func (d DefaultTenantResolver) Resolve(_ context.Context) TenantResolveResult {
	res := TenantResolveResult{}
	trCtx := TenantResolveContext{}
	for _, resolver := range d.o.resolvers {
		resolver.Resolve(&trCtx)
		res.AppliedResolvers=append(res.AppliedResolvers, resolver.Name())
		if trCtx.HasResolved() {
			//set
			res.TenantIdOrName=trCtx.TenantIdOrName
		}
	}
	return res
}

