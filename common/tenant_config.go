package common

type ConnectionStrings map[string]string

const defaultKey = "Default"

type TenantConfig struct {
	//Id of this tenant
	Id string
	Name string
	//Connection string map
	ConnectionString ConnectionStrings
}

func (cfg *TenantConfig)SetDefault(value string)  {
	cfg.ConnectionString[defaultKey] = value
}

func (cfg *TenantConfig)Default() string {
	return cfg.ConnectionString[defaultKey]
}