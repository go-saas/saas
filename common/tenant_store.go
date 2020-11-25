package common

type TenantStore interface {
	
	getByNameOrId(nameOrId string)(*TenantConfig,error)

}
