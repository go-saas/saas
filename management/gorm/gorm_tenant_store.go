package gorm

import (
	"context"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/management/domain"
)

type GormTenantStore struct {
	tr domain.TenantRepo
}
func NewGormTenantStore(tr domain.TenantRepo)*GormTenantStore{
	return &GormTenantStore{
		tr: tr,
	}
}

func (g GormTenantStore) GetByNameOrId(_ context.Context,nameOrId string) (*common.TenantConfig, error) {
	//change to host side
	newCtx := common.NewCurrentTenant(context.Background(),"","")
	t,err :=g.tr.FindByIdOrName(newCtx,nameOrId)
	if err!=nil{
		return nil, err
	}
	ret:=common.NewTenantConfig(t.ID,t.Name,t.Region)
	for _, conn := range t.Conn {
		ret.Conn[conn.Key]=conn.Value
	}
	return ret,nil

}
