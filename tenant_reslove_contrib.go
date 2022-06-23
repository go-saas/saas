package saas

type TenantResolveContrib interface {
	// Name of resolver
	Name() string
	// Resolve tenant
	Resolve(ctx *Context) error
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

func (t *TenantNormalizerContrib) Resolve(ctx *Context) error {
	if len(ctx.TenantIdOrName) > 0 {
		tenant, err := t.ts.GetByNameOrId(ctx.Context(), ctx.TenantIdOrName)
		if err != nil {
			return err
		}
		ctx.TenantIdOrName = tenant.ID
		//store for cache
		ctx.WithContext(NewTenantConfigContext(ctx.Context(), tenant.ID, tenant))
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

func (c *ContextContrib) Resolve(ctx *Context) error {
	info, ok := FromCurrentTenant(ctx.Context())
	if ok {
		ctx.TenantIdOrName = info.GetId()
		//terminate
		ctx.HasHandled = true
	}
	return nil
}
