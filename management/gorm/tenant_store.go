package gorm

import (
	"context"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/management/domain"
)

type TenantStore struct {
	tr domain.TenantRepo
}

func NewGormTenantStore(tr domain.TenantRepo) *TenantStore {
	return &TenantStore{
		tr: tr,
	}
}

func (g TenantStore) GetByNameOrId(_ context.Context, nameOrId string) (*common.TenantConfig, error) {
	//change to host side
	newCtx := common.NewCurrentTenant(context.Background(), "", "")
	t, err := g.tr.FindByIdOrName(newCtx, nameOrId)
	if err != nil {
		return nil, err
	}
	ret := common.NewTenantConfig(t.ID, t.Name, t.Region)
	for _, conn := range t.Conn {
		ret.Conn[conn.Key] = conn.Value
	}
	return ret, nil

}
