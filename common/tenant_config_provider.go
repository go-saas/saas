package common

import "context"

type TenantConfigProvider interface {
	//get tenant config
	get(ctx context.Context) TenantConfig
}