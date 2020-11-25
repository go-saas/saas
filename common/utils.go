package common

import (
	"context"
	"github.com/jinzhu/copier"
)


func getMultiTenantSide(ctx context.Context,ct CurrentTenant)MultiTenancySide  {
	if ct.Id(ctx) ==""{
		return Host
	}else{
		return Tenant
	}
}

// Copy
func Copy(s, ts interface{}) error {
	return copier.Copy(ts, s)
}
