package http

import (
	"github.com/go-saas/saas"
	"net/http"
)

type QueryTenantResolveContrib struct {
	key     string
	request *http.Request
}

func NewQueryTenantResolveContrib(key string, r *http.Request) *QueryTenantResolveContrib {
	return &QueryTenantResolveContrib{
		key:     key,
		request: r,
	}
}

func (h *QueryTenantResolveContrib) Name() string {
	return "Query"
}

func (h *QueryTenantResolveContrib) Resolve(ctx *saas.Context) error {
	v := h.request.URL.Query().Get(h.key)
	if v == "" {
		return nil
	}
	ctx.TenantIdOrName = v
	return nil
}
