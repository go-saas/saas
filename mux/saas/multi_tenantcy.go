package saas

import (
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/common/http"
	"github.com/goxiaoy/go-saas/data"
	netHttp "net/http"
)

type MultiTenancy struct {
	hmtOpt *http.WebMultiTenancyOption
	trOptF []common.PatchTenantResolveOption
	ts     common.TenantStore
}

func (m *MultiTenancy) Middleware(next netHttp.Handler) netHttp.Handler {
	return netHttp.HandlerFunc(func(w netHttp.ResponseWriter, r *netHttp.Request) {
		hmtOpt := m.hmtOpt
		df := []common.TenantResolveContributor{
			//TODO route
			http.NewCookieTenantResolveContributor(*hmtOpt, r),
			http.NewFormTenantResolveContributor(*hmtOpt, r),
			http.NewHeaderTenantResolveContributor(*hmtOpt, r),
			http.NewQueryTenantResolveContributor(*hmtOpt, r),
		}
		if hmtOpt.DomainFormat != "" {
			df := append(df[:1], df[0:]...)
			df[0] = http.NewDomainTenantResolveContributor(*hmtOpt, r, hmtOpt.DomainFormat)
		}
		trOpt := common.NewTenantResolveOption(df...)
		for _, option := range m.trOptF {
			option(trOpt)
		}
		//get tenant config
		tenantConfigProvider := common.NewDefaultTenantConfigProvider(common.NewDefaultTenantResolver(*trOpt), m.ts)
		tenantConfig, trCtx, err := tenantConfigProvider.Get(r.Context(), true)
		if err != nil {
			//not found
			netHttp.Error(w, "Not Found", 404)
		}
		//set current tenant
		newContext := common.NewCurrentTenant(trCtx, tenantConfig.ID, tenantConfig.Name)
		//data filter
		dataFilterCtx := data.NewEnableMultiTenancyDataFilter(newContext)
		next.ServeHTTP(w, r.WithContext(dataFilterCtx))
	})
}
