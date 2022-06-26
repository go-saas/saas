package http

import (
	"github.com/go-saas/saas"
	"net/http"
)

type FormTenantResolveContrib struct {
	key     string
	request *http.Request
}

func NewFormTenantResolveContrib(key string, r *http.Request) *FormTenantResolveContrib {
	return &FormTenantResolveContrib{
		key:     key,
		request: r,
	}
}

func (h *FormTenantResolveContrib) Name() string {
	return "Form"
}

func (h *FormTenantResolveContrib) Resolve(ctx *saas.Context) error {
	v := h.request.FormValue(h.key)
	if v == "" {
		return nil
	}
	ctx.TenantIdOrName = v
	return nil
}
