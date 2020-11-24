package common

import "context"

type TenantResolver interface {
	Resolve(ctx context.Context) TenantResolveResult
}
