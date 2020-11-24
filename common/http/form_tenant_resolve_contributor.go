package http

import (
	"github.com/goxiaoy/go-saas/common"
	"net/http"
)

type FormTenantResolveContributor struct {
	opt     MultiTenancyOption
	request *http.Request
}

func NewFormTenantResolveContributor(opt MultiTenancyOption,r *http.Request) * FormTenantResolveContributor  {
	return &FormTenantResolveContributor{
		opt: opt,
		request: r,
	}
}

func (h *FormTenantResolveContributor) Name() string {
	return "Form"
}

func (h *FormTenantResolveContributor) Resolve(trCtx *common.TenantResolveContext) {
	v:=h.request.FormValue(h.opt.TenantKey)
	if v == ""{
		return
	}
	trCtx.TenantIdOrName = v
}

