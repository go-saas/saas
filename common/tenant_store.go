package common

import (
	"context"
	"errors"
)

var ErrTenantNotFound = errors.New("tenant not found")

type TenantStore interface {
	// GetByNameOrId return nil and ErrTenantNotFound if tenant not found
	GetByNameOrId(ctx context.Context, nameOrId string) (*TenantConfig, error)
}
