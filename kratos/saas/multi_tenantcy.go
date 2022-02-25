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

func Server(hmtOpt *shttp.WebMultiTenancyOption, ts common.TenantStore, trOptF ...common.PatchTenantResolveOption) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			trOpt := common.NewTenantResolveOption()
			if tr, ok := transport.FromServerContext(ctx); ok {
				if ht, ok := tr.(*http.Transport); ok {
					r := ht.Request()
					df := []common.TenantResolveContributor{
						//TODO route
						shttp.NewCookieTenantResolveContributor(hmtOpt.TenantKey, r),
						shttp.NewFormTenantResolveContributor(hmtOpt.TenantKey, r),
						shttp.NewHeaderTenantResolveContributor(hmtOpt.TenantKey, r),
						shttp.NewQueryTenantResolveContributor(hmtOpt.TenantKey, r),
					}
					if hmtOpt.DomainFormat != "" {
						df := append(df[:1], df[0:]...)
						df[0] = shttp.NewDomainTenantResolveContributor(hmtOpt.DomainFormat, r)
					}
					trOpt.AppendContributors(df...)
				} else {
					trOpt.AppendContributors(NewHeaderTenantResolveContributor(hmtOpt.TenantKey, tr))
				}
				for _, option := range trOptF {
					option(trOpt)
				}

				//get tenant config
				tenantConfigProvider := common.NewDefaultTenantConfigProvider(common.NewDefaultTenantResolver(*trOpt), ts)
				tenantConfig, trCtx, err := tenantConfigProvider.Get(ctx, true)
				if err != nil {
					//not found
					if errors.Is(err, common.ErrTenantNotFound) {
						return nil, errors.NotFound("TENANT", err.Error())
					}
					return nil, err
				}
				newContext := common.NewCurrentTenant(trCtx, tenantConfig.ID, tenantConfig.Name)
				//data filter
				dataFilterCtx := data.NewEnableMultiTenancyDataFilter(newContext)
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
