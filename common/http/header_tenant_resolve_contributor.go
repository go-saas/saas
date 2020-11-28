package http

import (
	"github.com/goxiaoy/go-saas/common"
	"net/http"
)

type HeaderTenantResolveContributor struct {
	opt     WebMultiTenancyOption
	request *http.Request
}

func NewHeaderTenantResolveContributor(opt WebMultiTenancyOption, r *http.Request) *HeaderTenantResolveContributor {
	return &HeaderTenantResolveContributor{
		opt:     opt,
		request: r,
	}
}

func (h *HeaderTenantResolveContributor) Name() string {
	return "Header"
}

func (h *HeaderTenantResolveContributor) Resolve(trCtx *common.TenantResolveContext) {
	v := h.request.Header.Get(h.opt.TenantKey)
	if v == "" {
		return
	}
	trCtx.TenantIdOrName = v
}
