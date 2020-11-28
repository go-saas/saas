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
	assert.Equal(t,TestDbProvider.createCounter,uint32(1))
	TestDbProvider.Get(ctx,data.Default)

	ctx = common.NewCurrentTenant(ctx,TenantId1,"Test1")

	TestDbProvider.Get(ctx,data.Default)


	assert.Equal(t,TestDbProvider.createCounter,uint32(1))

}
