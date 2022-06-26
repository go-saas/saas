package saas

import "github.com/go-saas/saas/data"

type TenantConfig struct {
	ID     string           `json:"id"`
	Name   string           `json:"name"`
	Region string           `json:"region"`
	Conn   data.ConnStrings `json:"conn"`
}

func NewTenantConfig(id string, name string, region string) *TenantConfig {
	return &TenantConfig{
		ID:     id,
		Name:   name,
		Region: region,
		Conn:   make(data.ConnStrings),
	}
}
