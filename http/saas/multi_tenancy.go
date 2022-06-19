package saas

import (
	"github.com/goxiaoy/go-saas/common"
	shttp "github.com/goxiaoy/go-saas/common/http"
	"github.com/goxiaoy/go-saas/data"
	"net/http"
)

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

type option struct {
	hmtOpt  *shttp.WebMultiTenancyOption
	ef      ErrorFormatter
	resolve []common.ResolveOption
}

type Option func(*option)

func WithMultiTenancyOption(opt *shttp.WebMultiTenancyOption) Option {
	return func(o *option) {
		o.hmtOpt = opt
	}
}

func WithErrorFormatter(e ErrorFormatter) Option {
	return func(o *option) {
		o.ef = e
	}
}

func WithResolveOption(opt ...common.ResolveOption) Option {
	return func(o *option) {
		o.resolve = opt
	}
}

func Middleware(ts common.TenantStore, options ...Option) func(next http.Handler) http.Handler {
	opt := &option{
		hmtOpt:  shttp.NewDefaultWebMultiTenancyOption(),
		ef:      DefaultErrorFormatter,
		resolve: nil,
	}
	for _, o := range options {
		o(opt)
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			df := []common.TenantResolveContributor{
				shttp.NewCookieTenantResolveContributor(opt.hmtOpt.TenantKey, r),
				shttp.NewFormTenantResolveContributor(opt.hmtOpt.TenantKey, r),
				shttp.NewHeaderTenantResolveContributor(opt.hmtOpt.TenantKey, r),
				shttp.NewQueryTenantResolveContributor(opt.hmtOpt.TenantKey, r),
			}

			if opt.hmtOpt.DomainFormat != "" {
				df = append(df, shttp.NewDomainTenantResolveContributor(opt.hmtOpt.DomainFormat, r))
			}
			df = append(df, common.NewTenantNormalizerContributor(ts))
			trOpt := common.NewTenantResolveOption(df...)
			for _, resolveOption := range opt.resolve {
				resolveOption(trOpt)
			}
			//get tenant config
			tenantConfigProvider := common.NewDefaultTenantConfigProvider(common.NewDefaultTenantResolver(trOpt), common.NewCachedTenantStore(ts))
			tenantConfig, trCtx, err := tenantConfigProvider.Get(r.Context())
			if err != nil {
				opt.ef(w, err)
				return
			}
			//set current tenant
			newContext := common.NewCurrentTenant(trCtx, tenantConfig.ID, tenantConfig.Name)
			//data filter
			newContext = data.NewMultiTenancyDataFilter(newContext)
			next.ServeHTTP(w, r.WithContext(newContext))
		})
	}

}
