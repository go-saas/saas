package common

type PatchTenantResolveOption func(resolveOption * TenantResolveOption)
type TenantResolveOption struct {
	Resolvers []TenantResolveContributor
}

func NewTenantResolveOption(c ...TenantResolveContributor) *TenantResolveOption {
	return &TenantResolveOption{
		Resolvers: c,
	}
}