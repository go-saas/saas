package middleware

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

func MultiTenancy(hmtOptF shttp.PatchHttpMultiTenancyOption, trOptF common.PatchTenantResolveOption, ts common.TenantStore) middleware.Middleware {
	hmtOpt := shttp.DefaultWebMultiTenancyOption()
	if hmtOptF != nil {
		//patch
		hmtOptF(hmtOpt)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			trOpt := common.NewTenantResolveOption()
			if tr, ok := transport.FromServerContext(ctx); ok {
				if ht, ok := tr.(*http.Transport); ok {
					r := ht.Request()
					df := []common.TenantResolveContributor{
						//TODO route
						shttp.NewCookieTenantResolveContributor(*hmtOpt, r),
						shttp.NewFormTenantResolveContributor(*hmtOpt, r),
						shttp.NewHeaderTenantResolveContributor(*hmtOpt, r),
						shttp.NewQueryTenantResolveContributor(*hmtOpt, r),
					}
					if hmtOpt.DomainFormat != "" {
						df := append(df[:1], df[0:]...)
						df[0] = shttp.NewDomainTenantResolveContributor(*hmtOpt, r, hmtOpt.DomainFormat)
					}
					trOpt.AppendContributors(df...)
				} else {
					trOpt.AppendContributors(NewHeaderTenantResolveContributor(*hmtOpt, tr))
				}
				if trOptF != nil {
					//patch
					trOptF(trOpt)
				}

				//get tenant config
				tenantConfigProvider := common.NewDefaultTenantConfigProvider(common.NewDefaultTenantResolver(*trOpt), ts)
				tenantConfig, trCtx, err := tenantConfigProvider.Get(ctx, true)
				if err != nil {
					//not found
					return nil, errors.NotFound("TENANT",err.Error())
				}
				//set current tenant
				currentTenant := common.ContextCurrentTenant{}
				newContext, _ := currentTenant.Change(trCtx, tenantConfig.ID, tenantConfig.Name)
				//data filter
				dataFilterCtx := data.NewEnableMultiTenancyDataFilter(newContext)

				return handler(dataFilterCtx, req)
			}
			return handler(ctx, req)
		}
	}
}
