package common

import "context"

type TenantResolver interface {
	Resolve(ctx context.Context) (TenantResolveResult, error)
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

func (d *DefaultTenantResolver) Resolve(ctx context.Context) (TenantResolveResult, error) {
	res := TenantResolveResult{}
	trCtx := TenantResolveContext{Context: ctx}
	for _, resolver := range d.o.Resolvers {
		if err := resolver.Resolve(&trCtx); err != nil {
			return res, err
		}
		res.AppliedResolvers = append(res.AppliedResolvers, resolver.Name())
		if trCtx.HasResolved() {
			break
		}
	}
	res.TenantIdOrName = trCtx.TenantIdOrName
	return res, nil
}
