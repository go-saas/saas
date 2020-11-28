package common

import "github.com/goxiaoy/go-saas/data"

type TenantConfig struct {
	//ID of this tenant
	ID     string
	Name   string
	Region string
	//Connection string map
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
