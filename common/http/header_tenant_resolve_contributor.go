package http

import (
	"github.com/goxiaoy/go-saas/common"
	"net/http"
)

type HeaderTenantResolveContributor struct {
	key     string
	request *http.Request
}

func NewHeaderTenantResolveContributor(key string, r *http.Request) *HeaderTenantResolveContributor {
	return &HeaderTenantResolveContributor{
		key:     key,
		request: r,
	}
}

func (h *HeaderTenantResolveContributor) Name() string {
	return "Header"
}

func (h *HeaderTenantResolveContributor) Resolve(trCtx *common.TenantResolveContext) error {
	v := h.request.Header.Get(h.key)
	if v == "" {
		return nil
	}
	trCtx.TenantIdOrName = v
	return nil
}
