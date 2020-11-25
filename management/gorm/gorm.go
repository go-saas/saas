package gorm

import (
	"context"
	"github.com/goxiaoy/go-saas/gorm"
	g "gorm.io/gorm"
)

const ConnKey = "Tenant"

func GetDb(ctx context.Context,provider gorm.DbProvider) *g.DB {
	return provider.Get(ctx,ConnKey)
}
