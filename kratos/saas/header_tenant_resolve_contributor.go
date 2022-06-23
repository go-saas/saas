package saas

import (
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goxiaoy/go-saas/common"
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

func (h *HeaderTenantResolveContrib) Resolve(ctx *common.Context) error {
	v := h.transporter.RequestHeader().Get(h.key)
	if v == "" {
		return nil
	}
	ctx.TenantIdOrName = v
	return nil
}
