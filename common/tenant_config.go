package common

import "github.com/goxiaoy/go-saas/data"

type TenantConfig struct {
	//Id of this tenant
	Id string
	Name string
	//Connection string map
	Conn data.ConnStrings
}

func NewTenantConfig(id string,name string) *TenantConfig {
	return &TenantConfig{
		Id:   id,
		Name: name,
		Conn: make(data.ConnStrings),
	}
}