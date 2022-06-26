package kratos

import (
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-saas/saas"
)

type HeaderTenantResolveContrib struct {
	key         string
	transporter transport.Transporter
}

func NewHeaderTenantResolveContrib(key string, transporter transport.Transporter) *HeaderTenantResolveContrib {
	return &HeaderTenantResolveContrib{
		key:         key,
		transporter: transporter,
	}
}
func (h *HeaderTenantResolveContrib) Name() string {
	return "KratosHeader"
}

func (h *HeaderTenantResolveContrib) Resolve(ctx *saas.Context) error {
	v := h.transporter.RequestHeader().Get(h.key)
	if v == "" {
		return nil
	}
	ctx.TenantIdOrName = v
	return nil
}
