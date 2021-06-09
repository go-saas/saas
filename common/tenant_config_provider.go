package common

import "context"

type TenantConfigProvider interface {
	// Get tenant config
	Get(ctx context.Context, store bool) (TenantConfig, context.Context, error)
}
