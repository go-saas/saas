package common

type TenantResolveContributor interface {
	// Name of resolver
	Name() string
	// Resolve tenant
	Resolve(trCtx *TenantResolveContext)
}
