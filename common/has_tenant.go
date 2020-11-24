package common

type HasTenant struct {
	//Mark tenant owner. zero value means Host
	TenantId string
}