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
		//store
		trCtx.WithContext(NewTenantConfigContext(trCtx.Context(), tenant.ID, tenant))
	}
	return nil
}
