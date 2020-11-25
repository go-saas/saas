package gorm

import (
	"context"
	"github.com/google/uuid"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/go-saas/management/domain"
	e "github.com/goxiaoy/go-saas/management/gorm/entity"
	g "gorm.io/gorm"
)

type GormTenantRepo struct {
	DbProvider gorm.DbProvider
}
func (g *GormTenantRepo)Db(ctx context.Context) *g.DB {
	return GetDb(ctx,g.DbProvider)
}

func (g *GormTenantRepo) FindByIdOrName(ctx context.Context, idOrName string) (*domain.Tenant,error) {
	var t *domain.Tenant
	var tDb *e.Tenant
	//parse
	if idOrName==""{
		return t,nil
	}
	//parse uuid
	id,err := uuid.Parse(idOrName)
	if err!=nil{
		//id
		g.Db(ctx).Where("id = ?",id.String()).First(&tDb)
	}else{
		g.Db(ctx).Where("name = ?",id.String()).First(&tDb)
	}
	common.Copy(tDb,t)
	return t,nil
}

func (g *GormTenantRepo) GetCount(ctx context.Context) (int,error) {
	panic("implement me")
}

func (g *GormTenantRepo) GetPaged(ctx context.Context, p common.Pagination) (c int64, t []*domain.Tenant,err error) {
	panic("implement me")
}

func (g *GormTenantRepo) Create(ctx context.Context, t domain.Tenant)error {
	var tDb *e.Tenant
	common.Copy(&t,tDb)
	ret := g.Db(ctx).Create(tDb)
	return ret.Error

}

func (g *GormTenantRepo) Update(ctx context.Context, id string, t domain.Tenant)error {
	panic("implement me")
}




