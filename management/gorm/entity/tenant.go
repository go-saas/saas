package entity

import "time"

type Tenant struct {
	// unique id
	ID string `gorm:"column:id;primary_key;size:36;"`
	//unique name. usually for domain name
	Name string `gorm:"column:name;index;size:255;"`
	//localed display name
	DisplayName string `gorm:"column:display_name;index;size:255;"`
	//region of this tenant
	Region    string     `gorm:"column:region;index;size:255;"`
	CreatedAt time.Time  `gorm:"column:created_at;index;"`
	UpdatedAt time.Time  `gorm:"column:updated_at;index;"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index;"`

	//connection
	Conn []TenantConn `gorm:"foreignKey:TenantId"`
	//edition
	Features []TenantFeature `gorm:"foreignKey:TenantId"`
}

type Tenants []*Tenant
