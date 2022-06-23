package saas

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/goxiaoy/go-saas/common"
	shttp "github.com/goxiaoy/go-saas/common/http"
	"github.com/goxiaoy/go-saas/data"
)

type ErrorFormatter func(err error) (interface{}, error)

var (
	DefaultErrorFormatter ErrorFormatter = func(err error) (interface{}, error) {
		//not found
		if errors.Is(err, common.ErrTenantNotFound) {
			return nil, errors.NotFound("TENANT", err.Error())
		}
		return nil, err
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

func Server(ts common.TenantStore, options ...Option) middleware.Middleware {
	opt := &option{
		hmtOpt:  shttp.NewDefaultWebMultiTenancyOption(),
		ef:      DefaultErrorFormatter,
		resolve: nil,
	}
	for _, o := range options {
		o(opt)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var trOpt []common.ResolveOption
			if tr, ok := transport.FromServerContext(ctx); ok {
				if ht, ok := tr.(*http.Transport); ok {
					r := ht.Request()
					df := []common.TenantResolveContrib{
						shttp.NewCookieTenantResolveContrib(opt.hmtOpt.TenantKey, r),
						shttp.NewFormTenantResolveContrib(opt.hmtOpt.TenantKey, r),
						shttp.NewHeaderTenantResolveContrib(opt.hmtOpt.TenantKey, r),
						shttp.NewQueryTenantResolveContrib(opt.hmtOpt.TenantKey, r),
					}
					if opt.hmtOpt.DomainFormat != "" {
						df = append(df, shttp.NewDomainTenantResolveContrib(opt.hmtOpt.DomainFormat, r))
					}
					df = append(df, common.NewTenantNormalizerContrib(ts))
					trOpt = append(trOpt, common.AppendContribs(df...))
				} else {
					trOpt = append(trOpt, common.AppendContribs(NewHeaderTenantResolveContrib(opt.hmtOpt.TenantKey, tr)))
				}
				trOpt = append(trOpt, opt.resolve...)

				//get tenant config
				tenantConfigProvider := common.NewDefaultTenantConfigProvider(common.NewDefaultTenantResolver(trOpt...), ts)
				tenantConfig, ctx, err := tenantConfigProvider.Get(ctx)
				if err != nil {
					return opt.ef(err)
				}
				newContext := common.NewCurrentTenant(ctx, tenantConfig.ID, tenantConfig.Name)
				//data filter
				dataFilterCtx := data.NewMultiTenancyDataFilter(newContext)
				return handler(dataFilterCtx, req)
			}
			return handler(ctx, req)
		}
	}
}

func Client(hmtOpt *shttp.WebMultiTenancyOption) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			ti, _ := common.FromCurrentTenant(ctx)
			if tr, ok := transport.FromClientContext(ctx); ok {
				if tr.Kind() == transport.KindHTTP {
					if ht, ok := tr.(*http.Transport); ok {
						ht.RequestHeader().Set(hmtOpt.TenantKey, ti.GetId())
					}
				} else if tr.Kind() == transport.KindGRPC {
					tr.RequestHeader().Set(hmtOpt.TenantKey, ti.GetName())
				}
			}
			return handler(ctx, req)
		}
	}
}
