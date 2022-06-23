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

func (h *QueryTenantResolveContrib) Resolve(trCtx *common.TenantResolveContext) error {
	v := h.request.URL.Query().Get(h.key)
	if v == "" {
		return nil
	}
	trCtx.TenantIdOrName = v
	return nil
}
