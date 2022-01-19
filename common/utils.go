package common

import (
	"context"
)

func GetMultiTenantSide(ctx context.Context) MultiTenancySide {
	tenantInfo := FromCurrentTenant(ctx)
	if tenantInfo.GetId() == "" {
		return Host
	} else {
		return Tenant
	}
}
