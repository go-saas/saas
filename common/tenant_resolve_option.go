package common

type TenantResolveOption struct {
	Resolvers []TenantResolveContributor
}

type ResolveOption func(resolveOption *TenantResolveOption)

func AppendContributors(c ...TenantResolveContributor) ResolveOption {
	return func(resolveOption *TenantResolveOption) {
		resolveOption.AppendContributors(c...)
	}
}

func RemoveContributors(c ...TenantResolveContributor) ResolveOption {
	return func(resolveOption *TenantResolveOption) {
		resolveOption.RemoveContributors(c...)
	}
}

func NewTenantResolveOption(c ...TenantResolveContributor) *TenantResolveOption {
	return &TenantResolveOption{
		Resolvers: c,
	}
}

func (opt *TenantResolveOption) AppendContributors(c ...TenantResolveContributor) {
	opt.Resolvers = append(opt.Resolvers, c...)
}

func (opt *TenantResolveOption) RemoveContributors(c ...TenantResolveContributor) {
	var r []TenantResolveContributor
	for _, resolver := range opt.Resolvers {
		if !contains(c, resolver) {
			r = append(r, resolver)
		}
	}
	opt.Resolvers = r
}

func contains(a []TenantResolveContributor, b TenantResolveContributor) bool {
	for i := 0; i < len(a); i++ {
		if a[i] == b {
			return true
		}
	}
	return false
}
