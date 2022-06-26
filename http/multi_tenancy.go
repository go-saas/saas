package http

import (
	"github.com/go-saas/saas"
	"github.com/go-saas/saas/data"
	"net/http"
)

type ErrorFormatter func(w http.ResponseWriter, err error)

var (
	DefaultErrorFormatter ErrorFormatter = func(w http.ResponseWriter, err error) {
		if err == saas.ErrTenantNotFound {
			//not found
			http.Error(w, "Not Found", 404)
		} else {
			http.Error(w, err.Error(), 500)
		}
	}
)

type option struct {
	hmtOpt  *WebMultiTenancyOption
	ef      ErrorFormatter
	resolve []saas.ResolveOption
}

type Option func(*option)

func WithMultiTenancyOption(opt *WebMultiTenancyOption) Option {
	return func(o *option) {
		o.hmtOpt = opt
	}
}

func WithErrorFormatter(e ErrorFormatter) Option {
	return func(o *option) {
		o.ef = e
	}
}

func WithResolveOption(opt ...saas.ResolveOption) Option {
	return func(o *option) {
		o.resolve = opt
	}
}

func Middleware(ts saas.TenantStore, options ...Option) func(next http.Handler) http.Handler {
	opt := &option{
		hmtOpt:  NewDefaultWebMultiTenancyOption(),
		ef:      DefaultErrorFormatter,
		resolve: nil,
	}
	for _, o := range options {
		o(opt)
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var trOpt []saas.ResolveOption
			df := []saas.TenantResolveContrib{
				NewCookieTenantResolveContrib(opt.hmtOpt.TenantKey, r),
				NewFormTenantResolveContrib(opt.hmtOpt.TenantKey, r),
				NewHeaderTenantResolveContrib(opt.hmtOpt.TenantKey, r),
				NewQueryTenantResolveContrib(opt.hmtOpt.TenantKey, r),
			}

			if opt.hmtOpt.DomainFormat != "" {
				df = append(df, NewDomainTenantResolveContrib(opt.hmtOpt.DomainFormat, r))
			}
			df = append(df, saas.NewTenantNormalizerContrib(ts))
			trOpt = append(trOpt, saas.AppendContribs(df...))
			trOpt = append(trOpt, opt.resolve...)

			//get tenant config
			tenantConfigProvider := saas.NewDefaultTenantConfigProvider(saas.NewDefaultTenantResolver(trOpt...), ts)
			tenantConfig, ctx, err := tenantConfigProvider.Get(r.Context())
			if err != nil {
				opt.ef(w, err)
				return
			}
			//set current tenant
			newContext := saas.NewCurrentTenant(ctx, tenantConfig.ID, tenantConfig.Name)
			//data filter
			newContext = data.NewMultiTenancyDataFilter(newContext)
			next.ServeHTTP(w, r.WithContext(newContext))
		})
	}

}
