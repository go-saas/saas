package gorm

import "github.com/goxiaoy/go-saas/common"

type GormTenantStore struct {

}

func (g GormTenantStore) GetByNameOrId(nameOrId string) (*common.TenantConfig, error) {
	//TODO
	panic("implement me")
}
