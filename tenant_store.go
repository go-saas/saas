package saas

import (
	"context"
	"errors"
)

var (
	ErrTenantNotFound = errors.New("tenant not found")
)

type TenantStore interface {
	// GetByNameOrId return nil and ErrTenantNotFound if tenant not found
	GetByNameOrId(ctx context.Context, nameOrId string) (*TenantConfig, error)
}

type MemoryTenantStore struct {
	TenantConfig []TenantConfig
}

var _ TenantStore = (*MemoryTenantStore)(nil)

func NewMemoryTenantStore(t []TenantConfig) *MemoryTenantStore {
	return &MemoryTenantStore{
		TenantConfig: t,
	}
}

func (m *MemoryTenantStore) GetByNameOrId(_ context.Context, nameOrId string) (*TenantConfig, error) {
	for _, config := range m.TenantConfig {
		if config.ID == nameOrId || config.Name == nameOrId {
			return &config, nil
		}
	}
	return nil, ErrTenantNotFound
}
