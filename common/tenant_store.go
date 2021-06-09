package common

import (
	"context"
	"errors"
)

var ErrTenantNotFound = errors.New("tenant not found")

type TenantStore interface {
	GetByNameOrId(ctx context.Context, nameOrId string) (*TenantConfig, error)
}
