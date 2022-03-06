package http

import (
	"github.com/goxiaoy/go-saas/common"
	"net/http"
)

type FormTenantResolveContributor struct {
	key     string
	request *http.Request
}

func NewFormTenantResolveContributor(key string, r *http.Request) *FormTenantResolveContributor {
	return &FormTenantResolveContributor{
		key:     key,
		request: r,
	}
}

func (h *FormTenantResolveContributor) Name() string {
	return "Form"
}

func (h *FormTenantResolveContributor) Resolve(trCtx *common.TenantResolveContext) error {
	v := h.request.FormValue(h.key)
	if v == "" {
		return nil
	}
	trCtx.TenantIdOrName = v
	return nil
}
