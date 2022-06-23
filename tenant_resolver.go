package saas

import "context"

type TenantResolver interface {
	Resolve(ctx context.Context) (TenantResolveResult, context.Context, error)
}

type DefaultTenantResolver struct {
	//options
	o *TenantResolveOption
}

func NewDefaultTenantResolver(opt ...ResolveOption) TenantResolver {
	o := NewTenantResolveOption(&ContextContrib{})
	for _, resolveOption := range opt {
		resolveOption(o)
	}
	return &DefaultTenantResolver{
		o: o,
	}
}

func (d *DefaultTenantResolver) Resolve(ctx context.Context) (TenantResolveResult, context.Context, error) {
	res := TenantResolveResult{}
	trCtx := NewTenantResolveContext(ctx)
	for _, resolver := range d.o.Resolvers {
		if err := resolver.Resolve(trCtx); err != nil {
			return res, trCtx.Context(), err
		}
		res.AppliedResolvers = append(res.AppliedResolvers, resolver.Name())
		if trCtx.HasResolved() {
			break
		}
	}
	res.TenantIdOrName = trCtx.TenantIdOrName
	ctx = NewTenantResolveRes(trCtx.Context(), &res)
	return res, ctx, nil
}
