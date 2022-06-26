package http

import (
	"github.com/go-saas/saas"
	"net/http"
)

type HeaderTenantResolveContrib struct {
	key     string
	request *http.Request
}

func NewHeaderTenantResolveContrib(key string, r *http.Request) *HeaderTenantResolveContrib {
	return &HeaderTenantResolveContrib{
		key:     key,
		request: r,
	}
}

func (h *HeaderTenantResolveContrib) Name() string {
	return "Header"
}

func (h *HeaderTenantResolveContrib) Resolve(ctx *saas.Context) error {
	v := h.request.Header.Get(h.key)
	if v == "" {
		return nil
	}
	ctx.TenantIdOrName = v
	return nil
}
