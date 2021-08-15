package entity

import "time"

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
