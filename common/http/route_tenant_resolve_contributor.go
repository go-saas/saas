package http

import (
	"github.com/goxiaoy/go-saas/common"
	"net/http"
)

type RouteTenantResolveContributor struct {
	opt     MultiTenancyOption
	request *http.Request
}

func NewRouteTenantResolveContributor(opt MultiTenancyOption,r *http.Request) * RouteTenantResolveContributor  {
	return &RouteTenantResolveContributor{
		opt: opt,
		request: r,
	}
}

func (h *RouteTenantResolveContributor) Name() string {
	return "Route"
}

func (h *RouteTenantResolveContributor) Resolve(trCtx *common.TenantResolveContext) {
	//TODO
	panic("NotImplemented")
}

