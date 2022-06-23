package common

type TenantResolveContrib interface {
	// Name of resolver
	Name() string
	// Resolve tenant
	Resolve(trCtx *TenantResolveContext) error
}

//TenantNormalizerContrib normalize tenant id or name into tenant id
type TenantNormalizerContrib struct {
	ts TenantStore
}

var _ TenantResolveContrib = (*TenantNormalizerContrib)(nil)

func NewTenantNormalizerContrib(ts TenantStore) *TenantNormalizerContrib {
	return &TenantNormalizerContrib{
		ts: ts,
	}
}
func (t *TenantNormalizerContrib) Name() string {
	return "TenantNormalizer"
}

func (t *TenantNormalizerContrib) Resolve(trCtx *TenantResolveContext) error {
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

// ContextContrib resolve from current context
type ContextContrib struct {
}

var _ TenantResolveContrib = (*ContextContrib)(nil)

func (c *ContextContrib) Name() string {
	return "ContextContrib"
}

func (c *ContextContrib) Resolve(trCtx *TenantResolveContext) error {
	info, ok := FromCurrentTenant(trCtx.Context())
	if ok {
		trCtx.TenantIdOrName = info.GetId()
		//terminate
		trCtx.HasHandled = true
	}
	return nil
}
