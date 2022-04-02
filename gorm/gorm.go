package gorm

import (
	"gorm.io/gorm"
)

type DialectFunc func(s string) gorm.Dialector
type EnsureDbExistFunc func(cfg *Config, s string) error

// Config 配置参数
type Config struct {
	Debug         bool
	Dialect       DialectFunc
	EnsureDbExist EnsureDbExistFunc
	Cfg           *gorm.Config
	MaxLifetime   *int
	MaxOpenConn   *int
	MaxIdleConn   *int
}

// MultiTenancy entity
type MultiTenancy struct {
	TenantId HasTenant `gorm:"index"`
}
