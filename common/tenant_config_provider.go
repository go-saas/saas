package common

import "context"

type TenantConfigProvider interface {
	//get tenant config
	Get(ctx context.Context) (TenantConfig,error)
}