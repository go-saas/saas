package http

import (
	"github.com/goxiaoy/go-saas/common"
	"net/http"
)

type CookieTenantResolveContributor struct {
	opt     WebMultiTenancyOption
	request *http.Request
}

func NewCookieTenantResolveContributor(opt WebMultiTenancyOption, r *http.Request) *CookieTenantResolveContributor {
	return &CookieTenantResolveContributor{
		opt:     opt,
		request: r,
	}
}

func (h *CookieTenantResolveContributor) Name() string {
	return "Cookie"
}

func (h *CookieTenantResolveContributor) Resolve(trCtx *common.TenantResolveContext) {
	v, err := h.request.Cookie(h.opt.TenantKey)
	if err != nil {
		//no cookie
		return
	}
	if v.Value == "" {
		return
	}
	trCtx.TenantIdOrName = v.Value
}
