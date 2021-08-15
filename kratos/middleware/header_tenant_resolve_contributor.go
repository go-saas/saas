package middleware

import (
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/common/http"
)

type HeaderTenantResolveContributor struct {
	opt         http.WebMultiTenancyOption
	transporter transport.Transporter
}

func NewHeaderTenantResolveContributor(opt http.WebMultiTenancyOption, transporter transport.Transporter) *HeaderTenantResolveContributor {
	return &HeaderTenantResolveContributor{
		opt:         opt,
		transporter: transporter,
	}
}
func (h *HeaderTenantResolveContributor) Name() string {
	return "KratosHeader"
}

func (h *HeaderTenantResolveContributor) Resolve(trCtx *common.TenantResolveContext) {
	v := h.transporter.RequestHeader().Get(h.opt.TenantKey)
	if v == "" {
		return
	}
	trCtx.TenantIdOrName = v
}
