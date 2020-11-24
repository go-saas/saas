package common

type BasicTenantInfo struct {
	Id string
	Name string
}

func NewBasicTenantInfo(id string, name string) *BasicTenantInfo {
	return &BasicTenantInfo{Id: id, Name: name}
}


