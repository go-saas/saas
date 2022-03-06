package saas

import (
	"github.com/goxiaoy/go-saas/common"
	shttp "github.com/goxiaoy/go-saas/common/http"
	"github.com/goxiaoy/go-saas/data"
	"net/http"
)

type MultiTenancy struct {
	hmtOpt *shttp.WebMultiTenancyOption
	trOptF []common.PatchTenantResolveOption
	ts     common.TenantStore
	ef     ErrorFormatter
}

type ErrorFormatter func(w http.ResponseWriter, err error)

var (
	DefaultErrorFormatter ErrorFormatter = func(w http.ResponseWriter, err error) {
		if err == common.ErrTenantNotFound {
			//not found
			http.Error(w, "Not Found", 404)
		} else {
			http.Error(w, err.Error(), 500)
		}
	}
)

func NewMultiTenancy(hmtOpt *shttp.WebMultiTenancyOption, ts common.TenantStore) *MultiTenancy {
	return &MultiTenancy{
		hmtOpt: hmtOpt,
		ts:     ts,
		ef:     DefaultErrorFormatter,
	}
}

func (m *MultiTenancy) WithErrorFormatter(ef ErrorFormatter) *MultiTenancy {
	m.ef = ef
	return m
}
func (m *MultiTenancy) WithOptions(options ...common.PatchTenantResolveOption) *MultiTenancy {
	m.trOptF = options
	return m
}

func (m *MultiTenancy) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hmtOpt := m.hmtOpt
		df := []common.TenantResolveContributor{
			shttp.NewCookieTenantResolveContributor(hmtOpt.TenantKey, r),
			shttp.NewFormTenantResolveContributor(hmtOpt.TenantKey, r),
			shttp.NewHeaderTenantResolveContributor(hmtOpt.TenantKey, r),
			shttp.NewQueryTenantResolveContributor(hmtOpt.TenantKey, r),
		}

		if hmtOpt.DomainFormat != "" {
			df = append(df, shttp.NewDomainTenantResolveContributor(hmtOpt.DomainFormat, r))
		}
		df = append(df, common.NewTenantNormalizerContributor(m.ts))
		trOpt := common.NewTenantResolveOption(df...)
		for _, option := range m.trOptF {
			option(trOpt)
		}
		//get tenant config
		tenantConfigProvider := common.NewDefaultTenantConfigProvider(common.NewDefaultTenantResolver(trOpt), common.NewCachedTenantStore(m.ts))
		tenantConfig, trCtx, err := tenantConfigProvider.Get(r.Context())
		if err != nil {
			m.ef(w, err)
			return
		}
		//set current tenant
		newContext := common.NewCurrentTenant(trCtx, tenantConfig.ID, tenantConfig.Name)
		//data filter
		newContext = data.NewEnableMultiTenancyDataFilter(newContext)
		next.ServeHTTP(w, r.WithContext(newContext))
	})
}
