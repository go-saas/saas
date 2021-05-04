package middleware

import (
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/common/http"
	"github.com/goxiaoy/go-saas/data"
	netHttp "net/http"
)

type MultiTenancy struct {
	hmtOptF http.PatchHttpMultiTenancyOption
	trOptF  common.PatchTenantResolveOption
	ts      common.TenantStore
}

func (m *MultiTenancy) Middleware(next netHttp.Handler) netHttp.Handler {
	return netHttp.HandlerFunc(func(w netHttp.ResponseWriter, r *netHttp.Request) {

		hmtOpt := http.DefaultWebMultiTenancyOption()
		if m.hmtOptF != nil {
			//patch
			m.hmtOptF(hmtOpt)
		}
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
		if m.trOptF != nil {
			//patch
			m.trOptF(trOpt)
		}
		//get tenant config
		tenantConfigProvider := common.NewDefaultTenantConfigProvider(common.NewDefaultTenantResolver(*trOpt), m.ts)
		tenantConfig, trCtx, err := tenantConfigProvider.Get(r.Context(), true)
		if err != nil {
			//not found
			netHttp.Error(w, "Not Found", 404)
		}
		//set current tenant
		currentTenant := common.ContextCurrentTenant{}
		newContext, cancel := currentTenant.Change(trCtx, tenantConfig.ID, tenantConfig.Name)
		//data filter
		dataFilterCtx := data.NewEnableMultiTenancyDataFilter(newContext)
		//cancel
		defer cancel(dataFilterCtx)
		next.ServeHTTP(w, r.WithContext(dataFilterCtx))
	})
}
