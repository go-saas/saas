package gorm

import (
	"context"
	"gorm.io/gorm"
)

type DbProvider interface {
	//get db
	Get(ctx context.Context, key string) *gorm.DB
}
