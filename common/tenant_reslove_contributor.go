package common

type TenantResolveContributor interface {
	//Name of resolver
	Name() string
	//Chain of responsibility pattern
	Resolve(trCtx *TenantResolveContext)
}
