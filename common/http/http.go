package http

import (
	"github.com/google/wire"
	"github.com/goxiaoy/go-saas/common"
)

var HttpSet = wire.NewSet(
	DefaultWebMultiTenancyOption,
	//TODO route
	wire.Struct(new(CookieTenantResolveContributor)),
	wire.Bind(new(common.TenantResolveContributor),new(*CookieTenantResolveContributor)),
	wire.Struct(new(FormTenantResolveContributor)),
	wire.Bind(new(common.TenantResolveContributor),new(*FormTenantResolveContributor)),
	wire.Struct(new(HeaderTenantResolveContributor)),
	wire.Bind(new(common.TenantResolveContributor),new(*HeaderTenantResolveContributor)),
	wire.Struct(new(QueryTenantResolveContributor)),
	wire.Bind(new(common.TenantResolveContributor),new(*QueryTenantResolveContributor)),
	)
