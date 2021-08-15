package common

import "context"

type TenantResolver interface {
	Resolve(ctx context.Context) TenantResolveResult
}

type DefaultTenantResolver struct {
	//options
	o TenantResolveOption
}

func NewDefaultTenantResolver(o TenantResolveOption) TenantResolver {
	return &DefaultTenantResolver{
		o: o,
	}
}

func (d DefaultTenantResolver) Resolve(_ context.Context) TenantResolveResult {
	res := TenantResolveResult{}
	trCtx := TenantResolveContext{}
	for _, resolver := range d.o.Resolvers {
		resolver.Resolve(&trCtx)
		res.AppliedResolvers = append(res.AppliedResolvers, resolver.Name())
		if trCtx.HasResolved() {
			//set
			res.TenantIdOrName = trCtx.TenantIdOrName
			break
		}
	}
	return res
}
