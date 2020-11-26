package common

import "context"

type TenantStore interface {
	
	GetByNameOrId(ctx context.Context,nameOrId string)(*TenantConfig,error)

}
