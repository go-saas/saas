package http

import (
	"github.com/goxiaoy/go-saas/common"
	"net/http"
)

type QueryTenantResolveContributor struct {
	key     string
	request *http.Request
}

func NewQueryTenantResolveContributor(key string, r *http.Request) *QueryTenantResolveContributor {
	return &QueryTenantResolveContributor{
		key:     key,
		request: r,
	}
}

func (h *QueryTenantResolveContributor) Name() string {
	return "Query"
}

func (h *QueryTenantResolveContributor) Resolve(trCtx *common.TenantResolveContext) error {
	v := h.request.URL.Query().Get(h.key)
	if v == "" {
		return nil
	}
	trCtx.TenantIdOrName = v
	return nil
}
