package gorm

import "gorm.io/gorm"

type DialectFunc func(s string)gorm.Dialector

// Config 配置参数
type Config struct {
	Debug        bool
	//Dialect
	Dialect      DialectFunc
	Cfg *gorm.Config
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
}