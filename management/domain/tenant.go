package domain

import (
	"context"
	"github.com/goxiaoy/go-saas/common"
	"time"
)

type Tenant struct {
	// unique id
	ID string
	//unique name. usually for domain name
	Name string
	//localed display name
	DisplayName string
	//region of this tenant. Useful for data storage location
	Region    string
	CreatedAt time.Time
	UpdatedAt time.Time
	//should apply soft delete
	DeletedAt *time.Time
	//connection
	Conn []TenantConn
	//edition
	Features []TenantFeature
}

type TenantRepo interface {
	FindByIdOrName(ctx context.Context, idOrName string) (*Tenant, error)
	GetCount(ctx context.Context) (int64, error)
	GetPaged(ctx context.Context, p common.Pagination) (c int64, t []*Tenant, err error)
	Create(ctx context.Context, t Tenant) error
	Update(ctx context.Context, id string, t Tenant) error
}
