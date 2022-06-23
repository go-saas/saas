package main

import (
	"context"
	"errors"
	"github.com/goxiaoy/go-saas"

	sgorm "github.com/goxiaoy/go-saas/gorm"
	"gorm.io/gorm"
	"time"
)

type Tenant struct {
	ID string `gorm:"type:char(36)" json:"id"`
	//unique name. usually for domain name
	Name string `gorm:"column:name;index;size:255;"`
	//localed display name
	DisplayName string `gorm:"column:display_name;index;size:255;"`
	//region of this tenant
	Region    string `gorm:"column:region;index;size:255;"`
	Logo      string
	CreatedAt time.Time      `gorm:"column:created_at;index;"`
	UpdatedAt time.Time      `gorm:"column:updated_at;index;"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index;"`

	//connection
	Conn []TenantConn `gorm:"foreignKey:TenantId"`
}

// TenantConn connection string info
type TenantConn struct {
	TenantId string `gorm:"column:tenant_id;primary_key;size:36;"`
	//key of connection string
	Key string `gorm:"column:key;primary_key;size:100;"`
	//connection string
	Value     string    `gorm:"column:value;size:1000;"`
	CreatedAt time.Time `gorm:"column:created_at;index;"`
	UpdatedAt time.Time `gorm:"column:updated_at;index;"`
}

type TenantStore struct {
	dbProvider sgorm.DbProvider
}

func (t *TenantStore) GetByNameOrId(ctx context.Context, nameOrId string) (*saas.TenantConfig, error) {
	//change to host side
	ctx = saas.NewCurrentTenant(ctx, "", "")
	db := t.dbProvider.Get(ctx, "")
	var tenant Tenant
	err := db.Model(&Tenant{}).Preload("Conn").Where("id = ? OR name = ?", nameOrId, nameOrId).First(&tenant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, saas.ErrTenantNotFound
		} else {
			return nil, err
		}
	}
	ret := saas.NewTenantConfig(tenant.ID, tenant.Name, tenant.Region)
	for _, conn := range tenant.Conn {
		ret.Conn[conn.Key] = conn.Value
	}
	return ret, nil
}

var _ saas.TenantStore = (*TenantStore)(nil)
