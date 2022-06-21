package common

type TenantResolveContributor interface {
	// Name of resolver
	Name() string
	// Resolve tenant
	Resolve(trCtx *TenantResolveContext) error
}

//TenantNormalizerContributor normalize tenant id or name into tenant id
type TenantNormalizerContributor struct {
	ts TenantStore
}

var _ TenantResolveContributor = (*TenantNormalizerContributor)(nil)

func NewTenantNormalizerContributor(ts TenantStore) *TenantNormalizerContributor {
	return &TenantNormalizerContributor{
		ts: ts,
	}
}
func (t *TenantNormalizerContributor) Name() string {
	return "TenantNormalizer"
}

func (t *TenantNormalizerContributor) Resolve(trCtx *TenantResolveContext) error {
	if len(trCtx.TenantIdOrName) > 0 {
		tenant, err := t.ts.GetByNameOrId(trCtx.Context(), trCtx.TenantIdOrName)
		if err != nil {
			return err
		}
		trCtx.TenantIdOrName = tenant.ID
		//store for cache
		trCtx.WithContext(NewTenantConfigContext(trCtx.Context(), tenant.ID, tenant))
	}
	return nil
}

// ContextContributor resolve from current context
type ContextContributor struct {
}

var _ TenantResolveContributor = (*ContextContributor)(nil)

func (c *ContextContributor) Name() string {
	return "ContextContributor"
}

func (c *ContextContributor) Resolve(trCtx *TenantResolveContext) error {
	info, ok := FromCurrentTenant(trCtx.Context())
	if ok {
		trCtx.TenantIdOrName = info.GetId()
		//terminate
		trCtx.HasHandled = true
	}
	return nil
}
