package saas

import (
	"context"
)

func GetMultiTenantSide(ctx context.Context) MultiTenancySide {
	tenantInfo, _ := FromCurrentTenant(ctx)
	if tenantInfo.GetId() == "" {
		return Host
	} else {
		return Tenant
	}
}
