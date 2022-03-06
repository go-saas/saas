package saas

import (
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goxiaoy/go-saas/common"
)

type HeaderTenantResolveContributor struct {
	key         string
	transporter transport.Transporter
}

func NewHeaderTenantResolveContributor(key string, transporter transport.Transporter) *HeaderTenantResolveContributor {
	return &HeaderTenantResolveContributor{
		key:         key,
		transporter: transporter,
	}
}
func (h *HeaderTenantResolveContributor) Name() string {
	return "KratosHeader"
}

func (h *HeaderTenantResolveContributor) Resolve(trCtx *common.TenantResolveContext) error {
	v := h.transporter.RequestHeader().Get(h.key)
	if v == "" {
		return nil
	}
	trCtx.TenantIdOrName = v
	return nil
}
