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

type MemoryTenantStore struct {
	TenantConfig []TenantConfig
}

var _ TenantStore = (*MemoryTenantStore)(nil)

func NewMemoryTenantStore(t []TenantConfig) *MemoryTenantStore {
	return &MemoryTenantStore{
		TenantConfig: t,
	}
}

func (m MemoryTenantStore) GetByNameOrId(_ context.Context, nameOrId string) (*TenantConfig, error) {
	for _, config := range m.TenantConfig {
		if config.ID == nameOrId || config.Name == nameOrId {
			return &config, nil
		}
	}
	return nil, ErrTenantNotFound
}

//AlwaysTrustedAsIdTenantStore get the nameOrId as trusted id. usually used behind gateway
type AlwaysTrustedAsIdTenantStore struct {
}

var _ TenantStore = (*AlwaysTrustedAsIdTenantStore)(nil)

func (t *AlwaysTrustedAsIdTenantStore) GetByNameOrId(ctx context.Context, nameOrId string) (*TenantConfig, error) {
	return &TenantConfig{
		ID: nameOrId,
	}, nil
}

type CachedTenantStore struct {
	up TenantStore
}

func NewCachedTenantStore(up TenantStore) *CachedTenantStore {
	return &CachedTenantStore{
		up: up,
	}
}

var _ TenantStore = (*CachedTenantStore)(nil)

func (c *CachedTenantStore) GetByNameOrId(ctx context.Context, nameOrId string) (*TenantConfig, error) {
	cfg, ok := FromTenantConfigContext(ctx, nameOrId)
	if ok {
		return cfg, nil
	}
	return c.up.GetByNameOrId(ctx, nameOrId)
}
