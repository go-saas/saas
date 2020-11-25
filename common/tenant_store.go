package common

type TenantStore interface {
	
	GetByNameOrId(nameOrId string)(*TenantConfig,error)

}
