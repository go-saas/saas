package gorm

import (
	"gorm.io/gorm"
)

const (
	MultiTenantBeforeCreateName = "multi_tenancy:before_create"
	MultiTenantQueryName        = "multi_tenancy:query"
)

//Auto set current tenant value before create
// Deprecated: application should set value self
func AutoSetTenant(db *gorm.DB) {
	//t := common.FromCurrentTenant(db.Statement.Context)
	//tId := t.ID
}

// Deprecated: Use custom data type
func AutoFilterCurrentTenant(db *gorm.DB) {
	//t := common.FromCurrentTenant(db.Statement.Context)
	//tId := t.ID
}
