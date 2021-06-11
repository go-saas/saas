package gorm

import (
	"context"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"github.com/stretchr/testify/assert"
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
func TestDefaultDbProvider_UOW_Get(t *testing.T) {
	ctx := context.Background()
	TestUnitOfWorkManager.WithNew(ctx, func(ctx context.Context) error {

		dbA1 := TestDbProvider.Get(ctx, data.Default)

		ctx = common.NewCurrentTenant(ctx, TenantId1, "Test1")

		dbA2 := TestDbProvider.Get(ctx, data.Default)
		assert.Equal(t, dbA1, dbA2)

		ctx = common.NewCurrentTenant(ctx, TenantId2, "Test2")

		dbA3 := TestDbProvider.Get(ctx, data.Default)

		assert.NotEqual(t, dbA2, dbA3)

		////nested
		//TestUnitOfWorkManager.WithNew(ctx, func(ctx context.Context) error {
		//
		//	dbB1 := TestDbProvider.Get(ctx, data.Default)
		//
		//	ctx = common.NewCurrentTenant(ctx, TenantId1, "Test1")
		//
		//	dbB2 := TestDbProvider.Get(ctx, data.Default)
		//
		//	assert.NotEqual(t, dbA1, dbB1)
		//	assert.Equal(t, dbB1, dbB2)
		//
		//	ctx = common.NewCurrentTenant(ctx, TenantId2, "Test2")
		//
		//	dbB3 := TestDbProvider.Get(ctx, data.Default)
		//	assert.NotEqual(t, dbB2, dbB3)
		//	return nil
		//})
		return nil
	})

}
