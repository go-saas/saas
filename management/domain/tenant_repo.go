package domain

import (
	"context"
	"github.com/goxiaoy/go-saas/common"
)

type TenantRepo interface {
	//Find by id or name
	FindByIdOrName(ctx context.Context,idOrName string)(*Tenant,error)
	//total count
	GetCount(ctx context.Context,)(int64,error)
	//get paged list
	GetPaged(ctx context.Context,p common.Pagination)(c int64,t []*Tenant,err error)
	//create a tenant
	Create(ctx context.Context,t Tenant)error
	//update a tenant
	Update(ctx context.Context,id string,t Tenant)error
}
