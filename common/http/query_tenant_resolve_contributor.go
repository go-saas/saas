package http

import (
	"github.com/goxiaoy/go-saas/common"
	"net/http"
)

type QueryTenantResolveContributor struct {
	opt     WebMultiTenancyOption
	request *http.Request
}

func NewQueryTenantResolveContributor(opt WebMultiTenancyOption, r *http.Request) *QueryTenantResolveContributor {
	return &QueryTenantResolveContributor{
		opt:     opt,
		request: r,
	}
}

func (h *QueryTenantResolveContributor) Name() string {
	return "Query"
}

func (h *QueryTenantResolveContributor) Resolve(trCtx *common.TenantResolveContext) {
	v := h.request.URL.Query().Get(h.opt.TenantKey)
	if v == "" {
		return
	}
	trCtx.TenantIdOrName = v
}
