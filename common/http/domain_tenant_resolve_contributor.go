package http

import (
	"github.com/goxiaoy/go-saas/common"
	"net/http"
	"regexp"
)

type DomainTenantResolveContributor struct {
	request *http.Request
	format  string
}

func NewDomainTenantResolveContributor(f string, r *http.Request) *DomainTenantResolveContributor {
	return &DomainTenantResolveContributor{
		request: r,
		format:  f,
	}
}

func (h *DomainTenantResolveContributor) Name() string {
	return "Domain"
}

func (h *DomainTenantResolveContributor) Resolve(trCtx *common.TenantResolveContext) error {
	host := h.request.Host
	r := regexp.MustCompile(h.format)
	f := r.FindAllStringSubmatch(host, -1)
	if f == nil {
		//no match
		return nil
	}
	trCtx.TenantIdOrName = f[0][1]
	return nil
}
