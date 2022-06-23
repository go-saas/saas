package saas

type TenantInfo interface {
	GetId() string
	GetName() string
}

type BasicTenantInfo struct {
	Id   string
	Name string
}

func (b *BasicTenantInfo) GetId() string {
	return b.Id
}

func (b *BasicTenantInfo) GetName() string {
	return b.Name
}

func NewBasicTenantInfo(id string, name string) *BasicTenantInfo {
	return &BasicTenantInfo{Id: id, Name: name}
}
