package http

import (
	"github.com/go-saas/saas"
	"net/http"
	"regexp"
)

type DomainTenantResolveContrib struct {
	request *http.Request
	format  string
}

func NewDomainTenantResolveContrib(f string, r *http.Request) *DomainTenantResolveContrib {
	return &DomainTenantResolveContrib{
		request: r,
		format:  f,
	}
}

func (h *DomainTenantResolveContrib) Name() string {
	return "Domain"
}

func (h *DomainTenantResolveContrib) Resolve(ctx *saas.Context) error {
	host := h.request.Host
	r := regexp.MustCompile(h.format)
	f := r.FindAllStringSubmatch(host, -1)
	if f == nil {
		//no match
		return nil
	}
	ctx.TenantIdOrName = f[0][1]
	return nil
}
