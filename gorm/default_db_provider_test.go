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
	assert.Equal(t, uint32(1), TestDbProvider.createCounter)
	TestDbProvider.Get(ctx, data.Default)

	ctx = common.NewCurrentTenant(ctx, TenantId1, "Test1")

	TestDbProvider.Get(ctx, data.Default)

	assert.Equal(t, uint32(1), TestDbProvider.createCounter)

	ctx = common.NewCurrentTenant(ctx, TenantId2, "Test2")

	TestDbProvider.Get(ctx, data.Default)
	assert.Equal(t, uint32(2), TestDbProvider.createCounter)

}
