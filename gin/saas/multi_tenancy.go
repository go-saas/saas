package saas

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/common/http"
	"github.com/goxiaoy/go-saas/data"
)

func MultiTenancy(hmtOpt *http.WebMultiTenancyOption, ts common.TenantStore, trOptF ...common.PatchTenantResolveOption) gin.HandlerFunc {
	return func(context *gin.Context) {
		df := []common.TenantResolveContributor{
			//TODO route
			http.NewCookieTenantResolveContributor(*hmtOpt, context.Request),
			http.NewFormTenantResolveContributor(*hmtOpt, context.Request),
			http.NewHeaderTenantResolveContributor(*hmtOpt, context.Request),
			http.NewQueryTenantResolveContributor(*hmtOpt, context.Request),
		}
		if hmtOpt.DomainFormat != "" {
			df := append(df[:1], df[0:]...)
			df[0] = http.NewDomainTenantResolveContributor(*hmtOpt, context.Request, hmtOpt.DomainFormat)
		}
		trOpt := common.NewTenantResolveOption(df...)
		for _, option := range trOptF {
			option(trOpt)
		}
		//get tenant config
		tenantConfigProvider := common.NewDefaultTenantConfigProvider(common.NewDefaultTenantResolver(*trOpt), ts)
		tenantConfig, trCtx, err := tenantConfigProvider.Get(context, true)
		if err != nil {
			//not found
			if errors.Is(err, common.ErrTenantNotFound) {
				context.AbortWithError(404, err)
			} else {
				context.AbortWithError(500, err)
			}
		}
		previousTenant := common.FromCurrentTenant(trCtx)
		//set current tenant
		newContext := common.NewCurrentTenant(trCtx, tenantConfig.ID, tenantConfig.Name)
		//data filter
		dataFilterCtx := data.NewEnableMultiTenancyDataFilter(newContext)

		//with newContext
		context.Request = context.Request.WithContext(dataFilterCtx)
		//next
		context.Next()
		//cancel
		context.Request = context.Request.WithContext(common.NewCurrentTenantInfo(dataFilterCtx, previousTenant))
	}
}
