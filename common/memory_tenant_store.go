package common

import "errors"

type MemoryTenantStore struct {
	TenantConfig []TenantConfig
}

func NewMemoryTenantStore(t []TenantConfig) *MemoryTenantStore {
	return &MemoryTenantStore{
		TenantConfig: t,
	}
}

func (m MemoryTenantStore) getByNameOrId(nameOrId string) (*TenantConfig, error) {
	for _, config := range m.TenantConfig {
		if config.Id==nameOrId||config.Name==nameOrId {
			return &config,nil
		}
	}
	return nil,errors.New("tenant not found")
}

