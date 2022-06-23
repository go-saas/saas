package http

import (
	"github.com/goxiaoy/go-saas/common"
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

func (h *QueryTenantResolveContrib) Resolve(ctx *common.Context) error {
	v := h.request.URL.Query().Get(h.key)
	if v == "" {
		return nil
	}
	ctx.TenantIdOrName = v
	return nil
}
