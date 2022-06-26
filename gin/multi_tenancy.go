package gin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-saas/saas"

	"github.com/go-saas/saas/data"
	"github.com/go-saas/saas/http"
)

type ErrorFormatter func(context *gin.Context, err error)

var (
	DefaultErrorFormatter ErrorFormatter = func(context *gin.Context, err error) {
		if errors.Is(err, saas.ErrTenantNotFound) {
			context.AbortWithError(404, err)
		} else {
			context.AbortWithError(500, err)
		}
	}
)

type option struct {
	hmtOpt  *http.WebMultiTenancyOption
	ef      ErrorFormatter
	resolve []saas.ResolveOption
}

type Option func(*option)

func WithMultiTenancyOption(opt *http.WebMultiTenancyOption) Option {
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

func MultiTenancy(ts saas.TenantStore, options ...Option) gin.HandlerFunc {
	opt := &option{
		hmtOpt:  http.NewDefaultWebMultiTenancyOption(),
		ef:      DefaultErrorFormatter,
		resolve: nil,
	}
	for _, o := range options {
		o(opt)
	}
	return func(context *gin.Context) {
		var trOpt []saas.ResolveOption
		df := []saas.TenantResolveContrib{
			http.NewCookieTenantResolveContrib(opt.hmtOpt.TenantKey, context.Request),
			http.NewFormTenantResolveContrib(opt.hmtOpt.TenantKey, context.Request),
			http.NewHeaderTenantResolveContrib(opt.hmtOpt.TenantKey, context.Request),
			http.NewQueryTenantResolveContrib(opt.hmtOpt.TenantKey, context.Request)}
		if opt.hmtOpt.DomainFormat != "" {
			df = append(df, http.NewDomainTenantResolveContrib(opt.hmtOpt.DomainFormat, context.Request))
		}
		df = append(df, saas.NewTenantNormalizerContrib(ts))
		trOpt = append(trOpt, saas.AppendContribs(df...))
		trOpt = append(trOpt, opt.resolve...)

		//get tenant config
		tenantConfigProvider := saas.NewDefaultTenantConfigProvider(saas.NewDefaultTenantResolver(trOpt...), ts)
		tenantConfig, ctx, err := tenantConfigProvider.Get(context)
		if err != nil {
			opt.ef(context, err)
			return
		}
		//set current tenant
		newContext := saas.NewCurrentTenant(ctx, tenantConfig.ID, tenantConfig.Name)
		//data filter
		newContext = data.NewMultiTenancyDataFilter(newContext)

		//with newContext
		context.Request = context.Request.WithContext(newContext)
		//next
		context.Next()

	}
}
