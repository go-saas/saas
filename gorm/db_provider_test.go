package gorm

import (
	"context"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"testing"
)

func TestDefaultDbProvider_Get(t *testing.T) {

	ctx := context.Background()

	TestDbProvider.Get(ctx, data.Default)

	ctx = common.NewCurrentTenant(ctx, TenantId1, "Test1")

	TestDbProvider.Get(ctx, data.Default)

	ctx = common.NewCurrentTenant(ctx, TenantId2, "Test2")

	TestDbProvider.Get(ctx, data.Default)

}
