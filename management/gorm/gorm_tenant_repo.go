package gorm

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/go-saas/management/domain"
	e "github.com/goxiaoy/go-saas/management/gorm/entity"
	gg "gorm.io/gorm"
)

type GormTenantRepo struct {
	DbProvider gorm.DbProvider
}

func (g *GormTenantRepo) Db(ctx context.Context, preload bool) *gg.DB {
	ret := GetDb(ctx, g.DbProvider)
	if preload {
		ret = ret.Preload("Conn").Preload("Features")
	}
	return ret
}

func (g *GormTenantRepo) FindByIdOrName(ctx context.Context, idOrName string) (*domain.Tenant, error) {
	var t = new(domain.Tenant)
	var tDb e.Tenant
	//parse
	if idOrName == "" {
		return t, nil
	}
	//parse uuid
	id, err := uuid.Parse(idOrName)
	if err == nil {
		//id
		err = g.Db(ctx, true).Where("id = ?", id.String()).First(&tDb).Error
	} else {
		err = g.Db(ctx, true).Where("name = ?", idOrName).First(&tDb).Error
	}
	if err != nil {
		if errors.Is(err, gg.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	err = common.Copy(&tDb, t)
	return t, err
}

func (g *GormTenantRepo) GetCount(ctx context.Context) (int64, error) {
	var count int64
	//check count
	tx := g.Db(ctx, false).Model(&e.Tenant{}).Count(&count)
	return count, tx.Error
}

func (g *GormTenantRepo) GetPaged(ctx context.Context, p common.Pagination) (c int64, t []*domain.Tenant, err error) {
	err = g.Db(ctx, false).Model(&e.Tenant{}).Count(&c).Error
	var tDb e.Tenants
	if err != nil {
		return c, nil, err
	}
	err = gorm.BuildPage(g.Db(ctx, false), p).Find(&tDb).Error
	//copy
	common.Copy(tDb, &t)
	return
}

func (g *GormTenantRepo) Create(ctx context.Context, t domain.Tenant) error {
	var tDb = new(e.Tenant)
	common.Copy(&t, tDb)
	d := g.Db(ctx, true)
	ret := d.Create(tDb)
	return ret.Error

}

func (g *GormTenantRepo) Update(ctx context.Context, id string, t domain.Tenant) error {
	var tDb = new(e.Tenant)
	common.Copy(&t, tDb)
	d := g.Db(ctx, true)
	return d.Model(&e.Tenant{}).Where("id = ?", id).Updates(tDb).Error
}
