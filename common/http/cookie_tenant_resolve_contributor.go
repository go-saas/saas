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

func (h *CookieTenantResolveContributor) Resolve(trCtx *common.TenantResolveContext) error {
	v, err := h.request.Cookie(h.key)
	if err != nil {
		//no cookie
		return nil
	}
	if v.Value == "" {
		return nil
	}
	trCtx.TenantIdOrName = v.Value
	return nil
}
