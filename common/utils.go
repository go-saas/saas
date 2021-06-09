package common

import (
	"context"
	"github.com/jinzhu/copier"
)

func GetMultiTenantSide(ctx context.Context, ct CurrentTenant) MultiTenancySide {
	if ct.Id(ctx) == "" {
		return Host
	} else {
		return Tenant
	}
}

// Copy deeply
func Copy(s, ts interface{}) error {
	return copier.Copy(ts, s)
}
