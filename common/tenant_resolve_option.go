package common

type PatchTenantResolveOption func(resolveOption *TenantResolveOption)

type TenantResolveOption struct {
	Resolvers []TenantResolveContributor
}

func NewTenantResolveOption(c ...TenantResolveContributor) *TenantResolveOption {
	return &TenantResolveOption{
		Resolvers: c,
	}
}

func (opt *TenantResolveOption) AppendContributors(c ...TenantResolveContributor) *TenantResolveOption {
	opt.Resolvers = append(opt.Resolvers, c...)
	return opt
}

func (opt *TenantResolveOption) RemoveContributors(c ...TenantResolveContributor) *TenantResolveOption {
	var r []TenantResolveContributor
	for _, resolver := range opt.Resolvers {
		if !contains(c, resolver) {
			r = append(r, resolver)
		}
	}
	opt.Resolvers = r
	return opt
}

func contains(a []TenantResolveContributor, b TenantResolveContributor) bool {
	for i := 0; i < len(a); i++ {
		if a[i] == b {
			return true
		}
	}
	return false
}
