package common

import "context"

func getMultiTenantSide(ctx context.Context,ct CurrentTenant)MultiTenancySide  {
	if ct.Id(ctx) ==""{
		return Host
	}else{
		return Tenant
	}
}