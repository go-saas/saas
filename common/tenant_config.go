package common

import "github.com/goxiaoy/go-saas/data"

type TenantConfig struct {
	ID     string
	Name   string
	Region string
	Conn data.ConnStrings
}

func NewTenantConfig(id string, name string, region string) *TenantConfig {
	return &TenantConfig{
		ID:     id,
		Name:   name,
		Region: region,
		Conn:   make(data.ConnStrings),
	}
}
