package http

import (
	"github.com/goxiaoy/go-saas/common"
	"net/http"
	"regexp"
)

type DomainTenantResolveContributor struct {
	opt     WebMultiTenancyOption
	request *http.Request
	format  string
}

func NewDomainTenantResolveContributor(opt WebMultiTenancyOption,r *http.Request,f string) * DomainTenantResolveContributor  {
	return &DomainTenantResolveContributor{
		opt: opt,
		request: r,
		format: f,
	}
}

func (h *DomainTenantResolveContributor) Name() string {
	return "Domain"
}

func (h *DomainTenantResolveContributor) Resolve(trCtx *common.TenantResolveContext) {
	host := h.request.Host
	r:=regexp.MustCompile(h.format)
	f := r.FindAllStringSubmatch(host, -1)
	if f ==nil{
		//no match
		return
	}
	trCtx.TenantIdOrName=f[0][1]
}

