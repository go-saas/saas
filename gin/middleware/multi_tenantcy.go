package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/common/http"
)

func MultiTenancy(hmtOptF http.PatchHttpMultiTenancyOption,trOptF common.PatchTenantResolveOption,ts common.TenantStore) gin.HandlerFunc {
	return func(context *gin.Context) {

		hmtOpt := http.DefaultWebMultiTenancyOption()
		if hmtOptF != nil{
			//patch
			hmtOptF(hmtOpt)
		}
		trOpt := common.NewTenantResolveOption(
			//TODO route
			http.NewCookieTenantResolveContributor(*hmtOpt,context.Request),
			http.NewFormTenantResolveContributor(*hmtOpt,context.Request),
			http.NewHeaderTenantResolveContributor(*hmtOpt,context.Request),
			http.NewQueryTenantResolveContributor(*hmtOpt,context.Request),
			)
		if trOptF != nil{
			//patch
			trOptF(trOpt)
		}
		//get tenant config
		tenantConfigProvider := common.NewDefaultTenantConfigProvider(common.NewDefaultTenantResolver(*trOpt),ts)
		tenantConfig,trCtx,err := tenantConfigProvider.Get(context,true)
		if err!=nil{
			//not found
			context.AbortWithError(404,err)
		}
		//set current tenant
		currentTenant :=common.ContextCurrentTenant{}
		newContext,cancel := currentTenant.Change(trCtx,tenantConfig.Id,tenantConfig.Name)
		//cancel
		defer cancel()
		//with newContext
		context.Request=context.Request.WithContext(newContext)
		//next
		context.Next()
	}
}
