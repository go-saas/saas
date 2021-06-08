package gorm

import (
	"context"
	"gorm.io/gorm"
)

type DbProvider interface {
	// Get gorm db instance by key
	Get(ctx context.Context, key string) *gorm.DB
}
