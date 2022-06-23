package saas

type TenantResolveOption struct {
	Resolvers []TenantResolveContrib
}

type ResolveOption func(resolveOption *TenantResolveOption)

func AppendContribs(c ...TenantResolveContrib) ResolveOption {
	return func(resolveOption *TenantResolveOption) {
		resolveOption.AppendContribs(c...)
	}
}

func RemoveContribs(c ...TenantResolveContrib) ResolveOption {
	return func(resolveOption *TenantResolveOption) {
		resolveOption.RemoveContribs(c...)
	}
}

func NewTenantResolveOption(c ...TenantResolveContrib) *TenantResolveOption {
	return &TenantResolveOption{
		Resolvers: c,
	}
}

func (opt *TenantResolveOption) AppendContribs(c ...TenantResolveContrib) {
	opt.Resolvers = append(opt.Resolvers, c...)
}

func (opt *TenantResolveOption) RemoveContribs(c ...TenantResolveContrib) {
	var r []TenantResolveContrib
	for _, resolver := range opt.Resolvers {
		if !contains(c, resolver) {
			r = append(r, resolver)
		}
	}
	opt.Resolvers = r
}

func contains(a []TenantResolveContrib, b TenantResolveContrib) bool {
	for i := 0; i < len(a); i++ {
		if a[i] == b {
			return true
		}
	}
	return false
}
