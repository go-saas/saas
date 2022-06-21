package saas

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/common/http"
	"github.com/goxiaoy/go-saas/data"
)

type ErrorFormatter func(context *gin.Context, err error)

var (
	DefaultErrorFormatter ErrorFormatter = func(context *gin.Context, err error) {
		if errors.Is(err, common.ErrTenantNotFound) {
			context.AbortWithError(404, err)
		} else {
			context.AbortWithError(500, err)
		}
	}
)

type option struct {
	hmtOpt  *http.WebMultiTenancyOption
	ef      ErrorFormatter
	resolve []common.ResolveOption
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

func WithResolveOption(opt ...common.ResolveOption) Option {
	return func(o *option) {
		o.resolve = opt
	}
}

func MultiTenancy(ts common.TenantStore, options ...Option) gin.HandlerFunc {
	opt := &option{
		hmtOpt:  http.NewDefaultWebMultiTenancyOption(),
		ef:      DefaultErrorFormatter,
		resolve: nil,
	}
	for _, o := range options {
		o(opt)
	}
	return func(context *gin.Context) {
		var trOpt []common.ResolveOption
		df := []common.TenantResolveContributor{
			http.NewCookieTenantResolveContributor(opt.hmtOpt.TenantKey, context.Request),
			http.NewFormTenantResolveContributor(opt.hmtOpt.TenantKey, context.Request),
			http.NewHeaderTenantResolveContributor(opt.hmtOpt.TenantKey, context.Request),
			http.NewQueryTenantResolveContributor(opt.hmtOpt.TenantKey, context.Request)}
		if opt.hmtOpt.DomainFormat != "" {
			df = append(df, http.NewDomainTenantResolveContributor(opt.hmtOpt.DomainFormat, context.Request))
		}
		df = append(df, common.NewTenantNormalizerContributor(ts))
		trOpt = append(trOpt, common.AppendContributors(df...))
		trOpt = append(trOpt, opt.resolve...)

		//get tenant config
		tenantConfigProvider := common.NewDefaultTenantConfigProvider(common.NewDefaultTenantResolver(trOpt...), ts)
		tenantConfig, trCtx, err := tenantConfigProvider.Get(context)
		if err != nil {
			opt.ef(context, err)
			return
		}
		//set current tenant
		newContext := common.NewCurrentTenant(trCtx, tenantConfig.ID, tenantConfig.Name)
		//data filter
		newContext = data.NewMultiTenancyDataFilter(newContext)

		//with newContext
		context.Request = context.Request.WithContext(newContext)
		//next
		context.Next()

	}
}
