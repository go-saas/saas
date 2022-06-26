package kratos

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-saas/saas"
	"github.com/go-saas/saas/data"
	shttp "github.com/go-saas/saas/http"
)

type ErrorFormatter func(err error) (interface{}, error)

var (
	DefaultErrorFormatter ErrorFormatter = func(err error) (interface{}, error) {
		//not found
		if errors.Is(err, saas.ErrTenantNotFound) {
			return nil, errors.NotFound("TENANT", err.Error())
		}
		return nil, err
	}
)

type option struct {
	hmtOpt  *shttp.WebMultiTenancyOption
	ef      ErrorFormatter
	resolve []saas.ResolveOption
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

func WithResolveOption(opt ...saas.ResolveOption) Option {
	return func(o *option) {
		o.resolve = opt
	}
}

func Server(ts saas.TenantStore, options ...Option) middleware.Middleware {
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
			var trOpt []saas.ResolveOption
			if tr, ok := transport.FromServerContext(ctx); ok {
				if ht, ok := tr.(*http.Transport); ok {
					r := ht.Request()
					df := []saas.TenantResolveContrib{
						shttp.NewCookieTenantResolveContrib(opt.hmtOpt.TenantKey, r),
						shttp.NewFormTenantResolveContrib(opt.hmtOpt.TenantKey, r),
						shttp.NewHeaderTenantResolveContrib(opt.hmtOpt.TenantKey, r),
						shttp.NewQueryTenantResolveContrib(opt.hmtOpt.TenantKey, r),
					}
					if opt.hmtOpt.DomainFormat != "" {
						df = append(df, shttp.NewDomainTenantResolveContrib(opt.hmtOpt.DomainFormat, r))
					}
					df = append(df, saas.NewTenantNormalizerContrib(ts))
					trOpt = append(trOpt, saas.AppendContribs(df...))
				} else {
					trOpt = append(trOpt, saas.AppendContribs(NewHeaderTenantResolveContrib(opt.hmtOpt.TenantKey, tr)))
				}
				trOpt = append(trOpt, opt.resolve...)

				//get tenant config
				tenantConfigProvider := saas.NewDefaultTenantConfigProvider(saas.NewDefaultTenantResolver(trOpt...), ts)
				tenantConfig, ctx, err := tenantConfigProvider.Get(ctx)
				if err != nil {
					return opt.ef(err)
				}
				newContext := saas.NewCurrentTenant(ctx, tenantConfig.ID, tenantConfig.Name)
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
			ti, _ := saas.FromCurrentTenant(ctx)
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
