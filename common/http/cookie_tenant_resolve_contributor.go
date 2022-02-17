package http

import (
	"github.com/goxiaoy/go-saas/common"
	"net/http"
)

type CookieTenantResolveContributor struct {
	key     string
	request *http.Request
}

func NewCookieTenantResolveContributor(key string, r *http.Request) *CookieTenantResolveContributor {
	return &CookieTenantResolveContributor{
		key:     key,
		request: r,
	}
}

func (h *CookieTenantResolveContributor) Name() string {
	return "Cookie"
}

func (h *CookieTenantResolveContributor) Resolve(trCtx *common.TenantResolveContext) {
	v, err := h.request.Cookie(h.key)
	if err != nil {
		//no cookie
		return
	}
	if v.Value == "" {
		return
	}
	trCtx.TenantIdOrName = v.Value
}
