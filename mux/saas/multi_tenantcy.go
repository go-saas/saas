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
			http.NewCookieTenantResolveContributor(hmtOpt.TenantKey, r),
			http.NewFormTenantResolveContributor(hmtOpt.TenantKey, r),
			http.NewHeaderTenantResolveContributor(hmtOpt.TenantKey, r),
			http.NewQueryTenantResolveContributor(hmtOpt.TenantKey, r),
		}
		if hmtOpt.DomainFormat != "" {
			df := append(df[:1], df[0:]...)
			df[0] = http.NewDomainTenantResolveContributor(hmtOpt.DomainFormat, r)
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
