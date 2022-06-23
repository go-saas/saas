package http

import (
	"github.com/goxiaoy/go-saas/common"
	"net/http"
)

type CookieTenantResolveContrib struct {
	key     string
	request *http.Request
}

func NewCookieTenantResolveContrib(key string, r *http.Request) *CookieTenantResolveContrib {
	return &CookieTenantResolveContrib{
		key:     key,
		request: r,
	}
}

func (h *CookieTenantResolveContrib) Name() string {
	return "Cookie"
}

func (h *CookieTenantResolveContrib) Resolve(ctx *common.Context) error {
	v, err := h.request.Cookie(h.key)
	if err != nil {
		//no cookie
		return nil
	}
	if v.Value == "" {
		return nil
	}
	ctx.TenantIdOrName = v.Value
	return nil
}
